-- 000029_drop_context_engine_remnants.down.sql
-- Recrea las tablas eliminadas con su DDL original (000019 / 000022). Los datos
-- (seed few-shot) no se restauran; volver a correr el seed si se necesitan.

-- response_examples (de 000019)
CREATE TABLE IF NOT EXISTS response_examples (
    id               BIGSERIAL PRIMARY KEY,
    organization_id  UUID REFERENCES organizations(id),
    mode             VARCHAR(50) NOT NULL,
    context_snapshot JSONB NOT NULL DEFAULT '{}',
    response         TEXT,
    label            VARCHAR(50) NOT NULL,
    tags             TEXT[],
    source           VARCHAR(50),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_response_examples_mode_label ON response_examples(mode, label);
CREATE INDEX IF NOT EXISTS idx_response_examples_tags ON response_examples USING GIN (tags);

-- integradora_assignments (de 000022)
CREATE TABLE IF NOT EXISTS integradora_assignments (
    id              BIGSERIAL PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    student_id      BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, student_id)
);
CREATE INDEX IF NOT EXISTS idx_integradora_assignments_org ON integradora_assignments(organization_id);
CREATE INDEX IF NOT EXISTS idx_integradora_assignments_user ON integradora_assignments(user_id);
CREATE INDEX IF NOT EXISTS idx_integradora_assignments_student ON integradora_assignments(student_id);
