# ADR 004: Multi-tenancy por organization_id

**Date**: 2026-05-22
**Status**: accepted

## Context

La plataforma sirve a múltiples organizaciones educativas (escuelas, distritos, ONGs).
Cada organización debe ver exclusivamente sus propios datos — estudiantes, dispositivos,
adaptaciones — sin posibilidad de acceso cruzado entre tenants.

## Decision

Multi-tenancy mediante `organization_id` (UUID) como columna en todas las tablas de negocio:

- **TenantMiddleware** extrae `org_uuid` del `audience` claim del JWT y lo inyecta en `context.Context`
- **Todos los repositories** reciben `orgID` como parámetro explícito y filtran con `WHERE organization_id = ?`
- No existen queries cross-tenant en el código de aplicación
- Schema compartido (shared schema, shared database) — una sola instancia de PostgreSQL

```
JWT audience: ["org:uuid-123"] → middleware extrae → ctx.Value("org_id") → repo filtra
```

No se implementó row-level security en PostgreSQL (potencial mejora futura).

## Consequences

**Positivas**:
- Aislamiento de datos a nivel de aplicación — enforced automáticamente en cada request
- Schema único compartido: bajo costo operacional (una DB, un set de migraciones)
- Audit trail natural: `organization_id` presente en todos los registros

**Negativas**:
- Riesgo si un desarrollador agrega un nuevo repository sin filtrar por `org_id`; mitigado
  con code review y potencialmente con tests de integración
- No hay aislamiento a nivel de base de datos — un bug de aplicación podría exponer datos
  cross-tenant (sin RLS de PostgreSQL como safety net)
- Queries de reporting global requieren omitir el filtro explícitamente (solo para admins)
