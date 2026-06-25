.PHONY: build test test-cover test-integration vet lint docker migrate run seed mocks

build:
	CGO_ENABLED=0 go build -o alizia-inclusion-api ./cmd

test:
	go test -race ./...

test-cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Repository integration tests against a real Postgres (testcontainers). Requires Docker.
test-integration:
	go test -tags=integration -race ./src/repositories/... ./src/testutil/pgtest/...

vet:
	go vet ./...

lint:
	golangci-lint run

docker:
	docker compose up -d

migrate:
	./scripts/migrate.sh

seed:
	./scripts/seed.sh

mocks:
	mockery

run:
	air
