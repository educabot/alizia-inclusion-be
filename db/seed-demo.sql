-- =============================================================================
-- Seed de demo: Valija Adaptativa - Alizia Inclusión
-- Basado en el documento "Valijas adaptativa - Flujo de demo.pdf"
-- =============================================================================

BEGIN;

-- Org ID constante
-- a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

-- =============================================================================
-- 1. Perfiles faltantes: Sofía Álvarez (id=4) y Mateo García (id=5)
-- =============================================================================

INSERT INTO student_profiles (student_id, is_transitory, difficulties, free_description)
VALUES
  (4, false,
   '{acceso_tecnologico,motricidad_gruesa}',
   'Sofía presenta dificultades de motricidad gruesa que afectan su postura y manipulación de dispositivos. Necesita soportes ergonómicos y acceso tecnológico adaptado para participar en actividades digitales.'),
  (5, true,
   '{distraccion_constante,no_inicia_tarea}',
   'Mateo muestra dificultad para iniciar tareas y sostener la atención. Requiere consignas segmentadas, apoyos visuales de organización temporal y validaciones frecuentes. Condición transitoria en seguimiento.')
ON CONFLICT (student_id) DO NOTHING;

-- =============================================================================
-- 2. Adaptaciones: 10 casos basados en los escenarios del PDF
-- Docente María González (id=2) y Carlos López (id=3)
-- =============================================================================

INSERT INTO adaptations (organization_id, student_id, teacher_id, device_id, subject, activity_description, adaptation_strategy, adaptation_type, outcome, notes, status)
VALUES
-- Caso 1: Lucas - Dictado de texto (motricidad fina) - Pinzas de escritura
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 2, 10,
 'Prácticas del Lenguaje',
 'Dictado de un texto de 10 minutos. Los alumnos escriben lo que escuchan.',
 'Se entrega el adaptador de escritura (pinzas) para que Lucas pueda sujetar la lapicera de forma adaptada durante el dictado. Se reduce la extensión del dictado a frases cortas y se valida cada párrafo antes de continuar.',
 'planificada',
 'Lucas completó el 80% del dictado sin mostrar fatiga en la mano. Mejoró la legibilidad respecto a la clase anterior.',
 'Caso del flujo de demo: modo planificador. El adaptador permite participar sin restricción, reduce frustración y da autonomía inmediata.',
 'completada'),

-- Caso 2: Valentina - Comprensión lectora (acceso visual + bloqueo lectoescritor) - Pen reader
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 2, 6,
 'Prácticas del Lenguaje',
 'Lectura comprensiva de un cuento corto y respuesta a preguntas sobre el texto.',
 'Se utiliza el pen reader (lápiz lector) para que Valentina pueda escuchar el texto en voz alta mientras lo sigue visualmente. Se complementa con la regla de lectura con ventana para evitar saltos de línea.',
 'planificada',
 'Valentina pudo acceder al texto completo y respondió 4 de 5 preguntas correctamente. Mostró mayor autonomía que en clases anteriores.',
 'Se combinaron dos dispositivos: pen reader + regla de lectura. La combinación fue muy efectiva.',
 'completada'),

-- Caso 3: Tomás - Evaluación de Matemática (distracción + sobreestimulación) - Auriculares cancelación
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 2, 4,
 'Matemática',
 'Evaluación individual de resolución de problemas. Duración: 30 minutos.',
 'Se ofrecen auriculares de cancelación de ruido para reducir la sobrecarga sensorial durante la evaluación. Se segmenta la evaluación en 3 bloques de 10 minutos con el Time Timer visible, permitiendo micro-pausas entre bloques.',
 'planificada',
 'Tomás completó los 3 bloques. En el segundo bloque pidió una pausa breve. Resolvió correctamente 6 de 8 problemas.',
 'La combinación de cancelación de ruido + Time Timer fue clave para que sostenga la actividad.',
 'completada'),

-- Caso 4: Sofía - Actividad digital (acceso tecnológico) - Tablet + Soporte flexible
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, 3, 11,
 'Ciencias Naturales',
 'Investigación guiada sobre el ciclo del agua usando recursos digitales.',
 'Se configura la tablet con opciones de accesibilidad (texto grande, lector de pantalla). Se coloca en el soporte flexible para liberar las manos de Sofía y mantener postura cómoda. Se prepara la actividad antes de entregarla.',
 'planificada',
 'Sofía navegó los recursos de forma autónoma y completó el cuestionario digital. Pidió ayuda solo una vez para ajustar el soporte.',
 'El soporte flexible fue fundamental para que pueda trabajar sin fatiga postural.',
 'completada'),

