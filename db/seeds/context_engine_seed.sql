-- ============================================================
-- SEED — Context Engine (catálogos GLOBALES, sin PII)
-- organization_id = NULL  -> catálogo global definido por Educabot.
-- Idempotente: anti-join por NOT EXISTS (ON CONFLICT no dedupe con org NULL).
-- Corre DESPUÉS de las migraciones 000015-000021.
-- ============================================================

-- ---------- situations_catalog: ~15 situaciones observables de aula ----------
-- Entrada pedagógica primaria (lo observable, no el diagnóstico).
INSERT INTO situations_catalog (organization_id, code, name, description, phase, sort_order)
SELECT NULL, v.code, v.name, v.description, v.phase, v.sort_order
FROM (VALUES
  ('no_inicia_tarea',        'No inicia la tarea',                 'Cuesta arrancar; se queda sin empezar aunque entienda la consigna.',        'durante',     1),
  ('se_distrae',             'Se distrae constantemente',          'Pierde el foco con estímulos del entorno; atención intermitente.',          'durante',     2),
  ('se_desregula',           'Se desregula emocionalmente',        'Frustración, llanto o enojo que interrumpen el trabajo.',                   'durante',     3),
  ('no_sostiene_atencion',   'No sostiene la atención',            'Mantiene el foco poco tiempo y abandona antes de terminar.',                'durante',     4),
  ('dificultad_consignas',   'Dificultad para comprender consignas','No interpreta lo que se pide; necesita reformulación o apoyo visual.',       'preventiva',  5),
  ('no_termina_tarea',       'No termina la tarea',                'Comienza pero no llega al cierre dentro del tiempo de clase.',              'cierre',      6),
  ('ritmo_lento',            'Trabaja a ritmo muy lento',          'Avanza por debajo del ritmo del grupo; necesita más tiempo.',               'durante',     7),
  ('dificultad_lectura',     'Dificultad en la lectura',           'Decodificación lenta o con errores; impacta la comprensión.',               'durante',     8),
  ('dificultad_escritura',   'Dificultad en la escritura',         'Trazo, ortografía u organización de ideas por escrito.',                    'durante',     9),
  ('dificultad_matematica',  'Dificultad en matemática',           'Dificultad con conteo, operaciones o razonamiento numérico.',               'durante',    10),
  ('no_interactua_pares',    'No interactúa con pares',            'Se aísla o evita el trabajo en grupo.',                                     'preventiva', 11),
  ('conducta_disruptiva',    'Conducta disruptiva',                'Interrumpe la clase; comportamientos que afectan al grupo.',                'durante',    12),
  ('fatiga_sensorial',       'Sobrecarga sensorial',               'Ruido, luz o estímulos que generan saturación y bloqueo.',                  'preventiva', 13),
  ('dependencia_adulto',     'Dependencia del adulto',             'No avanza sin acompañamiento individual constante.',                        'durante',    14),
  ('dificultad_transiciones','Dificultad en las transiciones',     'Le cuesta el cambio de actividad o de espacio.',                            'preventiva', 15)
) AS v(code, name, description, phase, sort_order)
WHERE NOT EXISTS (
  SELECT 1 FROM situations_catalog s WHERE s.organization_id IS NULL AND s.code = v.code
);

-- ---------- diagnoses_catalog: capa secundaria, opcional ----------
INSERT INTO diagnoses_catalog (organization_id, name, category)
SELECT NULL, v.name, v.category
FROM (VALUES
  ('Trastorno del Espectro Autista (TEA)',                'neurodesarrollo'),
  ('Trastorno por Déficit de Atención e Hiperactividad (TDAH)', 'neurodesarrollo'),
  ('Dislexia',                                            'aprendizaje'),
  ('Discalculia',                                         'aprendizaje'),
  ('Disgrafía',                                           'aprendizaje'),
  ('Discapacidad intelectual',                            'intelectual'),
  ('Discapacidad visual',                                 'sensorial'),
  ('Discapacidad auditiva',                               'sensorial'),
  ('Discapacidad motriz',                                 'motora'),
  ('Trastorno específico del lenguaje',                   'lenguaje')
) AS v(name, category)
WHERE NOT EXISTS (
  SELECT 1 FROM diagnoses_catalog d WHERE d.organization_id IS NULL AND d.name = v.name
);

