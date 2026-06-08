-- Context Engine — Capa B: memoria de conversación (resumen + vínculos N a entidades)

CREATE TABLE IF NOT EXISTS conversation_summaries (
    conversation_id BIGINT PRIMARY KEY REFERENCES conversations(id) ON DELETE CASCADE,
    summary         TEXT,
    topic_keywords  TEXT[],
    token_count     INTEGER NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_conversation_summaries_topics
    ON conversation_summaries USING GIN (topic_keywords);

CREATE TABLE IF NOT EXISTS conversation_summary_students (
    conversation_id BIGINT NOT NULL REFERENCES conversation_summaries(conversation_id) ON DELETE CASCADE,
    student_id      BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    PRIMARY KEY (conversation_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_conv_summary_students_student
    ON conversation_summary_students(student_id);

CREATE TABLE IF NOT EXISTS conversation_summary_devices (
    conversation_id BIGINT NOT NULL REFERENCES conversation_summaries(conversation_id) ON DELETE CASCADE,
    device_id       BIGINT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    PRIMARY KEY (conversation_id, device_id)
);

CREATE INDEX IF NOT EXISTS idx_conv_summary_devices_device
    ON conversation_summary_devices(device_id);
