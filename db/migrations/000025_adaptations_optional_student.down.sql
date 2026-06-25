-- 000025_adaptations_optional_student.down.sql
-- Revierte: vuelve student_id a NOT NULL. Falla si existen filas con student_id NULL
-- (recursos de situación sin alumno) — limpiarlas antes de revertir.
ALTER TABLE adaptations ALTER COLUMN student_id SET NOT NULL;
