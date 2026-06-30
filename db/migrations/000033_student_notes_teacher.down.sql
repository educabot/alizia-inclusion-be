-- 000030_student_notes_teacher.down.sql
DROP INDEX IF EXISTS idx_student_notes_student_user;
ALTER TABLE student_notes DROP COLUMN IF EXISTS user_id;
