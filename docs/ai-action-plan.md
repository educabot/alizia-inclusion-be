# AI Action Plan — Alizia Inclusion

**Owner:** Sebastian Luser
**Started:** 2026-05-26
**Status:** Fase 1 en curso

Plan para llevar a Alizia de "asistente que responde texto" a "asistente interactuable, persistente, agéntico y exportable".

---

## Estado base (snapshot 2026-05-26)

### Ya funciona
- `POST /recommend-device` — recomendación de dispositivo dado alumno + contexto pedagógico
- `POST /assist-classroom` — chat en aula, modos `normal` y `guided`
- `GET /chat-history/:contextId` — historial por modo
- Azure OpenAI client (`Chat`) + stub client para dev sin API key
- Extracción de marcadores `[DEVICE_ID]`, `[STUDENT_ID]`, `[ADAPTATION_JSON:{}]`
- Adaptaciones persisten como entidad con CRUD completo

### Gaps
1. Sin streaming — request/response bloqueante
2. Conversaciones no se persisten automáticamente (cliente manda `History` en cada call)
3. `ChatWithTools` sin implementar (segundo arg ignorado)
4. Sin exportación de documentos (PDF/DOCX)
5. Sin truncation de historial largo
6. Sin tracking de tokens/costo
7. Sin rate limiting
8. Sin multimodal (imágenes/audio)

---

## Fase 0 — Verificación post-deploy

**Objetivo:** confirmar que el deploy a Railway dejó el sistema operativo.

- [ ] `GET /health` responde 200 con `db: ok`
- [ ] `POST /recommend-device` con JWT real responde (no cae al stub)
- [ ] `POST /assist-classroom` igual
- [ ] Logs estructurados visibles en Railway dashboard
- [ ] Azure OpenAI consume tokens (validar en portal Azure)

---

## Fase 1 — Chat conversacional real

**Objetivo:** que sienta vivo como ChatGPT — streaming + persistencia + reanudación.

### 1a — Streaming SSE
- Agregar `ChatStream(ctx, msgs) (<-chan ChatChunk, error)` al `AIClient`
- Implementar en `azure_client.go` con `stream: true` y parser de SSE de Azure
- Adaptar `assist_classroom` y `recommend_device` a streaming
- Handlers SSE que escriben `text/event-stream` chunk a chunk
- Stub client devuelve chunks fake para tests
- Front recibe progresivamente (UI feel "está escribiendo")

### 1b — Persistencia automática
- Extender `ConversationProvider`:
  - `AppendTurn(ctx, orgID, userID, mode, contextID, userMsg, assistantMsg, metadata) error`
  - `GetOrCreate(ctx, orgID, userID, mode, contextID) (*Conversation, error)`
- `assist_classroom` y `recommend_device` llaman `AppendTurn` después de cada respuesta del AI
- Metadata guarda: marcadores extraídos (device_id, student_id, adaptation), tokens consumidos
- Regenerar mocks vía `make mocks`

### 1c — Reanudación
- `GET /chat-history/:contextId` ya existe — validar que devuelve mensajes en orden
- Agregar `GET /chat-history/:contextId/:conversation_id` para retomar conversación específica
- Front llama esto al cargar y hidrata el chat

### 1d — Truncation de historial largo
- Helper `truncateHistory(msgs []ChatMessage, maxTokens int) []ChatMessage`
- Cuenta aproximada de tokens (regla 4 chars ≈ 1 token)
- Si excede, conserva system + últimos N + resumen comprimido de los viejos
- Resumen vía llamada al AI cuando excede umbral

---

## Fase 2 — Generación de documentos

**Objetivo:** docente baja la planificación como archivo compartible/imprimible.

### 2a — Endpoint
- `POST /adaptations/:id/export?format=pdf|docx|md`
- Auth + tenant check vía middleware existente

### 2b — Generadores
- Markdown: trivial, string templating
- PDF: `gofpdf` (Go puro, sin headless Chrome)
- DOCX: `unioffice` o `gooxml`

### 2c — Template pedagógico
- Header con branding Educabot
- Datos del alumno (con opción de anonimizar)
- Adaptación: título, tipo, estrategia
- Dispositivos sugeridos con fundamento + cómo usar
- Notas para el docente
- Footer con fecha + ID adaptación

