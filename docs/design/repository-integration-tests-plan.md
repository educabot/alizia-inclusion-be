# Plan de acción — Tests de integración de repositories (Postgres testcontainers)

**Estado:** propuesto · **Fecha:** 2026-06-12 · **Owner:** —

## Objetivo

Llevar `src/repositories/{auth,catalog,inclusion,management}` de **0%** a cobertura
real, validando el **SQL contra un Postgres de verdad** (no SQLite), porque varios
repos dependen de features Postgres-only que SQLite no puede ejecutar.

## Por qué Postgres real y no SQLite

Los repos no son homogéneos. Hay SQL que **solo** corre en Postgres:

| Repo | Feature Postgres-only |
|------|------------------------|
| `inclusion/pedagogical_content.go` | `to_tsvector` / `plainto_tsquery` / `setweight` / `array_to_string` (full-text search ponderado) |
| `inclusion/adaptation.go:45` | `ILIKE` |
| `inclusion/ai_usage.go`, `conversation.go` | `datatypes.JSON` (JSONB) |
| varios | arrays Postgres (`keywords`, `tags`), `uuid` |

Testear con SQLite daría **falsos OK** (dialecto distinto) y dejaría el buscador RAG
sin cubrir. Por eso: un Postgres efímero por test-run vía `testcontainers-go`.

## Inventario de repos a cubrir

- **auth:** `user.go`
- **catalog:** `ramp.go`, `device.go`
- **management:** `classroom.go`
- **inclusion (15):** `student.go`, `student_profile.go`, `teacher_profile.go`,
  `adaptation.go`, `adaptation_resource.go`, `conversation.go`,
  `conversation_summary.go`, `ai_usage.go`, `diagnosis.go`, `ppi.go`, `situation.go`,
  `integradora_assignment.go`, `pedagogical_content.go`

## Infraestructura

### Dependencias nuevas
- `github.com/testcontainers/testcontainers-go`
- `github.com/testcontainers/testcontainers-go/modules/postgres`
- (driver `github.com/lib/pq` ya está; no hace falta golang-migrate: las migraciones
  se aplican leyendo los `.sql` igual que `scripts/dbmigrate`)

### Requisito de entorno
- **Docker** disponible en dev y en CI. Sin Docker, estos tests se saltean (ver build tag).

## Diseño del harness

Paquete nuevo `src/testutil/pgtest` (o `src/repositories/internal/dbtest`):

1. **`StartPostgres(t)`** — levanta **un** contenedor Postgres por paquete de test
   (`TestMain` + `sync.Once`), no uno por test (caro). Devuelve el `*gorm.DB`.
2. **Migraciones** — aplica en orden los `db/migrations/*.up.sql` contra el contenedor,
   reutilizando la lógica de `scripts/dbmigrate` (leer archivo → `db.Exec`). Extraer esa
   lógica a una función reutilizable para no duplicar.
3. **Aislamiento entre tests** — cada test corre dentro de una **transacción que se
   hace rollback** al final (`db.Begin()` → test → `tx.Rollback()`), o `TRUNCATE ... CASCADE`
   en `t.Cleanup`. Preferir rollback: más rápido y sin estado compartido (cumple
   test-standards: "integration tests must clean up — transactional rollback").
4. **Seed/fixtures** — extender `src/testutil/fixtures.go` con builders de filas
   (org, classroom, student, device, adaptation…) que inserten vía GORM.
5. **Scope multi-tenant** — todos los asserts verifican el filtro por `organization_id`
   (insertar 2 orgs y confirmar que el repo nunca cruza tenants).

### Build tag
Marcar los tests con `//go:build integration` para que el `go test ./...` por defecto
(unit, sin Docker) siga siendo rápido. CI corre `go test -tags=integration ./...` en un
job con servicio Docker.

## Layout de archivos

```
src/testutil/pgtest/pgtest.go          # StartPostgres, applyMigrations, WithTx
src/repositories/catalog/ramp_integration_test.go      //go:build integration
src/repositories/catalog/device_integration_test.go
src/repositories/management/classroom_integration_test.go
src/repositories/auth/user_integration_test.go
src/repositories/inclusion/*_integration_test.go
```

## Fases (incremental, mergeable por fase)

**Fase 0 — Harness (1 PR):**
añadir deps, escribir `pgtest`, extraer la lógica de migraciones de `scripts/dbmigrate`
a una función compartida, smoke test que levanta el contenedor y migra. Configurar el
job de CI con Docker + `-tags=integration`.

**Fase 1 — Repos CRUD simples (validar el patrón):**
`catalog/ramp`, `catalog/device`, `management/classroom`, `auth/user`.
Cubre Create/Get/List/Update/Delete + filtro por org + `ErrNotFound`.

**Fase 2 — Inclusion CRUD:**
`student`, `student_profile`, `teacher_profile`, `adaptation` (CRUD + `ILIKE` search),
`adaptation_resource`, `diagnosis`, `ppi`, `situation`, `integradora_assignment`.

**Fase 3 — Inclusion con estado/JSON:**
`conversation` (AppendTurn + metadata JSONB), `conversation_summary`, `ai_usage`
(ContextSnapshot JSONB).

**Fase 4 — El difícil: `pedagogical_content`:**
full-text search ponderado. Tests con corpus sembrado que verifican ranking y que
`plainto_tsquery` matchea por palabra suelta (el `OR` query de `orTSQuery`). Es el de
mayor valor porque es SQL que hoy nadie valida.

## Gotchas / riesgos

- **Arranque del contenedor:** ~2-5s; por eso uno por paquete, no por test.
- **CI sin Docker:** el job debe declarar el servicio Docker; si no, los tests se saltean
  por el build tag (no rompen el pipeline default).
- **tsvector/idioma `'spanish'`:** la imagen Postgres debe tener la config de texto en
  español (la imagen oficial la trae). Verificar en Fase 0.
- **Paralelismo:** con rollback por test se puede usar `t.Parallel()` solo si cada test
  tiene su propia tx; cuidado con el contenedor compartido.
- **Migraciones divergentes:** usar SIEMPRE los `.sql` reales de `db/migrations`, nunca
  un esquema paralelo, para que el test valide lo que corre en prod.

## Estimación

- Fase 0: ~media sesión (harness + CI).
- Fases 1-3: ~1 sesión c/u (mecánico una vez que el harness anda).
- Fase 4: ~media sesión (cuidado con el ranking).

## Decisiones tomadas

- **DB de test:** Postgres real vía testcontainers (no SQLite, no híbrido).
- **Aislamiento:** transacción con rollback por test.
- **Migraciones:** los `.up.sql` reales, aplicados por lógica compartida con `scripts/dbmigrate`.
