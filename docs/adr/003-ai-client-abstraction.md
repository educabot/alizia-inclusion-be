# ADR 003: Abstracción del Cliente de IA

**Date**: 2026-05-22
**Status**: accepted

## Context

El sistema integra IA generativa para recomendar dispositivos adaptativos personalizados.
Azure OpenAI requiere credenciales de producción que no están disponibles en entornos de
desarrollo local ni en pipelines de CI, lo que bloqueaba el ciclo de desarrollo.

## Decision

Definimos una interface `AIClient` en `src/core/providers/` con dos implementaciones:

- **AzureClient** — llama a Azure OpenAI; usado en producción
- **StubClient** — retorna respuestas hardcodeadas; usado en desarrollo y CI

La selección es automática en `cmd/repositories.go` basada en env vars:

```go
if cfg.AzureOpenAIKey != "" {
    aiClient = azure.NewClient(cfg)
} else {
    aiClient = stub.NewClient()
}
```

Agregar un nuevo provider (OpenAI directo, Anthropic, etc.) requiere solo implementar la
interface sin modificar usecases ni handlers.

## Consequences

**Positivas**:
- Desarrollo y CI funcionan sin credenciales de Azure
- Testing de usecases sin llamadas externas (determinístico, rápido)
- Intercambio de provider sin modificar lógica de negocio (Open/Closed Principle)
- Facilita A/B testing entre modelos en el futuro

**Negativas**:
- StubClient no valida el formato de prompts enviados — errores de prompt engineering
  pueden pasar desapercibidos hasta llegar a producción
- Las respuestas del StubClient no reflejan variabilidad real del modelo
