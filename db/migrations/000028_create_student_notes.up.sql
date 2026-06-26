-- 000028_create_student_notes.up.sql
-- Notas / seguimiento del alumno (Epic ALZ-246). Historial del alumno en el
-- perfil: diagnóstico u observables del alta (internal=true, no visible en el
-- front) y seguimiento en el tiempo. Idempotente.
CREATE TABLE IF NOT EXISTS student_notes (
    id              BIGSERIAL PRIMARY KEY,
    student_id      BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    organization_id uuid NOT NULL,
    content         text NOT NULL,
    type            text NOT NULL DEFAULT 'seguimiento', -- diagnostico | observable | seguimiento
    internal        boolean NOT NULL DEFAULT true,
    created_at      timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_student_notes_student ON student_notes (student_id, created_at DESC);
