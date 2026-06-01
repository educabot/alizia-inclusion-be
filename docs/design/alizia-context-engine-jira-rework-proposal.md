# Contenido literal para Jira — Context Engine (para validar antes de crear)

> **Estado:** PROPUESTA, no ejecutada. Es el texto EXACTO que irá a cada campo `Título` (summary) y `Descripción`.
> **Fecha:** 2026-06-01 · Proyecto ALZ.
> Jerarquía y códigos: **EP-{n}** Épica · **HU-{n}.{m}** Historia · **T-{n}.{k}** Tarea (corrida dentro de la épica `n`).
> 6 Épicas (ya creadas: ALZ-266..271) · 10 Historias · 30 Subtareas (a crear al aprobar).

---

# EP-1 · Traza  (ALZ-266)
**Descripción:** Registrar, en cada turno, **qué prompt + qué contexto → qué resultado**, sin cambiar el comportamiento. Es la base de todo lo demás: sin traza no hay datos para mejorar. Spec: `docs/design/alizia-context-engine.md` §5, §6.4, §7.

## HU-1.1 · Trazar cada turno y la aceptación implícita
**Descripción:**
Objetivo: correlacionar prompt+contexto con el resultado y capturar si el docente aceptó la sugerencia, sin pedirle nada extra.
Criterios de aceptación:
- Cada turno guarda en `ai_usage`: modelo, latencia, tool calls, `conversation_id`, `message_id` y un snapshot de contexto con **IDs (no PII en claro)**.
- Una adaptación creada desde una sugerencia guarda su origen y queda marcada según se editó o no antes de guardar.
- El tablero del Director sigue funcionando; filas viejas con los campos nuevos en `NULL`.
- Registro best-effort: si falla, no bloquea la respuesta.

- **T-1.1 · Enriquecer `ai_usage` con columnas de traza**
  `ALTER ai_usage`: `latency_ms`, `model`, `tool_calls`, `conversation_id`, `message_id`, `context_snapshot` (IDs, no PII). Filas viejas en `NULL`.
- **T-1.2 · Ligar adaptaciones a su origen + señal de aceptación implícita**
  La adaptación persiste `source_conversation_id` y `source_message_id`; `was_edited=false` si se guardó tal cual, `true` si se modificó.

---

# EP-2 · Contexto del alumno y del docente  (ALZ-267)
**Descripción:** Que Alizia **sepa con quién habla y de quién habla** (docente; alumno con perfil, situaciones observables, diagnósticos opcionales, PPI, entorno) y que un **Context Assembler** lo junte todo en orden cacheable. Spec: §6.2, §7, §8.

## HU-2.1 · Alizia usa el contexto del alumno en su respuesta
**Descripción:**
Objetivo: que el prompt de cada turno incluya el contexto disponible del alumno y del docente, y que degrade sin romper cuando falta.
Criterios de aceptación:
- Con contexto cargado, la respuesta/el prompt lo refleja.
- Campos ausentes (todos opcionales) no imprimen "N/A" ni rompen.
- No se filtra PII a logs.

- **T-2.1 · Perfil del docente (`teacher_profiles`)**
  `teacher_profiles` 1:1 con `users`: `age_range`/`birthdate` (nullable), `years_experience`, `specialization`, `subjects[]`, `tone_preference`, `bio`. Aislamiento por organización.
- **T-2.2 · Enriquecer alumno (`students` / `student_profiles`)**
  `students`: `birthdate`, `age_range`, `grade_level`, `preferred_name` (nullable). `student_profiles`: `support_level`, `strengths[]`, `interests[]`, `triggers[]`, `effective_strategies[]`, `ineffective_strategies[]`, `situation_codes[]`, `has_therapeutic_companion`, `environment_notes`. Todo opcional.
- **T-2.3 · Catálogo de situaciones observables (`situations_catalog`)**
  Tabla global (con `organization_id` para per-org futuro), `phase` nullable. El seed de las ~15 situaciones se carga por **script de entorno**, no en la migración.
- **T-2.4 · Diagnósticos estructurados opcionales**
  `diagnoses_catalog` (global) + `student_diagnoses` (M2M con `severity`). Opcional y secundario a las situaciones; Alizia puede sugerir, nunca exigir. Seed por script de entorno.
- **T-2.5 · Proyecto Pedagógico Individual (`ppi`)**
  `ppi` 1 por alumno, todos los campos nullable. Cuando existe, es contexto de primera línea.
- **T-2.6 · Rol maestra integradora + asignación**
  Sumar `maestra_integradora` al enum `member_role` y modelar la asignación integradora↔alumno. RBAC fuera de scope (otro equipo).
- **T-2.7 · Context Assembler (`BuildPromptContext`)**
  Usecase que junta todo el contexto en un struct tipado, ordenado para caching (prefijo invariante adelante, variable atrás). Punto de convergencia de las tareas anteriores.

