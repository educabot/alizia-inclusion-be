CREATE TABLE IF NOT EXISTS ai_usage (
    id                BIGSERIAL PRIMARY KEY,
    organization_id   UUID NOT NULL,
    user_id           BIGINT NOT NULL,
    mode              VARCHAR(50) NOT NULL DEFAULT '',
    prompt_tokens     INTEGER NOT NULL DEFAULT 0,
    completion_tokens INTEGER NOT NULL DEFAULT 0,
    total_tokens      INTEGER NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_usage_org_created ON ai_usage(organization_id, created_at);
