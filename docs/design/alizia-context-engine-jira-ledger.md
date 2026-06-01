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
**91 subtareas creadas** (ALZ-302..392), el 2026-06-01. Una `Subtarea` Jira por cada 🔧 del backlog,
con summary prefijado por CE para identificarlas en el board. Reparto:

| Historia | Subtareas (rango ALZ) |
|---|---|
| CE-0.1 | 302–305 | 
| CE-0.2 | 306–309 |
| CE-1.1 | 310–313 |
| CE-1.2 | 314–316 |
| CE-1.3 | 317–320 |
| CE-1.4 | 321–324 |
| CE-1.5 | 325–327 |
| CE-1.6 | 328–330 |
| CE-1.7 | 331–335 |
| CE-1.8 | 336–338 |
| CE-2.1 | 339–342 |
| CE-2.2 | 343–347 |
| CE-2.3 | 348–350 |
| CE-2.4 | 351–354 |
| CE-2.5 | 355–358 |
| CE-2.6 | 359–361 |
| CE-3.1 | 362–365 |
| CE-3.2 | 366–369 |
| CE-3.3 | 370–372 |
| CE-4.1 | 373–376 |
| CE-4.2 | 377–380 |
| CE-4.3 | 381–383 |
| CE-4.4 | 384–387 |
| CE-4.5 | 388–392 |

> Las CE-5.x (contenido) no llevan subtareas técnicas — son autoría.

## Issue links (Blocks)
**37 links creados** (`Blocks`), el 2026-06-01. Semántica: el que bloquea → el bloqueado.

| Bloqueada | La bloquean |
|---|---|
| ALZ-280 (CE-1.7) | 274, 275, 276, 277, 278, 279 |
| ALZ-281 (CE-1.8) | 280, 289 |
| ALZ-283 (CE-2.2) | 282, 280 |
| ALZ-284 (CE-2.3) | 282, 283 |
| ALZ-285 (CE-2.4) | 283 |
| ALZ-286 (CE-2.5) | 282, 283, 284, 285 |
| ALZ-287 (CE-2.6) | 282, 296, 297, 298, 299 |
| ALZ-289 (CE-3.2) | 288 |
| ALZ-290 (CE-3.3) | 288 |
| ALZ-291 (CE-4.1) | 301 |
| ALZ-292 (CE-4.2) | 291, 283 |
| ALZ-293 (CE-4.3) | 272, 273, 282 |
| ALZ-294 (CE-4.4) | 282, 291, 293 |
| ALZ-295 (CE-4.5) | 291, 293, 289 |
| ALZ-276 (CE-1.3) | 300 |
