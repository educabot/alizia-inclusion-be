# Backlog Context Engine — DOCUMENTO DE TAREAS (fuente de verdad del backlog)

> **Spec / diseño (el "normal"):** [`alizia-context-engine.md`](./alizia-context-engine.md) — el *cómo* a fondo (modelo de datos, DBML, validaciones).
> **Este doc:** el *qué* y el *en qué orden* — las historias y subtareas tal como están en Jira.
>
> **Decisión final (Juan, 2026-06-02):** **NO hay épicas intermedias.** Todo cuelga de la épica **AlizIA Inclusión - Chubut (ALZ-246)**.
> **Conteo definitivo: 10 Historias · 30 Subtareas.** (El "4" y el "6" de versiones anteriores eran cantidades de *épicas* de modelos descartados, no de historias.)

## Jerarquía (3 niveles, lo que soporta Jira)

```
ALZ-246 · AlizIA Inclusión - Chubut        (ÉPICA — la única)
 └─ HU-n   (Historia)   "Como… quiero… para…" + criterios de aceptación
      └─ T-n.k  (Subtarea)   trabajo técnico
```

- Códigos: **HU-{n}** Historia · **T-{n}.{k}** Subtarea (la `k` corre dentro de la historia `n`).
- Historias en formato caso de uso testeable: *"Como {user} quiero {resultado} para {motivo}"*, lenguaje común, con todos los criterios de aceptación.

## Índice de las 10 Historias (key Jira → título)

| HU | Jira | Título | Subtareas |
|---|---|---|---|
| HU-1 | ALZ-449 | Trazar cada turno y la aceptación implícita | T-1.1 · T-1.2 (2) |
| HU-2 | ALZ-394 | Alizia usa el contexto del alumno en su respuesta | T-2.1 … T-2.7 (7) |
| HU-3 | ALZ-395 | Profundizar bajo demanda con tools agénticas | T-3.1 (1) |
| HU-4 | ALZ-396 | Editar y publicar un prompt desde DB sin deploy (validación + fallback) | T-4.1 … T-4.6 (6) |
| HU-5 | ALZ-397 | La conversación larga conserva memoria (resumen) | T-5.1 · T-5.2 (2) |
| HU-6 | ALZ-398 | Memoria viva del alumno entre sesiones (insights) | T-6.1 (1) |
| HU-7 | ALZ-450 | Ejemplos golden alimentan el few-shot (cold-start) | T-7.1 · T-7.2 (2) |
| HU-8 | ALZ-400 | Medir si el sistema mejora (win-rate y métricas) | T-8.1 (1) |
| HU-9 | ALZ-399 | Promover golden y correr A/B entre versiones | T-9.1 · T-9.2 (2) |
| HU-10 | ALZ-402 | Contenido pedagógico base del prompt | T-10.1 … T-10.6 (6) |

**Totales: 10 Historias · 30 Subtareas · 1 sola épica (Chubut).**

---

# HU-1 · Trazar cada turno y la aceptación implícita  (ALZ-449)

**Como** equipo de Alizia **quiero** registrar en cada interacción qué se pidió, qué se respondió y si el docente lo aprovechó, **para** medir y mejorar las recomendaciones con datos reales.

**Criterios de aceptación**
- Cada turno guarda en `ai_usage`: modelo, latencia, tool calls, `conversation_id`, `message_id` y un snapshot de contexto con **IDs (no PII en claro)**.
- Una adaptación creada desde una sugerencia guarda su origen y queda marcada según se editó o no antes de guardar.
- El tablero del Director sigue funcionando; filas viejas con los campos nuevos en `NULL`.
- Registro best-effort: si falla, no bloquea la respuesta.

- **T-1.1 · Enriquecer `ai_usage` con columnas de traza** (ALZ-403)
  `ALTER ai_usage`: `latency_ms`, `model`, `tool_calls`, `conversation_id`, `message_id`, `context_snapshot` (IDs, no PII). Filas viejas en `NULL`.
- **T-1.2 · Ligar adaptaciones a su origen + señal de aceptación implícita** (ALZ-404)
  La adaptación persiste `source_conversation_id` y `source_message_id`; `was_edited=false` si se guardó tal cual, `true` si se modificó.

