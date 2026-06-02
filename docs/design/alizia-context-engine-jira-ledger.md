# Ledger de creación en Jira — proyecto ALZ (Alizia)

> Estructura del Context Engine en Jira. Sitio: aula-educabot.atlassian.net · cloudId `b171db1b-26f3-4903-98d1-0dcfca599382`.
> **Modelo Producto (refactor 2026-06-01):** épicas = objetivos de producto · historias = casos de uso testeables
> en formato "Como… quiero… para…" · subtareas = trabajo técnico. Códigos `EP-{n}` · `HU-{n}.{m}` · `T-{n}.{k}`.

## EP-1 · Alizia entiende a cada alumno (ALZ-266)
| Historia | Jira | Tareas |
|---|---|---|
| HU-1.1 · Alizia tiene en cuenta el contexto del alumno | ALZ-394 | T-1.1 ALZ-405 · T-1.2 ALZ-406 · T-1.3 ALZ-407 · T-1.4 ALZ-408 · T-1.5 ALZ-409 · T-1.6 ALZ-410 · T-1.7 ALZ-411 |
| HU-1.2 · Alizia profundiza con más datos del alumno cuando hace falta | ALZ-395 | T-1.8 ALZ-412 |

## EP-2 · Alizia con memoria (ALZ-267)
| Historia | Jira | Tareas |
|---|---|---|
| HU-2.1 · Alizia no pierde el hilo en conversaciones largas | ALZ-397 | T-2.1 ALZ-419 · T-2.2 ALZ-420 |
| HU-2.2 · Alizia recuerda entre sesiones qué funcionó con mi alumno | ALZ-398 | T-2.3 ALZ-421 |

## EP-3 · Alizia mejora con el uso (ALZ-268)
| Historia | Jira | Tareas |
|---|---|---|
| HU-3.1 · Medir si las recomendaciones de Alizia mejoran | ALZ-400 | T-3.1 ALZ-403 · T-3.2 ALZ-404 · T-3.3 ALZ-424 |
| HU-3.2 · Ajustar el comportamiento de Alizia sin deploy | ALZ-396 | T-3.4 ALZ-413 · T-3.5 ALZ-414 · T-3.6 ALZ-415 · T-3.7 ALZ-416 · T-3.8 ALZ-417 · T-3.9 ALZ-418 |
| HU-3.3 · Alizia mejora sola promoviendo buenos ejemplos | ALZ-399 | T-3.10 ALZ-422 · T-3.11 ALZ-423 · T-3.12 ALZ-425 · T-3.13 ALZ-426 |

## EP-4 · Alizia con criterio pedagógico (ALZ-269)
| Historia | Jira | Tareas |
|---|---|---|
| HU-4.1 · Las respuestas de Alizia siguen los lineamientos de inclusión | ALZ-402 | T-4.1 ALZ-427 · T-4.2 ALZ-428 · T-4.3 ALZ-429 · T-4.4 ALZ-430 · T-4.5 ALZ-431 · T-4.6 ALZ-432 |

**Totales:** 4 EP (ALZ-266–269) · 8 HU · 30 T. Todo asignado a Sebastian. Las historias llevan el "Como… quiero… para…" como primera línea de la descripción; los criterios de aceptación viven en la historia.

## PENDIENTE DE BORRAR (manual — el MCP no expone delete)
1. **Modelo CE original** — ALZ-272 … ALZ-392 (30 historias + 91 subtareas), label **`context-engine`**.
   Filtro UI: `project = ALZ AND labels = context-engine` → Bulk delete.
2. **Sobrantes del modelo Plataforma** — épicas ALZ-270, ALZ-271 e historias ALZ-393, ALZ-401, todas en estado **Cancelado** y con prefijo `[BORRAR]`.
   Filtro UI: `project = ALZ AND statusCategory = Done AND summary ~ "BORRAR"` → Bulk delete.

Las nuevas EP/HU/T (ALZ-266–269, 394–432, y 403–432) **no** tienen label `context-engine` ni `[BORRAR]`, así que ninguno de los dos filtros las toca.
