-- 000025_adaptations_optional_student.up.sql
-- Permite guardar una adaptación/recurso asociada a una SITUACIÓN, sin exigir un
-- alumno registrado (ni aula ni dispositivo). student_id pasa a ser nullable.
-- Idempotente: DROP NOT NULL es no-op si ya está nullable.
ALTER TABLE adaptations ALTER COLUMN student_id DROP NOT NULL;
