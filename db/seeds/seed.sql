-- Demo organization
INSERT INTO organizations (id, name) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Escuela Demo Inclusión')
ON CONFLICT (id) DO NOTHING;

-- Demo users (password: demo123)
INSERT INTO users (organization_id, email, name, password_hash, role) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@demo.edu', 'Admin Demo', '$2a$10$SYOE0umlZfWbGbMVYBROauPfTtaIMF9YTbiAa1KjB/fhWM6gySu8u', 'admin'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'docente1@demo.edu', 'María González', '$2a$10$SYOE0umlZfWbGbMVYBROauPfTtaIMF9YTbiAa1KjB/fhWM6gySu8u', 'teacher'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'docente2@demo.edu', 'Carlos López', '$2a$10$SYOE0umlZfWbGbMVYBROauPfTtaIMF9YTbiAa1KjB/fhWM6gySu8u', 'teacher'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'ministerio@demo.edu', 'Laura Ministerio', '$2a$10$SYOE0umlZfWbGbMVYBROauPfTtaIMF9YTbiAa1KjB/fhWM6gySu8u', 'ministerio'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'psico@demo.edu', 'Ana Psicopedagoga', '$2a$10$SYOE0umlZfWbGbMVYBROauPfTtaIMF9YTbiAa1KjB/fhWM6gySu8u', 'psicopedagogo')
ON CONFLICT DO NOTHING;

-- Demo classrooms
INSERT INTO classrooms (organization_id, name, grade, section) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '3ro A', '3', 'A'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '4to B', '4', 'B')
ON CONFLICT DO NOTHING;

-- Demo students
INSERT INTO students (organization_id, classroom_id, name) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Lucas Martínez'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Valentina Ruiz'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Tomás Fernández'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Sofía Álvarez'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Mateo García')
ON CONFLICT DO NOTHING;

-- Ramps (5 categories matching frontend icon/color mappings)
INSERT INTO ramps (organization_id, name, description, short_description, video_url, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Regulación sensorial y motriz', 'Materiales para regulación sensorial, propiocepción y descarga motriz controlada.', 'Regulación y movimiento', NULL, 1),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Atención y foco', 'Herramientas que reducen distractores y ayudan a sostener la atención en la tarea.', 'Concentración y reducción de estímulos', NULL, 2),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Organización', 'Dispositivos para gestión del tiempo, planificación y organización temporal.', 'Gestión del tiempo y planificación', NULL, 3),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Acceso a la lectura, la escritura y la producción', 'Herramientas que apoyan la lectura, escritura y producción manual.', 'Lectura, escritura y producción', NULL, 4),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Acceso tecnológico adaptado', 'Dispositivos tecnológicos adaptados para facilitar el acceso digital en el aula.', 'Acceso tecnológico adaptado', NULL, 5)
ON CONFLICT DO NOTHING;

-- Devices: Regulación sensorial y motriz (ramp_id = 1)
INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use, recommendations, rationale, classroom_benefit, needs_description, useful_when, evaluation_criteria, quantity, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Elástico para silla', 'Banda elástica entre patas de silla para estimulación sensorial con los pies.', 'DEVICE-ELA-001',
     'Colocar entre patas delanteras de la silla. El estudiante empuja con los pies mientras trabaja.',
     'Usar durante actividades que requieran estar sentado por períodos prolongados. Verificar que el elástico esté bien sujeto. Combinar con pausas activas.',
     'Proporciona canal de descarga motriz controlado que no interrumpe la clase.',
     'Canaliza necesidad de movimiento sin levantarse ni interrumpir.',
     'Estudiantes que necesitan moverse constantemente o presentan inquietud motora.',
     'Durante actividades prolongadas en el banco, cuando el estudiante muestra signos de inquietud o necesita moverse para concentrarse.',
     'Observar si sostiene atención más tiempo y reduce interrupciones.',
     4, 1),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Pelota antiestrés de gel', 'Pelota flexible rellena de gel para regulación sensorial.', 'DEVICE-PEL-001',
     'Entregar cuando se detecte necesidad de regulación. Puede usarse mientras escucha o espera.',
     'Ideal para momentos de espera o transiciones. No usar como juguete. Establecer acuerdo de uso con el estudiante. Tener una de repuesto.',
     'Proporciona estímulo propioceptivo que ayuda a regulación emocional y concentración.',
     'Permite regularse sin interrumpir la dinámica de la clase.',
     'Estudiantes que necesitan regulación sensorial, ansiedad o dificultad de atención.',
     'En momentos de espera, transiciones entre actividades o cuando el estudiante muestra signos de ansiedad o desregulación emocional.',
     'Observar si se muestra más regulado y la concentración mejora.',
     1, 2)
