# Backlog Context Engine — DOCUMENTO DE TAREAS (vista de producto + subtareas)

> **Spec / diseño (el "normal"):** [`alizia-context-engine.md`](./alizia-context-engine.md) — el *cómo* a fondo (modelo de datos, DBML, RAG, prompts, validaciones, guardrails).
> **Este doc:** el *qué* en clave de producto — **historias concisas y testeables** (para demo y QA) — más, por historia, las **subtareas técnicas (T)** con su estado, para guiar la implementación.

## Scope del MVP (modelo fusionado)

- **Apertura con router (un solo modo)**: Alizia recibe y pregunta de qué hablar (alumno / valija / tema); con eso decide qué contexto cargar. **No hay elección de modo**: una sola forma de responder (directo si puede, pregunta/busca si le falta).
- **Valija en contexto**: el catálogo de dispositivos viaja en el prompt (no usa RAG).
- **RAG sólo de contenido pedagógico** (libros / papers / material que nos brinden), por **keywords / full-text** contra Postgres. Los **embeddings quedan inertes** (columna lista, sin índice; se activan en Futuro). Es transversal a cualquier rama.
- **Tools** de lectura + **acción** (crear alumno, crear recurso/adaptación, vincular) + RAG. **Loop agéntico acotado**: 1 pasada (2 máx) en el primer release.
- **Guardrails**: validación por código de la respuesta antes de mostrarla + **off-ramp** (Alizia intenta responder con lo que tiene; paso al costado solo para casos fuera de alcance).
- **Output**: dispositivos desde catálogo cerrado (`DEVICE_ID`); **adaptación de escritura libre** (el front ofrece opciones predefinidas, pero el docente siempre puede escribir libre).
- **Recursos privados por docente**; **guardado siempre confirmado** por el docente.
- **Prompts en código** (paquete aparte, un solo modo). **Historial compactado en DB** (2 momentos: al cerrar y al reabrir).
- **Migraciones con golang-migrate**.
- **Quedan a Futuro**: embeddings activos (pgvector), memoria viva (insights), flywheel (golden / win-rate / A/B), prompts editables desde backoffice, archivado en buckets GCP, presets de tono.

Todo cuelga de la épica **AlizIA Inclusión - Chubut (ALZ-246)**.

---

## Estado de la fundación (DB) — ✅ hecho

| ID | Subtarea | Estado |
|---|---|---|
| F-1 | Migraciones `000015`–`000021` (schema context-engine, Capas A–E) | ✅ aplicado a Railway |
| F-2 | Seed de catálogos globales (15 situaciones, 10 diagnósticos, few-shot curated, contenido pedagógico) | ✅ aplicado |
| F-3 | Adopción de **golang-migrate** (tracking en `schema_migrations`) | ✅ Railway en v21 |
| F-4 | Colisión de versión `000011` resuelta (→ `000014`) | ✅ |

> La capa de datos ya existe en Railway. Lo que sigue (entities/repos/usecases/tools/prompts) es código de aplicación, desglosado por historia abajo.

---

## Resumen

### 🟢 MVP — 6 historias

| HU | Título (producto) | Se demuestra mostrando… |
|---|---|---|
| HU-1 | Apertura guiada: elijo de qué hablar | la bienvenida + la pregunta de dimensión que enruta el contexto |
| HU-2 | Respuestas pertinentes a mi alumno | una sugerencia que usa el perfil/situaciones/PPI del chico |
| HU-3 | Consultar contenido pedagógico (temas) | Alizia trayendo info de un libro/paper cargado, sin inventar |
| HU-4 | Recomendar la valija y guardar adaptaciones | la ficha del recurso al frente, confirmada, guardada y vinculada al alumno |
| HU-5 | La conversación conserva memoria | una charla larga que no pierde el hilo y retoma al volver |
| HU-6 | Alizia responde dentro de su marco (+ guardrails) | respuestas accionables con límites verificados por código + traza por turno |

### 🔵 Futuro (post-MVP)

| Tema | En una línea |
|---|---|
| Embeddings activos (pgvector) | búsqueda semántica que entiende sinónimos/paráfrasis, además de keywords |
| Memoria viva del alumno | recuerda entre sesiones qué funcionó con cada chico (insights) |
| Aprendizaje por resultado (flywheel) | golden few-shot + win-rate + A/B para mejorar sola |
| Prompts editables desde backoffice | producto edita/publica el comportamiento sin deploy (DB versionada) |
| Presets de tono | modular el tono según el momento (clase vs. planificación) |
| Archivado en buckets GCP | conversaciones viejas a almacenamiento barato para analytics/ML |

