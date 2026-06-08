-- Context Engine — Capa A (parte 1): rol maestra integradora + perfil docente + PPI

ALTER TYPE member_role ADD VALUE IF NOT EXISTS 'maestra_integradora';

CREATE TABLE IF NOT EXISTS teacher_profiles (
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id  UUID NOT NULL REFERENCES organizations(id),
    birthdate        DATE,
    age_range        VARCHAR(50),
    years_experience INTEGER,
    specialization   VARCHAR(255),
    subjects         TEXT[],
    tone_preference  VARCHAR(50),
    bio              TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id)
);

CREATE INDEX IF NOT EXISTS idx_teacher_profiles_org ON teacher_profiles(organization_id);

CREATE TABLE IF NOT EXISTS ppi (
    id                     BIGSERIAL PRIMARY KEY,
    organization_id        UUID NOT NULL REFERENCES organizations(id),
    student_id             BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    objectives             TEXT[],
    curricular_adaptations TEXT,
    follow_up              TEXT,
    status                 VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_by             BIGINT REFERENCES users(id),
    validated_by           BIGINT REFERENCES users(id),
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (student_id)
);

CREATE INDEX IF NOT EXISTS idx_ppi_org ON ppi(organization_id);
