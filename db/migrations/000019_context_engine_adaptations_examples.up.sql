-- Context Engine — Capa C: traza de origen IA en adaptaciones + banco de few-shot

-- adaptations: liga la sugerencia de IA con su resultado real + aceptación implícita
ALTER TABLE adaptations ADD COLUMN IF NOT EXISTS source_conversation_id BIGINT REFERENCES conversations(id);
ALTER TABLE adaptations ADD COLUMN IF NOT EXISTS source_message_id      BIGINT REFERENCES conversation_messages(id);
ALTER TABLE adaptations ADD COLUMN IF NOT EXISTS was_edited             BOOLEAN NOT NULL DEFAULT FALSE;

-- response_examples: few-shot golden/bad + set de evaluación
CREATE TABLE IF NOT EXISTS response_examples (
    id               BIGSERIAL PRIMARY KEY,
    organization_id  UUID REFERENCES organizations(id),  -- NULL = global (Educabot)
    mode             VARCHAR(50) NOT NULL,
    context_snapshot JSONB NOT NULL DEFAULT '{}',
    response         TEXT,
    label            VARCHAR(50) NOT NULL,               -- golden / bad
    tags             TEXT[],
    source           VARCHAR(50),                        -- curated / from_outcome
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_response_examples_mode_label ON response_examples(mode, label);
CREATE INDEX IF NOT EXISTS idx_response_examples_tags ON response_examples USING GIN (tags);
