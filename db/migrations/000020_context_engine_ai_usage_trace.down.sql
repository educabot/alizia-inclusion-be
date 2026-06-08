-- Revierte 000020

ALTER TABLE ai_usage DROP COLUMN IF EXISTS context_snapshot;
ALTER TABLE ai_usage DROP COLUMN IF EXISTS tool_calls;
ALTER TABLE ai_usage DROP COLUMN IF EXISTS latency_ms;
ALTER TABLE ai_usage DROP COLUMN IF EXISTS model;
ALTER TABLE ai_usage DROP COLUMN IF EXISTS prompt_version_id;
ALTER TABLE ai_usage DROP COLUMN IF EXISTS message_id;
ALTER TABLE ai_usage DROP COLUMN IF EXISTS conversation_id;
