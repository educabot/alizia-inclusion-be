.PHONY: build test test-cover vet lint docker migrate run seed mocks

build:
	CGO_ENABLED=0 go build -o alizia-inclusion-api ./cmd

test:
	go test -race ./...

test-cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

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