---

# HU-2 · Alizia usa el contexto del alumno en su respuesta  (ALZ-394)

**Como** docente **quiero** que Alizia tenga en cuenta el contexto de mi alumno (perfil, situaciones de aula, PPI, entorno) al responder, **para** que las sugerencias sean pertinentes a ese chico y no genéricas.

**Criterios de aceptación**
- Con contexto cargado, la respuesta/el prompt lo refleja.
- Campos ausentes (todos opcionales) no imprimen "N/A" ni rompen.
- No se filtra PII a logs.

- **T-2.1 · Perfil del docente (`teacher_profiles`)** (ALZ-405)
  `teacher_profiles` 1:1 con `users`: `age_range`/`birthdate` (nullable), `years_experience`, `specialization`, `subjects[]`, `tone_preference`, `bio`. Aislamiento por organización.
- **T-2.2 · Enriquecer alumno (`students` / `student_profiles`)** (ALZ-406)
  `students`: `birthdate`, `age_range`, `grade_level`, `preferred_name` (nullable). `student_profiles`: `support_level`, `strengths[]`, `interests[]`, `triggers[]`, `effective_strategies[]`, `ineffective_strategies[]`, `situation_codes[]`, `has_therapeutic_companion`, `environment_notes`. Todo opcional.
- **T-2.3 · Catálogo de situaciones observables (`situations_catalog`)** (ALZ-407)
  Tabla global (con `organization_id` para per-org futuro), `phase` nullable. El seed de las ~15 situaciones se carga por **script de entorno**, no en la migración.
- **T-2.4 · Diagnósticos estructurados opcionales** (ALZ-408)
  `diagnoses_catalog` (global) + `student_diagnoses` (M2M con `severity`). Opcional y secundario a las situaciones; Alizia puede sugerir, nunca exigir. Seed por script de entorno.
- **T-2.5 · Proyecto Pedagógico Individual (`ppi`)** (ALZ-409)
  `ppi` 1 por alumno, todos los campos nullable. Cuando existe, es contexto de primera línea.
- **T-2.6 · Rol maestra integradora + asignación** (ALZ-410)
  Sumar `maestra_integradora` al enum `member_role` y modelar la asignación integradora↔alumno. RBAC fuera de scope (otro equipo).
- **T-2.7 · Context Assembler (`BuildPromptContext`)** (ALZ-411)
  Usecase que junta todo el contexto en un struct tipado, ordenado para caching (prefijo invariante adelante, variable atrás). Punto de convergencia de las tareas anteriores.

---

# HU-3 · Profundizar bajo demanda con tools agénticas  (ALZ-395)

**Como** docente **quiero** que Alizia traiga más datos del alumno cuando hace falta (historial, adaptaciones, aprendizajes), **para** que profundice sin que yo le repita todo.

**Criterio de aceptación**
- Disponibles `get_student_history`, `get_past_adaptations`, `get_student_insights`; el dispatcher vive en el usecase (clean architecture); el modelo las usa bajo demanda.

- **T-3.1 · Tools `get_student_history` / `get_past_adaptations` / `get_student_insights`** (ALZ-412)
  Implementar las 3 tools y su dispatcher en el usecase. `get_student_insights` degrada hasta que exista la memoria viva.

---

# HU-4 · Editar y publicar un prompt desde DB sin deploy (validación + fallback)  (ALZ-396)

**Como** equipo de producto **quiero** ajustar y publicar el comportamiento de Alizia sin esperar un deploy (con validación previa y fallback), **para** iterar rápido y sin riesgo de romper producción.

**Criterios de aceptación**
- Edito un `body`, lo publico (pasa la validación) y el runtime usa la nueva versión; las versiones no se pisan.
- Una versión inválida no se activa (queda en `draft`, sigue la activa).
- Si la versión activa falla en runtime, cae a la última versión buena o al prompt de código.
- El comportamiento observable no cambia respecto al prompt actual.

- **T-4.1 · Modelo de datos versionado (`prompt_templates` + `prompt_versions`)** (ALZ-413)
  `prompt_templates(key unique)` con `key ∈ {recommend, assist, guided}`; `prompt_versions(template_id, version, body, model, params, status∈{draft,active,archived})`, unique `(template_id, version)`. `ai_usage.prompt_version_id` referencia esta tabla.
