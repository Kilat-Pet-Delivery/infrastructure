# Changelog

All notable changes to the Kilat Pet Delivery infrastructure repository are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `docs/run-from-source.md`: runbook for bringing up `api-gateway`,
  `service-identity` and `service-booking` against the shared dev-infra stack,
  with the register -> login -> create booking smoke test, the exact ports that
  work around the Desa Murni Batik collision, and cleanup steps. (KPD-3)
- `CHANGELOG.md`: this file. Partially advances KPD-52.

### Changed

- `README.md`: new "The two ways to run" section -- the daily loop (shared
  dev-infra stack plus `go run`) versus the full-stack smoke (this repo's
  `docker-compose.yml`), what differs between them (Postgres 5432 vs 5433,
  `localhost` vs container hostnames), and the 8001-8009 port collision with
  Desa Murni Batik. (KPD-6)
