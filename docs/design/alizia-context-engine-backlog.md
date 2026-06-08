# Backlog Context Engine — DOCUMENTO DE TAREAS (vista de producto)

> **Spec / diseño (el "normal"):** [`alizia-context-engine.md`](./alizia-context-engine.md) — el *cómo* a fondo (modelo de datos, DBML, RAG, prompts, validaciones).
> **Este doc:** el *qué* en clave de producto — **historias concisas y testeables** (para demo de producto y para QA). El detalle técnico vive en el spec; acá no se enumeran subtareas internas.

## Scope del MVP (modelo fusionado, 2026-06-04)

- **Apertura con router**: Alizia recibe, ofrece modo (recommend / assist / guided) y pregunta de qué hablar (alumno / valija / tema); con eso decide qué contexto cargar.
- **Valija en contexto**: el catálogo de dispositivos viaja en el prompt (no usa RAG).
- **RAG sólo de contenido pedagógico** (libros / papers / material que nos brinden), por keywords contra Postgres, sin vectores. Es transversal a cualquier rama.
- **Tools** de lectura + **acción** (crear alumno, crear recurso/adaptación, vincular) + RAG.
- **Prompts en código** (paquete aparte), no en DB. **Historial compactado en DB** (2 momentos: al cerrar y al reabrir).
- **Quedan a Futuro**: memoria viva (insights), flywheel (golden / win-rate / A/B), prompts editables desde backoffice, archivado en buckets GCP.

Todo cuelga de la épica **AlizIA Inclusión - Chubut (ALZ-246)**.

---

## Resumen

### 🟢 MVP — 6 historias

| HU | Título (producto) | Se demuestra mostrando… |
|---|---|---|
| HU-1 | Apertura guiada: elijo modo y de qué hablar | la bienvenida con 3 modos + la pregunta de dimensión que enruta |
| HU-2 | Respuestas pertinentes a mi alumno | una sugerencia que usa el perfil/situaciones/PPI del chico |
| HU-3 | Consultar contenido pedagógico (temas) | Alizia trayendo info de un libro/paper cargado, sin inventar |
| HU-4 | Recomendar la valija y guardar adaptaciones | la ficha del recurso al frente, guardada y vinculada al alumno |
| HU-5 | La conversación conserva memoria | una charla larga que no pierde el hilo y retoma al volver |
| HU-6 | Alizia responde dentro de su marco pedagógico | respuestas accionables con límites respetados + traza por turno |

### 🔵 Futuro (post-MVP)

| Tema | En una línea |
|---|---|
| Memoria viva del alumno | recuerda entre sesiones qué funcionó con cada chico (insights) |
| Aprendizaje por resultado (flywheel) | golden few-shot + win-rate + A/B para mejorar sola |
| Prompts editables desde backoffice | producto edita/publica el comportamiento sin deploy (DB versionada) |
| Archivado en buckets GCP | conversaciones viejas a almacenamiento barato para analytics/ML |

---

# 🟢 MVP

## HU-1 · Apertura guiada: elijo modo y de qué hablar

**Como** docente **quiero** que Alizia me reciba, me deje elegir cómo trabajar y me pregunte de qué quiero hablar **para** entrar directo a lo que necesito sin configurar nada.

**Se demuestra (producto):** abro una sesión → veo el saludo y **tres opciones de modo** (recommend / assist / guided) → elijo una → Alizia pregunta **"¿de qué querés hablar?"** → elijo **un alumno**, **la valija** o **un tema** → Alizia arranca enfocada en eso.

**Criterios de aceptación (QA):**
- Al iniciar una sesión se muestran la **bienvenida** y las **3 opciones de modo**.
- Tras elegir el modo, Alizia **pregunta la dimensión** (alumno / valija / tema).
- Si la respuesta es ambigua, **repregunta** en lugar de asumir.
- Lo elegido **determina qué contexto se carga**: si elijo "tema" no se carga perfil de un alumno; si elijo "alumno" se carga su perfil; la valija siempre está disponible en el catálogo.
- Si reabro una conversación previa, la apertura **retoma de qué veníamos hablando**.