ON CONFLICT DO NOTHING;

-- Devices: Atención y foco (ramp_id = 2)
INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use, recommendations, rationale, classroom_benefit, needs_description, useful_when, evaluation_criteria, quantity, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Auriculares con micrófono', 'Mejoran aislamiento del ruido y permiten acceso por comando de voz.', 'DEVICE-AUR-001',
     'Conectar al dispositivo. Verificar volumen. Activar micrófono para dictado por voz.',
     'Verificar volumen antes de entregar. Limpiar con toallita desinfectante entre usos. Combinar con actividades que requieran concentración individual.',
     'Reducen distractores auditivos y habilitan vías alternativas de producción.',
     'Facilitan concentración y permiten trabajar sin interrumpir al grupo.',
     'Estudiantes que se distraen con ruido ambiental o necesitan producir texto por voz.',
     'Durante evaluaciones, actividades de lectura individual o cuando el estudiante necesita dictar texto en lugar de escribirlo a mano.',
     'Observar si mejora concentración y completa la tarea con menor frustración.',
     3, 1),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Auriculares de cancelación de ruido', 'Auriculares pasivos que reducen ruido ambiental sin reproducir audio.', 'DEVICE-ACR-001',
     'Colocar cuando necesite reducir estímulos auditivos. No requiere conexión.',
     'No forzar el uso. Permitir que el estudiante decida cuándo los necesita. Ideal para evaluaciones o actividades de alta concentración.',
     'Reduce sobrecarga sensorial auditiva sin aislar completamente del grupo.',
     'Permite concentrarse en ambientes ruidosos sin retirarse del aula.',
     'Estudiantes que se sobreestimulan con ruido o necesitan reducir estímulos sensoriales.',
     'En ambientes ruidosos, durante evaluaciones o cuando el estudiante muestra signos de sobrecarga sensorial auditiva.',
     'Observar si muestra mayor calma y mejor rendimiento en tareas de concentración.',
     1, 2)
ON CONFLICT DO NOTHING;

-- Devices: Organización (ramp_id = 3)
INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use, recommendations, rationale, classroom_benefit, needs_description, useful_when, evaluation_criteria, quantity, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 'Time Timer (temporizador visual)', 'Muestra visualmente el paso del tiempo con disco de color que se reduce.', 'DEVICE-TIM-001',
     'Configurar tiempo girando el disco. Posicionar visible para el estudiante.',
     'Usar para segmentar actividades largas en bloques manejables. Combinar con recompensa al completar. No usar como presión sino como apoyo.',
     'Externaliza el paso del tiempo, reduciendo ansiedad y mejorando organización temporal.',
     'Facilita gestión del tiempo sin intervención constante del docente.',
     'Estudiantes que se desorganizan con el tiempo o presentan ansiedad ante plazos.',
     'Al inicio de actividades con tiempo limitado, durante evaluaciones o cuando se necesita segmentar tareas largas en bloques manejables.',
     'Observar si completa tareas dentro del tiempo y con menor ansiedad.',
     1, 1)
ON CONFLICT DO NOTHING;

