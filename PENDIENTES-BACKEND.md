# Pendientes Backend — Alizia Inclusión

Documento generado a partir de la auditoría BE↔FE del 2026-05-15.
Última actualización: 2026-05-19.

---

## 1. ~~Migración: agregar roles `ministerio` y `psicopedagogo` al ENUM~~ ✅

**Estado:** Resuelto (2026-05-19)
- Migración `000007_add_roles_to_enum.up.sql` creada y aplicada.
- Seed actualizado con usuarios demo `ministerio@demo.edu` y `psico@demo.edu`.
- DB verificada: ENUM `member_role` tiene 5 valores (`teacher`, `coordinator`, `admin`, `ministerio`, `psicopedagogo`).

**Mapeo de roles BE→FE (para referencia):**

| BE (DB) | FE (UI) |
|---|---|
| `teacher` | `docente` |
| `coordinator` | `integradora` |
| `admin` | `director` |
| `ministerio` | `ministerio` |
| `psicopedagogo` | `psicopedagogo` |

---

## 2. ~~Validar valores de `status` en adaptaciones~~ ✅

**Estado:** Resuelto (2026-05-19)
- Validación implementada en `src/core/usecases/inclusion/update_adaptation.go`.
- Error `errInvalidStatus` agregado en `src/core/usecases/inclusion/errors.go`.
- Valores válidos: `en_curso`, `probado`, `funciono`, `para_ajustar`.

---

## 3. ~~Recompilar el binario~~ ✅

**Estado:** Resuelto (2026-05-19)
- Binario recompilado vía Docker: `docker run --rm -v ... golang:1.26.3-alpine go build -o alizia-inclusion-api.exe ./cmd`
- Tamaño: ~40.5 MB, target: `windows/amd64`.

---

## 4. ~~Agregar campo `useful_when` a devices~~ ✅

**Estado:** Resuelto (2026-05-19)
- Migración `000008_add_useful_when_to_devices.up.sql` creada y aplicada.
- Campo `UsefulWhen` agregado a entity `device.go` y DTO `catalog_devices.go`.
- Seed actualizado con contenido diferenciado: `needs_description` describe a quién le sirve, `useful_when` describe cuándo usarlo.
- 15 dispositivos actualizados en la DB.
- Binario recompilado.

---

## 5. ~~Endpoint `/health` — verificar respuesta~~ ✅

**Estado:** Verificado (2026-05-19)
- El endpoint en `cmd/app.go:44` devuelve `{ "status": "ok" }` — coincide con el formato `{ status: string }` que espera el FE.
- No requiere cambios.

---

## Resumen

| Archivo | Cambio | Estado |
|---|---|---|
| `db/migrations/000007_add_roles_to_enum.up.sql` | Roles al ENUM | ✅ Hecho |
| `db/migrations/000007_add_roles_to_enum.down.sql` | Marcar irreversible | ✅ Hecho |
| `db/seeds/seed.sql` | Usuarios demo con roles nuevos + `useful_when` | ✅ Hecho |
| `alizia-inclusion-api.exe` | Recompilar con `go build` | ✅ Hecho |
| `src/core/usecases/inclusion/update_adaptation.go` | Validar status | ✅ Hecho |
| `src/core/usecases/inclusion/errors.go` | `errInvalidStatus` | ✅ Hecho |
| `db/migrations/000008_add_useful_when_to_devices.up.sql` | Campo `useful_when` en devices | ✅ Hecho |
| `db/migrations/000008_add_useful_when_to_devices.down.sql` | Rollback del campo | ✅ Hecho |
| `src/core/entities/device.go` | Campo `UsefulWhen` | ✅ Hecho |
| `src/entrypoints/catalog_devices.go` | DTO + mapper | ✅ Hecho |
| `cmd/app.go` | Endpoint `/health` verificado | ✅ Sin cambios |