*Spec: §6.0 (router de apertura), §8 (carga dirigida).*

---

## HU-2 · Respuestas pertinentes a mi alumno

**Como** docente **quiero** que, cuando hable de un alumno, las sugerencias estén pensadas para **ese** chico (perfil, situaciones de aula, PPI, entorno) **para** que no sean genéricas.

**Se demuestra (producto):** con un alumno que tiene perfil cargado (situación "no inicia la tarea", fortalezas, PPI), pregunto cómo encarar una actividad → la respuesta **menciona y usa** esos datos. Comparado con un alumno sin datos, la respuesta sigue siendo útil pero más general.

**Criterios de aceptación (QA):**
- Con perfil cargado, la respuesta **refleja** al menos las situaciones / fortalezas / PPI del alumno.
- Los campos vacíos (todos opcionales) **no aparecen como "N/A"** ni rompen la respuesta.
- Alizia puede **sugerir completar** datos faltantes; nunca los exige.
- En los logs **no aparece PII** (nombres ni diagnósticos): sólo IDs.

*Spec: §6.2 (Context Assembler), §7 Capa A, §8.*

---

## HU-3 · Consultar contenido pedagógico (temas)

**Como** docente **quiero** preguntarle a Alizia sobre temas de inclusión (estrategias para autismo, técnicas de lectura para baja visión, etc.) y que responda con **material pedagógico real** (libros / papers cargados) **para** tener respuestas fundadas y no inventadas.

**Se demuestra (producto):** pregunto *"¿qué estrategias hay para un alumno con TEA que se desregula?"* → Alizia **busca en el contenido pedagógico** y responde con lo que encontró. Si pregunto algo que no está cargado, **lo aclara** y responde con los lineamientos base, sin inventar.

**Criterios de aceptación (QA):**
- Una pregunta cuyo tema **está en el corpus** devuelve contenido de ese material (verificable contra el documento sembrado).
- La búsqueda **tolera errores de tipeo** (reescribe la pregunta a keywords antes de buscar).
- Devuelve **los más relevantes primero** (ranking por coincidencia de keywords).
- Sin match → **no inventa**: usa los lineamientos base.
- Es **transversal**: funciona hablando de un alumno, de un tema o de la valija.
- La **valija NO** pasa por este buscador (va en el catálogo del prompt); el RAG es sólo para contenido pedagógico.

*Spec: §0.1, §6.2 (`search_content` / `get_content`), §7 Capa E.*

---

## HU-4 · Recomendar la valija y guardar adaptaciones

**Como** docente **quiero** que Alizia me recomiende dispositivos de la valija y que pueda **crear un alumno** o **guardar una adaptación/recurso** desde la conversación, vinculada al alumno **para** no rehacer el trabajo a mano.

**Se demuestra (producto):** pido una adaptación → Alizia propone usando **dispositivos de la valija** → me muestra **la ficha del recurso de frente** → la guardo → queda **vinculada al alumno** y registra de qué conversación salió.

**Criterios de aceptación (QA):**
- Alizia recomienda dispositivos del **catálogo de la valija** (identificados con su `DEVICE_ID`).
- Puedo **crear un alumno** desde la conversación.
- Al crear un recurso/adaptación, **se muestra la ficha al frente** y **se guarda con su origen** (conversación/mensaje) y con la marca de **si fue editada** antes de guardar.
- El recurso queda **vinculado al alumno**.

*Spec: §6.2 (tools de acción: `create_student` / `create_recurso` / `relate_student_recurso`), §7 Capa C.*

---

## HU-5 · La conversación conserva memoria

**Como** docente **quiero** que Alizia no pierda el hilo en charlas largas y que al volver recuerde de qué veníamos **para** seguir trabajando sin recontextualizar todo.

**Se demuestra (producto):** mantengo una conversación larga → sigue coherente, sin "olvidar" lo del principio → la cierro → vuelvo más tarde → Alizia **retoma con un resumen** de lo anterior.