-- Devices: Acceso a la lectura, la escritura y la producción (ramp_id = 4)
INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use, recommendations, rationale, classroom_benefit, needs_description, useful_when, evaluation_criteria, quantity, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, 'Pen reader (lápiz lector)', 'Escanea texto impreso y lo lee en voz alta.', 'DEVICE-PEN-001',
     'Encender. Pasar sobre el texto. El dispositivo reproduce en audio.',
     'Practicar primero con textos simples. Asegurarse de que el texto esté bien impreso. Combinar con seguimiento visual para reforzar la decodificación.',
     'Permite acceder al contenido escrito sin depender de decodificación lectora.',
     'Facilita comprensión de textos para estudiantes con dificultades de lectura.',
     'Estudiantes con dificultades de lectura, decodificación o comprensión lectora.',
     'Cuando el estudiante necesita acceder a textos impresos de forma autónoma, en actividades de comprensión lectora o consignas escritas.',
     'Observar si comprende mejor el texto y participa con mayor autonomía.',
     1, 1),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, 'Regla de lectura con ventana', 'Aísla una línea de texto para facilitar seguimiento visual.', 'DEVICE-REG-001',
     'Posicionar sobre el texto. Ajustar ventana. Desplazar línea por línea.',
     'Dejar que el estudiante regule la velocidad. Puede combinarse con finger focus. Ideal para textos largos o evaluaciones escritas.',
     'Reduce sobrecarga visual y ayuda a mantener foco en la línea correcta.',
     'Permite seguir el texto sin perderse entre líneas.',
     'Estudiantes con dificultades de seguimiento visual o atención a la lectura.',
     'Durante lecturas de textos largos, evaluaciones escritas o cuando el estudiante salta líneas o pierde el lugar al leer.',
     'Observar si lee con mayor fluidez y menos saltos de línea.',
     2, 2),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, 'Finger focus (señalador de dedo)', 'Se coloca en el dedo para guiar la lectura.', 'DEVICE-FFO-001',
     'Colocar en dedo índice. Señalar la palabra que se lee. Avanzar al ritmo del estudiante.',
     'No presionar para leer más rápido. Respetar el ritmo del estudiante. Puede usarse junto con regla de lectura para mayor apoyo.',
     'Proporciona referencia táctil y visual para seguimiento de lectura.',
     'Ayuda a mantener ritmo de lectura sin perder el lugar.',
     'Estudiantes que se pierden al leer o necesitan apoyo para mantener ritmo.',
     'En actividades de lectura en voz alta, lectura compartida o cuando el estudiante necesita seguir el texto a su propio ritmo.',
     'Observar si la lectura es más fluida y con menos interrupciones.',
     2, 3),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, 'Tijeras adaptadas', 'Mango ergonómico para facilitar corte en actividades manuales.', 'DEVICE-TIJ-001',
     'Seleccionar tipo según necesidad (zurdo, resorte, adaptador). Acompañar primeros usos.',
     'Tener opciones para zurdos y diestros. El modelo con resorte reduce el esfuerzo. Acompañar los primeros usos para verificar agarre.',
     'Reduce demanda motriz del corte.',
     'Permite completar tareas de recorte sin frustración.',
     'Estudiantes con dificultades de motricidad fina o fuerza en las manos.',
     'En actividades plásticas, manualidades o cualquier tarea que requiera recortar papel u otros materiales.',
     'Observar si completa la tarea de corte con mayor autonomía.',
     1, 4),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 4, 'Pinzas de escritura', 'Adaptadores ergonómicos para lápices que facilitan el agarre.', 'DEVICE-PIN-001',
     'Colocar en el lápiz. Ajustar posición de dedos. El adaptador guía el agarre correcto.',
     'Probar distintos modelos hasta encontrar el más cómodo. No forzar el uso si genera rechazo. Combinar con ejercicios de motricidad fina.',
     'Reduce fatiga y mejora control del trazo al escribir.',
     'Permite escribir por más tiempo y con menor esfuerzo.',
     'Estudiantes con dificultades de prensión, agarre o motricidad fina.',
     'Durante actividades de escritura prolongada, copia de textos o cuando el estudiante muestra fatiga o dolor al escribir.',
     'Observar si mejora legibilidad y escribe con menos fatiga.',
     3, 5)
ON CONFLICT DO NOTHING;

