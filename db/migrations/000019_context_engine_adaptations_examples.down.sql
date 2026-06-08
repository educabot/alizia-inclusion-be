-- Revierte 000019

DROP TABLE IF EXISTS response_examples;

ALTER TABLE adaptations DROP COLUMN IF EXISTS was_edited;
ALTER TABLE adaptations DROP COLUMN IF EXISTS source_message_id;
ALTER TABLE adaptations DROP COLUMN IF EXISTS source_conversation_id;
