-- Seed de prueba para el modo de auth `test` (ENV=test).
--
-- En ENV=test el middleware de auth inyecta una identidad fija (ver
-- src/entrypoints/middleware/auth.go): user_id=1 y org_uuid=
-- 00000000-0000-0000-0000-000000000001. Los datos sembrados por el seed
-- principal (context_engine_seed.sql) viven bajo OTRA organización, así que
-- el camino "alumno" de HU-1 (POST /api/v1/inclusion/open) no encontraría
-- ningún alumno al pegarle desde Postman sin token.
--
-- Este seed crea, bajo la org de ceros, una organización, un aula y 3 alumnos
-- con perfil para poder probar HU-1 end-to-end. Es idempotente: se puede
-- correr varias veces sin duplicar (ON CONFLICT por PK / unique).
--
-- Uso:
--   DB_URL="postgres://..." go run ./scripts/dbmigrate db/seeds/test_org_seed.sql
--
-- IDs reservados en rango alto (9001+) para no colisionar con datos reales.

-- Organización de prueba (la que inyecta el modo test).
INSERT INTO organizations (id, name)
VALUES ('00000000-0000-0000-0000-000000000001', 'Org de prueba (ENV=test)')
ON CONFLICT (id) DO NOTHING;

-- Aula de prueba bajo la org de ceros.
INSERT INTO classrooms (id, organization_id, name, grade, section)
VALUES (9001, '00000000-0000-0000-0000-000000000001', 'Aula de prueba', '4to', 'A')
ON CONFLICT (id) DO UPDATE
    SET organization_id = EXCLUDED.organization_id,
        name            = EXCLUDED.name,
        grade           = EXCLUDED.grade,
        section         = EXCLUDED.section;

-- Alumnos de prueba bajo la org de ceros, aula 9001.
INSERT INTO students (id, organization_id, classroom_id, name)
VALUES
    (9001, '00000000-0000-0000-0000-000000000001', 9001, 'Tomás Prueba'),
    (9002, '00000000-0000-0000-0000-000000000001', 9001, 'Lucía Demo'),
    (9003, '00000000-0000-0000-0000-000000000001', 9001, 'Mateo Test')
ON CONFLICT (id) DO UPDATE
    SET organization_id = EXCLUDED.organization_id,
        classroom_id    = EXCLUDED.classroom_id,
        name            = EXCLUDED.name;

-- Perfiles (1:1 con alumno, UNIQUE(student_id)).
INSERT INTO student_profiles (student_id, is_transitory, difficulties, free_description)
VALUES
    (9001, false, ARRAY['se_distrae_facilmente', 'impulsividad'],
        'Le cuesta sostener la atención en consignas largas; responde bien a pausas activas.'),
    (9002, true, ARRAY['dificultad_lectura'],
        'Dificultad transitoria de lectura tras cambio de escuela; en proceso de adaptación.'),
    (9003, false, ARRAY['hipersensibilidad_sensorial', 'rutinas_rigidas'],
        'Necesita anticipación de cambios y entornos con baja estimulación sensorial.')
ON CONFLICT (student_id) DO UPDATE
    SET is_transitory    = EXCLUDED.is_transitory,
        difficulties     = EXCLUDED.difficulties,
        free_description = EXCLUDED.free_description;
