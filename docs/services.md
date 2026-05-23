# Services Registry

Source of truth for service ports, databases, and bounded-context responsibilities. See `ARCHITECTURE.md` for the why behind the partitioning. Keep this table in sync with `docker-compose.yml`, `infra/init-databases.sql`, and `kafka/topics.json`.

## Service registry

| Service              | Port | Database            | Bounded context / responsibility |
|----------------------|------|---------------------|----------------------------------|
| service-booking      | 8001 | `kilat_booking`     | Booking aggregate, runner-flow state machine. Owns transport requests + shop-side substates (Plan C Phase 5). |
| service-payment      | 8002 | `kilat_payment`     | Escrow saga, Stripe ACL. Owns shop wallet + withdrawals (Plan C Phase 7). |
| service-runner       | 8003 | `kilat_runner`      | Runner availability + GPS potential. PostGIS proximity queries. |
| service-identity     | 8004 | `kilat_identity`    | OIDC/OAuth2 STS. Owns shop merchant roles (Plan C Phase 6). |
| service-tracking     | 8005 | `kilat_tracking`    | Real-time GPS feed for owner + runner over WebSocket. |
| service-notification | 8006 | `kilat_notification`| Multi-channel notifications. Owns shop categories (Plan C Phase 8). |
| service-review       | 8007 | `kilat_review`      | Post-trip review and rating. |
| service-shop         | 8012 | `kilat_shop`        | Shop / Product / Inventory / SalesAggregate aggregates (Plan C, scaffolded in Phase 1). |
| api-gateway          | 8080 | n/a                 | BFF for mobile + web clients. Plan C Phase 9 adds shop routes. |

`postgres` listens on host port `5433` (mapped from container port `5432`) and hosts every `kilat_*` database. `kafka` listens on `9092` (host) / `29092` (in-network).

## Kafka topics

See `kafka/topics.json` for the declarative registry (partitions, retention). Quick reference:

| Topic              | Partitions | Retention | Partition key  | Plan |
|--------------------|------------|-----------|----------------|------|
| `booking.events`   | 6          | 7 d       | `booking_id`   | A    |
| `payment.events`   | 6          | 14 d      | `payment_id` / `withdrawal_id` | A + C |
| `runner.events`    | 6          | 7 d       | `runner_id`    | A    |
| `tracking.events`  | 6          | 7 d       | `booking_id`   | A    |
| `chat.events`      | 6          | 7 d       | `thread_id`    | A    |
| `shop.events`      | 12         | 7 d       | `shop_id`      | C    |
| `inventory.events` | 6          | 7 d       | `product_id`   | C    |

## Adding a new service

1. Add a row to the table above.
2. Add a `CREATE DATABASE kilat_<service>;` + extension blocks to `infra/init-databases.sql`.
3. Add the service block to `docker-compose.yml` (port + env wiring).
4. Add the `UPSTREAM_<NAME>` env var to `api-gateway`.
5. If the service publishes new event types, add the topic to `kafka/topics.json`.