-- ---------- response_examples: few-shot golden 'curated' (cold-start) ----------
-- Ligados a su situación vía tags. mode='assist', label='golden', source='curated'.
INSERT INTO response_examples (organization_id, mode, context_snapshot, response, label, tags, source)
SELECT NULL, v.mode, '{}'::jsonb, v.response, 'golden', v.tags, 'curated'
FROM (VALUES
  ('assist',
   'Para que arranque: dividí la consigna en 1 solo paso visible a la vez y dale un inicio concreto ("escribí solo el título"). Usá un time timer de la valija para marcar 5 minutos de foco. Reforzá apenas empieza, no recién al terminar.',
   ARRAY['no_inicia_tarea','dependencia_adulto']),
  ('assist',
   'Anticipá los cambios: avisá la transición 2 minutos antes con apoyo visual (pictograma o cartel de "ahora / después"). Para la sobrecarga, ofrecé auriculares de cancelación de ruido y un rincón de calma. Mantené la rutina estable.',
   ARRAY['dificultad_transiciones','fatiga_sensorial','TEA']),
  ('assist',
   'Para la lectura: entregá el texto en fragmentos cortos con tipografía accesible y mayor interlineado. Permití lectura en voz alta o apoyo audio. Evaluá la comprensión de forma oral, separándola de la decodificación.',
   ARRAY['dificultad_lectura','dislexia'])
) AS v(mode, response, tags)
WHERE NOT EXISTS (
  SELECT 1 FROM response_examples r WHERE r.source = 'curated' AND r.response = v.response
);

-- ---------- pedagogical_content (RAG): documentos globales de ejemplo ----------
-- MVP: 1 chunk = documento entero; embedding queda NULL hasta fijar el modelo de Azure.
INSERT INTO pedagogical_content (parent_id, type, title, status, keywords, organization_id)
SELECT NULL, v.type, v.title, 'published', v.keywords, NULL
FROM (VALUES
  ('material', 'Estrategias para acompañar a estudiantes con TEA en el aula',
     ARRAY['TEA','autismo','autorregulacion','anticipacion','transiciones']),
  ('material', 'Técnicas de lectura para estudiantes con dislexia',
     ARRAY['dislexia','lectura','comprension','accesibilidad']),
  ('material', 'Pautas para sostener la atención y el inicio de tareas',
     ARRAY['TDAH','atencion','consignas','no_inicia_tarea'])
) AS v(type, title, keywords)
WHERE NOT EXISTS (
  SELECT 1 FROM pedagogical_content pc WHERE pc.title = v.title
);

-- chunk por documento (1 chunk = documento entero en el MVP)
INSERT INTO pedagogical_content_chunks (content_id, chunk_text, tags)
SELECT pc.id, v.chunk_text, v.tags
FROM (VALUES
  ('Estrategias para acompañar a estudiantes con TEA en el aula',
   'Anticipar los cambios con apoyos visuales, ofrecer un rincón de calma y materiales para regular estímulos (auriculares, time timer). Dividir las consignas en pasos cortos y mantener rutinas predecibles favorece la autorregulación.',
   ARRAY['TEA','autismo','autorregulacion']),
  ('Técnicas de lectura para estudiantes con dislexia',
   'Presentar textos en fragmentos cortos, con tipografía accesible e interlineado amplio. Combinar lectura con audio y evaluar la comprensión de forma oral, separada de la decodificación.',
   ARRAY['dislexia','lectura','comprension']),
  ('Pautas para sostener la atención y el inicio de tareas',
   'Mostrar un solo paso a la vez, dar un inicio concreto y usar marcadores de tiempo. Reforzar el comienzo de la tarea y reducir estímulos distractores en el entorno cercano.',
   ARRAY['TDAH','atencion','no_inicia_tarea'])
) AS v(title, chunk_text, tags)
JOIN pedagogical_content pc ON pc.title = v.title
WHERE NOT EXISTS (
  SELECT 1 FROM pedagogical_content_chunks c WHERE c.content_id = pc.id
);