-- Devices: Acceso tecnológico adaptado (ramp_id = 5)
INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use, recommendations, rationale, classroom_benefit, needs_description, useful_when, evaluation_criteria, quantity, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, 'Tablet educativa (10")', 'Tablet Android 10" con 8GB RAM y 128GB, pensada para acompañar el trabajo en el aula.', 'DEVICE-TAB-001',
     'Encendé la tablet. Ajustá accesibilidad según el estudiante. Abrí la actividad antes de entregarla.',
     'Configurar accesibilidad antes de la clase. Cargar la batería la noche anterior. Instalar solo apps aprobadas. Usar funda protectora siempre.',
     'La tablet no cambia la propuesta, cambia el acceso. Funciona como apoyo para comprender la consigna y expresar lo que sabe.',
     'Permite incluir sin separar, ofreciendo distintos caminos para el mismo aprendizaje.',
     'Estudiantes que necesitan apoyos para lectura/escritura, información visual/auditiva, alternativas al trabajo manual.',
     'Cuando el estudiante necesita una vía alternativa para acceder a la consigna, producir contenido o participar en actividades digitales.',
     'Observar si accede mejor a la consigna, sostiene la actividad más tiempo, aumenta autonomía.',
     3, 1),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, 'Mouse trackball', 'Mouse adaptado para dificultades de control motor. Compatible con contactores.', 'DEVICE-MTR-001',
     'Conectar al dispositivo. Ajustar sensibilidad. Puede combinarse con contactores.',
     'Ajustar sensibilidad del cursor antes de entregar. Posicionar en superficie estable. Verificar que el estudiante alcance el trackball cómodamente.',
     'Elimina necesidad de desplazar el ratón, reduciendo demanda motriz.',
     'Permite interactuar con la computadora sin depender de motricidad fina.',
     'Estudiantes con dificultades de control motor en extremidades superiores.',
     'En actividades con computadora cuando el mouse convencional resulta difícil de manejar o genera fatiga excesiva.',
     'Observar si navega con mayor autonomía y menor fatiga.',
     1, 2),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, 'Teclado CLEVY', 'Letras grandes y colores diferenciados por grupos de teclas.', 'DEVICE-TCL-001',
     'Conectar al dispositivo. No requiere configuración adicional.',
     'Ideal como primer teclado para estudiantes que están aprendiendo a escribir en computadora. Los colores ayudan a ubicar las letras por zona.',
     'Facilita reconocimiento de letras y reduce errores de tipeo.',
     'Permite escribir con mayor autonomía y confianza.',
     'Estudiantes con dificultades motoras, cognitivas o de reconocimiento de letras.',
     'Cuando el estudiante está aprendiendo a tipear, tiene dificultad para ubicar letras o necesita teclas más grandes por motricidad reducida.',
     'Observar si velocidad y precisión de escritura mejoran.',
     1, 3),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, 'Pulsador botón USB', 'Alternativa al teclado/ratón para severas dificultades de control motor.', 'DEVICE-PUL-001',
     'Conectar. Configurar función (clic, enter, espacio). Posicionar al alcance.',
     'Probar distintas posiciones hasta encontrar la más accesible. Puede fijarse a la mesa con cinta. Configurar la función según la actividad.',
     'Permite interacción mínima pero efectiva con dispositivos electrónicos.',
     'Habilita participación de estudiantes que no pueden usar teclado ni ratón convencional.',
     'Estudiantes con severas dificultades de control motor.',
     'Cuando el estudiante no puede usar teclado ni mouse convencional y necesita una forma simplificada de interactuar con el dispositivo.',
     'Observar si logra interactuar con la actividad de forma autónoma.',
     1, 4),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 5, 'Soporte flexible celular/tablet', 'Soporte articulado para posicionar dispositivos al alcance del estudiante.', 'DEVICE-SOP-001',
     'Ajustar altura y ángulo. Fijar dispositivo. Verificar estabilidad.',
     'Verificar estabilidad antes de colocar el dispositivo. Ajustar el ángulo según la postura del estudiante. Ideal para uso prolongado de tablet.',
     'Libera las manos y permite posicionar de forma ergonómica.',
     'Facilita acceso al dispositivo sin depender de compañero o adulto.',
     'Estudiantes con dificultades de manipulación o postura.',
     'Cuando el estudiante usa tablet o celular por períodos prolongados y necesita liberar las manos o mantener una postura cómoda.',
     'Observar si accede al dispositivo con mayor comodidad.',
     1, 5)
ON CONFLICT DO NOTHING;

-- Student inclusion profiles (demo data with valid FE difficulty values)
INSERT INTO student_profiles (student_id, is_transitory, difficulties, free_description) VALUES
    (1, false, ARRAY['motricidad_fina', 'expresion'], 'Lucas presenta dificultades en la motricidad fina que afectan su escritura y producción manual. También muestra dificultades para expresarse verbalmente en clase.'),
    (2, false, ARRAY['acceso_visual', 'bloqueo_lectoescritura'], 'Valentina tiene dificultad de acceso visual y presenta bloqueos en lectoescritura. Requiere dispositivos adaptados para acceder a los contenidos.'),
    (3, true, ARRAY['distraccion_constante', 'sobreestimulacion'], 'Tomás presenta distracción constante y se sobreestimula con facilidad ante ruidos y estímulos del aula. Condición transitoria en seguimiento.')
ON CONFLICT (student_id) DO NOTHING;
