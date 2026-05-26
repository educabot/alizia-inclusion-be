# ADR 001: Clean Architecture

**Date**: 2026-05-22
**Status**: accepted

## Context

El proyecto necesita una arquitectura mantenible que separe concerns y permita testing
independiente de cada capa. Un enfoque monolítico con lógica mezclada en handlers
dificulta el testing, la extensibilidad y el onboarding de nuevos desarrolladores.

## Decision

Adoptamos Clean Architecture con el siguiente flujo de dependencias estricto:

```
entrypoints → usecases → providers (interfaces) ← repositories
```

Capas:

- **entities** — structs de dominio (modelos GORM), sin lógica de negocio
- **providers** — interfaces (contratos) y errores sentinel; no importan implementaciones
- **usecases** — lógica de negocio, patrón `Execute(ctx, input) (output, error)` por operación
- **entrypoints** — handlers HTTP Gin; solo traducen HTTP ↔ usecases
- **repositories** — implementaciones GORM de los providers

Dependency Injection manual en `cmd/` (composition root). No se usa framework de DI.

## Consequences

**Positivas**:
- Cada capa es testeable de forma aislada (usecases con mocks, repositores con DB de test)
- Implementaciones intercambiables sin modificar lógica de negocio (e.g., AI client stub/real)
- Flujo de dependencias unidireccional — imposible crear dependencias circulares entre capas
- Onboarding predecible: cada feature sigue la misma estructura

**Negativas**:
- Más archivos y boilerplate que un enfoque monolítico (handler + usecase + provider + repo)
- DI manual en `cmd/` requiere wiring explícito al agregar nuevos componentes