---

# 🟢 MVP

## HU-1 · Apertura guiada: elijo de qué hablar

**Como** docente **quiero** que Alizia me reciba y me pregunte de qué quiero hablar **para** entrar directo a lo que necesito sin configurar nada.

**Se demuestra (producto):** abro una sesión → veo el saludo → Alizia pregunta **"¿de qué querés hablar?"** → elijo **un alumno**, **la valija** o **un tema** → Alizia arranca enfocada en eso. **No tengo que elegir un "modo"**: Alizia se adapta sola.

**Criterios de aceptación (QA):**
- Al iniciar una sesión se muestra la **bienvenida** y la **pregunta de dimensión** (alumno / valija / tema). **No** hay selección de modo.
- Si la respuesta es ambigua, **repregunta** en lugar de asumir.
- Lo elegido **determina qué contexto se carga**: "tema" no carga perfil de alumno; "alumno" carga su perfil; la valija siempre está en el catálogo.
- Si reabro una conversación previa, la apertura **retoma de qué veníamos hablando**.

**Subtareas (T):**
- ⬜ T-1.1 · Prompt 0 / router: bienvenida + pregunta de dimensión; repregunta si es ambiguo.
- ⬜ T-1.2 · El output del router (dimensión) **dirige la carga lazy** del contexto.
- ⬜ T-1.3 · Al abrir, identifica la entidad y **recupera el resumen previo** (engancha con HU-5).

*Spec: §6.0 (router de apertura), §8 (carga dirigida).*

---

## HU-2 · Respuestas pertinentes a mi alumno

**Como** docente **quiero** que, cuando hable de un alumno, las sugerencias estén pensadas para **ese** chico (perfil, situaciones de aula, PPI, entorno) **para** que no sean genéricas.

**Se demuestra (producto):** con un alumno con perfil cargado (situación "no inicia la tarea", fortalezas, PPI), pregunto cómo encarar una actividad → la respuesta **menciona y usa** esos datos. Con un alumno sin datos, sigue siendo útil pero más general.

**Criterios de aceptación (QA):**
- Con perfil cargado, la respuesta **refleja** al menos situaciones / fortalezas / PPI del alumno.
- Los campos vacíos (todos opcionales) **no aparecen como "N/A"** ni rompen la respuesta.
- Alizia puede **sugerir completar** datos faltantes; nunca los exige.
- En los logs **no aparece PII** (nombres ni diagnósticos): sólo IDs.

**Subtareas (T):**
- ⬜ T-2.1 · Entities + repos Capa A (`teacher_profiles`, `ppi`, `situations_catalog`, `diagnoses_catalog`, `student_diagnoses`; campos nuevos de `students`/`student_profiles`).
- ⬜ T-2.2 · Context Assembler `BuildPromptContext`: struct tipado, estático cacheado + dinámico lazy por dimensión.
- ⬜ T-2.3 · Tools de lectura: `get_student`, `get_student_history`, `get_past_adaptations`, `list_classroom_students`, `list_devices`.
- ⬜ T-2.4 · Degradar con elegancia ante campos vacíos + sugerir completar.

*Spec: §6.2 (Context Assembler), §7 Capa A, §8.*

---

## HU-3 · Consultar contenido pedagógico (temas)

**Como** docente **quiero** preguntarle a Alizia sobre temas de inclusión (estrategias para autismo, técnicas de lectura para baja visión, etc.) y que responda con **material pedagógico real** (libros / papers cargados) **para** tener respuestas fundadas y no inventadas.

**Se demuestra (producto):** pregunto *"¿qué estrategias hay para un alumno con TEA que se desregula?"* → Alizia **busca en el contenido pedagógico** y responde con lo que encontró. Si pregunto algo que no está cargado, **lo aclara** y responde con los lineamientos base, sin inventar.

**Criterios de aceptación (QA):**
- Una pregunta cuyo tema **está en el corpus** devuelve contenido de ese material (verificable contra el documento sembrado).
- La búsqueda **tolera errores de tipeo** (reescribe la pregunta a keywords antes de buscar).
- Devuelve **los más relevantes primero** (ranking por keyword / full-text).
- Sin match → **no inventa**: usa los lineamientos base.
- Es **transversal**: funciona hablando de un alumno, de un tema o de la valija.
- La **valija NO** pasa por este buscador (va en el catálogo del prompt); el RAG es sólo para contenido pedagógico.

