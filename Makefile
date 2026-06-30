.PHONY: build test test-cover vet lint docker migrate run seed mocks eval trace

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

# Eval de comportamiento del modelo real (Azure / gpt-5.4): verifica que ante inputs
# que deberían disparar búsqueda, el modelo SÍ llama la tool correcta. No mockea el AI,
# cuesta tokens y es no-determinístico → on-demand / nightly, no en el CI bloqueante.
# Lee credenciales del .env. Tuneable: EVAL_RUNS, EVAL_THRESHOLD.
eval:
	go test -tags eval -count=1 -v -run TestAgenticEval ./src/core/usecases/inclusion/

# Reconstruye conversaciones del agente desde los logs de producción (Railway) para
# auditar tool calls, prompt y respuesta. Solo lectura. Requiere la CLI railway logueada.
# Uso: make trace ARGS="last --user 5"  |  make trace ARGS="conversation 238"
#      make trace ARGS="student Francisco --format md"
trace:
	python3 scripts/trace.py $(ARGS)

run:
	air
