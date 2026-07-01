-- 000034_create_message_feedback.up.sql
-- Feedback de mensajes del asistente (manito arriba/abajo + comentario). Uso
-- interno: la especialista marca dislike y comenta qué falló, para revisar los
-- errores de Alizia. Guardamos conversation_id además del mensaje para poder
-- reconstruir el hilo completo y entender el contexto. Idempotente.
CREATE TABLE IF NOT EXISTS message_feedback (
    id                      BIGSERIAL PRIMARY KEY,
    conversation_message_id BIGINT NOT NULL REFERENCES conversation_messages(id) ON DELETE CASCADE,
    conversation_id         BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    organization_id         uuid NOT NULL,
    user_id                 BIGINT NOT NULL,
    rating                  varchar(10) NOT NULL, -- like | dislike
    comment                 text,
    created_at              timestamptz NOT NULL DEFAULT now(),
    updated_at              timestamptz NOT NULL DEFAULT now(),
    -- Un feedback por usuario por mensaje: permite upsert / toggle.
    UNIQUE (conversation_message_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_message_feedback_conversation_id ON message_feedback (conversation_id);
CREATE INDEX IF NOT EXISTS idx_message_feedback_rating ON message_feedback (rating, created_at DESC);
