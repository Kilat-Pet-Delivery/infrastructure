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
