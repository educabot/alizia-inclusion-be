-- 000026_adaptations_source_conversation.up.sql
-- Traza el origen de una adaptación creada desde el chat: permite que el FE ofrezca
-- "Ver conversación de origen" (GAP B del contrato BE/FE). Ambas columnas nullable:
-- una adaptación puede crearse fuera del chat. Idempotente (ADD COLUMN IF NOT EXISTS).
ALTER TABLE adaptations ADD COLUMN IF NOT EXISTS source_conversation_id BIGINT;
ALTER TABLE adaptations ADD COLUMN IF NOT EXISTS source_message_id BIGINT;
