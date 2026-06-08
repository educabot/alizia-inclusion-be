-- Context Engine — Capa A (parte 3): diagnósticos estructurados (capa secundaria, opcional)

CREATE TABLE IF NOT EXISTS diagnoses_catalog (
    id              BIGSERIAL PRIMARY KEY,
    organization_id UUID REFERENCES organizations(id),  -- NULL = global (Educabot)
    name            VARCHAR(255) NOT NULL,
    category        VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS student_diagnoses (
    id                 BIGSERIAL PRIMARY KEY,
    student_profile_id BIGINT NOT NULL REFERENCES student_profiles(id) ON DELETE CASCADE,
    diagnosis_id       BIGINT NOT NULL REFERENCES diagnoses_catalog(id),
    severity           VARCHAR(50),
    notes              TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (student_profile_id, diagnosis_id)
);

CREATE INDEX IF NOT EXISTS idx_student_diagnoses_profile ON student_diagnoses(student_profile_id);