**Subtareas (T):**
- ⬜ T-3.1 · Entities + repos Capa E (`pedagogical_content`, `pedagogical_content_chunks`).
- ⬜ T-3.2 · Tool `search_content`: reescribe a keywords → búsqueda **keyword / full-text** (`tsvector` / GIN). Embeddings **inertes** (columna lista, sin índice).
- ⬜ T-3.3 · Tool `get_content` (chunk / documento completo).
- ⬜ T-3.4 · Ranking por coincidencia; sin match → no inventa.
- ⬜ T-3.5 · `EMBEDDING_DIM` (default 1536) en un solo punto de config (insumo Futuro).

*Spec: §0.1, §6.2 (`search_content` / `get_content`), §7 Capa E.*

---

## HU-4 · Recomendar la valija y guardar adaptaciones

**Como** docente **quiero** que Alizia me recomiende dispositivos de la valija y que pueda **crear un alumno** o **guardar una adaptación/recurso** desde la conversación, vinculada al alumno **para** no rehacer el trabajo a mano.

**Se demuestra (producto):** pido una adaptación → Alizia propone usando **dispositivos de la valija** → me muestra **la ficha del recurso de frente** → **me pregunta si la guardo** → confirmo → queda **vinculada al alumno** (privada mía) y registra de qué conversación salió.

**Criterios de aceptación (QA):**
- Alizia recomienda dispositivos del **catálogo de la valija** (`DEVICE_ID`), sin inventar herramientas.
- La **adaptación** puede ser de escritura libre; el front ofrece opciones predefinidas, pero siempre se puede escribir libre.
- Ni alumnos ni recursos se guardan en silencio: Alizia **propone**, el docente **confirma**, recién ahí persiste.
- El recurso (= `adaptation`) es **privado del docente** y queda **vinculado al alumno**, con su **origen** (conversación/mensaje) y la marca **`was_edited`**.

**Subtareas (T):**
- ⬜ T-4.1 · Tools de acción: `create_student`, `create_recurso` (persiste `adaptation` vía endpoint existente), `relate_student_recurso`.
- ⬜ T-4.2 · Guardado confirmado (propone → confirma → persiste).
- ⬜ T-4.3 · Recurso = `adaptation`, scope **privado por docente**.
- ⬜ T-4.4 · Persistir origen (`source_conversation_id`/`source_message_id`), `was_edited` y materiales de valija usados.
- ⬜ T-4.5 · Front: ficha al frente + opciones predefinidas con escritura libre siempre.

*Spec: §6.2 (tools de acción), §6.7 (output cerrado/libre), §7 Capa C.*

---

## HU-5 · La conversación conserva memoria

**Como** docente **quiero** que Alizia no pierda el hilo en charlas largas y que al volver recuerde de qué veníamos **para** seguir trabajando sin recontextualizar todo.

**Se demuestra (producto):** mantengo una conversación larga → sigue coherente → la cierro → vuelvo más tarde → Alizia **retoma con un resumen** de lo anterior.

**Criterios de aceptación (QA):**
- Una conversación que **excede el presupuesto de tokens** conserva system + últimos turnos + un **resumen comprimido** de lo viejo (no se tira el historial).
- **Al cerrar** la sesión se genera/actualiza el **resumen compactado** en DB, con **tags a 3 dimensiones** (alumno / tema / valija).
- **Al abrir**, Alizia **recupera los resúmenes** ligados a la entidad (máx. 10 más recientes) — engancha con HU-1.

**Subtareas (T):**
- ⬜ T-5.1 · Entities + repos Capa B (`conversation_summaries` + cross-tables de alumnos/devices).
- ⬜ T-5.2 · Compactación al cerrar: resumen (par de párrafos) + tags/FKs a alumno/tema/valija; fallback máx 10.
- ⬜ T-5.3 · Recuperar resúmenes por entidad al abrir.
- ⬜ T-5.4 · `capMessages`: preservar system + últimos turnos + resumen (no tirar lo viejo).

*Spec: §6.4 (compactación en 2 momentos), §7 Capa B.*

---

## HU-6 · Alizia responde dentro de su marco pedagógico (+ guardrails)

**Como** equipo **quiero** que las respuestas sigan nuestros lineamientos, que sus límites se **verifiquen por código** (no solo en el prompt), y que cada interacción quede registrada **para** garantizar calidad y poder depurar.

**Se demuestra (producto):** cualquier recomendación llega con **1-3 acciones** ordenadas por impacto, con niveles de diferenciación, en lenguaje cálido. Si pido un diagnóstico, Alizia **intenta ayudar con lo que tiene** y, si el caso se va de su alcance (clínico/crisis), **da un paso al costado** sin diagnosticar.

