-- Demo organization
INSERT INTO organizations (id, name) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Escuela Demo Inclusión')
ON CONFLICT (id) DO NOTHING;

-- Demo users (password: demo123)
INSERT INTO users (organization_id, email, name, password_hash, role) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@demo.edu', 'Admin Demo', '$argon2id$v=19$m=65536,t=3,p=4$c2FsdHNhbHRzYWx0$hash', 'admin'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'docente1@demo.edu', 'María González', '$argon2id$v=19$m=65536,t=3,p=4$c2FsdHNhbHRzYWx0$hash', 'teacher'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'docente2@demo.edu', 'Carlos López', '$argon2id$v=19$m=65536,t=3,p=4$c2FsdHNhbHRzYWx0$hash', 'teacher')
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

-- Ramps (3 categories of adaptive devices)
INSERT INTO ramps (organization_id, name, description, short_description, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Rampa Digital', 'Dispositivos tecnológicos adaptados para facilitar el acceso digital en el aula.', 'Acceso tecnológico adaptado', 1),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Rampa Didáctico-Pedagógica', 'Herramientas no tecnológicas que apoyan la lectura, escritura y producción.', 'Lectura, escritura y producción', 2),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Rampa de Autorregulación Sensorial', 'Materiales para regulación sensorial, atención y organización temporal.', 'Regulación sensorial y atención', 3)
ON CONFLICT DO NOTHING;

-- Devices: Rampa Digital (6 devices)
INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use, rationale, classroom_benefit, needs_description, evaluation_criteria, quantity, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Tablet educativa (10")', 'Tablet Android 10" con 8GB RAM y 128GB, pensada para acompañar el trabajo en el aula.', 'DEVICE-TAB-001', 'Encendé la tablet. Ajustá accesibilidad según el estudiante. Abrí la actividad antes de entregarla.', 'La tablet no cambia la propuesta, cambia el acceso. Funciona como apoyo para comprender la consigna y expresar lo que sabe.', 'Permite incluir sin separar, ofreciendo distintos caminos para el mismo aprendizaje.', 'Estudiantes que necesitan apoyos para lectura/escritura, información visual/auditiva, alternativas al trabajo manual.', 'Observar si accede mejor a la consigna, sostiene la actividad más tiempo, aumenta autonomía.', 3, 1),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Auriculares con micrófono', 'Mejoran aislamiento del ruido y permiten acceso por comando de voz.', 'DEVICE-AUR-001', 'Conectar al dispositivo. Verificar volumen. Activar micrófono para dictado por voz.', 'Reducen distractores auditivos y habilitan vías alternativas de producción.', 'Facilitan concentración y permiten trabajar sin interrumpir al grupo.', 'Estudiantes que se distraen con ruido ambiental o necesitan producir texto por voz.', 'Observar si mejora concentración y completa la tarea con menor frustración.', 3, 2),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Mouse trackball', 'Mouse adaptado para dificultades de control motor. Compatible con contactores.', 'DEVICE-MTR-001', 'Conectar al dispositivo. Ajustar sensibilidad. Puede combinarse con contactores.', 'Elimina necesidad de desplazar el ratón, reduciendo demanda motriz.', 'Permite interactuar con la computadora sin depender de motricidad fina.', 'Estudiantes con dificultades de control motor en extremidades superiores.', 'Observar si navega con mayor autonomía y menor fatiga.', 1, 3),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Teclado CLEVY', 'Letras grandes y colores diferenciados por grupos de teclas.', 'DEVICE-TCL-001', 'Conectar al dispositivo. No requiere configuración adicional.', 'Facilita reconocimiento de letras y reduce errores de tipeo.', 'Permite escribir con mayor autonomía y confianza.', 'Estudiantes con dificultades motoras, cognitivas o de reconocimiento de letras.', 'Observar si velocidad y precisión de escritura mejoran.', 1, 4),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Pulsador botón USB', 'Alternativa al teclado/ratón para severas dificultades de control motor.', 'DEVICE-PUL-001', 'Conectar. Configurar función (clic, enter, espacio). Posicionar al alcance.', 'Permite interacción mínima pero efectiva con dispositivos electrónicos.', 'Habilita participación de estudiantes que no pueden usar teclado ni ratón convencional.', 'Estudiantes con severas dificultades de control motor.', 'Observar si logra interactuar con la actividad de forma autónoma.', 1, 5),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1, 'Soporte flexible celular/tablet', 'Soporte articulado para posicionar dispositivos al alcance del estudiante.', 'DEVICE-SOP-001', 'Ajustar altura y ángulo. Fijar dispositivo. Verificar estabilidad.', 'Libera las manos y permite posicionar de forma ergonómica.', 'Facilita acceso al dispositivo sin depender de compañero o adulto.', 'Estudiantes con dificultades de manipulación o postura.', 'Observar si accede al dispositivo con mayor comodidad.', 1, 6)
ON CONFLICT DO NOTHING;

