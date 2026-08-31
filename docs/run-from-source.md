# Running the platform from source (KPD-3)

How to bring up `api-gateway` + `service-identity` + `service-booking` against the shared dev-infra
stack and prove the path works end to end. Verified 2026-08-31.

## 1. Data stores

```bash
cd ~/Documents/dev-infra
./dev.ps1 up kilat
```

Gives Postgres+PostGIS on `localhost:5432` with the eight `kilat_*` databases (`kilat` / `kilat_secret`),
Kafka on `localhost:9092`, Redis, MinIO and Mailpit.

## 2. Ports — read this before you start

The Desa Murni Batik services claim **8001–8009**, and Docker holds **8080** for `niaga-dev-gateway`.
KPD's documented ports are the same range, so the two products cannot both run as documented. Whichever
starts second dies with `Only one usage of each socket address`.

Until **KPD-65** settles a permanent port split, run KPD on the `18xxx` range. Every port below is
overridable by environment variable, so nothing needs editing in code.

| Service | Documented | Used here |
|---|---|---|
| `api-gateway` | 8080 | **18080** |
| `service-booking` | 8001 | **18001** |
| `service-identity` | 8004 | **18004** |

Check first:

```bash
netstat -ano | grep LISTENING | grep -E ":1800[0-9]|:18080"
```

## 3. Start the services

Each service applies its own SQL migrations at startup, in every environment, so there is no separate
migrate step — see each repo's README. `JWT_SECRET` **must be identical** across all three, or the
gateway and booking will reject tokens that identity issued.

```bash
# service-identity
cd ~/Documents/kilat-pet-delivery/service-identity
export DB_HOST=localhost DB_PORT=5432 DB_USER=kilat DB_PASSWORD=kilat_secret DB_SSL_MODE=disable
export JWT_SECRET=dev-smoke-secret-change-me
DB_NAME=kilat_identity SERVICE_PORT=18004 APP_ENV=development \
  JWT_ACCESS_EXPIRY=15m JWT_REFRESH_EXPIRY=168h go run ./cmd/server

# service-booking
cd ~/Documents/kilat-pet-delivery/service-booking
DB_NAME=kilat_booking SERVICE_PORT=18001 APP_ENV=development \
  KAFKA_BROKERS=localhost:9092 KAFKA_GROUP_PREFIX=kilat- go run ./cmd/server

# api-gateway
cd ~/Documents/kilat-pet-delivery/api-gateway
GATEWAY_PORT=18080 APP_ENV=development \
  UPSTREAM_IDENTITY=http://localhost:18004 \
  UPSTREAM_BOOKING=http://localhost:18001 go run ./cmd/server
```

Confirm all three are up. The aggregator reports the services you did not start as `unhealthy` — that is
expected, not a failure:

```bash
curl -s http://localhost:18080/health
```

## 4. Smoke test — register, login, create a booking

All three calls go through the gateway, which is the point: it proves routing, JWT issuance and
cross-service auth, not just that the binaries boot.

```bash
EMAIL="smoke-$(date +%s)@kilat.test"

# register -> 201
curl -s -X POST http://localhost:18080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"phone\":\"+60123456789\",\"full_name\":\"Smoke Owner\",\"password\":\"SmokeTest123!\",\"role\":\"owner\"}"

# login -> 200, returns data.access_token
TOKEN=$(curl -s -X POST http://localhost:18080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"SmokeTest123!\"}" | jq -r .data.access_token)

# create booking -> 201 (owner role required)
curl -s -X POST http://localhost:18080/api/v1/bookings \
  -H 'Content-Type: application/json' -H "Authorization: Bearer $TOKEN" \
  -d '{
    "pet_spec": {"pet_type":"dog","breed":"Golden Retriever","name":"Rex","weight_kg":18.5,"age_months":36},
    "pickup_address": {"line1":"12 Jalan Ampang","city":"Kuala Lumpur","state":"Wilayah Persekutuan","postal_code":"50450","country":"MY","latitude":3.1590,"longitude":101.7123},
    "dropoff_address": {"line1":"88 Jalan Universiti","city":"Petaling Jaya","state":"Selangor","postal_code":"46200","country":"MY","latitude":3.1073,"longitude":101.6420},
    "notes": "smoke test"
  }'
```

Expected: `201`, a `booking_number` like `BK-7MC56E`, `status` `requested`, and a computed
`estimated_price_cents` in `MYR` (the run on 2026-08-31 priced this route at 4423).

Confirm it persisted:

```bash
docker exec dev-postgres psql -U kilat -d kilat_booking \
  -tAc "SELECT booking_number, status, estimated_price_cents FROM bookings;"
```

Clean up after yourself — the smoke test writes real rows:

```bash
docker exec dev-postgres psql -U kilat -d kilat_booking -tAc "DELETE FROM bookings WHERE booking_number='<number>';"
docker exec dev-postgres psql -U kilat -d kilat_identity -tAc "DELETE FROM refresh_tokens; DELETE FROM users WHERE email LIKE 'smoke-%@kilat.test';"
```

## 5. Known gap — the booking event is silently dropped

The booking is created and returns `201`, but its `booking.requested` event does **not** reach Kafka:

```
failed to publish message {"topic": "booking.events",
  "error": "[3] Unknown Topic Or Partition: ... does not exist on this broker"}
```

`lib-common/kafka/producer.go` builds its `kafka.Writer` without `AllowAutoTopicCreation`, so kafka-go
will not create a missing topic even though the broker allows it. `payment.events` exists only because
service-booking's *consumer* subscribes to it, and consumers do trigger creation. The producer then logs
the failure and carries on, so the request still succeeds and the event is lost — `service-notification`
and `service-payment` never hear about the booking.

Tracked as **KPD-64**. The HTTP path is genuinely working; the event path is not.
