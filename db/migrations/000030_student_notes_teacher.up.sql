-- 000030_student_notes_teacher.up.sql
-- Notas privadas del docente: cada nota pasa a tener dueño (user_id). Al ver el
-- perfil de un alumno, un docente solo ve SUS notas. Nullable: las filas legacy
-- (sin dueño) quedan invisibles bajo el filtro por docente. Idempotente.
ALTER TABLE student_notes ADD COLUMN IF NOT EXISTS user_id BIGINT REFERENCES users(id);
CREATE INDEX IF NOT EXISTS idx_student_notes_student_user ON student_notes (student_id, user_id, created_at DESC);