-- Devices: Rampa Didáctico-Pedagógica (5 devices)
INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use, rationale, classroom_benefit, needs_description, evaluation_criteria, quantity, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Pen reader (lápiz lector)', 'Escanea texto impreso y lo lee en voz alta.', 'DEVICE-PEN-001', 'Encender. Pasar sobre el texto. El dispositivo reproduce en audio.', 'Permite acceder al contenido escrito sin depender de decodificación lectora.', 'Facilita comprensión de textos para estudiantes con dificultades de lectura.', 'Estudiantes con dificultades de lectura, decodificación o comprensión lectora.', 'Observar si comprende mejor el texto y participa con mayor autonomía.', 1, 1),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Regla de lectura con ventana', 'Aísla una línea de texto para facilitar seguimiento visual.', 'DEVICE-REG-001', 'Posicionar sobre el texto. Ajustar ventana. Desplazar línea por línea.', 'Reduce sobrecarga visual y ayuda a mantener foco en la línea correcta.', 'Permite seguir el texto sin perderse entre líneas.', 'Estudiantes con dificultades de seguimiento visual o atención a la lectura.', 'Observar si lee con mayor fluidez y menos saltos de línea.', 2, 2),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Finger focus (señalador de dedo)', 'Se coloca en el dedo para guiar la lectura.', 'DEVICE-FFO-001', 'Colocar en dedo índice. Señalar la palabra que se lee. Avanzar al ritmo del estudiante.', 'Proporciona referencia táctil y visual para seguimiento de lectura.', 'Ayuda a mantener ritmo de lectura sin perder el lugar.', 'Estudiantes que se pierden al leer o necesitan apoyo para mantener ritmo.', 'Observar si la lectura es más fluida y con menos interrupciones.', 2, 3),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Tijeras adaptadas', 'Mango ergonómico para facilitar corte en actividades manuales.', 'DEVICE-TIJ-001', 'Seleccionar tipo según necesidad (zurdo, resorte, adaptador). Acompañar primeros usos.', 'Reduce demanda motriz del corte.', 'Permite completar tareas de recorte sin frustración.', 'Estudiantes con dificultades de motricidad fina o fuerza en las manos.', 'Observar si completa la tarea de corte con mayor autonomía.', 1, 4),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2, 'Pinzas de escritura', 'Adaptadores ergonómicos para lápices que facilitan el agarre.', 'DEVICE-PIN-001', 'Colocar en el lápiz. Ajustar posición de dedos. El adaptador guía el agarre correcto.', 'Reduce fatiga y mejora control del trazo al escribir.', 'Permite escribir por más tiempo y con menor esfuerzo.', 'Estudiantes con dificultades de prensión, agarre o motricidad fina.', 'Observar si mejora legibilidad y escribe con menos fatiga.', 3, 5)
ON CONFLICT DO NOTHING;

-- Devices: Rampa de Autorregulación Sensorial (4 devices)
INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use, rationale, classroom_benefit, needs_description, evaluation_criteria, quantity, sort_order) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 'Elástico para silla', 'Banda elástica entre patas de silla para estimulación sensorial con los pies.', 'DEVICE-ELA-001', 'Colocar entre patas delanteras de la silla. El estudiante empuja con los pies mientras trabaja.', 'Proporciona canal de descarga motriz controlado que no interrumpe la clase.', 'Canaliza necesidad de movimiento sin levantarse ni interrumpir.', 'Estudiantes que necesitan moverse constantemente o presentan inquietud motora.', 'Observar si sostiene atención más tiempo y reduce interrupciones.', 4, 1),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 'Time Timer (temporizador visual)', 'Muestra visualmente el paso del tiempo con disco de color que se reduce.', 'DEVICE-TIM-001', 'Configurar tiempo girando el disco. Posicionar visible para el estudiante.', 'Externaliza el paso del tiempo, reduciendo ansiedad y mejorando organización temporal.', 'Facilita gestión del tiempo sin intervención constante del docente.', 'Estudiantes que se desorganizan con el tiempo o presentan ansiedad ante plazos.', 'Observar si completa tareas dentro del tiempo y con menor ansiedad.', 1, 2),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 'Pelota antiestrés de gel', 'Pelota flexible rellena de gel para regulación sensorial.', 'DEVICE-PEL-001', 'Entregar cuando se detecte necesidad de regulación. Puede usarse mientras escucha o espera.', 'Proporciona estímulo propioceptivo que ayuda a regulación emocional y concentración.', 'Permite regularse sin interrumpir la dinámica de la clase.', 'Estudiantes que necesitan regulación sensorial, ansiedad o dificultad de atención.', 'Observar si se muestra más regulado y la concentración mejora.', 1, 3),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 'Auriculares de cancelación de ruido', 'Auriculares pasivos que reducen ruido ambiental sin reproducir audio.', 'DEVICE-ACR-001', 'Colocar cuando necesite reducir estímulos auditivos. No requiere conexión.', 'Reduce sobrecarga sensorial auditiva sin aislar completamente del grupo.', 'Permite concentrarse en ambientes ruidosos sin retirarse del aula.', 'Estudiantes que se sobreestimulan con ruido o necesitan reducir estímulos sensoriales.', 'Observar si muestra mayor calma y mejor rendimiento en tareas de concentración.', 1, 4)
ON CONFLICT DO NOTHING;

-- Student inclusion profiles (demo data)
INSERT INTO student_profiles (student_id, is_transitory, difficulties, free_description) VALUES
    (1, false, ARRAY['motricidad_fina', 'debilidad_muscular_manos'], 'Lucas presenta dificultades en la motricidad fina con debilidades musculares en la mano. Tiene retrasos madurativos documentados.'),
    (2, false, ARRAY['acceso_tecnologia_digital'], 'Valentina tiene dificultad para acceder a la tecnología digital estándar y requiere dispositivos adaptados.'),
    (3, true, ARRAY['atencion', 'regulacion_emocional'], 'Tomás presenta dificultad para mantener la atención y regular emociones. Condición transitoria en seguimiento.')
ON CONFLICT (student_id) DO NOTHING;
