.PHONY: all build test lint docker-up docker-down migrate-up

SERVICES = service-booking service-payment service-runner service-identity service-tracking api-gateway

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
