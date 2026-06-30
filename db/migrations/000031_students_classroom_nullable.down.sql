-- Revierte: vuelve classroom_id a NOT NULL. Falla si hay alumnos sin aula (classroom_id
-- NULL); en ese caso hay que reasignarlos antes de revertir.
ALTER TABLE students ALTER COLUMN classroom_id SET NOT NULL;
