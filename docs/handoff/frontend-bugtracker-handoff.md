# Handoff a Front — Bug Tracker de Alizia Inclusion

> Contexto: análisis y fixes de backend sobre los 7 bugs del tracker. El trabajo previo de
> "context engine" atacó calidad conversacional (no estos bugs), así que ninguno se resolvió
> de rebote. Acá está qué se arregló en backend y qué queda del lado del front.

Fecha: 2026-06-17 · Branch backend: `fix/context-engine-hardening`

---

## 1. Cambios de contrato nuevos a consumir (ya en backend)

### a) `GET /api/v1/conversation/:id` — **NUEVO** (resuelve BUG-005)
Trae una conversación puntual con sus mensajes, scopeada a la org. Sirve para "Retomar con tu
asistente": leer `adaptation.source_conversation_id` del recurso y cargar **esa** conversación.

Respuesta:
```json
{
  "id": 42,
  "mode": "assist",
  "messages": [
    { "role": "user", "content": "...", "created_at": "2026-06-17T..." },
    { "role": "assistant", "content": "...", "created_at": "2026-06-17T..." }
  ],
  "created_at": "2026-06-17T..."
}
```
- `404` si no existe o es de otra org.
- Flujo "Retomar": `GET /adaptations/:id` → `source_conversation_id` → `GET /conversation/:id`.

### b) `referenced_content[]` en la respuesta del assist — **NUEVO** (resuelve BUG-006)
`POST /api/v1/inclusion/assist` ahora puede incluir los materiales que Ada citó en el turno,
con su id real, para que el chip deep-linkee al material específico (no a `/materiales`).

```json
{
  "response": "Te dejo este material ...",
  "conversation_id": 42,
  "referenced_content": [
    { "id": 2, "title": "Guía de acceso a la lectura" }
  ]
}
```
- Se omite (`omitempty`) cuando Ada no citó ningún material.
- Los ids inexistentes o de otra org se descartan silenciosamente (un chip roto nunca corta el chat).
- **Acción FE:** renderizar el chip desde `referenced_content[]` (id + title) y navegar a
  `/materiales/:id` (o el detalle correspondiente), en vez de parsear el texto.

### c) `adaptation.student_id` — **NUEVO** (ayuda a BUG-003)
El objeto `adaptation` de la respuesta del assist ahora trae `student_id` cuando se identificó
un alumno en el turno. Antes había que leerlo del campo hermano `identified_student`.

```json
{
  "adaptation": {
    "title": "Pasos cortos",
    "type": "estrategia_aula",
    "strategy": "...",
    "device_ids": [1],
    "device_names": ["Pictogramas"],
    "student_id": 7
  }
}
```

---

## 2. BUG-003 — "Guardar recurso" falla con 400

**Causa:** el FE arma el `POST /api/v1/adaptations` sin `student_id` (no lo extraía del texto).
**El backend ya lo entrega** de dos formas:
- `identified_student` (campo de nivel superior de la respuesta del assist) — ya existía.
- `adaptation.student_id` — nuevo (ver 1.c).

**Acción FE:** tomar `student_id` de `adaptation.student_id` (o `identified_student`) y mandarlo
en el body de `/adaptations`. No hace falta parsear `[STUDENT_ID:X]` del texto.

---

## 3. BUG-004 — No inicia conversación nueva desde "¿Cómo usarlo con mi alumno?"

**Es de FE (estado de navegación).** El endpoint `assist` respeta `conversation_id`:
- Para **iniciar una conversación nueva** (con el contexto del material), mandar
  `conversation_id: 0` (o no mandarlo). El backend crea una conversación nueva.
- Para **continuar** una existente, mandar el `conversation_id` real.

Hoy el FE estaría reusando el `conversation_id` previo al entrar desde el material, por eso abre
la conversación vieja en vez de una nueva con "Cómo uso X material".

---

## 4. BUG-006 (recordatorio de comportamiento)

Ada marca los materiales citados con `[CONTENT_ID:X]` en el texto y el backend los resuelve a
`referenced_content[]`. El FE **no** debería mostrar el tag crudo `[CONTENT_ID:X]` al usuario:
limpiar esos tags del texto al renderizar (igual que se hace —o debería— con `[STUDENT_ID:X]` y
`[DEVICE_ID:X]`).

---

## 5. BUG-001 / BUG-002 — Materiales sin título / sin descargables

**Diagnóstico: es dato del catálogo, no código de backend.**

El módulo "Materiales" consume el catálogo de la valija (`GET /api/v1/devices` y
`/devices/:id`). El backend ya:
- expone `name` de cada device, y
- preloadea y expone los descargables en `downloads[]` (de la tabla `device_resources`,
  con `title` / `file_url` / `file_type`).

Por lo tanto:
- **BUG-001 (sin título en "Acceso a la lectura"):** los devices de ese ramp tienen `name`
  vacío en la base (seed incompleto). El backend siempre mapea `name`. → corregir datos del
  catálogo.
- **BUG-002 (sin recursos para descargar):** no hay filas en `device_resources` para esos
  devices, por eso `downloads[]` viene vacío. → cargar los recursos descargables en el catálogo.

**Acción:** queda para el equipo de contenido/catálogo (completar `name` y poblar
`device_resources`). No requiere cambio de backend ni de FE.

---

## 6. BUG-007 — Latencia 5–25s sin feedback de carga

**Mitigaciones de backend ya aplicadas:** loop agéntico acotado a `maxAgenticIterations = 2` y
techo del buscador RAG (`defaultContentSearchLimit`). El tiempo restante es inherente a la
llamada al LLM (Azure OpenAI).

**El fix de fondo del "parece colgado" es de FE + contrato:**
- Corto plazo (solo FE): mostrar un spinner / estado "Ada está pensando…" mientras está pendiente
  el `POST /assist`.
- Mediano plazo (FE + back): **streaming SSE** de la respuesta token a token. Esto cambia el
  contrato del endpoint (de JSON único a `text/event-stream`), así que requiere trabajo
  coordinado. Si lo priorizan, lo encaramos como tarea aparte.

---

## Resumen accionable para Front

| Bug | Acción FE |
|---|---|
| 005 | Usar `GET /conversation/:id` para "Retomar" (vía `source_conversation_id`). |
| 006 | Renderizar chip desde `referenced_content[]` y deep-linkear por `id`. Limpiar tags `[..._ID:X]` del texto. |
| 003 | Mandar `student_id` (de `adaptation.student_id` / `identified_student`) en el POST de `/adaptations`. |
| 004 | Mandar `conversation_id: 0` para iniciar conversación nueva desde un material. |
| 007 | Spinner mientras carga; evaluar SSE streaming (coordinado con back). |
| 001/002 | No es FE ni back: completar datos del catálogo (`name` y `device_resources`). |