- **T-4.2 · Renderer de templates (placeholders + cache)** (ALZ-414)
  Motor estilo Mustache/Handlebars (lib probada, p. ej. `raymond`): soporta bloque `{{x}}`, campo `{{x.y}}` y flag `{{#x}}…{{/x}}`. Cache en memoria de la versión activa; invalida al publicar.
- **T-4.3 · Validación al publicar (los 4 checks)** (ALZ-415)
  Al publicar verifica: (1) cada `{{x}}` existe en el catálogo; (2) flags balanceados; (3) `{{output_contract}}` intacto; (4) nada dinámico antes del corte de cache. Si falla, no se activa.
- **T-4.4 · Fallback a última versión buena** (ALZ-416)
  Si la versión activa falla al renderizar/ejecutar, usa la última versión buena conocida; si no hay, el prompt de código. Loguea el incidente.
- **T-4.5 · Migrar la capa 1 fuera de `prompts.go`** (ALZ-417)
  Los 3 builders dejan de hardcodear la capa 1; el `{{output_contract}}` (capa 2) queda en código. Sin cambio de comportamiento observable.
- **T-4.6 · Seed inicial de los 3 `body` (recommend / assist / guided)** (ALZ-418)
  Una `version` `active` por modo con el `body` real (no ilustrativo) que pasa la validación. Seed por script de entorno.

---

# HU-5 · La conversación larga conserva memoria (resumen)  (ALZ-397)

**Como** docente **quiero** que Alizia no pierda el hilo en conversaciones largas, **para** seguir trabajando sin recontextualizar todo.

**Criterio de aceptación**
- Cuando el historial excede el presupuesto de tokens, se conserva system + últimos N + un resumen comprimido de lo viejo, en vez de tirar los mensajes.

- **T-5.1 · Resumen de conversación (`conversation_summaries` + job)** (ALZ-419)
  `conversation_summaries` 1:1 con `conversations`; un job genera/actualiza el resumen; el bloque `{{conversation_summary}}` lo consume.
- **T-5.2 · Usar el resumen en `capMessages`** (ALZ-420)
  `capMessages` usa el resumen en lugar de descartar historial cuando excede el presupuesto.

---

# HU-6 · Memoria viva del alumno entre sesiones (insights)  (ALZ-398)

**Como** docente **quiero** que Alizia recuerde entre sesiones qué funcionó con mi alumno, **para** no explicárselo cada vez.

**Criterio de aceptación**
- `student_insights` (resumen + aprendizajes clave) 1 por alumno, regenerado por job batch desde sesiones y adaptaciones; el bloque `{{insights}}` lo consume.

- **T-6.1 · Memoria viva del alumno (`student_insights`) + job batch** (ALZ-421)
  Crear `student_insights` (`summary` + `key_learnings[]`) y el job batch que lo regenera. Habilita `get_student_insights`.

---

# HU-7 · Ejemplos golden alimentan el few-shot (cold-start)  (ALZ-450)

**Como** equipo de Alizia **quiero** alimentarla con ejemplos de buenas respuestas, **para** que dé recomendaciones de calidad desde el día 1 aunque no haya historial.

**Criterios de aceptación**
- Con ejemplos cargados se inyectan top-3 por relevancia **tras el corte de cache**.
- Sin ejemplos, solo lineamientos.

- **T-7.1 · `response_examples` + seed curated (~15 casos)** (ALZ-422)
  `response_examples` (origen curated/promovido, relevancia, modo) + seed curated de ~15 casos por script de entorno.
- **T-7.2 · Selección de few-shot golden e inyección tras el corte** (ALZ-423)
  El bloque `{{few_shot}}` se rellena con top-3 por relevancia y se ubica después del marcador de corte de cache.

---

# HU-8 · Medir si el sistema mejora (win-rate y métricas)  (ALZ-400)

**Como** equipo de Alizia **quiero** medir si las recomendaciones mejoran versión a versión, **para** decidir con datos qué cambios mantener.

**Criterio de aceptación**
- Se computan win-rate (% `funcionó`), tasa de aceptación implícita, cobertura de contexto, costo por turno y % cache-hit, y la mejora por versión. Requiere la traza de HU-1.

