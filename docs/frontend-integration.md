# Frontend Integration — Alizia Inclusión

Qué debe construir el frontend para consumir las funcionalidades del backend.
Contratos verificados contra el código en `src/entrypoints/` al 2026-05-26.

Base URL: `/api/v1`. Todas las rutas requieren autenticación.

---

## 0. Autenticación y multi-tenancy (transversal)

Todas las llamadas a `/api/v1/**` requieren:

```
Authorization: Bearer <JWT_RS256>
```

El token lo emite el `auth-service`. El backend deriva `org_id` y `user_id`
del token (no se envían en el body). Sin token válido → `401`.

**A construir en el FE:**
- Interceptor HTTP que adjunte el header `Authorization` en cada request.
- Manejo de `401` → redirigir a login / refrescar token.

---

## 1. Endpoints de IA conversacional (LISTOS en backend)

### `POST /api/v1/inclusion/assist`

Chat del aula. Request:

```jsonc
{
  "conversation_id": 0,          // 0 para conversación nueva; el id devuelto para continuar
  "classroom_id": 12,
  "student_id": 5,               // opcional (puede ser null)
  "message": "Lucas no puede concentrarse",
  "mode": "",                    // "" = libre, "guided" = modo guiado paso a paso
  "history": [                   // opcional: turnos previos si el FE los mantiene en memoria
    { "role": "user", "content": "..." },
    { "role": "assistant", "content": "..." }
  ]
}
```

Response `200`:

```jsonc
{
  "response": "Podrías usar pictogramas [DEVICE_ID:1] ...",
  "conversation_id": 42,
  "identified_student": 5,       // omitido si no aplica
  "recommended_device": 1,       // omitido si no aplica
  "adaptation": {                // omitido si la IA no generó una adaptación
    "title": "Timer para fracciones",
    "type": "actividad_adaptada",
    "strategy": "Usar timer visual",
    "device_ids": [1],
    "device_names": ["Timer Visual"]
  }
}
```

### `POST /api/v1/inclusion/recommend`

Recomendación de dispositivo para un alumno. Request:

```jsonc
{
  "conversation_id": 0,
  "student_id": 5,
  "subject": "Matemáticas",
  "objective": "Sumar fracciones",
  "duration": "45 min",   // opcional
  "dynamic": "grupal",    // opcional
  "materials": "...",     // opcional
  "history": []           // opcional
}
```

Response `200`: igual a assist pero con `device_id` en vez de `recommended_device`.

### Marcadores embebidos en `response`

El texto de `response` puede contener marcadores que el FE **debe parsear y
renderizar como chips/links**, no mostrarlos crudos:

| Marcador | Significado | Acción FE sugerida |
|----------|-------------|--------------------|
| `[DEVICE_ID:N]` | Referencia a un dispositivo | Chip clickable → ficha del dispositivo |
| `[STUDENT_ID:N]` | Referencia a un alumno | Chip clickable → perfil del alumno |
| `[ADAPTATION_JSON:{...}]` | Adaptación generada | NO mostrar el JSON; usar el objeto `adaptation` ya parseado de la respuesta |

**A construir en el FE:**
- Componente de chat que mantenga `conversation_id` entre turnos.
- Parser de marcadores → chips/enlaces.
- Si llega `adaptation`, botón "Guardar adaptación" (`POST /adaptations`) y
  "Exportar" (ver §3).
- Toggle de `mode` libre/guiado.

### Manejo de `429 Too Many Requests` (NUEVO — rate limiting por organización)

Los endpoints `/inclusion/assist` y `/inclusion/recommend` están limitados por
organización. Al excederse:

```json
{ "code": "rate_limited", "message": "rate limit exceeded" }
```

**A construir en el FE:**
- Detectar `429` en estos dos endpoints → mostrar aviso "Demasiadas consultas,
  esperá un momento" y deshabilitar el botón de enviar temporalmente.
- No reintentar automáticamente en loop.

### Manejo de `503 service_unavailable` (NUEVO — circuit breaker)

Si el proveedor de IA está caído, el circuit breaker corta y devuelve `503`
`service_unavailable` de forma inmediata (sin esperar timeout).

**A construir en el FE:**
- Detectar `503` → mensaje "El asistente no está disponible por el momento".

---

## 2. Historial de chat (LISTO en backend)

### `GET /api/v1/chat/history/:contextId`

`contextId` = el `mode` ("assist", "recommend", "guided", etc.). Response `200`:

```jsonc
[
  {
    "id": 42,
    "mode": "assist",
    "created_at": "2026-05-26T10:00:00Z",
    "messages": [
      { "role": "user", "content": "...", "created_at": "..." },
      { "role": "assistant", "content": "...", "created_at": "..." }
    ]
  }
]
```

**A construir en el FE:** panel/lista de conversaciones previas; al abrir una,
cargar `messages` y setear `conversation_id` para continuarla.