-- Caso 5: Mateo - No inicia la tarea (distracción) - Time Timer + Pelota antiestrés
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, 3, 5,
 'Ciencias Sociales',
 'Lectura de un texto sobre pueblos originarios y elaboración de un mapa conceptual.',
 'Se desbloquea el inicio combinando: ubicación estratégica cerca del docente, consigna segmentada en 3 pasos visibles, Time Timer configurado a 8 minutos por bloque. Se entrega pelota antiestrés para regulación durante la lectura. Validación de cada etapa antes de avanzar.',
 'planificada',
 'Mateo inició la tarea dentro de los primeros 2 minutos (antes tardaba 10+). Completó 2 de 3 bloques dentro del tiempo.',
 'La clave fue la consigna segmentada + Time Timer. La pelota la usó sobre todo en el segundo bloque.',
 'completada'),

-- Caso 6: Lucas - Recorte en plástica (motricidad fina) - Tijeras adaptadas
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 2, 9,
 'Educación Artística',
 'Actividad de collage: recortar figuras de revistas y componer una escena.',
 'Se entregan tijeras adaptadas con resorte para reducir el esfuerzo de corte. Se acompaña los primeros cortes para verificar agarre. Se simplifica el diseño permitiendo figuras más grandes.',
 'durante',
 'Lucas completó el collage con 5 figuras recortadas. Mostró menos frustración que en actividades anteriores de recorte.',
 'Modelo con resorte fue el más adecuado. Verificar que haya tijeras para zurdos también.',
 'completada'),

-- Caso 7: Valentina - Trabajo con computadora (acceso visual) - Teclado CLEVY
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 2, 13,
 'Tecnología',
 'Escritura de un texto corto en procesador de texto. Primera experiencia con teclado.',
 'Se conecta el teclado CLEVY con letras grandes y colores diferenciados. No requiere configuración adicional. Se acompaña el inicio explicando las zonas de colores.',
 'planificada',
 'Valentina escribió 3 oraciones en 15 minutos. Los colores la ayudaron a ubicar las letras sin pedir ayuda.',
 'Buen primer contacto con la escritura digital. Continuar en próximas clases.',
 'completada'),

-- Caso 8: Tomás - Se mueve constantemente en Lengua - Elástico para silla
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 2, 1,
 'Prácticas del Lenguaje',
 'Lectura en voz alta grupal. Cada alumno lee un párrafo.',
 'Se coloca el elástico entre las patas delanteras de la silla de Tomás para que pueda empujar con los pies mientras espera su turno. Se le asigna un rol de "seguidor" con el finger focus para mantener la atención en el texto.',
 'durante',
 NULL,
 'Observar si sostiene la atención durante los turnos de los compañeros. Próxima clase evaluar resultado.',
 'en_curso'),

-- Caso 9: Sofía - Mouse adaptado para actividad digital - Mouse trackball
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, 3, 12,
 'Matemática',
 'Actividad interactiva de geometría en la computadora: arrastrar figuras y clasificarlas.',
 'Se conecta el mouse trackball y se ajusta la sensibilidad del cursor antes de entregar. Se posiciona en superficie estable. Se verifica que Sofía alcance el trackball cómodamente.',
 'planificada',
 NULL,
 'Actividad programada para la próxima clase. Preparar la computadora con anticipación.',
 'en_curso'),

-- Caso 10: Mateo - Se desregula emocionalmente - Pelota antiestrés + Auriculares
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, 3, 2,
 'Educación Física (teórica)',
 'Clase teórica sobre reglas de un deporte. Ambiente con mucho ruido por la clase vecina.',
 'Se anticipan señales de desregulación. Se entrega pelota antiestrés desde el inicio. Ante primeros signos de incomodidad, se ofrecen auriculares de cancelación. Consigna simplificada: "anotá solo las 3 reglas más importantes".',
 'durante',
 NULL,
 'Mateo mostró signos de incomodidad a los 10 minutos. Se le ofrecieron los auriculares y se reguló. Monitorear.',
 'en_curso');

-- =============================================================================
-- 3. Adaptation devices (many-to-many) - dispositivos adicionales por adaptación
-- =============================================================================