**Criterios de aceptación (QA):**
- Una conversación que **excede el presupuesto de tokens** conserva el system + los últimos turnos + un **resumen comprimido** de lo viejo (no se tira el historial).
- **Al cerrar** la sesión se genera/actualiza el **resumen compactado** en DB.
- **Al abrir**, Alizia **recupera el resumen** de la conversación/entidad correspondiente (momento que dispara la apertura, HU-1).

*Spec: §6.4 (compactación en 2 momentos), §7 Capa B.*

---

## HU-6 · Alizia responde dentro de su marco pedagógico

**Como** equipo **quiero** que las respuestas de Alizia sigan nuestros lineamientos de inclusión, respeten sus límites duros, y que cada interacción quede registrada **para** garantizar calidad y poder depurar (y, más adelante, medir).

**Se demuestra (producto):** cualquier recomendación llega con **1-3 acciones** ordenadas por impacto, con **niveles de diferenciación**, en lenguaje cálido y claro. Si pido un diagnóstico, Alizia **se corre y deriva**, no diagnostica.

**Criterios de aceptación (QA):**
- Las propuestas traen **1-3 acciones**, **≥3 niveles de diferenciación**, útiles en **<1 min**.
- Alizia **nunca** diagnostica, **no** reemplaza al docente, **no** produce informes clínicos: ante el pedido, lo aclara y deriva sin cortar la charla.
- Tono **cálido, español rioplatense, sin jerga clínica**.
- Los **prompts viven en un paquete de código** (capa editable separada del motor); cambiar el texto no rompe el formato de salida.
- **Cada turno deja traza**: modelo, latencia, tokens, tool calls y contexto por **IDs (sin PII)**; si el registro falla, **no bloquea** la respuesta.

*Spec: §4 (principios), §6.1 (capas), §9.4 (bodies) y Apéndice A; traza §6.4 / §7 Capa C.*

---

# 🔵 Futuro (post-MVP)

> Oportunidades de mejora una vez validado el MVP. Se describen en clave de valor; el detalle técnico está en el spec (§9, §10, §12).

- **Memoria viva del alumno (insights).** Alizia recuerda entre sesiones qué funcionó con cada chico, sin que el docente se lo repita. *(spec §6.2, §7 Capa B — `student_insights`).*
- **Aprendizaje por resultado (flywheel).** Las adaptaciones que funcionaron en el aula alimentan ejemplos golden; se mide win-rate por versión y se corren A/B para mejorar sola. *(spec §10, §13).*
- **Prompts editables desde backoffice.** Producto/pedagogía ajustan y publican el comportamiento de Alizia sin esperar un deploy, con validación y fallback. *(spec §5, §9 — prompts en DB versionada).*
- **Archivado en buckets GCP.** Conversaciones viejas a almacenamiento barato (file system) para liberar la DB y guardar datos recuperables para analytics/ML. *(spec §0, §12).*

---

## Dependencias (a nivel historia)

- **HU-2 (contexto del alumno)** habilita HU-3 (el RAG usa la situación/`difficulties[]` como keywords de tema) y HU-4 (recomendaciones pertinentes).
- **HU-1 (apertura/router)** dirige la carga de contexto de HU-2/HU-3 y consume el resumen de HU-5.
- **HU-6 (marco + traza)** es transversal: aplica a las respuestas de todas las demás.
- **HU-5 (memoria)** y **HU-6 (traza)** dejan la base sobre la que más adelante se montan memoria viva y flywheel (Futuro).

## Sincronización con Jira (pendiente)

> Esta versión reagrupa el backlog en **6 historias de producto**; la estructura anterior tenía 10 historias técnicas (épica `ALZ-246`). La conciliación de claves Jira (qué issues existentes se mantienen, fusionan o reetiquetan, y el reparenteo de subtareas) se hace en la tarea de **sync a Jira**, tomando este documento como fuente del *qué*.
