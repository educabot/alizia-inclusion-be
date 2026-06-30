-- El aula deja de ser obligatoria para un alumno: se puede crear sin aula (creación
-- no-bloqueante desde el chat) y completar el aula después. La FK a classrooms se mantiene
-- (NULL no la viola). Cambio aditivo y no destructivo.
ALTER TABLE students ALTER COLUMN classroom_id DROP NOT NULL;
