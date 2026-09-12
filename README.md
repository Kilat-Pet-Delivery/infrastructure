# infrastructure

Docker Compose orchestration and shared infrastructure for Kilat Pet Runner microservices architecture.

**Organization:** github.com/Kilat-Pet-Delivery

## Repository Contents

- **docker-compose.yml** — Orchestration of 8 containers (PostgreSQL, PostGIS, Zookeeper, Kafka, 5 microservices)
- **infra/init-databases.sql** — Database initialization with 5 databases, PostGIS, and uuid-ossp extensions
- **docs/ARCHITECTURE.md** — Complete Domain-Driven Design architecture specification
- **Makefile** — Build, test, and Docker command utilities
- **.env.example** — Environment variable configuration template

## Quick Start

```bash
git clone https://github.com/Kilat-Pet-Delivery/infrastructure.git
cd infrastructure

# Clone all service repositories as siblings
git clone https://github.com/Kilat-Pet-Delivery/service-identity.git ../
git clone https://github.com/Kilat-Pet-Delivery/service-booking.git ../
git clone https://github.com/Kilat-Pet-Delivery/service-payment.git ../
git clone https://github.com/Kilat-Pet-Delivery/service-runner.git ../
git clone https://github.com/Kilat-Pet-Delivery/service-tracking.git ../

# Start all services
docker-compose up -d
```

## Microservices

| Service | Port | Repository |
|---------|------|------------|
| service-identity | 8004 | Kilat-Pet-Delivery/service-identity |
| service-booking | 8001 | Kilat-Pet-Delivery/service-booking |
| service-payment | 8002 | Kilat-Pet-Delivery/service-payment |
| service-runner | 8003 | Kilat-Pet-Delivery/service-runner |
| service-tracking | 8005 | Kilat-Pet-Delivery/service-tracking |

## Infrastructure Services

- **PostgreSQL 16 + PostGIS 3.4** (port 5433) — Geospatial database with UUID support
- **Apache Kafka** via Confluent (port 9092) — Event streaming and message broker
- **Zookeeper** (port 2181) — Distributed coordination service

## The two ways to run

There are exactly two supported setups. Mixing them is what causes port clashes and
"connection refused" against the wrong Postgres.

### 1. Daily loop — shared stack + `go run`

What you want almost always. The data stores come from the shared dev-infra stack; the Go
services run from source on the host, so a change is one `go run` away.

```bash
cd ~/Documents/dev-infra && ./dev.ps1 up kilat   # Postgres+PostGIS, Kafka, Redis, MinIO, Mailpit
cd ~/Documents/kilat-pet-delivery/service-identity
cp .env.example .env                              # already points at localhost
go run ./cmd/migrate                              # optional: cmd/server applies the same migrations at startup
go run ./cmd/server
```

Every service repo ships a `.env.example` already pointed at the shared stack —
`DB_HOST=localhost`, `DB_PORT=5432`, `KAFKA_BROKERS=localhost:9092`. Copy it to `.env` and
edit only what you actually need. `.env` is gitignored in every repo; never commit one.

`JWT_SECRET` must be **identical** in every service, or tokens issued by `service-identity` are
rejected by the others. The gateway does not read it: it only proxies, and each service checks
the token itself.

Step-by-step, including a working `register → login → create booking` smoke test:
[`docs/run-from-source.md`](docs/run-from-source.md).

### 2. Full-stack smoke — `docker-compose.yml` in this repo

Everything in containers, for checking the whole platform end to end rather than iterating.

```bash
cd ~/Documents/kilat-pet-delivery/infrastructure
cp .env.example .env
make up
```

**Stop the shared stack first** — this compose file publishes its own Postgres and Kafka and
will fight the shared one for ports. Note the difference: this file uses Postgres on **5433**
and hostnames like `db` and `service-identity`, where the daily loop uses **5432** on
`localhost`. That is the single biggest source of confusion between the two modes.

### Ports, and why yours may not bind

The documented service ports are 8001–8007 plus 8012 and the gateway on 8080. On the shared
development laptop the Desa Murni Batik services already occupy 8001–8009, and Docker holds
8080. Whichever product starts second fails with
`Only one usage of each socket address (protocol/network address/port) is normally permitted`.

Until **KPD-65** settles a permanent split, override per service — `SERVICE_PORT=18004`,
`GATEWAY_PORT=18080`. Every port is environment-driven, so nothing in code needs editing.

## Configuration

Copy `.env.example` to `.env` and adjust environment variables for your deployment:

```bash
cp .env.example .env
```

## Development Commands

```bash
make build       # Build all Docker images
make up          # Start containers
make down        # Stop containers
make logs        # View service logs
make test        # Run tests
```

## Architecture

Refer to `docs/ARCHITECTURE.md` for the complete Domain-Driven Design specification, including bounded contexts, aggregates, and service responsibilities.
