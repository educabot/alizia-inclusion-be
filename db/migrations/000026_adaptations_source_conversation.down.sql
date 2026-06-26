-- 000026_adaptations_source_conversation.down.sql
ALTER TABLE adaptations DROP COLUMN IF EXISTS source_message_id;
ALTER TABLE adaptations DROP COLUMN IF EXISTS source_conversation_id;