## HU-2.2 · Profundizar bajo demanda con tools agénticas
**Descripción:**
Objetivo: que el modelo pueda pedir más datos de un alumno bajo demanda sin inflar el prompt base.
Criterio de aceptación: disponibles `get_student_history`, `get_past_adaptations`, `get_student_insights`; el dispatcher vive en el usecase (clean architecture).

- **T-2.8 · Tools `get_student_history` / `get_past_adaptations` / `get_student_insights`**
  Implementar las 3 tools y su dispatcher en el usecase. `get_student_insights` degrada hasta que exista la memoria viva.

---

# EP-3 · Prompts versionados en DB  (ALZ-268)
**Descripción:** Sacar la **capa 1** del prompt (persona, lineamientos, few-shot, params) del código a una tabla versionada, y rellenarla en runtime con un renderer seguro (cache + validación + fallback). El contrato de salida y el motor quedan en código. Spec: §9, §6.1, §6.3.

## HU-3.1 · Editar y publicar un prompt desde DB sin deploy (validación + fallback)
**Descripción:**
Objetivo: iterar los prompts sin deploy, comparando versiones y sin riesgo de tumbar producción.
Criterios de aceptación:
- Edito un `body`, lo publico (pasa la validación) y el runtime usa la nueva versión; las versiones no se pisan.
- Una versión inválida no se activa (queda en `draft`, sigue la activa).
- Si la versión activa falla en runtime, cae a la última versión buena o al prompt de código.
- El comportamiento observable no cambia respecto al prompt actual.

- **T-3.1 · Modelo de datos versionado (`prompt_templates` + `prompt_versions`)**
  `prompt_templates(key unique)` con `key ∈ {recommend, assist, guided}`; `prompt_versions(template_id, version, body, model, params, status∈{draft,active,archived})`, unique `(template_id, version)`. `ai_usage.prompt_version_id` referencia esta tabla.
- **T-3.2 · Renderer de templates (placeholders + cache)**
  Motor estilo Mustache/Handlebars (lib probada, p. ej. `raymond`): soporta bloque `{{x}}`, campo `{{x.y}}` y flag `{{#x}}…{{/x}}`. Cache en memoria de la versión activa; invalida al publicar.
- **T-3.3 · Validación al publicar (los 4 checks)**
  Al publicar verifica: (1) cada `{{x}}` existe en el catálogo; (2) flags balanceados; (3) `{{output_contract}}` intacto; (4) nada dinámico antes del corte de cache. Si falla, no se activa.
- **T-3.4 · Fallback a última versión buena**
  Si la versión activa falla al renderizar/ejecutar, usa la última versión buena conocida; si no hay, el prompt de código. Loguea el incidente.
- **T-3.5 · Migrar la capa 1 fuera de `prompts.go`**
  Los 3 builders dejan de hardcodear la capa 1; el `{{output_contract}}` (capa 2) queda en código. Sin cambio de comportamiento observable.
- **T-3.6 · Seed inicial de los 3 `body` (recommend / assist / guided)**
  Una `version` `active` por modo con el `body` real (no ilustrativo) que pasa la validación. Seed por script de entorno.

---

# EP-4 · Memoria  (ALZ-269)
**Descripción:** Dejar de **descartar** el historial viejo por presupuesto de tokens; resumir y mantener una memoria viva por alumno que se inyecte en el prompt. Spec: §7 Capa B, §9.1.

## HU-4.1 · La conversación larga conserva memoria (resumen)
**Descripción:**
Objetivo: no perder contexto en conversaciones largas.
Criterio de aceptación: cuando el historial excede el presupuesto, se conserva system + últimos N + un resumen comprimido de lo viejo, en vez de tirar los mensajes.

- **T-4.1 · Resumen de conversación (`conversation_summaries` + job)**
  `conversation_summaries` 1:1 con `conversations`; un job genera/actualiza el resumen; el bloque `{{conversation_summary}}` lo consume.
- **T-4.2 · Usar el resumen en `capMessages`**
  `capMessages` usa el resumen en lugar de descartar historial cuando excede el presupuesto.

## HU-4.2 · Memoria viva del alumno entre sesiones (insights)
**Descripción:**
Objetivo: anclar el saludo y personalizar sin re-preguntar.
Criterio de aceptación: `student_insights` (resumen + aprendizajes clave) 1 por alumno, regenerado por job batch desde sesiones y adaptaciones; el bloque `{{insights}}` lo consume.

- **T-4.3 · Memoria viva del alumno (`student_insights`)**
  Crear `student_insights` (`summary` + `key_learnings[]`) y el job batch que lo regenera. Habilita `get_student_insights`.

---

# EP-5 · Flywheel (mejora continua)  (ALZ-270)
**Descripción:** Cerrar el loop: usar resultado de aula + aceptación implícita para promover ejemplos golden, medir win-rate por versión y correr A/B + eval. Spec: §10, §13.

