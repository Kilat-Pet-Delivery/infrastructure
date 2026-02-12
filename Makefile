.PHONY: all build test test-integration lint docker-up docker-down migrate-up

SERVICES = service-booking service-payment service-runner service-identity service-tracking service-notification api-gateway

build:
	@for %%s in ($(SERVICES)) do ( \
		echo Building %%s... && \
		cd %%s && go build -o server.exe ./cmd/server && cd .. \
	)

test:
	@for %%s in ($(SERVICES)) do ( \
		echo Testing %%s... && \
		cd %%s && go test ./... -v && cd .. \
	)

test-integration:
	cd service-payment && go test -tags=integration -v -timeout 120s -count=1 .
	cd service-booking && go test -tags=integration -v -timeout 120s -count=1 .

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down -v

docker-infra:
	docker-compose up -d postgres zookeeper kafka

tidy:
	@for %%s in ($(SERVICES)) do ( \
		echo Tidying %%s... && \
		cd %%s && go mod tidy && cd .. \
	)
	cd lib-common && go mod tidy && cd ..
	cd lib-proto && go mod tidy && cd ..
	cd api-gateway && go mod tidy && cd ..
