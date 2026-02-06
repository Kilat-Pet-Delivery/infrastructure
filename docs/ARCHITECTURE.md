# Kilat Pet Runner — Architecture Document

> A "Grab for Pets" platform built with Domain-Driven Design (DDD) and Clean Architecture in Go.

---

## Table of Contents

1. [Business Logic (The "Why")](#1-business-logic-the-why)
2. [Microservices Architecture (The "How")](#2-microservices-architecture-the-how)
3. [User Experience & Workflow (The "Flow")](#3-user-experience--workflow-the-flow)
4. [Technical Stack & Integration](#4-technical-stack--integration)
5. [Risk Management & Safety (The "Critical Part")](#5-risk-management--safety-the-critical-part)
6. [Domain Event Catalog](#6-domain-event-catalog)
7. [API Endpoint Reference](#7-api-endpoint-reference)
8. [Kafka Topic Mapping](#8-kafka-topic-mapping)

---

## 1. Business Logic (The "Why")

In Domain-Driven Design (DDD), we must first define the **Domain**, which represents the specific subject matter and business expertise of the software. For a "Grab for Pets" platform, the **Core Domain** centers on the safe and reliable transfer of live animals between locations, which is a specialized form of logistics.

### Target Audience

The primary users are:
- **Pet Owners** (Resource Owners) requiring specialized transport
- **Pet Runners** (independent contractors) providing the capability

### Unique Value Proposition

Unlike standard logistics, this platform treats pets as **"Live Aggregates"** requiring specific environmental conditions (climate control, hydration, and medical oversight) rather than mere cargo. By using a **Ubiquitous Language** shared between developers and pet experts, the platform ensures that terms like "Crate Specification" or "Stress Monitoring" are rigorously defined in the code to eliminate ambiguity.

### Incentive Structure

Runners are incentivized through dynamic pricing models based on the **"Capability" layer** — matching the runner's vehicle equipment to the pet's specific needs. High-quality service is encouraged via a ranking system that directly impacts job allocation priority in the **Decision Support layer**.

---

## 2. Microservices Architecture (The "How")

To build a scalable system, we partition the application into **Bounded Contexts**, ensuring each microservice owns its domain data and logic. This prevents a "Big Ball of Mud" and enforces microservice autonomy.

### Bounded Contexts

| Service | Port | Database | Responsibility |
|---------|------|----------|---------------|
| **service-booking** | 8001 | kilat_booking | Manages the Order Aggregate. Responsible for transport request state and ensuring invariants (e.g., validating runner has correct crate size before accepting). |
| **service-payment** | 8002 | kilat_payment | Handles financial transactions using a Saga Pattern. Acts as escrow, holding funds until a "DeliveryConfirmed" Domain Event is received. |
| **service-runner** | 8003 | kilat_runner | Tracks runner availability and vehicle "Potential". Manages runner session lifecycle and proximity to pickup points via PostGIS. |
| **service-identity** | 8004 | kilat_identity | Centralized Security Token Service (STS) using OpenID Connect for authentication and OAuth 2.0 for authorizing API access. |
| **service-tracking** | 8005 | kilat_tracking | Consumes Integration Events from Kafka to provide real-time GPS updates to owners via WebSocket. |

### Why Clean Architecture?

We implement the **Handlers -> Services -> Repositories** pattern to isolate business logic from infrastructure concerns:

```
┌─────────────────────────────────────┐
│           Handler Layer             │  ← HTTP/WebSocket (Gin)
├─────────────────────────────────────┤
│         Application Layer           │  ← Use Cases, DTOs, Orchestration
├─────────────────────────────────────┤
│           Domain Layer              │  ← Entities, VOs, Events, Policies (PURE)
├─────────────────────────────────────┤
│        Infrastructure Layer         │  ← DB (GORM/PostgreSQL), Kafka, External APIs
└─────────────────────────────────────┘
```

The **Domain Layer** (Entities/Services) remains pure, while the **Infrastructure Layer** handles database persistence (PostgreSQL + PostGIS) and external API calls.

### Service Internal Structure

```
service-xxx/
├── cmd/server/main.go          # Entry point, DI wiring, graceful shutdown
├── internal/
│   ├── domain/                 # Pure domain: entities, VOs, events, repo interfaces, policies
│   ├── application/            # Use cases, application services, DTOs
│   ├── handler/                # HTTP handlers (Gin), request/response mapping
│   ├── repository/             # Concrete PostgreSQL/GORM implementations
│   ├── adapter/                # Anti-Corruption Layer (Maps, Stripe, Kafka)
│   ├── config/                 # Viper-based config loading
│   ├── middleware/             # Service-specific middleware
│   ├── events/                 # Kafka publisher/consumer wiring
│   └── saga/                   # Saga orchestrator (service-payment primarily)
├── migrations/                 # Numbered SQL migration files
├── Dockerfile                  # Multi-stage Alpine build
└── go.mod
```

---

## 3. User Experience & Workflow (The "Flow")

The user journey is driven by the state changes of the **Aggregate Root**.

### Booking State Machine

```
┌──────────┐    accept     ┌──────────┐    pickup    ┌────────────┐
│ REQUESTED ├──────────────►│ ACCEPTED ├─────────────►│ IN_PROGRESS │
└─────┬─────┘              └─────┬─────┘              └──────┬──────┘
      │                          │                           │
      │ cancel                   │ cancel                    │ deliver
      ▼                          ▼                           ▼
┌──────────┐              ┌──────────┐              ┌──────────┐
│ CANCELLED │              │ CANCELLED │              │ DELIVERED │
└──────────┘              └──────────┘              └─────┬─────┘
                                                          │ confirm
                                                          ▼
                                                    ┌──────────┐
                                                    │ COMPLETED │
                                                    └──────────┘
```

### Pet Owner Journey (Johor to Kuala Lumpur)

1. **Registration:** User authenticates via a BFF (Backend for Frontend) that delegates to the STS.
2. **Booking:** Owner enters the "Pet Specification" (Cat, breed, health status, vaccination records).
3. **Quote:** The Decision Support Layer uses a **Strategy Pattern** to calculate the best route and price.
4. **Monitoring:** Once the runner picks up the cat, the owner receives a stream of Domain Events via WebSockets or push notifications.

### Pet Runner Journey

1. **Availability:** Runner toggles their status to "Active," updating the Potential Layer.
2. **Job Selection:** The system matches the runner's vehicle capability to the cat's crate requirement.
3. **Execution:** The runner follows the Itinerary generated by the Routing Service.

### Critical Data Collection

To ensure safety, the app must collect a **Route Specification** that includes:

| Data | Purpose | Storage |
|------|---------|---------|
| **Pet Health Status** | Validation against transport safety rules | JSONB in bookings table |
| **Crate Size** | "Container Specification" fulfillment | JSONB in bookings table |
| **Vaccination Records** | Compliance checks | Immutable Value Objects |

---

## 4. Technical Stack & Integration

### Technology Choices

| Concern | Technology | Go Package |
|---------|-----------|------------|
| HTTP Framework | Gin | `github.com/gin-gonic/gin` |
| ORM | GORM | `gorm.io/gorm`, `gorm.io/driver/postgres` |
| UUID | google/uuid | `github.com/google/uuid` |
| Logging | Zap | `go.uber.org/zap` |
| Config | Viper | `github.com/spf13/viper` |
| JWT | golang-jwt | `github.com/golang-jwt/jwt/v5` |
| Event Bus | Kafka | `github.com/segmentio/kafka-go` |
| Geospatial | PostGIS | ST_DWithin, ST_Distance, geography type |
| WebSocket | gorilla/websocket | `github.com/gorilla/websocket` |
| Testing | testify | `github.com/stretchr/testify` |
| Migration | golang-migrate | `github.com/golang-migrate/migrate/v4` |

### Real-time GPS Tracking

Integrate the **Google Maps API** within an Adapter to translate external coordinates into our internal domain model. Runner locations are stored using PostGIS `GEOMETRY(Point, 4326)` with GIST spatial indexes for efficient proximity queries.

### Event-Driven Communication

Use **Apache Kafka** for event streaming between microservices to ensure eventual consistency. Each service publishes Domain Events to dedicated topics and consumes Integration Events from other services.

### Payment Gateway

Use an **Anti-Corruption Layer** to wrap Stripe, ensuring our domain logic remains independent of the payment provider's API. The escrow pattern uses Stripe's manual capture (`CaptureMethod=manual`) to hold funds until delivery confirmation.

### Repository Pattern (Go Example)

```go
// PetRepository defines the interface for persisting pet data.
type PetRepository interface {
    GetByID(ctx context.Context, id string) (*domain.Pet, error)
    Save(ctx context.Context, pet *domain.Pet) error
}

// SQLPetRepository implements persistence using a SQL database.
type SQLPetRepository struct {
    db *sql.DB
}

func (r *SQLPetRepository) Save(ctx context.Context, pet *domain.Pet) error {
    // Implement logic here...
    return nil
}
```

**Why?** Interfaces allow us to swap the database for a mock during unit testing or move from SQL to a NoSQL store without breaking the service layer.

---

## 5. Risk Management & Safety (The "Critical Part")

### Race Conditions in Job Booking

If two runners attempt to accept the same job simultaneously, the system may face a **Race Condition**, leading to data inconsistency.

**Solution:** Implement **Optimistic Locking** in the database. Each Aggregate Root has a `Version` field. If the version in the database differs from the version in memory during an update, the transaction fails and the user is prompted to retry.

```sql
-- Example optimistic lock check
UPDATE bookings
SET runner_id = $1, status = 'accepted', version = version + 1, updated_at = NOW()
WHERE id = $2 AND version = $3;
-- If rows affected = 0, another transaction already modified this record
```

### Safety Protocols

| Protocol | Implementation |
|----------|---------------|
| **Health Checks** | Every microservice implements `/health` and `/readiness` endpoints to ensure dependencies are ready before accepting traffic. |
| **Validation Invariants** | Entities are "always valid." A Booking must throw an error if the VaccinationDate is older than one year. |
| **Audit Logging** | All state changes are recorded in `booking_status_history` for forensic analysis in case of a transport incident. |

### Rating/Review System — Policy Pattern

```go
// QualityPolicy determines if a runner is eligible for high-tier jobs.
type QualityPolicy interface {
    IsEligible(runner *domain.Runner) bool
}

type FiveStarPolicy struct{}

func (p *FiveStarPolicy) IsEligible(runner *domain.Runner) bool {
    return runner.Rating >= 4.8 && runner.CompletionRate > 0.95
}
```

This isolates the rules for "Quality" from the runner's identity data, allowing the business to refine the rating algorithm without modifying the core Runner Entity.

---

## 6. Domain Event Catalog

### Booking Events (Topic: `booking.events`)

| Event | Trigger | Raised By | Consumed By |
|-------|---------|-----------|-------------|
| `booking.requested` | Owner creates booking | service-booking | notification |
| `booking.accepted` | Runner accepts booking | service-booking | service-tracking |
| `booking.pet_picked_up` | Runner picks up pet | service-booking | service-tracking |
| `booking.delivery_confirmed` | Owner confirms delivery | service-booking | service-payment |
| `booking.completed` | Payment released | service-booking | analytics |
| `booking.cancelled` | Owner/Runner cancels | service-booking | service-payment |

### Payment Events (Topic: `payment.events`)

| Event | Trigger | Raised By | Consumed By |
|-------|---------|-----------|-------------|
| `payment.escrow_created` | Payment initiated | service-payment | service-booking |
| `payment.escrow_held` | Stripe authorize success | service-payment | service-booking |
| `payment.escrow_released` | Delivery confirmed + capture | service-payment | service-booking |
| `payment.escrow_refunded` | Cancellation + refund | service-payment | service-booking |
| `payment.failed` | Stripe failure | service-payment | service-booking |

### Runner Events (Topic: `runner.events`)

| Event | Trigger | Raised By | Consumed By |
|-------|---------|-----------|-------------|
| `runner.online` | Runner activates session | service-runner | service-booking |
| `runner.offline` | Runner deactivates | service-runner | service-booking |
| `runner.location_update` | GPS ping (every 5-10s) | service-runner | service-tracking |

### Tracking Events (Topic: `tracking.events`)

| Event | Trigger | Raised By | Consumed By |
|-------|---------|-----------|-------------|
| `tracking.started` | Booking accepted | service-tracking | future: analytics |
| `tracking.updated` | Waypoint added | service-tracking | future: analytics |
| `tracking.completed` | Delivery confirmed | service-tracking | future: analytics |

---

## 7. API Endpoint Reference

### service-identity (Port 8004)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/register` | No | Register owner or runner |
| POST | `/api/v1/auth/login` | No | Login with email/password |
| POST | `/api/v1/auth/refresh` | No | Refresh JWT token pair |
| POST | `/api/v1/auth/logout` | Yes | Revoke refresh tokens |
| GET | `/api/v1/auth/profile` | Yes | Get current user profile |
| PUT | `/api/v1/auth/profile` | Yes | Update profile |

### service-runner (Port 8003)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/runners` | Runner | Register runner profile |
| GET | `/api/v1/runners/me` | Runner | Get own profile |
| PUT | `/api/v1/runners/me` | Runner | Update profile |
| POST | `/api/v1/runners/me/online` | Runner | Go online (with GPS coords) |
| POST | `/api/v1/runners/me/offline` | Runner | Go offline |
| POST | `/api/v1/runners/me/location` | Runner | Update GPS location |
| POST | `/api/v1/runners/me/crates` | Runner | Add crate specification |
| GET | `/api/v1/runners/nearby` | Internal | Find nearby active runners |

### service-booking (Port 8001)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/bookings` | Owner | Create new booking |
| GET | `/api/v1/bookings` | Owner/Runner | List bookings (role-filtered) |
| GET | `/api/v1/bookings/:id` | Owner/Runner | Get booking details |
| POST | `/api/v1/bookings/:id/accept` | Runner | Accept booking |
| POST | `/api/v1/bookings/:id/pickup` | Runner | Mark pet picked up |
| POST | `/api/v1/bookings/:id/deliver` | Runner | Mark delivered |
| POST | `/api/v1/bookings/:id/confirm` | Owner | Confirm delivery |
| POST | `/api/v1/bookings/:id/cancel` | Owner/Runner | Cancel booking |
| GET | `/api/v1/bookings/:id/runners` | Owner | List available runners |

### service-payment (Port 8002)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/payments/initiate` | Owner | Initiate escrow for booking |
| POST | `/api/v1/payments/webhook` | Stripe Sig | Stripe webhook handler |
| GET | `/api/v1/payments/:id` | Owner/Runner | Get payment details |
| GET | `/api/v1/payments/booking/:bookingId` | Owner/Runner | Get payment for booking |
| POST | `/api/v1/payments/:id/refund` | Admin | Manual refund |

### service-tracking (Port 8005)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/tracking/:bookingId` | Owner/Runner | Get tracking data |
| GET | `/api/v1/tracking/:bookingId/route` | Owner/Runner | Get full route as GeoJSON |
| WS | `/ws/tracking/:bookingId` | Owner (token) | Real-time WebSocket feed |

---

## 8. Kafka Topic Mapping

### Complete Event Flow — Booking Lifecycle

```
1. Owner creates booking
   service-booking ──► Kafka[booking.events] : BookingRequested

2. Owner initiates payment (escrow)
   service-payment ──► Stripe: CreatePaymentIntent (manual capture)
   service-payment ──► Kafka[payment.events] : EscrowHeld

3. Runner accepts booking
   service-booking ──► Kafka[booking.events] : BookingAccepted

4. Tracking starts
   service-tracking ◄── Kafka[booking.events] : BookingAccepted → creates TripTrack

5. Runner sends GPS updates
   service-runner ──► Kafka[runner.events] : RunnerLocationUpdate
   service-tracking ◄── Kafka[runner.events] → adds Waypoint + broadcasts via WebSocket

6. Runner marks pet picked up
   service-booking ──► Kafka[booking.events] : PetPickedUp

7. Runner marks delivered
   service-booking ──► Kafka[booking.events] : DeliveryConfirmed

8. Payment released to runner
   service-payment ◄── Kafka[booking.events] : DeliveryConfirmed
   service-payment ──► Stripe: CapturePaymentIntent
   service-payment ──► Kafka[payment.events] : EscrowReleased

9. Booking finalized
   service-booking ◄── Kafka[payment.events] : EscrowReleased → completes booking
```