**Criterios de aceptación (QA):**
- Las propuestas traen **1-3 acciones**, **≥3 niveles de diferenciación**, útiles en **<1 min**.
- **Guardrail por código** antes de mostrar: el `DEVICE_ID` existe en catálogo, el `ADAPTATION_JSON` parsea, no se cruzan límites duros. Si falla, reintenta o cae al off-ramp; **nunca** muestra salida inválida.
- **Off-ramp**: primero intenta responder con lo que tiene; el paso al costado es el último recurso, solo para lo fuera de alcance (clínico/crisis/diagnóstico).
- Alizia **nunca** diagnostica, **no** reemplaza al docente, **no** produce informes clínicos.
- Tono **cálido, español rioplatense, sin jerga clínica**.
- Los **prompts viven en un paquete de código** (capa editable separada del motor, un solo modo); cambiar el texto no rompe el formato de salida.
- **Cada turno deja traza**: modelo, latencia, tokens, tool calls y contexto por **IDs (sin PII)**; si el registro falla, **no bloquea** la respuesta.

**Subtareas (T):**
- ⬜ T-6.1 · Refactor `prompts.go` → paquete `prompts/` (un solo modo; capa 1 estática + capa 2 dinámica; few-shot estático cacheado + dinámico por alumno).
- ⬜ T-6.2 · Guardrail post-respuesta (valida `DEVICE_ID` / `ADAPTATION_JSON` / límites).
- ⬜ T-6.3 · Off-ramp: comportamiento + wording por defecto en constante editable.
- ⬜ T-6.4 · Loop agéntico acotado a 1 pasada (2 máx) para el MVP.
- ⬜ T-6.5 · Traza por turno en `ai_usage` (model, latency, tokens, tool_calls, context_snapshot por IDs); best-effort.
- ⬜ T-6.6 · Formato de salida: 1-3 acciones, ≥3 niveles de diferenciación, <1 min, tono rioplatense.

*Spec: §4 (principios), §6.1 (capas), §6.7 (guardrails), §9.4 (bodies) y Apéndice A; traza §6.4 / §7 Capa C.*

---

# 🔵 Futuro (post-MVP)

> Oportunidades una vez validado el MVP. En clave de valor; el detalle técnico está en el spec (§9, §10, §12).

- **Embeddings activos (pgvector).** Búsqueda semántica (sinónimos/paráfrasis) sobre el contenido pedagógico, además de keywords. La estructura ya está (columna `embedding`); falta fijar la dimensión real del modelo y crear el índice vectorial. *(spec §0.1, §7 Capa E).*
- **Memoria viva del alumno (insights).** Alizia recuerda entre sesiones qué funcionó con cada chico. *(spec §6.2, §7 Capa B — `student_insights`).*
- **Aprendizaje por resultado (flywheel).** Las adaptaciones que funcionaron alimentan ejemplos golden; win-rate por versión + A/B. *(spec §10, §13).*
- **Prompts editables desde backoffice.** Producto/pedagogía ajustan y publican sin deploy, con validación y fallback. *(spec §5, §9).*
- **Presets de tono.** Modular el tono según el momento (clase vs. planificación) sin volver a 3 modos. *(spec §0).*
- **Archivado en buckets GCP.** Conversaciones viejas a almacenamiento barato para analytics/ML. *(spec §0, §12).*

---

## Dependencias (a nivel historia)

- **Fundación DB (F-1…F-4)** ya está → habilita todas las HU (entities/repos sobre tablas existentes).
- **HU-2 (contexto del alumno)** habilita HU-3 (el RAG usa situación/`difficulties[]` como keywords de tema) y HU-4 (recomendaciones pertinentes).
- **HU-1 (apertura/router)** dirige la carga de contexto de HU-2/HU-3 y consume el resumen de HU-5.
- **HU-6 (marco + guardrails + traza)** es transversal: aplica a las respuestas de todas las demás.
- **HU-5 (memoria)** y **HU-6 (traza)** dejan la base sobre la que se montan memoria viva y flywheel (Futuro).

## Sincronización con Jira (pendiente)

> Esta versión mantiene **6 historias de producto** (épica `ALZ-246`) y agrega **subtareas (T) por historia** + la **fundación DB (F)** ya completada. La conciliación de claves Jira (qué issues se mantienen/fusionan/reetiquetan, reparenteo de subtareas y carga de las nuevas T) se hace en la tarea de **sync a Jira**, tomando este documento como fuente del *qué*.
