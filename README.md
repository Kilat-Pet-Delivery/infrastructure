# infrastructure

Docker Compose orchestration and shared infrastructure for Kilat Pet Runner microservices architecture.

**Organization:** github.com/Kilat-Pet-Delivery

## Repository Contents

- **docker-compose.yml** — Orchestration for Postgres/PostGIS, Kafka, MinIO, current microservices, and the new runner-service stubs
- **infra/init-databases.sql** — Database initialization with existing service databases plus chat, incident, loyalty, and zones databases
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

# Create the shared object-storage bucket
make minio-init
```

## Microservices

| Service | Port | Repository |
|---------|------|------------|
| service-identity | 8004 | Kilat-Pet-Delivery/service-identity |
| service-booking | 8001 | Kilat-Pet-Delivery/service-booking |
| service-payment | 8002 | Kilat-Pet-Delivery/service-payment |
| service-runner | 8003 | Kilat-Pet-Delivery/service-runner |
| service-tracking | 8005 | Kilat-Pet-Delivery/service-tracking |
| service-notification | 8006 | Kilat-Pet-Delivery/service-notification |
| service-review | 8007 | Kilat-Pet-Delivery/service-review |
| service-chat | 8008 | Kilat-Pet-Delivery/service-chat |
| service-incident | 8009 | Kilat-Pet-Delivery/service-incident |
| service-loyalty | 8010 | Kilat-Pet-Delivery/service-loyalty |
| service-zones | 8011 | Kilat-Pet-Delivery/service-zones |

## Infrastructure Services

- **PostgreSQL 16 + PostGIS 3.4** (port 5433) — Geospatial database with UUID support
- **Apache Kafka** via Confluent (port 9092) — Event streaming and message broker
- **Zookeeper** (port 2181) — Distributed coordination service
- **MinIO** (API port 9000, console port 9001 by default) — S3-compatible local object storage using the shared `kilat-runner` bucket

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
make minio-init  # Create the kilat-runner bucket
```

If another local stack already owns ports 9000/9001, run MinIO with overrides:

```bash
MINIO_API_PORT=19000 MINIO_CONSOLE_PORT=19001 make minio-init
```

## Architecture

Refer to `docs/ARCHITECTURE.md` for the complete Domain-Driven Design specification, including bounded contexts, aggregates, and service responsibilities.
