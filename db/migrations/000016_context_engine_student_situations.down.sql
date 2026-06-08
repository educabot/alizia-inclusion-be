-- Revierte 000016

DROP TABLE IF EXISTS situations_catalog;

ALTER TABLE student_profiles DROP COLUMN IF EXISTS environment_notes;
ALTER TABLE student_profiles DROP COLUMN IF EXISTS has_therapeutic_companion;
ALTER TABLE student_profiles DROP COLUMN IF EXISTS situation_codes;
ALTER TABLE student_profiles DROP COLUMN IF EXISTS ineffective_strategies;
ALTER TABLE student_profiles DROP COLUMN IF EXISTS effective_strategies;
ALTER TABLE student_profiles DROP COLUMN IF EXISTS triggers;
ALTER TABLE student_profiles DROP COLUMN IF EXISTS interests;
ALTER TABLE student_profiles DROP COLUMN IF EXISTS strengths;
ALTER TABLE student_profiles DROP COLUMN IF EXISTS support_level;

ALTER TABLE students DROP COLUMN IF EXISTS preferred_name;
ALTER TABLE students DROP COLUMN IF EXISTS grade_level;
ALTER TABLE students DROP COLUMN IF EXISTS age_range;
ALTER TABLE students DROP COLUMN IF EXISTS birthdate;
