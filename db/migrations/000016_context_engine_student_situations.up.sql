-- Context Engine — Capa A (parte 2): enriquecer alumno + catálogo de situaciones observables

-- students: doble granularidad y nombre preferido (todo nullable)
ALTER TABLE students ADD COLUMN IF NOT EXISTS birthdate       DATE;
ALTER TABLE students ADD COLUMN IF NOT EXISTS age_range       VARCHAR(50);
ALTER TABLE students ADD COLUMN IF NOT EXISTS grade_level     VARCHAR(50);
ALTER TABLE students ADD COLUMN IF NOT EXISTS preferred_name  VARCHAR(255);

-- student_profiles: capa rica de necesidades (todo nullable)
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS support_level             VARCHAR(50);
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS strengths                 TEXT[];
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS interests                 TEXT[];
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS triggers                  TEXT[];
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS effective_strategies      TEXT[];
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS ineffective_strategies    TEXT[];
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS situation_codes           TEXT[];
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS has_therapeutic_companion BOOLEAN;
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS environment_notes         TEXT;

-- situations_catalog: ~15 situaciones observables de aula (entrada pedagógica primaria)
CREATE TABLE IF NOT EXISTS situations_catalog (
    id              BIGSERIAL PRIMARY KEY,
    organization_id UUID REFERENCES organizations(id),  -- NULL = global (Educabot)
    code            VARCHAR(100) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    phase           VARCHAR(50),
    sort_order      INTEGER NOT NULL DEFAULT 0,
    UNIQUE (organization_id, code)
);
