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

-- Usuario de prueba = la identidad que inyecta el modo test (user_id=1). Necesario
-- para los FKs de adaptations.teacher_id / ppi.created_by, etc.
INSERT INTO users (id, organization_id, email, name, role)
VALUES (1, '00000000-0000-0000-0000-000000000001', 'test@educabot.com', 'Test User', 'teacher')
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

-- Perfiles (1:1 con alumno, UNIQUE(student_id)). Capa rica de HU-2: situation_codes
-- usa el vocabulario controlado de situations_catalog; el resto es opcional.
INSERT INTO student_profiles (
    student_id, is_transitory, difficulties, free_description,
    support_level, strengths, interests, triggers,
    effective_strategies, ineffective_strategies, situation_codes,
    has_therapeutic_companion, environment_notes
)
VALUES
    (9001, false, ARRAY['se_distrae_facilmente', 'impulsividad'],
        'Le cuesta sostener la atención en consignas largas; responde bien a pausas activas.',
        'medio', ARRAY['memoria visual', 'creatividad'], ARRAY['dinosaurios', 'dibujo'],
        ARRAY['ruidos fuertes', 'consignas extensas'],
        ARRAY['pausas activas', 'consignas cortas paso a paso'], ARRAY['retos en público'],
        ARRAY['no_inicia_tarea', 'se_distrae'],
        true, 'Acompañante terapéutico 3 veces por semana; familia muy presente.'),
    (9002, true, ARRAY['dificultad_lectura'],
        'Dificultad transitoria de lectura tras cambio de escuela; en proceso de adaptación.',
        'bajo', ARRAY['oralidad'], ARRAY['fútbol'], NULL,
        ARRAY['lectura compartida'], NULL, ARRAY['dificultad_consignas'],
        false, NULL),
    (9003, false, ARRAY['hipersensibilidad_sensorial', 'rutinas_rigidas'],
        'Necesita anticipación de cambios y entornos con baja estimulación sensorial.',
        'alto', ARRAY['atención al detalle'], ARRAY['trenes', 'mapas'],
        ARRAY['cambios imprevistos', 'luces brillantes'],
        ARRAY['anticipación con pictogramas', 'rincón de calma'], ARRAY['improvisar la rutina'],
        ARRAY['se_desregula', 'no_sostiene_atencion'],
        true, 'Sensibilidad sensorial alta; trabaja con fonoaudióloga externa.')
ON CONFLICT (student_id) DO UPDATE
    SET is_transitory             = EXCLUDED.is_transitory,
        difficulties              = EXCLUDED.difficulties,
        free_description          = EXCLUDED.free_description,
        support_level             = EXCLUDED.support_level,
        strengths                 = EXCLUDED.strengths,
        interests                 = EXCLUDED.interests,
        triggers                  = EXCLUDED.triggers,
        effective_strategies      = EXCLUDED.effective_strategies,
        ineffective_strategies    = EXCLUDED.ineffective_strategies,
        situation_codes           = EXCLUDED.situation_codes,
        has_therapeutic_companion = EXCLUDED.has_therapeutic_companion,
        environment_notes         = EXCLUDED.environment_notes;

-- Enriquecer alumnos con doble granularidad de edad / grado / nombre preferido.
UPDATE students SET age_range = '8-9', grade_level = '4to', preferred_name = 'Tomi'  WHERE id = 9001;
UPDATE students SET age_range = '9-10', grade_level = '4to', preferred_name = 'Lu'    WHERE id = 9002;
UPDATE students SET age_range = '8-9', grade_level = '4to', preferred_name = 'Mate'   WHERE id = 9003;

-- Diagnóstico de prueba (catálogo global) + asignación al perfil de 9001 (capa
-- secundaria opcional; 9002 y 9003 quedan sin diagnóstico para mostrar degradación).
INSERT INTO diagnoses_catalog (id, organization_id, name, category)
VALUES (9001, NULL, 'TDAH (prueba)', 'neurodesarrollo')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, category = EXCLUDED.category;

INSERT INTO student_diagnoses (student_profile_id, diagnosis_id, severity, notes)
SELECT sp.id, 9001, 'leve', 'Diagnóstico de prueba para validar el Context Assembler.'
FROM student_profiles sp
WHERE sp.student_id = 9001
ON CONFLICT (student_profile_id, diagnosis_id) DO UPDATE
    SET severity = EXCLUDED.severity, notes = EXCLUDED.notes;

-- PPI de prueba para 9001 (1:1 con alumno). 9002/9003 sin PPI a propósito.
INSERT INTO ppi (id, organization_id, student_id, objectives, curricular_adaptations, follow_up, status)
VALUES (9001, '00000000-0000-0000-0000-000000000001', 9001,
    ARRAY['Sostener la atención 15 min con apoyos', 'Iniciar la tarea de forma autónoma'],
    'Consignas segmentadas y tiempos extendidos en evaluaciones.',
    'Revisión mensual con la maestra integradora.',
    'active')
ON CONFLICT (student_id) DO UPDATE
    SET objectives             = EXCLUDED.objectives,
        curricular_adaptations = EXCLUDED.curricular_adaptations,
        follow_up              = EXCLUDED.follow_up,
        status                 = EXCLUDED.status;