- **T-8.1 · Win-rate por versión + métricas de éxito** (ALZ-424)
  Implementar el cálculo de las métricas anteriores por versión de prompt.

---

# HU-9 · Promover golden y correr A/B entre versiones  (ALZ-399)

**Como** equipo de Alizia **quiero** probar dos versiones en paralelo y promover la mejor automáticamente, **para** que Alizia mejore sola sin degradar la calidad.

**Criterios de aceptación**
- 2 versiones activas con split de tráfico; antes de promover corre el set de eval y compara; la promoción queda condicionada por win-rate; el job es idempotente y aislado del request path.

- **T-9.1 · A/B entre versiones + eval antes de promover** (ALZ-425)
  Soporta 2 versiones activas con split; corre `response_examples` como eval y compara; promoción condicionada por win-rate.
- **T-9.2 · Job batch del flywheel (cron interno)** (ALZ-426)
  Lógica en `RunBatch(ctx)`, disparada por cron interno. Config por env (`FLYWHEEL_ENABLED` default `false`, `FLYWHEEL_CRON`). Idempotente (marca de agua), aislado (goroutine), observable (loguea procesadas/promovidos/duración).

---

# HU-10 · Contenido pedagógico base del prompt  (ALZ-402)

**Como** equipo pedagógico **quiero** dejar redactado el contenido base de Alizia (persona, marco, límites, formato, situaciones, ejemplos), **para** que sus respuestas sigan nuestros lineamientos de inclusión.

**Criterio de aceptación**
- Cada bloque existe y pasa el checklist del Apéndice A: persona, marco pedagógico, límites, formato de salida, situaciones y ejemplos golden de referencia.

- **T-10.1 · Identidad y persona** (ALZ-427)
  Quién es Alizia, a quién le habla, tono/registro, qué NO es y estilo de salida. Tono base literal.
- **T-10.2 · Marco pedagógico (3 ejes + Criterios de AliZia + DUA)** (ALZ-428)
  Texto real de los 3 ejes, Criterios de AliZia y marco DUA. Incluye "entrada pedagógica, no clínica".
- **T-10.3 · Límites duros / guardrails** (ALZ-429)
  Redacción oficial de los 3 límites (no diagnostica · no reemplaza al docente · no produce informes clínicos) y protocolos ante pedidos fuera de scope y situaciones de riesgo.
- **T-10.4 · Reglas de formato de salida** (ALZ-430)
  1–3 acciones, ≥3 niveles de diferenciación, "útil en <1 min" y la estructura por modo. Coherente con el `{{output_contract}}`.
- **T-10.5 · Las ~15 situaciones observables (seed)** (ALZ-431)
  Lista oficial de ~15 situaciones de aula con `code` + `name` (y `phase` si aplica), en formato consumible por el script de seed.
- **T-10.6 · Few-shot golden curated (~15 casos)** (ALZ-432)
  ~15 ejemplos golden (contexto → respuesta ideal), etiquetados por modo/situación, listos para el seed de `response_examples`.

---

## Dependencias (a nivel Historia)

- **HU-2 Contexto** habilita: HU-3 (tools), el renderer de **HU-4** y el few-shot de **HU-7**.
- **HU-1 Traza** + **HU-4 Prompts** habilitan **HU-8** (medir mejora).
- **HU-7** + **HU-8** habilitan **HU-9** (promover y A/B).
- **HU-10 Contenido** habilita el seed real de **T-4.6** (3 body) y **T-7.1** (curated).

## Estado en Jira (2026-06-02)

- Las 10 Historias cuelgan de **Chubut (ALZ-246)**, en **AI Sprint 36**, asignadas a Sebastian, con título HU-1…HU-10.
- Las 30 Subtareas cuelgan de su historia (keys ALZ-403…432, según el mapeo de arriba).
- **Pendiente cosmético:** los títulos de las 30 subtareas en Jira todavía arrastran numeración vieja (T-1.x/T-3.x); este documento define la numeración correcta T-1.1…T-10.6.
- **Pendiente de limpieza (UI, bulk-delete):** issues cancelados con `[BORRAR]` y el modelo CE viejo (label `context-engine`).