-- Caso 2: Valentina también usó regla de lectura con ventana (device 7)
INSERT INTO adaptation_devices (adaptation_id, device_id) VALUES
  (2, 6),   -- pen reader (principal)
  (2, 7),   -- regla de lectura con ventana
  (3, 4),   -- auriculares cancelación
  (3, 5),   -- Time Timer
  (4, 11),  -- tablet
  (4, 15),  -- soporte flexible
  (5, 5),   -- Time Timer
  (5, 2),   -- pelota antiestrés
  (8, 1),   -- elástico para silla
  (8, 8),   -- finger focus
  (10, 2),  -- pelota antiestrés
  (10, 4);  -- auriculares cancelación

-- =============================================================================
-- 4. Adaptation resources (fichas pedagógicas y materiales de apoyo)
-- =============================================================================

INSERT INTO adaptation_resources (adaptation_id, title, file_url, file_type)
VALUES
  (1, 'Ficha pedagógica: Pinzas de escritura', '/resources/fichas/pinzas-escritura.pdf', 'pdf'),
  (1, 'Guía de actividad: Dictado adaptado', '/resources/guias/dictado-adaptado.pdf', 'pdf'),
  (2, 'Ficha pedagógica: Pen reader (lápiz lector)', '/resources/fichas/pen-reader.pdf', 'pdf'),
  (2, 'Ficha pedagógica: Regla de lectura con ventana', '/resources/fichas/regla-lectura.pdf', 'pdf'),
  (3, 'Guía de evaluación adaptada: Matemática', '/resources/guias/evaluacion-adaptada-matematica.pdf', 'pdf'),
  (3, 'Ficha pedagógica: Auriculares de cancelación', '/resources/fichas/auriculares-cancelacion.pdf', 'pdf'),
  (4, 'Ficha pedagógica: Tablet educativa', '/resources/fichas/tablet-educativa.pdf', 'pdf'),
  (4, 'Guía de accesibilidad: Configuración de tablet', '/resources/guias/config-tablet-accesibilidad.pdf', 'pdf'),
  (5, 'Ficha pedagógica: Time Timer', '/resources/fichas/time-timer.pdf', 'pdf'),
  (5, 'Estrategia: Desbloquear el inicio de la tarea', '/resources/estrategias/desbloquear-inicio.pdf', 'pdf'),
  (6, 'Ficha pedagógica: Tijeras adaptadas', '/resources/fichas/tijeras-adaptadas.pdf', 'pdf'),
  (7, 'Ficha pedagógica: Teclado CLEVY', '/resources/fichas/teclado-clevy.pdf', 'pdf'),
  (8, 'Ficha pedagógica: Elástico para silla', '/resources/fichas/elastico-silla.pdf', 'pdf'),
  (8, 'Estrategia: Rol activo durante espera en clase', '/resources/estrategias/rol-activo-espera.pdf', 'pdf'),
  (9, 'Ficha pedagógica: Mouse trackball', '/resources/fichas/mouse-trackball.pdf', 'pdf'),
  (9, 'Guía de configuración: Accesibilidad en computadora', '/resources/guias/accesibilidad-computadora.pdf', 'pdf'),
  (10, 'Ficha pedagógica: Pelota antiestrés de gel', '/resources/fichas/pelota-antiestres.pdf', 'pdf'),
  (10, 'Protocolo de desregulación emocional en el aula', '/resources/estrategias/protocolo-desregulacion.pdf', 'pdf');

-- =============================================================================
-- 5. Conversaciones de demo (flujo "Modo Planificador" del PDF)
-- =============================================================================

-- Conversación 1: María planifica el dictado para Lucas (modo planificador)
INSERT INTO conversations (organization_id, user_id, student_id, mode)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 1, 'plan');

