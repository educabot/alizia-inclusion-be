# Ledger de creación en Jira — proyecto ALZ (Alizia)

> Mapeo CE → clave Jira de la carga del backlog `alizia-context-engine-backlog.md`.
> Sitio: aula-educabot.atlassian.net · cloudId `b171db1b-26f3-4903-98d1-0dcfca599382`.
> Última reconciliación: 2026-06-01 (verificado contra Jira vía JQL).

## Epics
| CE | Jira | Título |
|---|---|---|
| EPIC 0 | ALZ-266 | Traza |
| EPIC 1 | ALZ-267 | Contexto del alumno y del docente |
| EPIC 2 | ALZ-268 | Prompts versionados en DB |
| EPIC 3 | ALZ-269 | Memoria |
| EPIC 4 | ALZ-270 | Flywheel (self-improvement) |
| EPIC 5 | ALZ-271 | Contenido pedagógico base |

## Historias
| CE | Jira | Epic |
|---|---|---|
| CE-0.1 | ALZ-272 | ALZ-266 |
| CE-0.2 | ALZ-273 | ALZ-266 |
| CE-1.1 | ALZ-274 | ALZ-267 |
| CE-1.2 | ALZ-275 | ALZ-267 |
| CE-1.3 | ALZ-276 | ALZ-267 |
| CE-1.4 | ALZ-277 | ALZ-267 |
| CE-1.5 | ALZ-278 | ALZ-267 |
| CE-1.6 | ALZ-279 | ALZ-267 |
| CE-1.7 | ALZ-280 | ALZ-267 |
| CE-1.8 | ALZ-281 | ALZ-267 |
| CE-2.1 | ALZ-282 | ALZ-268 |
| CE-2.2 | ALZ-283 | ALZ-268 |
| CE-2.3 | ALZ-284 | ALZ-268 |
| CE-2.4 | ALZ-285 | ALZ-268 |
| CE-2.5 | ALZ-286 | ALZ-268 |
| CE-2.6 | ALZ-287 | ALZ-268 |
| CE-3.1 | ALZ-288 | ALZ-269 |
| CE-3.2 | ALZ-289 | ALZ-269 |
| CE-3.3 | ALZ-290 | ALZ-269 |
| CE-4.1 | ALZ-291 | ALZ-270 |
| CE-4.2 | ALZ-292 | ALZ-270 |
| CE-4.3 | ALZ-293 | ALZ-270 |
| CE-4.4 | ALZ-294 | ALZ-270 |
| CE-4.5 | ALZ-295 | ALZ-270 |
| CE-5.1 | ALZ-296 | ALZ-271 |
| CE-5.2 | ALZ-297 | ALZ-271 |
| CE-5.3 | ALZ-298 | ALZ-271 |
| CE-5.4 | ALZ-299 | ALZ-271 |
| CE-5.5 | ALZ-300 | ALZ-271 |
| CE-5.6 | ALZ-301 | ALZ-271 |

> **30/30 historias creadas.** (CE-5.5 y CE-5.6 — ALZ-300/301 — creadas el 2026-06-01;
> el resto fue creado el 2026-05-29 y no se había registrado en este ledger.)

## Convenciones aplicadas en las historias
- **issuetype:** `Historia` · **priority:** `Medium` · **parent:** el epic correspondiente.
- **Descripción:** objetivo (`Como… quiero… para…`) + criterios + `Depende de` / `Paralelizable` / `Spec`
  como **texto** (apunta a la sección del design doc).
- **Labels:** `context-engine` + component (`data-model`/`context-assembler`/`renderer`/`memory`/`flywheel`/`content`)
  + `migration-000xxx` (si aplica) + `risk-low|med|high` + flags (`parallelizable`, `cold-start`).
- El **component** del backlog va como **label**, no en el campo Components.

## Subtareas
**No creadas** — las 🔧 sub-tareas del backlog quedan como detalle dentro de cada historia/spec,
no como issues separados (patrón aplicado uniformemente a las 30 historias).

## Issue links (Blocks)
**No creados** — las dependencias (`🔗 Depende de`) están capturadas como **texto** en la
descripción de cada historia, no como links nativos de Jira.
