-- Revierte 000015. El valor 'maestra_integradora' del enum member_role NO se revierte:
-- PostgreSQL no soporta DROP VALUE de un ENUM (mismo criterio que 000007).

DROP TABLE IF EXISTS ppi;
DROP TABLE IF EXISTS teacher_profiles;