### 2d — Plan de clase agregado
- `POST /classrooms/:id/export-plan?format=pdf` — exporta todas las adaptaciones activas de un aula

### 2e — Tests
- Golden files: comparar bytes de PDFs generados contra fixtures
- Skip si dependencia no instalada en CI

---

## Fase 3 — Function calling / Tools

**Objetivo:** la IA actúa sobre el dominio, no solo describe.

### 3a — Cliente
- Implementar `ChatWithTools` real en `azure_client.go`
- Soportar tool_calls de Azure OpenAI
- Loop: si respuesta es tool_call, ejecutar tool, devolver resultado, continuar

### 3b — Tools del dominio
- `create_adaptation(student_id, title, type, strategy, device_ids[])` — graba en BD
- `update_adaptation_status(id, status)` — para feedback de aula
- `search_devices(query, limit)` — búsqueda en catálogo (eventualmente embeddings)
- `list_students_with_difficulty(classroom_id, difficulty)` — query estructurada
- `mark_resource_used(student_id, resource_id)` — tracking de uso

### 3c — Dispatcher
- En el usecase (no en handler) — clean architecture
- Map de tool name → handler function
- Cada handler usa los providers ya inyectados

### 3d — Confirmación humana
- Tools de escritura devuelven `requires_confirmation: true` en metadata
- Front muestra preview + botón confirmar antes de ejecutar
- Configurable por tool

### 3e — Tests
- Mock AI client que devuelve tool_calls predefinidos
- Verificar side effects en BD (entity creado/actualizado)

---

## Fase 4 — Multimodal + observabilidad

### 4a — Imágenes (vision)
- `assist_classroom` acepta `image_urls []string` o `image_b64 []string`
- Azure OpenAI GPT-4o soporta image input nativo
- Use case: foto del cuaderno de un alumno, screenshot de actividad

### 4b — Audio STT (opcional)
- Endpoint `POST /transcribe` que recibe audio + devuelve texto
- Azure Speech SDK o Whisper vía Azure OpenAI

### 4c — Token tracking
- Tabla `ai_usage(org_id, user_id, mode, prompt_tokens, completion_tokens, cost_usd, created_at)`
- Cada call al AI registra uso
- Endpoint admin: `GET /admin/usage?org_id=...&from=...&to=...`

### 4d — Rate limiting
- Middleware con token bucket por `org_id` (in-memory inicialmente, Redis después)
- Configurable por env: `AI_RATE_LIMIT_PER_HOUR`
- Headers `X-RateLimit-Remaining`, `X-RateLimit-Reset`

### 4e — Dashboard
- `GET /admin/ai-usage` — uso por org, top users, costo total
- Solo accesible a rol admin

---

## Fase 5 — Hardening

### 5a — Prompt caching
- Validar que system prompts grandes (catálogo) golpean el prompt cache de Azure
- Reordenar prompts: parte estable arriba, parte variable abajo
- Monitorear `cached_tokens` en respuesta

### 5b — Circuit breaker
- Si Azure devuelve 5xx consecutivos, abrir circuito por N segundos
- Fallback: respuesta canned ("estoy con dificultades, intentá de nuevo en un momento")
- Logging de incidentes

### 5c — E2E smoke contra Azure real
- Workflow nocturno en GitHub Actions
- Prompts fijos con expected shape de respuesta (no contenido literal)
- Alerta si shape rompe (regresión de modelo o config)

---

## Orden de ejecución recomendado

```
Fase 0 (verif) → 1b (persist) → 1a (stream) → 1d (truncate) → 2 (export)
                                                                  ↓
                                                         3 (tools) → 4 → 5
```

**Mínimo viable "chat interactuable + docs":** 0 + 1b + 1a + 2 (~4-5 días)
**Diferenciación IA agéntica:** + 3 (~2-3 días extra)
**Producción robusta:** + 4 + 5 (~3-4 días extra)

---

## Decisiones abiertas

- ¿`contextId` en chat-history identifica conversación o solo modo? → Hoy es modo. Para Fase 1c hay que agregar `conversation_id` separado.
- ¿Confirmación humana de tools por defecto on/off? → Default ON para writes, OFF para reads.
- ¿Storage de PDFs generados? → Inicialmente stream directo en response. Si pide reuso, S3/Railway volume.
- ¿Embeddings para búsqueda de dispositivos? → Fase 3b empieza con LIKE, eventualmente pgvector.
