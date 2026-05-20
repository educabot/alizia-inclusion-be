CREATE TABLE IF NOT EXISTS adaptation_resources (
    id            BIGSERIAL PRIMARY KEY,
    adaptation_id BIGINT NOT NULL REFERENCES adaptations(id) ON DELETE CASCADE,
    title         VARCHAR(255) NOT NULL,
    file_url      TEXT NOT NULL,
    file_type     VARCHAR(50) DEFAULT 'pdf',
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_adaptation_resources_adaptation ON adaptation_resources(adaptation_id);
