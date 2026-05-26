# ADR 006: Estrategia de Testing

**Date**: 2026-05-22
**Status**: accepted

## Context

El proyecto partió con 0 tests sobre 88 archivos de código fuente. Se necesitaba una
estrategia que maximizara cobertura de lógica crítica con el menor esfuerzo posible,
priorizando la capa donde los bugs tienen mayor impacto.

## Decision

Pirámide de testing en 4 niveles, implementados en orden de prioridad:

**1. Unit tests de usecases** (prioridad máxima)
- Mocks de todos los providers usando `testify/mock`
- Cobertura de happy path + casos de error por usecase
- Sin DB, sin red — tests en milisegundos

**2. Tests de extractores y prompts** (white-box en package `inclusion`)
- Validación de parseo de respuestas de IA
- Garantiza que cambios en prompts no rompan el formato esperado

**3. Tests del AI client con `httptest`**
- `httptest.NewServer` simula el endpoint de Azure OpenAI
- Valida manejo de errores HTTP (timeout, 429, 500)

**4. CI pipeline con GitHub Actions**
- Jobs: `lint` (golangci-lint) + `vet` + `test` + `build`
- Ejecuta en cada push a `main` y en pull requests
- Falla el merge si algún job falla

Integration tests con DB real (testcontainers) están planificados como siguiente paso.

## Consequences

**Positivas**:
- Cobertura de lógica de negocio sin dependencia de infraestructura
- Tests rápidos y determinísticos — feedback inmediato en desarrollo local
- CI automatizado previene regresiones en merges a `main`
- Mocks generados de providers facilitan agregar tests a nuevos usecases

**Negativas**:
- No hay integration tests con DB real — queries GORM no se validan hasta staging
- Handlers HTTP no tienen tests end-to-end aún (gap conocido)
- StubClient en tests de usecases no detecta errores en prompt engineering
