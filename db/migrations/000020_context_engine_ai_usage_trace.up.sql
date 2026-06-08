-- Context Engine — Capa D: evolucionar ai_usage in-place con columnas de traza (todas nullable).
-- Sin FK (consistente con ai_usage actual, que tampoco referencia organizations). El tablero del
-- Director (GROUP BY mode) ignora estas columnas y sigue igual. prompt_version_id apuntará a
-- prompt_versions cuando esa tabla exista (Futuro); por ahora es solo un BIGINT de traza.

ALTER TABLE ai_usage ADD COLUMN IF NOT EXISTS conversation_id   BIGINT;
ALTER TABLE ai_usage ADD COLUMN IF NOT EXISTS message_id        BIGINT;
ALTER TABLE ai_usage ADD COLUMN IF NOT EXISTS prompt_version_id BIGINT;
ALTER TABLE ai_usage ADD COLUMN IF NOT EXISTS model             VARCHAR(100);
ALTER TABLE ai_usage ADD COLUMN IF NOT EXISTS latency_ms        INTEGER NOT NULL DEFAULT 0;
ALTER TABLE ai_usage ADD COLUMN IF NOT EXISTS tool_calls        INTEGER NOT NULL DEFAULT 0;
ALTER TABLE ai_usage ADD COLUMN IF NOT EXISTS context_snapshot  JSONB NOT NULL DEFAULT '{}';
