-- Context Engine — Capa A (parte 4): asignación maestra integradora ↔ alumno.
-- El motor de contexto necesita saber qué alumnos tiene asignados una maestra
-- integradora para cargar su contexto. El RBAC en sí queda fuera de scope (otro
-- equipo); acá solo habilitamos que el dato exista y llegue al prompt.

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
