# Backlog — Motor de Contexto, Trazabilidad y Self-Improvement de Alizia

> **Documento de especificación (la fuente de verdad):** [`alizia-context-engine.md`](./alizia-context-engine.md)
>
> Este backlog se desprende de ese design doc. **Toda duda de "cómo se desarrolla a fondo" se
> resuelve ahí** — cada épica e historia linkea a la sección exacta del spec donde está el
> detalle (modelo de datos, DBML, orden de bloques, validaciones, etc.). Acá está el *qué* y el
> *en qué orden*; el *cómo* está en el spec.

---

## Cómo leer este documento

- **Épica** = una **fase** del [§12 — Fases / roadmap](./alizia-context-engine.md#12-fases--roadmap-insumo-para-el-goal). Cada fase entrega valor sola y respeta migraciones inmutables/incrementales.
- **Historia** = un slice con valor demostrable (`Como… quiero… para…` + criterios de aceptación).
- **Sub-tareas** = trabajo técnico concreto, con referencia a migración/archivo del spec.
- **Mapeo a Jira:** Épica → Story → Sub-task. `Component`, `Label` y `Fix Version` sugeridos abajo.

### Convenciones para Jira

| Campo Jira | Valor sugerido |
|---|---|
| **Fix Version / Release** | La fase (`Fase 0`, `Fase 1`, …). El roadmap del §12 *es* el plan de releases. |
| **Component** | `data-model` · `context-assembler` · `renderer` · `memory` · `flywheel` · `content` |
| **Label** | `migration-000xxx` · `risk:low/med/high` · `parallelizable` · `content` · `cold-start` |
| **Priority** | = orden de fase (ya viene dado; no reinventar). |
| **Estimación** | T-shirt (S/M/L). El riesgo del §12 va como label, no como puntos. |

### Leyenda de la ficha de historia

- **🔗 Depende de:** historias que deben estar listas antes.
- **⚡ Paralelizable:** puede avanzar en simultáneo con lo indicado.
- **📄 Spec:** sección del design doc con el detalle.
- **✅ Criterios de aceptación / 🔧 Sub-tareas.**

---

## Mapa de épicas

| Épica | Fase | Migraciones | Riesgo | Component principal |
|---|---|---|---|---|
| [EPIC 0 · Traza](#epic-0--traza) | Fase 0 | 000014–000015 | bajo | `data-model` |
| [EPIC 1 · Contexto](#epic-1--contexto-del-alumno-y-del-docente) | Fase 1 | 000017–000021 | medio | `data-model` + `context-assembler` |
| [EPIC 2 · Prompts en DB](#epic-2--prompts-versionados-en-db) | Fase 2 | 000022 | medio | `renderer` |
| [EPIC 3 · Memoria](#epic-3--memoria) | Fase 3 | 000023 | medio | `memory` |
| [EPIC 4 · Flywheel](#epic-4--flywheel-self-improvement) | Fase 4 | 000024 | alto | `flywheel` |
| [EPIC 5 · Contenido pedagógico](#epic-5--contenido-pedagógico-base) | Transversal | — (autoría) | — | `content` |

### Paralelización (dos streams desde el día 1)

```text
 STREAM BACKEND (código + migraciones, secuencial por dependencia)
 ┌────────┐   ┌──────────────────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
 │ EPIC 0 │──►│ EPIC 1               │──►│ EPIC 2   │──►│ EPIC 3   │──►│ EPIC 4   │
 │ Traza  │   │ Contexto             │   │ Prompts  │   │ Memoria  │   │ Flywheel │
 └────────┘   │ (tablas en paralelo) │   │ en DB    │   └──────────┘   └────┬─────┘
              └──────────────────────┘   └──────────┘                       │
                                                                            │
 STREAM CONTENIDO (autoría, SIN dependencia de código — arranca ya)         │
 ┌─────────────────────────────────────────────────────┐                   │
 │ EPIC 5 · 3 ejes · Criterios AliZia · límites duros · │── alimenta ──► seed Fase 2
 │ ~15 situaciones · few-shot golden curated            │── alimenta ──► EPIC 4 (curated)
 └─────────────────────────────────────────────────────┘                   ▲
                                                                            │
                                   (EPIC 4 necesita: traza F0 + versiones F2 + contenido C)
```

**Reglas de oro:**
- **EPIC 5 arranca primero** — lo hace producto/pedagogía, no espera código. Es el cuello de botella escondido (texto real de los `body` + ~15 casos golden).
- **EPIC 0** es chico y de bajo riesgo: ideal para abrir el stream backend mientras se diseña Fase 1.
- **Dentro de EPIC 1**, las migraciones de tablas son independientes entre sí → varios devs en paralelo. El **Context Assembler es el punto de convergencia** (va al final del epic).
- **EPIC 4** junta todas las dependencias (traza + versiones + contenido) → último y de mayor riesgo.

---

## EPIC 0 · Traza

> **Fase 0** · Riesgo **bajo** · Migraciones **000014–000015** · Component `data-model`
> **Objetivo:** registrar *qué prompt + qué contexto → qué resultado* **sin cambiar el comportamiento**. Es la base del self-improvement: sin traza no hay datos de los que aprender.
> **📄 Spec:** [§5 — Evolucionar `ai_usage` in-place](./alizia-context-engine.md#5-decisiones-de-diseño) · [§6.4 Agentic Run + traza](./alizia-context-engine.md#64-agentic-run--traza) · [§7 Modelo de datos](./alizia-context-engine.md#7-modelo-de-datos)

### CE-0.1 · Enriquecer `ai_usage` con columnas de traza
**Como** plataforma **quiero** registrar en cada turno la versión de prompt, modelo, latencia, tool calls y a qué conversación/mensaje pertenece **para** poder correlacionar prompt+contexto con el resultado.
`risk:low` · `migration-000014`

- **🔗 Depende de:** — (puede arrancar ya)
- **⚡ Paralelizable:** con todo EPIC 5 y con CE-0.2.
- **📄 Spec:** [§7 — ALTER `ai_usage`](./alizia-context-engine.md#7-modelo-de-datos) (lista de columnas nullable).
- **✅ Criterios de aceptación:**
  - Dado un turno que llama al modelo, cuando termina, entonces `ai_usage` guarda `latency_ms`, `model`, `tool_calls`, `conversation_id`, `message_id` y `context_snapshot`.
  - `context_snapshot` guarda **IDs, no PII en claro**.
  - El tablero del Director (`get_ai_usage.go`) sigue funcionando; filas viejas quedan con los campos nuevos en `NULL`.
  - El registro es **best-effort**: si falla, no bloquea la respuesta.
- **🔧 Sub-tareas:**
  - Migración `000014`: `ALTER ai_usage` agregando columnas de traza, **todas nullable** (`conversation_id`, `message_id`, `prompt_version_id` nullable sin FK aún, `model`, `latency_ms`, `tool_calls`, `context_snapshot jsonb`).
  - Actualizar entidad `ai_usage` + repositorio.
  - `recordAIUsage` (`usage.go`) escribe los nuevos campos.
  - Regenerar mocks (`make mocks`) + tests.

### CE-0.2 · Ligar adaptaciones a su origen + señal de aceptación implícita
**Como** plataforma **quiero** saber de qué conversación/mensaje salió una adaptación y si el docente la editó antes de guardar **para** capturar la señal de feedback sin pedirle nada extra al docente.
`risk:low` · `migration-000015`

- **🔗 Depende de:** —
- **⚡ Paralelizable:** con CE-0.1 y EPIC 5.
- **📄 Spec:** [§5 — Feedback por resultado + aceptación implícita](./alizia-context-engine.md#5-decisiones-de-diseño) · [§7 Capa C](./alizia-context-engine.md#7-modelo-de-datos).
- **✅ Criterios de aceptación:**
  - Dada una adaptación creada a partir de una sugerencia de IA, entonces se persisten `source_conversation_id` y `source_message_id`.
  - Cuando el docente guarda la adaptación tal cual la propuso Alizia, entonces `was_edited = false`; si la modificó, `was_edited = true`.
  - `status` (`en_curso`/`funcionó`/`para_ajustar`) y `outcome` siguen existiendo y son el resultado de aula.
- **🔧 Sub-tareas:**
  - Migración `000015`: `ALTER adaptations` (`source_conversation_id`, `source_message_id`, `was_edited boolean`).
  - Entidad + repositorio de `adaptations`.
  - Capturar `was_edited` en el flujo de guardado (comparar propuesta vs. lo guardado).
  - Tests.

---

## EPIC 1 · Contexto del alumno y del docente

> **Fase 1** · Riesgo **medio** · Migraciones **000017–000021** · Component `data-model` + `context-assembler`
> **Objetivo:** que Alizia **sepa con quién habla y de quién habla** — docente, alumno (perfil + situaciones observables + diagnósticos opcionales), PPI y entorno — y que un **Context Assembler** lo junte todo en orden cacheable.
> **📄 Spec:** [§6.2 Context Assembler](./alizia-context-engine.md#62-context-assembler) · [§7 Modelo de datos](./alizia-context-engine.md#7-modelo-de-datos) · [§8 El Context Assembler en detalle](./alizia-context-engine.md#8-el-context-assembler-en-detalle)

> **Nota de paralelización:** CE-1.1 a CE-1.6 son **migraciones independientes** → repartibles entre devs. CE-1.7 (Assembler) es el **punto de convergencia** y va al final.
> **Nota de catálogos/seeds:** los catálogos (`situations_catalog`, `diagnoses_catalog`) **no llevan seed en la migración** (regla del repo). El seed va por **script de entorno** — ver CE-5.5.

### CE-1.1 · Perfil del docente (`teacher_profiles`)
**Como** Alizia **quiero** conocer rango/edad, experiencia, materias y tono preferido del docente **para** adaptar el registro y las sugerencias a quién tengo enfrente.
`risk:med` · `migration-000017` · `parallelizable`

- **🔗 Depende de:** —
- **⚡ Paralelizable:** con CE-1.2…CE-1.6.
- **📄 Spec:** [§7 Capa A — `teacher_profiles`](./alizia-context-engine.md#7-modelo-de-datos).
- **✅ Criterios de aceptación:**
  - Existe `teacher_profiles` 1:1 con `users`, con `age_range` **y** `birthdate` (ambos nullable, doble granularidad), `years_experience`, `specialization`, `subjects text[]`, `tone_preference`, `bio`.
  - Aislamiento por organización respetado.
- **🔧 Sub-tareas:** migración `000017` · entidad + provider `TeacherProfileProvider` · path de carga/edición (sin UI, vía API existente) · mocks + tests.

### CE-1.2 · Enriquecer alumno (`students` + `student_profiles`)
**Como** Alizia **quiero** edad, grado, nombre preferido, fortalezas, intereses, disparadores, estrategias que funcionan/no, situaciones observables y entorno del alumno **para** personalizar la adaptación.
`risk:med` · `migration-000018` · `parallelizable`

- **🔗 Depende de:** —
- **⚡ Paralelizable:** con el resto de CE-1.x.
- **📄 Spec:** [§7 Capa A — `students` / `student_profiles` ALTER](./alizia-context-engine.md#7-modelo-de-datos).
- **✅ Criterios de aceptación:**
  - `students` suma `birthdate`, `age_range`, `grade_level`, `preferred_name` (**todos nullable**).
  - `student_profiles` suma `support_level`, `strengths[]`, `interests[]`, `triggers[]`, `effective_strategies[]`, `ineffective_strategies[]`, `situation_codes[]` (vocabulario del `situations_catalog`), `has_therapeutic_companion`, `environment_notes` (**todos nullable**); se mantiene `difficulties[]` libre.
  - Todo opcional: el sistema degrada sin romper si los campos vienen vacíos.
- **🔧 Sub-tareas:** migración(es) `000018` · entidades + providers · mocks + tests.

### CE-1.3 · Catálogo de situaciones observables (`situations_catalog`)
**Como** Alizia **quiero** un vocabulario controlado de situaciones de aula (~15: "no inicia la tarea", "se distrae constantemente", …) **para** partir de lo **observable**, no del diagnóstico (entrada pedagógica, no clínica).
`risk:med` · `migration-000019` · `parallelizable`

- **🔗 Depende de:** — (tabla); **el contenido de las ~15 situaciones lo provee CE-5.5.**
- **⚡ Paralelizable:** con el resto de CE-1.x.
- **📄 Spec:** [§7 — `situations_catalog`](./alizia-context-engine.md#7-modelo-de-datos) · [§9.3/§9.5 bloque `{{situations_catalog}}`](./alizia-context-engine.md#93-catálogo-de-variables-placeholders).
- **✅ Criterios de aceptación:**
  - Existe `situations_catalog` (global, con `organization_id` para per-org futuro), con `phase` nullable.
  - El seed de las situaciones se carga **por script de entorno**, no en la migración.
- **🔧 Sub-tareas:** migración `000019` (estructura, sin datos) · entidad + provider · **script de entorno** para el seed (consume CE-5.5) · tests.

### CE-1.4 · Diagnósticos estructurados opcionales (`diagnoses_catalog` + `student_diagnoses`)
**Como** Alizia **quiero** una capa estructurada de diagnósticos **opcional y secundaria** a las situaciones **para** sumar detalle solo cuando la escuela lo brinda.
`risk:med` · `migration-000020` · `parallelizable`

- **🔗 Depende de:** —
- **⚡ Paralelizable:** con el resto de CE-1.x.
- **📄 Spec:** [§7 — `diagnoses_catalog` / `student_diagnoses`](./alizia-context-engine.md#7-modelo-de-datos).
- **✅ Criterios de aceptación:**
  - `diagnoses_catalog` (global, `organization_id` para per-org) + `student_diagnoses` (M2M con severity).
  - Se llena solo si la escuela aporta diagnóstico; Alizia puede **sugerir** etiquetas, nunca exigirlas.
  - Seed del catálogo por script de entorno.
- **🔧 Sub-tareas:** migración `000020` · entidades + providers · script de entorno del catálogo · tests.

### CE-1.5 · Proyecto Pedagógico Individual (`ppi`)
**Como** docente **quiero** que Alizia conozca el PPI del alumno (objetivos, adaptaciones curriculares, seguimiento) **para** que las sugerencias mantengan coherencia con él.
`risk:med` · `migration-000021` · `parallelizable`

- **🔗 Depende de:** —
- **⚡ Paralelizable:** con el resto de CE-1.x.
- **📄 Spec:** [§7 — `ppi`](./alizia-context-engine.md#7-modelo-de-datos) · [§9.5 bloque `{{ppi}}`](./alizia-context-engine.md#95-cómo-renderiza-cada-bloque-dinámico).
- **✅ Criterios de aceptación:**
  - `ppi` 1 por alumno, **todos los campos nullable** (el PPI puede no existir).
  - Cuando existe, es contexto de primera línea (Alizia asiste en redacción de objetivos y mantiene coherencia).
- **🔧 Sub-tareas:** migración `000021` · entidad + provider · tests.

### CE-1.6 · Rol maestra integradora + asignación
**Como** plataforma **quiero** modelar el rol **maestra integradora** y su asignación a un alumno **para** que el motor de contexto sepa con quién habla y qué alumno tiene a cargo.
`risk:med` · `parallelizable`

- **🔗 Depende de:** —
- **⚡ Paralelizable:** con el resto de CE-1.x.
- **📄 Spec:** [§7 — Nota Roles (`member_role`)](./alizia-context-engine.md#7-modelo-de-datos).
- **✅ Criterios de aceptación:**
  - El enum `member_role` incluye `maestra_integradora` (hoy: `teacher`, `coordinator`, `admin`, `ministerio`, `psicopedagogo`).
  - Existe el modelado de la asignación **integradora ↔ alumno**.
  - **RBAC sigue fuera de scope** (otro equipo): acá solo se habilita que el dato exista y llegue al prompt.
- **🔧 Sub-tareas:** migración del enum + tabla/columna de asignación · entidad + provider · tests.

### CE-1.7 · Context Assembler (`BuildPromptContext`)  ⟵ convergencia
**Como** sistema **quiero** un usecase que junte **todo** el contexto disponible en un struct tipado, ordenado para caching **para** alimentar el prompt de cada turno.
`risk:med`

- **🔗 Depende de:** CE-1.1, CE-1.2, CE-1.3, CE-1.4, CE-1.5, CE-1.6 (necesita las tablas).
- **⚡ Paralelizable:** el **andamiaje** del struct puede empezar antes; la integración real espera las tablas.
- **📄 Spec:** [§6.2 Context Assembler](./alizia-context-engine.md#62-context-assembler) · [§8 orden de bloques + corte de cache](./alizia-context-engine.md#8-el-context-assembler-en-detalle).
- **✅ Criterios de aceptación:**
  - `BuildPromptContext(ctx, orgID, userID, classroomID, studentID, mode)` devuelve `PromptContext` con prefijo invariante adelante y variable atrás (orden del §8).
  - Cada campo nullable ausente **degrada sin romper** (no imprime "N/A").
  - No filtra PII a logs.
- **🔧 Sub-tareas:** usecase `BuildPromptContext` · struct `PromptContext` · orquestar providers de CE-1.x · respetar orden cacheable (§8) · tests con contexto pobre y rico.

### CE-1.8 · Tools agénticas para profundizar bajo demanda
**Como** modelo **quiero** poder pedir historial/adaptaciones/insights de un alumno bajo demanda **para** profundizar sin inflar el prompt base.
`risk:med`

- **🔗 Depende de:** CE-1.7 (Assembler). `get_student_insights` degrada hasta que exista CE-3.2.
- **📄 Spec:** [§6.2 — Tools nuevas (`agentic.go`)](./alizia-context-engine.md#62-context-assembler).
- **✅ Criterios de aceptación:**
  - Disponibles `get_student_history(student_id)`, `get_past_adaptations(student_id)`, `get_student_insights(student_id)`.
  - El dispatcher vive en el usecase (no en handler) — clean architecture.
- **🔧 Sub-tareas:** definir tools · dispatcher en usecase · tests con mock AI client.

---

## EPIC 2 · Prompts versionados en DB

> **Fase 2** · Riesgo **medio** · Migración **000022** · Component `renderer`
> **Objetivo:** sacar la **capa 1** (persona, lineamientos, few-shot, params) del código a una tabla versionada, y rellenarla en runtime con un renderer seguro (cache + validación + fallback). El contrato de salida (capa 2) y el motor (capa 3) quedan en código.
> **📄 Spec:** [§9 Prompts versionados en DB](./alizia-context-engine.md#9-prompts-versionados-en-db) · [§6.1 Las 3 capas](./alizia-context-engine.md#61-las-tres-capas-de-un-prompt-el-reframe-del-híbrido) · [§6.3 Prompt Renderer](./alizia-context-engine.md#63-prompt-renderer)

### CE-2.1 · Modelo de datos versionado (`prompt_templates` + `prompt_versions`)
**Como** plataforma **quiero** guardar el contenido editable del prompt versionado **para** iterarlo sin deploy y comparar versiones.
`risk:med` · `migration-000022`

- **🔗 Depende de:** — (puede arrancar apenas se decide el esquema).
- **📄 Spec:** [§7 Capa D — DBML `prompt_templates`/`prompt_versions`](./alizia-context-engine.md#7-modelo-de-datos).
- **✅ Criterios de aceptación:**
  - `prompt_templates(key unique)` con `key ∈ {recommend, assist, guided}`.
  - `prompt_versions(template_id, version, body, model, params, status)`, `status ∈ {draft, active, archived}`, unique `(template_id, version)`.
  - Las versiones **no se pisan**: se crea nueva y se promueve.
  - `ai_usage.prompt_version_id` ahora referencia `prompt_versions` (cerrar la FK pendiente de CE-0.1).
- **🔧 Sub-tareas:** migración `000022` · entidades + provider · cerrar FK `ai_usage.prompt_version_id` · tests.

### CE-2.2 · Renderer de templates (3 tipos de placeholder + cache)
**Como** sistema **quiero** un motor estilo Mustache/Handlebars que rellene el `body` activo con el `PromptContext` **para** producir el system prompt de cada turno.
`risk:med`

- **🔗 Depende de:** CE-2.1 + CE-1.7 (rellena desde `PromptContext`).
- **📄 Spec:** [§6.3 Prompt Renderer](./alizia-context-engine.md#63-prompt-renderer) · [§9.3 Catálogo de variables](./alizia-context-engine.md#93-catálogo-de-variables-placeholders) · [§9.5 Cómo renderiza cada bloque](./alizia-context-engine.md#95-cómo-renderiza-cada-bloque-dinámico).
- **✅ Criterios de aceptación:**
  - Soporta los **3 tipos**: bloque `{{x}}`, campo `{{x.y}}`, flag `{{#x}}…{{/x}}`.
  - Usa una **lib probada** (p. ej. `raymond`), no sintaxis inventada.
  - **Cache en memoria** de la versión activa (no lee DB por request); invalida al publicar.
- **🔧 Sub-tareas:** integrar lib · mapear catálogo §9.3 → campos de `PromptContext` · sub-formatos §9.5 · cache + invalidación · tests.

### CE-2.3 · Validación al publicar (los 4 checks)
**Como** plataforma **quiero** validar una versión antes de activarla **para** que un edit no rompa el runtime ni el parser.
`risk:med`

- **🔗 Depende de:** CE-2.1, CE-2.2.
- **📄 Spec:** [§6.3 — Validación al publicar](./alizia-context-engine.md#63-prompt-renderer) · [§9 — diagrama edición→validación→runtime](./alizia-context-engine.md#9-prompts-versionados-en-db).
- **✅ Criterios de aceptación:** al intentar publicar, se verifica: (1) cada `{{x}}` existe en el catálogo §9.3; (2) flags balanceados; (3) `{{output_contract}}` intacto (no editable); (4) **nada dinámico antes del corte de cache**. Si algo falla, **no se activa** (queda en `draft`, sigue la activa).
- **🔧 Sub-tareas:** validador con los 4 checks · feedback de error claro · tests de cada check (caso ok + caso que rechaza).

### CE-2.4 · Fallback a última versión buena + red en código
**Como** plataforma **quiero** caer a la última versión buena (o al prompt de código) si la activa falla en runtime **para** que un prompt malo nunca tumbe producción.
`risk:med`

- **🔗 Depende de:** CE-2.2.
- **📄 Spec:** [§6.3 — Fallback](./alizia-context-engine.md#63-prompt-renderer) · [§14 Riesgos](./alizia-context-engine.md#14-riesgos-y-mitigaciones).
- **✅ Criterios de aceptación:** si la versión activa falla al renderizar/ejecutar, el sistema usa la última versión buena conocida; si tampoco hay, el prompt de código como red final. Se loguea el incidente.
- **🔧 Sub-tareas:** detección de fallo de render · selección de "última versión buena" · red de seguridad en código · tests de fallo forzado.

### CE-2.5 · Migrar la capa 1 fuera de `prompts.go`
**Como** equipo **quiero** mover persona/lineamientos/few-shot/params de `prompts.go` a DB **para** dejar de necesitar deploy para iterar el prompt.
`risk:med`

- **🔗 Depende de:** CE-2.1…CE-2.4.
- **📄 Spec:** [§9.2 Descomposición de los prompts actuales](./alizia-context-engine.md#92-descomposición-de-los-prompts-actuales-en-las-3-capas).
- **✅ Criterios de aceptación:** los 3 builders (`buildRecommendSystemPrompt`, `buildAssistSystemPrompt`, `buildGuidedAssistPrompt`) dejan de hardcodear la capa 1; el `{{output_contract}}` (capa 2) queda en código; el comportamiento observable no cambia respecto al prompt actual.
- **🔧 Sub-tareas:** extraer capa 1 a `prompt_versions` · dejar capa 2 en código como `{{output_contract}}` · ruteo por `mode`→`template.key` · tests de paridad.

### CE-2.6 · Seed inicial de los 3 `body` (recommend / assist / guided)
**Como** producto **quiero** cargar la **v1 real** de los tres prompts **para** estrenar el sistema versionado con el contenido pedagógico definitivo.
`risk:med` · `cold-start`

- **🔗 Depende de:** CE-2.1 + **EPIC 5** (texto real de los `body`).
- **📄 Spec:** [§9.4 Los tres `body` versionados](./alizia-context-engine.md#94-los-tres-body-versionados-v1) · [Apéndice A — Checklist de contenido](./alizia-context-engine.md#apéndice-a--checklist-de-contenido-para-los-body-base) · [Apéndice B — Esqueletos](./alizia-context-engine.md#apéndice-b--esqueletos-de-los-tres-body-para-rellenar).
- **✅ Criterios de aceptación:** existe una `version` `active` por modo con el `body` real (no el ilustrativo), que pasa la validación de CE-2.3; el seed se carga **por script de entorno** (no en migración).
- **🔧 Sub-tareas:** rellenar los esqueletos del Apéndice B con la prosa real de EPIC 5 · script de entorno de seed · validar al publicar.

---

## EPIC 3 · Memoria

> **Fase 3** · Riesgo **medio** · Migración **000023** · Component `memory`
> **Objetivo:** dejar de **descartar** el historial viejo por presupuesto de tokens; resumir y mantener una memoria viva por alumno que se inyecte en el prompt.
> **📄 Spec:** [§7 Capa B — Memoria](./alizia-context-engine.md#7-modelo-de-datos) · [§9.1 etapa 0 (retomar memoria)](./alizia-context-engine.md#91-ciclo-de-vida-de-la-conversación-prompt-y-variables-por-etapa)

### CE-3.1 · Resumen de conversación (`conversation_summaries` + job)
**Como** Alizia **quiero** un resumen comprimido de cada conversación **para** no tirar el historial viejo cuando excede el presupuesto de tokens.
`risk:med` · `migration-000023`

- **🔗 Depende de:** — (modelo); se integra mejor con CE-1.7 ya listo.
- **📄 Spec:** [§7 — `conversation_summaries`](./alizia-context-engine.md#7-modelo-de-datos).
- **✅ Criterios de aceptación:** existe `conversation_summaries` 1:1 con `conversations`; un job genera/actualiza el resumen; el bloque `{{conversation_summary}}` lo consume.
- **🔧 Sub-tareas:** migración `000023` (parte) · entidad + provider · job de resumen · tests.

### CE-3.2 · Memoria viva del alumno (`student_insights`)
**Como** Alizia **quiero** una memoria acumulada por alumno (qué funciona / qué no) **para** anclar el saludo y personalizar sin re-preguntar.
`risk:med`

- **🔗 Depende de:** CE-3.1 (comparten job batch); habilita CE-1.8 `get_student_insights`.
- **📄 Spec:** [§7 — `student_insights`](./alizia-context-engine.md#7-modelo-de-datos) · [§9.5 bloque `{{insights}}`](./alizia-context-engine.md#95-cómo-renderiza-cada-bloque-dinámico).
- **✅ Criterios de aceptación:** `student_insights` 1 por alumno (`summary` + `key_learnings[]`), regenerado por job batch desde sesiones + adaptaciones; el bloque `{{insights}}` lo consume.
- **🔧 Sub-tareas:** migración `000023` (parte) · entidad + provider · regeneración en el job batch · tests.

### CE-3.3 · Usar el resumen en `capMessages`
**Como** sistema **quiero** que `capMessages` use el resumen en vez de descartar historial **para** no perder contexto en conversaciones largas.
`risk:med`

- **🔗 Depende de:** CE-3.1.
- **📄 Spec:** [§3 Estado actual (`history.go:capMessages`)](./alizia-context-engine.md#3-estado-actual-cómo-funciona-hoy).
- **✅ Criterios de aceptación:** cuando el historial excede el presupuesto, se conserva system + últimos N + resumen comprimido de lo viejo (en vez de tirar los mensajes).
- **🔧 Sub-tareas:** modificar `capMessages` · integrar resumen · tests de truncado.

---

## EPIC 4 · Flywheel (self-improvement)

> **Fase 4** · Riesgo **alto** · Migración **000024** · Component `flywheel`
> **Objetivo:** cerrar el loop — usar resultado de aula + aceptación implícita para promover ejemplos golden, medir win-rate por versión y correr A/B + eval.
> **📄 Spec:** [§10 Loop de self-improvement (flywheel)](./alizia-context-engine.md#10-loop-de-self-improvement-flywheel) · [§13 Métricas de éxito](./alizia-context-engine.md#13-métricas-de-éxito)
> **Junta todas las dependencias:** traza (EPIC 0) + versiones (EPIC 2) + contenido curated (EPIC 5).

### CE-4.1 · `response_examples` + seed `curated` de los ~15 casos
**Como** Alizia **quiero** un repositorio de ejemplos (curated + promovidos) **para** alimentar el few-shot y arrancar con buenos ejemplos desde el día 1 (cold-start).
`risk:high` · `migration-000024` · `cold-start`

- **🔗 Depende de:** **EPIC 5** (los ~15 casos golden); el seed va por script de entorno.
- **📄 Spec:** [§10 — `response_examples`](./alizia-context-engine.md#10-loop-de-self-improvement-flywheel) · [Apéndice A — few-shot golden](./alizia-context-engine.md#apéndice-a--checklist-de-contenido-para-los-body-base).
- **✅ Criterios de aceptación:** existe `response_examples` (origen `curated`/promovido, relevancia, modo); seed `curated` de ~15 casos por script de entorno.
- **🔧 Sub-tareas:** migración `000024` (parte) · entidad + provider · script de seed (consume CE-5.6) · tests.

### CE-4.2 · Selección de few-shot golden e inyección tras el corte
**Como** sistema **quiero** seleccionar top-3 ejemplos por relevancia y ubicarlos **tras el corte de cache** **para** personalizar por alumno sin romper el prefijo cacheable.
`risk:high`

- **🔗 Depende de:** CE-4.1, CE-2.2.
- **📄 Spec:** [§8 — few-shot tras el corte](./alizia-context-engine.md#8-el-context-assembler-en-detalle) · [§9.3 `{{few_shot}}`](./alizia-context-engine.md#93-catálogo-de-variables-placeholders).
- **✅ Criterios de aceptación:** el bloque `{{few_shot}}` se rellena con top-3 por relevancia y se ubica **después** del marcador de corte; sin ejemplos → solo lineamientos (cold-start).
- **🔧 Sub-tareas:** scoring de relevancia · selección top-3 · ubicación post-corte · tests.

### CE-4.3 · Win-rate por versión + métricas de éxito
**Como** equipo **quiero** medir win-rate, aceptación implícita, cobertura de contexto y costo/cache-hit **para** saber si el sistema mejora.
`risk:high`

- **🔗 Depende de:** EPIC 0 (traza), CE-2.1 (versiones).
- **📄 Spec:** [§13 Métricas de éxito](./alizia-context-engine.md#13-métricas-de-éxito).
- **✅ Criterios de aceptación:** se computan win-rate (% `funcionó`), tasa de aceptación implícita, cobertura de contexto, costo por turno y % cache-hit, y mejora por versión.
- **🔧 Sub-tareas:** queries de métricas · exposición (tablero/endpoint admin) · tests.

### CE-4.4 · A/B entre versiones + eval antes de promover
**Como** equipo **quiero** correr dos versiones activas con split de tráfico y un set de eval antes de promover **para** no degradar la calidad.
`risk:high`

- **🔗 Depende de:** CE-2.1, CE-4.1, CE-4.3.
- **📄 Spec:** [§9 — flujo de versiones (A/B, eval)](./alizia-context-engine.md#9-prompts-versionados-en-db).
- **✅ Criterios de aceptación:** soporta 2 versiones activas con split; antes de promover, corre el set de `response_examples` y compara; promoción condicionada por win-rate.
- **🔧 Sub-tareas:** split de tráfico por versión · runner de eval · criterio de promoción · tests.

### CE-4.5 · Job batch del flywheel (cron interno)
**Como** plataforma **quiero** un job batch que procese resultados y promueva golden de forma idempotente y aislada del request path **para** que el loop gire sin afectar la API.
`risk:high`

- **🔗 Depende de:** CE-4.1, CE-4.3 (y CE-3.x para insights).
- **📄 Spec:** [§10.1 Job batch (cron interno)](./alizia-context-engine.md#10-loop-de-self-improvement-flywheel).
- **✅ Criterios de aceptación:**
  - La lógica vive en `RunBatch(ctx)`; se dispara por cron interno.
  - **Config por env:** `FLYWHEEL_ENABLED` (default `false`) y `FLYWHEEL_CRON` — **sin intervalos hardcodeados**.
  - **Idempotente** (marca de agua: procesar solo lo nuevo); **aislado** del request path (goroutine; si falla, loguea y reintenta, no tumba el server).
  - **Observabilidad:** cada corrida loguea cuántas adaptaciones procesó, cuántos golden promovió y cuánto tardó.
  - Anotado el camino a **worker único** (advisory lock / `cmd/worker`) para cuando escale a N réplicas.
- **🔧 Sub-tareas:** `RunBatch(ctx)` · scheduler interno + config env · marca de agua de idempotencia · logging estructurado · tests.

---

## EPIC 5 · Contenido pedagógico base

> **Transversal · autoría (sin código)** · Component `content` · **arranca el día 1**
> **Objetivo:** producir el **texto real** que rellena los `body` (capa 1) y los ejemplos golden. Es el insumo de CE-2.6 (seed de prompts) y CE-4.1 (curated). Lo redacta producto/pedagogía.
> **📄 Spec (la checklist completa, ítem por ítem):** [Apéndice A — Checklist de contenido para los `body` base](./alizia-context-engine.md#apéndice-a--checklist-de-contenido-para-los-body-base) · [Apéndice B — Esqueletos](./alizia-context-engine.md#apéndice-b--esqueletos-de-los-tres-body-para-rellenar).

> **Nota:** cada historia de abajo corresponde a un bloque del Apéndice A. **El detalle ítem-por-ítem (los `- [ ]`) está en el spec** — acá no se duplica para no desincronizar.

### CE-5.1 · Identidad y persona
**Como** redactor de contenido **quiero** definir quién es Alizia, a quién le habla, tono/registro, qué NO es y estilo de salida **para** fijar la persona literal del prefijo invariante.
`content` · `parallelizable`
- **📄 Spec:** [Apéndice A · bloque A](./alizia-context-engine.md#apéndice-a--checklist-de-contenido-para-los-body-base).
- **✅ Criterio:** bloque A del checklist tildado; tono base **literal** (la preferencia del docente se aplica desde su bloque, no se interpola en la persona).

### CE-5.2 · Marco pedagógico (3 ejes + Criterios de AliZia + DUA)
**Como** redactor **quiero** el texto real de los 3 ejes, Criterios de AliZia y marco DUA **para** anclar el razonamiento pedagógico.
`content` · `parallelizable`
- **📄 Spec:** [Apéndice A · bloque B](./alizia-context-engine.md#apéndice-a--checklist-de-contenido-para-los-body-base).
- **✅ Criterio:** bloque B tildado, incluyendo "entrada pedagógica, no clínica".

### CE-5.3 · Límites duros / guardrails
**Como** redactor **quiero** la redacción oficial de los 3 límites (no diagnostica · no reemplaza al docente · no produce informes clínicos) y protocolos de riesgo **para** proteger el dominio (datos sensibles de menores).
`content` · `parallelizable`
- **📄 Spec:** [Apéndice A · bloque C](./alizia-context-engine.md#apéndice-a--checklist-de-contenido-para-los-body-base).
- **✅ Criterio:** bloque C tildado, incl. qué hace ante pedidos fuera de scope y situaciones de riesgo.

### CE-5.4 · Reglas de formato de salida
**Como** redactor **quiero** confirmar 1–3 acciones, ≥3 niveles de diferenciación, "útil en <1 min" y la estructura por modo **para** que la salida sea consistente y parseable.
`content` · `parallelizable`
- **📄 Spec:** [Apéndice A · bloque D](./alizia-context-engine.md#apéndice-a--checklist-de-contenido-para-los-body-base).
- **✅ Criterio:** bloque D tildado; coherente con el `{{output_contract}}` (capa 2).

### CE-5.5 · Las ~15 situaciones observables (seed de `situations_catalog`)
**Como** redactor **quiero** la lista oficial de ~15 situaciones de aula con código y nombre **para** alimentar el seed de `situations_catalog` (CE-1.3).
`content` · `parallelizable`
- **📄 Spec:** [Apéndice A](./alizia-context-engine.md#apéndice-a--checklist-de-contenido-para-los-body-base) · [§7 `situations_catalog`](./alizia-context-engine.md#7-modelo-de-datos).
- **✅ Criterio:** ~15 situaciones con `code`+`name` (y `phase` si la clasificación la define); entregadas en formato consumible por el script de seed.

### CE-5.6 · Few-shot golden curated (~15 casos)
**Como** redactor **quiero** ~15 ejemplos golden de respuesta (contexto → respuesta ideal) **para** el cold-start del flywheel (CE-4.1).
`content` · `parallelizable` · `cold-start`
- **📄 Spec:** [Apéndice A — few-shot golden curated](./alizia-context-engine.md#apéndice-a--checklist-de-contenido-para-los-body-base) · [§9.5 `{{few_shot}}`](./alizia-context-engine.md#95-cómo-renderiza-cada-bloque-dinámico).
- **✅ Criterio:** ~15 casos curados, etiquetados por modo/situación, listos para el seed de `response_examples`.

> **Cobertura completa:** el Apéndice A tiene **11 bloques (A–K)**. CE-5.1…CE-5.6 cubren los de mayor impacto para el seed; los bloques restantes (tipos de adaptación, marco normativo, etc.) se tratan como ítems dentro de estas historias o se abren como sub-tareas según haga falta — **el inventario canónico es el Apéndice A**.

---

## Resumen de dependencias (vista rápida)

```text
EPIC 5  ───────────────────────────────────► alimenta CE-2.6 (seed prompts) y CE-4.1 (curated)
EPIC 0  ──► (independiente, base de traza) ─► habilita métricas de EPIC 4
EPIC 1  ──► CE-1.1..1.6 (paralelo) ──► CE-1.7 Assembler ──► CE-1.8 tools
EPIC 2  ──► CE-2.1 ──► CE-2.2 ──► CE-2.3 / CE-2.4 ──► CE-2.5 ──► CE-2.6   (CE-2.2 depende de CE-1.7)
EPIC 3  ──► CE-3.1 ──► CE-3.2 ──► CE-3.3                                  (mejor con CE-1.7 listo)
EPIC 4  ──► CE-4.1 ──► CE-4.2 / CE-4.3 ──► CE-4.4 / CE-4.5               (necesita EPIC 0 + 2 + C)
```

> **Recordatorio:** este backlog es el *qué* y el *orden*. Para el *cómo* (DBML, columnas exactas,
> orden de bloques, validaciones, formato de cada placeholder) la fuente es siempre
> [`alizia-context-engine.md`](./alizia-context-engine.md) — **tiene literalmente la data ahí**.
