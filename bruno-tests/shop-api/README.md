# Kilat Shop API Bruno Collection

Run after the local docker-compose stack is up and seeded:

```bash
bru run bruno-tests/shop-api --env local
```

Set `ownerToken`, `staffToken`, `customerToken`, and `runnerToken` in `_env/local.bru` before running. The first create calls return IDs that can be copied back into the environment for the later lifecycle requests.
