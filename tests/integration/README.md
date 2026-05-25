# Kilat Cross-Service Integration Tests

These tests exercise the runner-app flows through the real docker-compose stack. They are intentionally gated behind `KILAT_RUN_INTEGRATION=1` so local unit-test runs do not try to boot Docker.

## Run

```bash
make integration-up
make seed
make integration-test
make integration-down
```

Set these fixture variables when running the full suite:

- `KILAT_CUSTOMER_TOKEN` or `KILAT_CUSTOMER_EMAIL` + `KILAT_CUSTOMER_PASSWORD`
- `KILAT_RUNNER_TOKEN` or the seeded runner credentials from `seed/runner-test-user.sql`
- `KILAT_AGENT_TOKEN` or `KILAT_AGENT_EMAIL` + `KILAT_AGENT_PASSWORD`
- `KILAT_CHAT_THREAD_ID`
- `KILAT_ACTIVE_BOOKING_ID`
- `KILAT_REFEREE_TOKEN`
- `KILAT_REFERRAL_ID`
- `KILAT_SURGE_FIXTURE_READY=1` after the zone demand fixture is seeded

Defaults:

- `KILAT_GATEWAY_URL=http://localhost:8080`
- `KILAT_GATEWAY_WS_URL=ws://localhost:8080`
- `KILAT_KAFKA_BROKERS=localhost:9092`
