# ADR 005: Logging Estructurado con log/slog

**Date**: 2026-05-22
**Status**: accepted

## Context

El sistema usaba `log.Printf` con formato de texto manual. Esto dificulta la búsqueda
de logs en producción, la correlación entre requests, la creación de alertas basadas
en campos específicos, y la integración con herramientas como CloudWatch o Datadog.

## Decision

Migramos a `log/slog` de la stdlib de Go (disponible desde Go 1.21):

- **Producción**: `slog.NewJSONHandler` — output JSON por línea, parseable por cualquier
  agregador de logs
- **Desarrollo local**: `slog.NewTextHandler` — output human-readable en terminal
- Selección automática basada en la env var `APP_ENV`

El middleware de request logging agrega los siguientes campos en cada request:

```json
{
  "request_id": "uuid-v4",
  "method": "POST",
  "path": "/v1/students",
  "status": 201,
  "duration_ms": 42,
  "org_id": "uuid-org",
  "user_id": "uuid-user"
}
```

`request_id` se genera en el middleware y se propaga via `context.Context` para
correlacionar logs de un mismo request a través de múltiples capas.

## Consequences

**Positivas**:
- Logs parseables por CloudWatch Logs Insights, Datadog, Grafana Loki sin configuración extra
- Correlación de todos los logs de un request via `request_id`
- Zero dependencias adicionales — `log/slog` es stdlib de Go
- Campos tipados: búsquedas exactas por `org_id`, `status`, `duration_ms`

**Negativas**:
- Requiere actualización de todos los call sites que usan `log.Printf` (refactor incremental)
- Logs JSON son menos legibles para humanos sin herramientas (mitigado con text handler en dev)