---

## 3. Exportar adaptación (LISTO en backend)

### `GET /api/v1/adaptations/:id/export?format=pdf|md`

Devuelve un archivo binario con `Content-Disposition: attachment`.
`format` por defecto: `pdf`.

**A construir en el FE:** botones "Descargar PDF" / "Descargar Markdown".
Como requiere header `Authorization`, no usar un `<a href>` directo: hacer
`fetch` con el token, obtener el blob y disparar la descarga con
`URL.createObjectURL`.

---

## 4. Dashboard de métricas (LISTO en backend)

### `GET /api/v1/dashboard/metrics`

```jsonc
{
  "total_students": 120,
  "students_with_profiles": 80,
  "total_adaptations": 45,
  "adaptations_by_status": { "draft": 10, "active": 35 },
  "adaptations_by_type": { "actividad_adaptada": 20 },
  "top_used_devices": [
    { "device_id": 1, "device_name": "Timer Visual", "count": 12 }
  ],
  "adaptations_this_week": 7,
  "classroom_count": 6
}
```

**A construir en el FE:** tarjetas de resumen + gráfico de barras de
`top_used_devices` y de `adaptations_by_status/type`.

---

## 5. PENDIENTE en backend — el FE debe prepararse

Estas features aún no están en el backend. Se documentan acá para que el FE
reserve el diseño y, cuando el backend las exponga, la integración sea directa.

### 5.1 Streaming SSE del chat (Fase 1a)

Hoy `assist`/`recommend` responden el texto completo de una sola vez. Cuando se
implemente streaming, el contrato propuesto será:

- Mismo endpoint, con header `Accept: text/event-stream`.
- Respuesta `Content-Type: text/event-stream`, eventos:
  - `event: token` / `data: {"delta": "texto parcial"}`
  - `event: done` / `data: {<el mismo objeto JSON de la respuesta no-streaming>}`

**A construir en el FE (cuando esté disponible):**
- Cliente SSE (`EventSource` no soporta headers → usar `fetch` + `ReadableStream`
  por el `Authorization`).
- Render incremental del texto; al recibir `done`, parsear marcadores y
  `adaptation` como en §1.

### 5.2 Multimodal / visión (Fase 4a)

Permitir adjuntar imágenes (ej. foto de una actividad o material) al chat.
Contrato propuesto: `assist` aceptará un campo `images: [{ "url" | "base64", "mime" }]`.

**A construir en el FE (cuando esté disponible):**
- Input de carga de imagen en el chat (con preview y validación de tamaño/tipo).
- Envío de la imagen en el formato que defina el backend.

### 5.3 Indicador de herramientas / function calling (Fase 3)

Cuando el backend pase al loop agéntico, la IA podrá ejecutar acciones
(buscar alumno, listar dispositivos, crear adaptación) durante una misma
respuesta. El contrato no-streaming no cambia, pero conviene que el FE muestre
estado intermedio ("Alizia está buscando el perfil de Lucas…").

**A construir en el FE (cuando esté disponible):**
- Estado de "pensando/ejecutando acción" mientras el backend corre el loop.
- Opcional: render de los pasos de herramientas si el backend los expone.

### 5.4 Dashboard de uso de IA / tokens — ✅ LISTO

`GET /api/v1/dashboard/ai-usage?days=30` (`days` opcional, default 30, máx 365).

```jsonc
{
  "window_days": 30,
  "total_requests": 42,
  "prompt_tokens": 12000,
  "completion_tokens": 4500,
  "total_tokens": 16500,
  "by_mode": [
    { "mode": "assist", "requests": 30, "prompt_tokens": 9000,
      "completion_tokens": 3000, "total_tokens": 12000 },
    { "mode": "recommend", "requests": 12, "prompt_tokens": 3000,
      "completion_tokens": 1500, "total_tokens": 4500 }
  ]
}
```

**A construir en el FE:** vista admin con consumo de tokens del período y
desglose por modo (assist/recommend); selector de ventana (`days`).

---

## Resumen de estados

| Feature | Backend | FE debe construir |
|---------|---------|-------------------|
| Chat assist/recommend | ✅ Listo | Componente de chat + parser de marcadores |
| Rate limit `429` / CB `503` | ✅ Listo | Manejo de errores en chat |
| Historial de chat | ✅ Listo | Panel de conversaciones |
| Export PDF/MD | ✅ Listo | Botones de descarga (fetch + blob) |
| Dashboard métricas | ✅ Listo | Tarjetas + gráficos |
| Dashboard de tokens | ✅ Listo | Vista admin de uso (`/dashboard/ai-usage`) |
| Streaming SSE | ⏳ Pendiente | Cliente SSE + render incremental |
| Multimodal/visión | ⏳ Pendiente | Carga de imágenes |
| Function calling UI | ⚙️ Backend listo (flag off) | Estado de ejecución de acciones |