## HU-5.1 · Ejemplos golden alimentan el few-shot (cold-start)
**Descripción:**
Objetivo: arrancar con buenos ejemplos desde el día 1 y personalizar por alumno sin romper el cache.
Criterios de aceptación: con ejemplos cargados se inyectan top-3 por relevancia **tras el corte de cache**; sin ejemplos, solo lineamientos.

- **T-5.1 · `response_examples` + seed curated (~15 casos)**
  `response_examples` (origen curated/promovido, relevancia, modo) + seed curated de ~15 casos por script de entorno.
- **T-5.2 · Selección de few-shot golden e inyección tras el corte**
  El bloque `{{few_shot}}` se rellena con top-3 por relevancia y se ubica después del marcador de corte de cache.

## HU-5.2 · Medir si el sistema mejora (win-rate y métricas)
**Descripción:**
Objetivo: saber si el sistema mejora versión a versión.
Criterio de aceptación: se computan win-rate (% `funcionó`), tasa de aceptación implícita, cobertura de contexto, costo por turno y % cache-hit, y la mejora por versión.

- **T-5.3 · Win-rate por versión + métricas de éxito**
  Implementar el cálculo de las métricas anteriores por versión de prompt.

## HU-5.3 · Promover golden y correr A/B entre versiones
**Descripción:**
Objetivo: que el loop gire solo, sin afectar la API y sin degradar calidad.
Criterios de aceptación: 2 versiones activas con split de tráfico; antes de promover corre el set de eval y compara; la promoción queda condicionada por win-rate; el job es idempotente y aislado del request path.

- **T-5.4 · A/B entre versiones + eval antes de promover**
  Soporta 2 versiones activas con split; corre `response_examples` como eval y compara; promoción condicionada por win-rate.
- **T-5.5 · Job batch del flywheel (cron interno)**
  Lógica en `RunBatch(ctx)`, disparada por cron interno. Config por env (`FLYWHEEL_ENABLED` default `false`, `FLYWHEEL_CRON`). Idempotente (marca de agua), aislado (goroutine), observable (loguea procesadas/promovidos/duración).

---

# EP-6 · Contenido pedagógico base  (ALZ-271)
**Descripción:** Producir el **texto real** que rellena los `body` (capa 1) y los ejemplos golden. Es el insumo del seed de prompts y del curated del flywheel. Lo redacta producto/pedagogía; no espera código. Spec: Apéndice A (checklist), Apéndice B (esqueletos).

## HU-6.1 · Contenido pedagógico base del prompt
**Descripción:**
Objetivo: dejar listo y revisado todo el contenido literal que alimenta el prompt.
Criterio de aceptación: cada bloque existe y pasa el checklist del Apéndice A.

- **T-6.1 · Identidad y persona**
  Quién es Alizia, a quién le habla, tono/registro, qué NO es y estilo de salida. Tono base literal.
- **T-6.2 · Marco pedagógico (3 ejes + Criterios de AliZia + DUA)**
  Texto real de los 3 ejes, Criterios de AliZia y marco DUA. Incluye "entrada pedagógica, no clínica".
- **T-6.3 · Límites duros / guardrails**
  Redacción oficial de los 3 límites (no diagnostica · no reemplaza al docente · no produce informes clínicos) y protocolos ante pedidos fuera de scope y situaciones de riesgo.
- **T-6.4 · Reglas de formato de salida**
  1–3 acciones, ≥3 niveles de diferenciación, "útil en <1 min" y la estructura por modo. Coherente con el `{{output_contract}}`.
- **T-6.5 · Las ~15 situaciones observables (seed)**
  Lista oficial de ~15 situaciones de aula con `code` + `name` (y `phase` si aplica), en formato consumible por el script de seed.
- **T-6.6 · Few-shot golden curated (~15 casos)**
  ~15 ejemplos golden (contexto → respuesta ideal), etiquetados por modo/situación, listos para el seed de `response_examples`.

---

## Dependencias (a nivel Épica/Historia)
- **EP-2 Contexto** habilita: HU-2.2 (tools), el renderer de **EP-3** y el few-shot de **EP-5**.
- **EP-1 Traza** + **EP-3 Prompts** habilitan **HU-5.2** (medir mejora).
- **HU-5.1** + **HU-5.2** habilitan **HU-5.3** (promover y A/B).
- **EP-6 Contenido** habilita el seed real de **T-3.6** (3 body) y **T-5.1** (curated).

## Pendiente al aprobar
1. Crear las 10 Historias (HU-*) + 30 Tareas (T-*) con este texto exacto, con su código en el título.
2. Limpiar las descripciones de las 6 Épicas (hoy todavía dicen "Fase N / CE / EPIC" en el cuerpo) y ponerles el código `EP-{n}` en el título.
3. Marcar `[BORRAR]` las 30 historias viejas (ALZ-272..301) + 91 subtareas (ALZ-302..392) para bulk-delete.
4. Asignarte todo y actualizar el ledger.
