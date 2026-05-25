.PHONY: all build test test-integration integration-up integration-test integration-test-fast integration-down lint docker-up docker-down up down seed seed-chat seed-zones seed-loyalty minio-init minio-prune docker-infra tidy

COMPOSE ?= docker compose
REPO_ROOT := $(abspath ..)
SERVICES = service-booking service-payment service-runner service-identity service-tracking service-notification service-review service-chat service-incident service-loyalty service-zones api-gateway
LIBS = lib-common lib-proto

build:
	@for s in $(SERVICES); do \
		echo "Building $$s..."; \
		(cd "$(REPO_ROOT)/$$s" && go build -o server ./cmd/server); \
	done

test:
	@for s in $(LIBS) $(SERVICES); do \
		echo "Testing $$s..."; \
		(cd "$(REPO_ROOT)/$$s" && go test ./... -v); \
	done

test-integration:
	cd "$(REPO_ROOT)/service-payment" && go test -tags=integration -v -timeout 120s -count=1 .
	cd "$(REPO_ROOT)/service-booking" && go test -tags=integration -v -timeout 120s -count=1 .

integration-up:
	$(COMPOSE) up -d --build
	$(COMPOSE) run --rm minio-init

integration-test:
	cd tests/integration && KILAT_RUN_INTEGRATION=1 go test -v -timeout 30m -count=1 ./...

integration-test-fast:
	cd tests/integration && go test -v -timeout 2m -count=1 ./...

integration-down:
	$(COMPOSE) down -v

docker-up:
	$(COMPOSE) up -d --build

docker-down:
	$(COMPOSE) down -v

up: docker-up

down: docker-down

seed:
	docker exec -i kilat-postgres psql -U kilat -f - < seed/runner-test-user.sql

seed-chat:
	@echo "service-chat seed data lands with Phase 1."

seed-zones:
	@echo "service-zones seeds KL zone polygons through its startup migrations."

seed-loyalty:
	@echo "service-loyalty seeds quest definitions on startup in development."

minio-init:
	$(COMPOSE) run --rm minio-init

minio-prune:
	$(COMPOSE) down
	docker volume rm infrastructure_miniodata 2>/dev/null || true

docker-infra:
	$(COMPOSE) up -d postgres zookeeper kafka minio

tidy:
	@for s in $(LIBS) $(SERVICES); do \
		echo "Tidying $$s..."; \
		(cd "$(REPO_ROOT)/$$s" && go mod tidy); \
	done
