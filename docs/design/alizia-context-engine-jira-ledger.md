# Ledger de creación en Jira — proyecto ALZ (Alizia)

> Estructura del Context Engine en Jira. Sitio: aula-educabot.atlassian.net · cloudId `b171db1b-26f3-4903-98d1-0dcfca599382`.
> **Refactor 2026-06-01:** jerarquía EP → HU → T (épica → historia testeable → subtarea técnica).
> Códigos en el título: `EP-{n}` · `HU-{n}.{m}` · `T-{n}.{k}` (T corrida dentro de la épica). Sin "CE".

## Épicas (EP)
| Código | Jira | Título |
|---|---|---|
| EP-1 | ALZ-266 | Traza |
| EP-2 | ALZ-267 | Contexto del alumno y del docente |
| EP-3 | ALZ-268 | Prompts versionados en DB |
| EP-4 | ALZ-269 | Memoria |
| EP-5 | ALZ-270 | Flywheel (mejora continua) |
| EP-6 | ALZ-271 | Contenido pedagógico base |

## Historias (HU) y Tareas (T)
| Código | Jira | Épica | Tareas (Jira) |
|---|---|---|---|
| HU-1.1 Trazar cada turno y la aceptación implícita | ALZ-393 | EP-1 | T-1.1 ALZ-403 · T-1.2 ALZ-404 |
| HU-2.1 Alizia usa el contexto del alumno en su respuesta | ALZ-394 | EP-2 | T-2.1 ALZ-405 · T-2.2 ALZ-406 · T-2.3 ALZ-407 · T-2.4 ALZ-408 · T-2.5 ALZ-409 · T-2.6 ALZ-410 · T-2.7 ALZ-411 |
| HU-2.2 Profundizar bajo demanda con tools agénticas | ALZ-395 | EP-2 | T-2.8 ALZ-412 |
| HU-3.1 Editar y publicar un prompt desde DB sin deploy | ALZ-396 | EP-3 | T-3.1 ALZ-413 · T-3.2 ALZ-414 · T-3.3 ALZ-415 · T-3.4 ALZ-416 · T-3.5 ALZ-417 · T-3.6 ALZ-418 |
| HU-4.1 La conversación larga conserva memoria (resumen) | ALZ-397 | EP-4 | T-4.1 ALZ-419 · T-4.2 ALZ-420 |
| HU-4.2 Memoria viva del alumno entre sesiones (insights) | ALZ-398 | EP-4 | T-4.3 ALZ-421 |
| HU-5.1 Ejemplos golden alimentan el few-shot (cold-start) | ALZ-399 | EP-5 | T-5.1 ALZ-422 · T-5.2 ALZ-423 |
| HU-5.2 Medir si el sistema mejora (win-rate y métricas) | ALZ-400 | EP-5 | T-5.3 ALZ-424 |
| HU-5.3 Promover golden y correr A/B entre versiones | ALZ-401 | EP-5 | T-5.4 ALZ-425 · T-5.5 ALZ-426 |
| HU-6.1 Contenido pedagógico base del prompt | ALZ-402 | EP-6 | T-6.1 ALZ-427 · T-6.2 ALZ-428 · T-6.3 ALZ-429 · T-6.4 ALZ-430 · T-6.5 ALZ-431 · T-6.6 ALZ-432 |

**Totales:** 6 EP (ALZ-266–271) · 10 HU (ALZ-393–402) · 30 T (ALZ-403–432). Todo asignado a Sebastian.
El contenido literal (descripciones) está en `alizia-context-engine-jira-rework-proposal.md`.

## Estructura vieja — PENDIENTE DE BORRAR
La carga anterior (historias CE-x.y + subtareas) quedó obsoleta. Son **ALZ-272 … ALZ-392** (30 historias + 91 subtareas), colgando de ALZ-246, todas con label **`context-engine`**.

**Borrado (manual, el MCP no expone delete):** en la UI, filtrar
`project = ALZ AND labels = context-engine`
→ selecciona las 121 viejas (las nuevas EP/HU/T **no** tienen ese label) → Bulk change → Delete.
Borrar las historias arrastra sus subtareas en cascada.