INSERT INTO conversation_messages (conversation_id, role, content, metadata) VALUES
(1, 'assistant', '¡Hola María! Soy Alizia, tu asistente de aula. ¿En qué te ayudo hoy?', '{"suggestions": ["Adaptar una actividad para un alumno", "Tengo una situación difícil en el aula", "Crear un material o recurso", "No sé por dónde empezar"]}'),
(1, 'user', 'Quiero planificar una actividad de dictado para mañana. Tengo un alumno con dificultades de motricidad fina.', '{}'),
(1, 'assistant', 'Perfecto, vamos a preparar esa actividad. Contame un poco más:

**Sobre la actividad:**
- ¿Cuánto dura la clase?
- ¿Qué asignatura y eje temático?
- ¿Van a trabajar con algún material específico?', '{"step": "relevamiento_actividad"}'),
(1, 'user', 'Prácticas del Lenguaje. Se va a realizar un dictado de texto y los alumnos van a tener que escribir lo que oyen. La actividad dura 10 minutos.', '{}'),
(1, 'assistant', 'Entendido. Veo que vas a trabajar con **Lucas Martínez**. Su perfil indica dificultades en la motricidad fina que afectan su escritura.

**¿La condición es transitoria o permanente?**
→ Según el perfil: permanente.

**Situación que mejor describe al estudiante:**
→ Tiene dificultad para mover o controlar sus manos.', '{"step": "relevamiento_alumno", "student_id": 1}'),
(1, 'user', 'Sí, es Lucas. ¿Qué me recomendás?', '{}'),
(1, 'assistant', 'Por lo que me contás, Lucas necesita un insumo que lo ayude a sujetar la lapicera de forma adaptada para la actividad de dictado.

**Te recomiendo comenzar con las Pinzas de escritura**, un recurso de la Rampa "Acceso a la lectura, la escritura y la producción".

**¿Por qué?**
- Permite participar de la actividad sin ninguna restricción
- Reduce frustración al escribir
- Da autonomía inmediata
- Se integra fácilmente al trabajo grupal

**Estrategia sugerida:**
1. Colocá el adaptador en la lapicera antes de que empiece la actividad
2. Reducí la extensión del dictado a frases cortas
3. Validá cada párrafo antes de continuar
4. Observá si sostiene la escritura sin fatiga

¿Querés que guarde esta adaptación?', '{"step": "recomendacion", "devices": [10], "ramp": "Acceso a la lectura, la escritura y la producción"}'),
(1, 'user', 'Sí, guardala.', '{}'),
(1, 'assistant', 'Listo, guardé la adaptación para Lucas en Prácticas del Lenguaje. Después de la clase contame cómo le fue para hacer seguimiento.

¿Necesitás algo más?', '{"step": "guardado", "adaptation_id": 1}');

-- Conversación 2: María pide ayuda en tiempo real con Tomás (modo asistencia)
INSERT INTO conversations (organization_id, user_id, student_id, mode)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 3, 'assist');

INSERT INTO conversation_messages (conversation_id, role, content, metadata) VALUES
(2, 'assistant', '¡Hola María! ¿En qué te ayudo?', '{"suggestions": ["Adaptar una actividad para un alumno", "Tengo una situación difícil en el aula", "Crear un material o recurso"]}'),
(2, 'user', 'Tomás no para de moverse y está interrumpiendo la clase de Lengua. Estamos haciendo lectura en voz alta.', '{}'),
(2, 'assistant', 'Entiendo, María. Antes de intervenir, observá brevemente: ¿la conducta aparece porque la tarea es difícil, porque está sin actividad (esperando su turno), o porque hay demasiado estímulo?

Eso nos orienta la respuesta.', '{"step": "observacion"}'),
(2, 'user', 'Creo que es porque está esperando su turno y se aburre.', '{}'),
(2, 'assistant', 'Perfecto, tiene sentido. Te sugiero estas acciones inmediatas:

**1. Elástico para silla** (Rampa: Regulación sensorial y motriz)
Colocalo entre las patas delanteras de la silla. Tomás puede empujar con los pies mientras espera, canalizando la necesidad de movimiento sin levantarse.

**2. Finger focus (señalador de dedo)**
Asignale el rol de "seguidor": que use el finger focus para señalar por dónde va la lectura. Eso le da una tarea activa mientras los compañeros leen.

**3. Turno estratégico**
Dale su turno relativamente pronto para que no tenga que esperar demasiado.

¿Probamos con esto?', '{"step": "asistencia_inmediata", "devices": [1, 8]}'),
(2, 'user', 'Le puse el elástico y el finger focus. Está más tranquilo.', '{}'),
(2, 'assistant', '¡Muy bien! Observá si sostiene la atención durante los turnos de los compañeros. Si funciona, podemos incorporar estas herramientas como parte de su rutina para actividades de espera.

¿Querés que guarde esta adaptación para Tomás?', '{"step": "seguimiento"}');

-- Conversación 3: Carlos planifica actividad digital para Sofía
INSERT INTO conversations (organization_id, user_id, student_id, mode)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 4, 'plan');

INSERT INTO conversation_messages (conversation_id, role, content, metadata) VALUES
(3, 'assistant', '¡Hola Carlos! Soy Alizia. ¿En qué te ayudo hoy?', '{"suggestions": ["Adaptar una actividad para un alumno", "Tengo una situación difícil en el aula"]}'),
(3, 'user', 'Quiero planificar una actividad de investigación digital sobre el ciclo del agua para Sofía. Tiene dificultades de motricidad gruesa.', '{}'),
(3, 'assistant', 'Veo el perfil de **Sofía Álvarez**: dificultades de motricidad gruesa que afectan su postura y manipulación de dispositivos.

Para una actividad digital, te recomiendo:

**Tablet educativa (10")** + **Soporte flexible**
- Configurá las opciones de accesibilidad: texto grande, lector de pantalla
- Colocá la tablet en el soporte flexible para liberar sus manos
- Abrí la actividad antes de entregarla
- Acompañá el inicio y dejá que avance con autonomía

**Enfoque pedagógico:**
La tablet no cambia la propuesta, cambia el acceso. Sofía puede comprender la consigna, sostener la actividad y expresar lo que sabe sin quedar afuera.

**Tips:**
- Revisá la carga antes de la clase
- Usá funda protectora
- Probá la configuración antes de la actividad

¿Guardamos esta adaptación?', '{"step": "recomendacion", "devices": [11, 15]}'),
(3, 'user', 'Sí, perfecto. Guardala para la clase de Ciencias Naturales del jueves.', '{}'),
(3, 'assistant', 'Guardada la adaptación para Sofía en Ciencias Naturales. Recordá probar la tablet y el soporte antes de la clase.

**Qué observar después:**
- Si accede mejor a la consigna
- Si sostiene la actividad por más tiempo
- Si aumenta su autonomía
- Si logra expresar lo que aprendió

¡Éxitos el jueves! Contame cómo le fue.', '{"step": "guardado", "adaptation_id": 4}');

-- Conversación 4: Carlos con Mateo - no inicia la tarea
INSERT INTO conversations (organization_id, user_id, student_id, mode)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 5, 'assist');

INSERT INTO conversation_messages (conversation_id, role, content, metadata) VALUES
(4, 'assistant', 'La última vez estuviste trabajando con Mateo en Ciencias Sociales. ¿Cómo te fue con la segmentación de consignas?', '{"suggestions": ["Necesito ajustar algo de lo que probé", "Quiero planificar otra actividad", "Tengo una situación nueva"]}'),
(4, 'user', 'Funcionó bastante bien pero hoy no arranca con nada. Tiene que leer un texto y hacer un mapa conceptual.', '{}'),
(4, 'assistant', 'Entiendo. Para desbloquear el inicio te sugiero esta combinación:

**1. Ubicación estratégica:** sentalo cerca tuyo para poder validar rápido.

**2. Consigna segmentada visible:**
- Paso 1: "Leé solo el primer párrafo"
- Paso 2: "Escribí 2 palabras clave"
- Paso 3: "Conectalas con una flecha"

**3. Time Timer a 8 minutos** por bloque — así Mateo ve cuánto falta.

**4. Pelota antiestrés** para que regule mientras lee.

**Clave:** validá cada paso antes de dar el siguiente. "¡Bien, ya tenés el paso 1!" es más efectivo que "te falta todo esto".

¿Probamos?', '{"step": "asistencia_inmediata", "devices": [5, 2]}'),
(4, 'user', 'Le puse el Time Timer y le di la consigna en 3 pasos. Arrancó a los 2 minutos, antes tardaba más de 10.', '{}'),
(4, 'assistant', '¡Excelente progreso! De 10+ minutos a 2 minutos para arrancar es un avance enorme. La segmentación visible + Time Timer le está funcionando muy bien a Mateo.

Para las próximas clases podemos ir sumando complejidad gradualmente. ¿Querés que guarde esta estrategia?', '{"step": "seguimiento"}');

-- =============================================================================
-- 6. Actualizar ramps con video_url de ejemplo (placeholders)
-- =============================================================================

UPDATE ramps SET video_url = '/videos/ramps/regulacion-sensorial-motriz.mp4' WHERE id = 1;
UPDATE ramps SET video_url = '/videos/ramps/atencion-foco.mp4' WHERE id = 2;
UPDATE ramps SET video_url = '/videos/ramps/organizacion.mp4' WHERE id = 3;
UPDATE ramps SET video_url = '/videos/ramps/lectura-escritura-produccion.mp4' WHERE id = 4;
UPDATE ramps SET video_url = '/videos/ramps/acceso-tecnologico-adaptado.mp4' WHERE id = 5;

-- =============================================================================
-- 7. image_url de devices
-- =============================================================================
-- Las imágenes de la valija se setean por product_code en
-- db/seeds/valija_chubut_images.sql (fuente única, servidas por el BE embebido
-- bajo /images/devices/<product_code>.<ext>). Se quitaron los placeholders por
-- id de acá para no competir como segunda fuente de verdad.

COMMIT;
