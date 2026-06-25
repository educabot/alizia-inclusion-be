-- valija_chubut_align.sql
-- Datos para alinear el catálogo de devices con la Valija Chubut (NO schema).
-- Requiere la migración 000023 (columnas product_code, product_family, stage, is_active).
-- Idempotente y multi-statement; lo corre scripts/dbmigrate con lib/pq, por eso NO usa
-- meta-comandos de psql (\set, \gset): el org y los ramp_id se resuelven por subquery.
-- Org objetivo: 'Escuela Demo Inclusión'. Los ramp_id se resuelven por (name, org).
--
-- (a) Backfill de product_code a 6 devices existentes (match por name exacto).
-- (b) Desactivación (is_active=FALSE) de 9 devices discontinuados.
-- (c) Upsert idempotente de los 18 devices nuevos (ON CONFLICT por org+product_code).

-- ============================================================
-- (a) BACKFILL: asigna product_code y reactiva 6 devices existentes
-- ============================================================
UPDATE devices SET product_code = 'ETE-S10816-EB', is_active = TRUE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Auriculares con micrófono';

UPDATE devices SET product_code = 'ETE-I10820-EB', is_active = TRUE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Auriculares de cancelación de ruido';

UPDATE devices SET product_code = 'ETE-I10817-EB', is_active = TRUE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Mouse trackball';

UPDATE devices SET product_code = 'ETE-I10795-EB', is_active = TRUE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Pulsador botón USB';

UPDATE devices SET product_code = 'ETE-I10821-EB', is_active = TRUE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Pen reader (lápiz lector)';

UPDATE devices SET product_code = 'ETE-I10818-EB', is_active = TRUE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Elástico para silla';

-- ============================================================
-- (b) DESACTIVAR: marca is_active=FALSE en 9 devices discontinuados
-- ============================================================
UPDATE devices SET is_active = FALSE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Time Timer (temporizador visual)';

UPDATE devices SET is_active = FALSE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Tablet educativa (10")';

UPDATE devices SET is_active = FALSE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Finger focus (señalador de dedo)';

UPDATE devices SET is_active = FALSE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Teclado CLEVY';

UPDATE devices SET is_active = FALSE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Soporte flexible celular/tablet';

UPDATE devices SET is_active = FALSE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Pelota antiestrés de gel';

UPDATE devices SET is_active = FALSE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Regla de lectura con ventana';

UPDATE devices SET is_active = FALSE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Pinzas de escritura';

UPDATE devices SET is_active = FALSE, updated_at = NOW()
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')
  AND name = 'Tijeras adaptadas';

-- ============================================================
-- (c) UPSERT: 18 devices nuevos de la Valija Chubut
--     qr_code = product_code; ramp_id resuelto por (ramp_name, org)
-- ============================================================
INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Soporte para lápiz 4 - etapa 1',
    'Adaptador de agarre inicial que aloja los tres dedos y guía la prensión correcta del lápiz.',
    'ETE-I10827-EB',
    'Colocá el soporte en el lápiz. Ubicá cada dedo en su hueco marcado. Acompañá los primeros trazos verificando que la posición se mantenga.',
    'Usalo en la etapa más temprana del trazo, cuando todavía no hay prensión estable. Acompañá cada uso al principio. Combinalo con ejercicios de motricidad fina. Retiralo recién cuando el agarre se sostenga solo.',
    'Al fijar la posición de cada dedo, reduce la carga motriz y deja que el estudiante registre por repetición cómo se sostiene un lápiz.',
    'Permite iniciar la escritura sin frustrarse por no poder sostener el lápiz.',
    'Estudiantes que recién inician el trazo o tienen prensión muy inestable y necesitan máxima guía.',
    'Observar si logra sostener el lápiz en la posición guiada y completa trazos simples sin abandonar.',
    'Cuando el estudiante toma el lápiz con el puño, cambia de dedos constantemente o no logra sostenerlo para empezar a escribir.',
    5,
    'ETE-I10827-EB',
    'soporte_lapiz',
    1,
    TRUE,
    100
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Soporte para lápiz 3 - etapa 2',
    'Adaptador de agarre que marca la posición de los dedos pero da algo más de libertad de movimiento que el inicial.',
    'ETE-I10826-EB',
    'Colocá el soporte en el lápiz. Apoyá los dedos en las zonas marcadas. Dejá que el estudiante escriba acompañando solo si pierde la posición.',
    'Usalo cuando ya hay una prensión incipiente pero todavía irregular. Reducí el acompañamiento de a poco. Si la posición se sostiene bien, probá avanzar a la etapa siguiente.',
    'Ofrece menos guía que el inicial: el estudiante corrige por sí mismo dentro de un marco que sigue presente, consolidando el patrón motor.',
    'Permite afianzar el agarre correcto mientras gana algo de autonomía en el trazo.',
    'Estudiantes que ya sostienen el lápiz pero con prensión irregular y necesitan apoyo para consolidarla.',
    'Observar si mantiene la prensión correcta durante más tiempo y con menos correcciones del docente.',
    'Cuando el estudiante escribe pero pierde la posición de los dedos a mitad de la tarea o se cansa al sostener el trazo.',
    5,
    'ETE-I10826-EB',
    'soporte_lapiz',
    2,
    TRUE,
    101
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Soporte para lápiz 2 - etapa 3',
    'Adaptador de agarre liviano que sugiere la posición de los dedos sin fijarla por completo.',
    'ETE-I10825-EB',
    'Colocá el soporte en el lápiz. Dejá que el estudiante acomode los dedos con la referencia del adaptador. Intervení solo si lo pide.',
    'Usalo cuando la prensión ya es estable y se busca afinar el control del trazo. Es un paso intermedio antes de escribir sin soporte. Observá si todavía lo necesita o puede prescindir de él.',
    'Brinda una referencia táctil sutil que recuerda la posición correcta sin imponerla, favoreciendo el control fino y la autorregulación.',
    'Permite trabajar la calidad del trazo y la escritura sostenida con mínima dependencia del apoyo.',
    'Estudiantes con prensión ya estable que necesitan refinar el control y la resistencia del trazo.',
    'Observar si mejora la legibilidad y sostiene la escritura sin volver a posiciones incorrectas.',
    'Cuando el estudiante escribe con autonomía pero su trazo pierde precisión en tareas largas o de copia.',
    5,
    'ETE-I10825-EB',
    'soporte_lapiz',
    3,
    TRUE,
    102
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Soporte para lápiz 1 - etapa 4',
    'Adaptador de agarre mínimo y discreto que aporta confort sin condicionar la posición de los dedos.',
    'ETE-I10824-EB',
    'Colocá el soporte en el lápiz. Dejá que el estudiante escriba con total autonomía. Usalo como apoyo de confort, no de guía.',
    'Usalo en la etapa más avanzada, cuando el agarre ya es autónomo y solo se busca reducir fatiga en escritura prolongada. Es el último paso antes de escribir sin ningún soporte.',
    'Aporta amortiguación y confort sin guiar la prensión, sosteniendo la autonomía ya lograda y reduciendo la fatiga en tareas extensas.',
    'Permite escribir por más tiempo con menor esfuerzo, sin depender de una guía de agarre.',
    'Estudiantes con agarre autónomo y correcto que solo necesitan reducir fatiga en escritura prolongada.',
    'Observar si escribe durante más tiempo sin fatiga ni molestias, manteniendo la calidad del trazo.',
    'Cuando el estudiante escribe bien de forma independiente pero muestra cansancio o molestia en jornadas largas de escritura.',
    5,
    'ETE-I10824-EB',
    'soporte_lapiz',
    4,
    TRUE,
    103
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Ayuda para la lectura - tamaño ajustable - etapa 1',
    'Apoyo visual de ventana ajustable que aísla y agranda un fragmento de texto para iniciar el seguimiento de la lectura.',
    'ETE-I10829-EB',
    'Apoyá la guía sobre el texto. Ajustá el tamaño de la ventana hasta dejar visible una sola línea o palabra. Acompañá al estudiante desplazándola con él, renglón por renglón.',
    'Empezá con textos cortos y de letra grande. Ajustá la ventana al fragmento más pequeño que el estudiante tolere sin perderse. Acompañá el desplazamiento en los primeros usos y dejá que después lo intente solo. Usá en lecturas guiadas, no en evaluaciones todavía.',
    'Reduce la sobrecarga visual al limitar el campo de lectura, lo que permite anclar la atención en un punto concreto del texto y construir el hábito de seguimiento.',
    'Permite acompañar la lectura inicial sin leer por el estudiante, sosteniendo el foco en una porción manejable del texto.',
    'Estudiantes que recién inician el seguimiento lector, que se pierden entre líneas o se abruman ante la cantidad de texto en la página.',
    'Observar si logra mantener el foco en la línea señalada y empieza a seguir el texto con menos pérdidas, al principio acompañado y luego con ayuda mínima.',
    'En lecturas guiadas con apoyo del docente, cuando el estudiante pierde el lugar constantemente o se bloquea ante una página llena de texto.',
    3,
    'ETE-I10829-EB',
    'ayuda_lectura',
    1,
    TRUE,
    104
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Ayuda para la lectura - Reglas de lectura guiada - etapa 3',
    'Regla de lectura guiada que marca el renglón en curso para sostener el seguimiento con menor acompañamiento del docente.',
    'ETE-I10830-EB',
    'Colocá la regla bajo el renglón que se está leyendo. Pedile al estudiante que la desplace hacia abajo a medida que avanza. Intervení solo si pierde el lugar.',
    'Apta para textos de mayor extensión que en etapas iniciales. Dejá que el estudiante regule su propia velocidad de desplazamiento. Combinala con lecturas en voz alta o compartidas. Reducí progresivamente tu acompañamiento para favorecer la autonomía.',
    'Marcar el renglón en curso libera recursos de atención que antes se usaban para no perder el lugar, lo que mejora la fluidez y la continuidad de la lectura.',
    'Habilita lecturas más largas con seguimiento propio del estudiante, liberando al docente de guiar línea por línea.',
    'Estudiantes que ya iniciaron el seguimiento lector y necesitan un apoyo intermedio para sostenerlo en textos más extensos con creciente autonomía.',
    'Observar si desplaza la regla solo, lee con mayor fluidez y comete menos saltos de línea sin necesidad de que el docente lo guíe.',
    'Durante lecturas de textos largos o lectura compartida, cuando el estudiante todavía salta líneas pero puede manejar la regla por su cuenta.',
    6,
    'ETE-I10830-EB',
    'ayuda_lectura',
    3,
    TRUE,
    105
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Ayuda para la lectura - Reglas de lectura transparente con renglón - etapa 4',
    'Regla transparente con renglón resaltado que guía la lectura sin tapar el texto, para un uso autónomo y discreto.',
    'ETE-I10831-EB',
    'Apoyá la regla sobre el texto alineando la línea resaltada con el renglón a leer. Deslizala hacia abajo al avanzar. El estudiante la usa por su cuenta, sin intervención.',
    'Ideal para uso autónomo en textos largos y en evaluaciones escritas, ya que la transparencia deja ver el contexto y resulta discreta. Dejala disponible para que el estudiante la tome cuando la necesite. Sirve también para revisar la propia producción escrita.',
    'El renglón resaltado sobre material transparente mantiene el foco sin aislar la línea del resto del texto, lo que sostiene el seguimiento a la vez que preserva la comprensión del contexto.',
    'Permite que el estudiante lea de forma autónoma y discreta en cualquier actividad, incluidas las evaluaciones, sin distinguirse del grupo.',
    'Estudiantes con seguimiento lector consolidado que necesitan un apoyo de baja asistencia para sostener la lectura de forma independiente.',
    'Observar si toma y usa la regla por iniciativa propia, sostiene la lectura de textos largos sin perder el lugar y mantiene la comprensión del contexto.',
    'En lectura individual, estudio y evaluaciones escritas, cuando el estudiante puede gestionar su lectura solo pero aún se beneficia de una guía visual discreta.',
    10,
    'ETE-I10831-EB',
    'ayuda_lectura',
    4,
    TRUE,
    106
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Tijera para zurdos',
    'Tijera con hojas y mango invertidos, pensada para cortar con la mano izquierda.',
    'ETE-I10835-EB',
    'Entregá la tijera al estudiante zurdo. Verificá que tome con la mano izquierda. Acompañá el primer corte para chequear que la línea quede a la vista.',
    'Reservala para estudiantes zurdos: en una tijera común las hojas tapan la línea de corte. Guardala identificada para no mezclarla con las de diestros. Acompañá los primeros usos para verificar el agarre.',
    'La inversión de las hojas deja la línea de corte visible para la mano izquierda, evitando que el estudiante tuerza la muñeca o el papel para compensar.',
    'Permite que el estudiante zurdo recorte con la misma precisión y autonomía que el resto, sin adaptaciones improvisadas.',
    'Estudiantes zurdos que cortan con dificultad o tuercen la mano al usar tijeras de diestros.',
    'Observar si recorta sobre la línea con mayor precisión y sin torcer la muñeca ni el papel.',
    'En actividades plásticas, manualidades o tareas de recorte donde el estudiante zurdo lucha con una tijera común o pide ayuda para cortar derecho.',
    1,
    'ETE-I10835-EB',
    'tijera',
    NULL,
    TRUE,
    107
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Tijera adaptada - etapa 1',
    'Tijera con resorte que vuelve a abrir sola, para iniciarse en el corte con mínimo esfuerzo.',
    'ETE-I10836-EB',
    'Colocá la tijera en la mano del estudiante. Indicá que solo apriete para cortar: el resorte la abre sola. Acompañá guiando la mano en los primeros cortes.',
    'Usala como primer paso: el estudiante solo aprende el movimiento de cerrar, no a reabrir. Si cuesta sostenerla, acompañá la mano por encima. Empezá con cortes cortos en papel firme antes de seguir líneas.',
    'El resorte automatiza la apertura, así el estudiante concentra toda la energía en un único gesto de cierre y reduce la demanda motriz y la frustración del aprendizaje inicial.',
    'Permite que estudiantes que aún no dominan el corte participen de las actividades de recorte desde el primer día.',
    'Estudiantes con dificultades marcadas de motricidad fina o poca fuerza en las manos que se inician en el corte.',
    'Observar si logra cerrar la tijera para cortar de forma sostenida y mantiene el interés en la actividad.',
    'En las primeras experiencias de recorte, o cuando el estudiante no logra abrir y cerrar una tijera común y abandona la tarea.',
    1,
    'ETE-I10836-EB',
    'tijera_adaptada',
    1,
    TRUE,
    108
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Tijera adaptada - etapa 2',
    'Tijera con lazo y resorte suave que acompaña el corte dando más control al estudiante.',
    'ETE-I10837-EB',
    'Pedí al estudiante que ubique los dedos dentro del lazo. Dejá que abra y cierre por su cuenta. Intervení solo si pierde la línea de corte.',
    'Usala cuando el estudiante ya domina el gesto de cierre de la etapa 1. El lazo da apoyo pero exige más control que el resorte automático. Proponé cortar siguiendo líneas rectas y luego curvas para ganar precisión.',
    'Al exigir que el estudiante participe en la apertura y guíe el corte, transfiere progresivamente el control desde el dispositivo hacia la propia motricidad, consolidando la autonomía.',
    'Permite avanzar hacia un corte más preciso y autónomo, acercando al estudiante al uso de una tijera convencional.',
    'Estudiantes con dificultades de motricidad fina que ya cortan con apoyo y están listos para mayor control y autonomía.',
    'Observar si recorta siguiendo la línea con creciente precisión y menor necesidad de acompañamiento.',
    'Cuando el estudiante recorta con la tijera de etapa 1 sin dificultad y necesita un desafío mayor para seguir líneas o cortar con más precisión.',
    1,
    'ETE-I10837-EB',
    'tijera_adaptada',
    2,
    TRUE,
    109
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Sacapuntas para zurdos',
    'Sacapuntas con giro invertido pensado para la mano izquierda, para sacar punta sin torcer la muñeca.',
    'ETE-I10833-EB',
    'Entregalo al estudiante zurdo. Sostené el lápiz con la mano izquierda y girá el sacapuntas en el sentido que indica el modelo. Vaciá el depósito cuando se llene.',
    'Diferenciá el modelo para zurdos del de diestros para que no se mezclen. Mostrá el sentido de giro la primera vez. Tené a mano para actividades de escritura y dibujo.',
    'El giro invertido respeta el gesto natural de la mano izquierda, evitando la torsión de muñeca que genera el sacapuntas común.',
    'Permite que el estudiante zurdo prepare su material de forma autónoma, sin pedir ayuda ni forzar la postura.',
    'Estudiantes zurdos que tienen dificultad o incomodidad al usar el sacapuntas convencional.',
    'Observar si saca punta sin torcer la muñeca y prepara su material de forma autónoma.',
    'Cuando el estudiante zurdo tuerce la muñeca, se frustra o pide ayuda para sacar punta al lápiz.',
    1,
    'ETE-I10833-EB',
    'zurdos',
    NULL,
    TRUE,
    110
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Lapicera para zurdos',
    'Lapicera con punta y apoyo de dedos diseñados para la mano izquierda, para escribir sin emborronar.',
    'ETE-I10834-EB',
    'Entregala al estudiante zurdo. Ubicá los dedos en la zona de agarre marcada. Escribí de izquierda a derecha apoyando la mano sin arrastrar la tinta.',
    'Verificá que la tinta seque rápido para evitar manchas. Acompañá los primeros usos para corregir la posición de la mano. Reservala para el estudiante zurdo, no la uses como lapicera común.',
    'La punta y el apoyo adaptados reducen el emborronamiento y la tensión que la mano izquierda sufre con lapiceras pensadas para diestros.',
    'Permite escribir con prolijidad y menos fatiga, sin manchar el cuaderno ni la mano.',
    'Estudiantes zurdos con dificultades de prensión, manchado o fatiga al escribir.',
    'Observar si escribe con mayor prolijidad, menos manchas y menor fatiga.',
    'Durante actividades de escritura prolongada, cuando el estudiante zurdo mancha el texto o muestra incomodidad al escribir.',
    1,
    'ETE-I10834-EB',
    'zurdos',
    NULL,
    TRUE,
    111
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso a la lectura, la escritura y la producción' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Pesas para lápices',
    'Pesas que se colocan en el lápiz para dar feedback propioceptivo y mejorar el control del trazo.',
    'ETE-I10828-EB',
    'Colocá la pesa en el cuerpo del lápiz cerca de la punta. El estudiante escribe sintiendo el peso extra. Retirala si el trazo se vuelve cansador.',
    'Probá distintos pesos hasta encontrar el más cómodo. Usala en sesiones cortas al principio para evitar fatiga. Combinala con pinzas de escritura si el estudiante también necesita apoyo de agarre.',
    'El peso adicional aumenta la información propioceptiva que recibe la mano, ayudando a regular la presión y estabilizar el trazo.',
    'Permite que el estudiante gane control y firmeza en la escritura sin intervención constante del docente.',
    'Estudiantes con trazo inestable, presión irregular o baja conciencia del movimiento de la mano al escribir.',
    'Observar si el trazo se vuelve más firme y regular, y si controla mejor la presión al escribir.',
    'Durante actividades de escritura o dibujo, cuando el estudiante presenta trazo tembloroso, presión despareja o poca firmeza.',
    5,
    'ETE-I10828-EB',
    'pesas_lapiz',
    NULL,
    TRUE,
    112
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso tecnológico adaptado' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Teclado admouse con covertor acrilico',
    'Teclado adaptado con cobertor acrílico que guía la digitación y evita pulsar teclas no deseadas.',
    'ETE-I10793-EB',
    'Conectá el teclado al dispositivo. Colocá el cobertor acrílico sobre las teclas. Verificá que cada orificio quede alineado con su tecla. Acompañá los primeros usos.',
    'Asegurate de que el cobertor esté bien encastrado antes de entregarlo. Limpiá la superficie acrílica con paño seco. Ideal para estudiantes que pulsan varias teclas a la vez por movimientos involuntarios. Acompañá los primeros usos hasta que ubique las teclas con confianza.',
    'El cobertor obliga a apoyar el dedo dentro de cada orificio, reduciendo pulsaciones accidentales y dando una referencia táctil estable para cada tecla.',
    'Permite escribir en computadora con mayor precisión y menos errores, sin depender de ayuda constante del adulto.',
    'Estudiantes con motricidad reducida, temblor o movimientos involuntarios que pulsan teclas no deseadas al escribir.',
    'Observar si comete menos errores de pulsación, escribe con mayor precisión y sostiene la tarea con más autonomía.',
    'Cuando el estudiante usa el teclado pero apoya la mano sobre varias teclas, repite letras o se frustra por errores de tipeo involuntarios.',
    1,
    'ETE-I10793-EB',
    NULL,
    NULL,
    TRUE,
    113
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Acceso tecnológico adaptado' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Mouse Admouse',
    'Mouse adaptado de la línea AdMouse para acceso con motricidad reducida.',
    'ETE-I10794-EB',
    'Conectá el mouse al dispositivo. Ajustá la sensibilidad del cursor. Posicionalo sobre una superficie estable al alcance del estudiante. Verificá que alcance los botones con comodidad.',
    'Ajustá la velocidad del cursor antes de entregarlo según el control motor del estudiante. Ubicalo sobre una superficie firme y a la altura adecuada. Puede combinarse con pulsadores externos. Acompañá los primeros usos para encontrar la posición más cómoda.',
    'Su diseño adaptado reduce la demanda motriz fina del clic y el desplazamiento, ofreciendo un control más estable para quienes tienen poca fuerza o precisión en la mano.',
    'Permite navegar e interactuar con la computadora de forma autónoma, sin depender de motricidad fina precisa.',
    'Estudiantes con motricidad reducida o poca fuerza y precisión en las extremidades superiores que tienen dificultad con el mouse convencional.',
    'Observar si navega y selecciona con mayor autonomía, menor fatiga y más precisión en el clic.',
    'En actividades con computadora cuando el mouse común resulta difícil de manejar, genera fatiga o el estudiante no logra hacer clic con precisión.',
    1,
    'ETE-I10794-EB',
    NULL,
    NULL,
    TRUE,
    114
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Regulación sensorial y motriz' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Patas de silla x4',
    'Juego de cuatro topes para las patas de la silla que estabilizan el asiento o permiten un balanceo controlado.',
    'ETE-I10819-EB',
    'Calzá un tope en cada pata de la silla. Elegí la posición fija para estabilizar el asiento o la posición con balanceo para habilitar un movimiento leve. Verificá que la silla quede firme antes de que el estudiante se siente.',
    'Revisá periódicamente que los cuatro topes estén bien calzados. Usá la opción de balanceo solo cuando el movimiento ayude a la regulación, no como juego. Combiná con pausas activas en actividades largas.',
    'El balanceo controlado ofrece un canal de descarga motriz e input propioceptivo desde el asiento, mientras que la opción estable evita el ruido y el vaivén que distraen al estudiante y al grupo.',
    'Permite que el estudiante regule su necesidad de movimiento sentado en su lugar, sin levantarse ni desordenar la silla.',
    'Estudiantes que necesitan moverse para concentrarse o que se inquietan al estar sentados por períodos prolongados.',
    'Observar si permanece sentado con mayor comodidad, reduce el balanceo descontrolado y sostiene la atención por más tiempo.',
    'Durante actividades prolongadas en el banco, cuando el estudiante balancea la silla, mueve las piernas o muestra inquietud motora.',
    1,
    'ETE-I10819-EB',
    NULL,
    NULL,
    TRUE,
    115
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Regulación sensorial y motriz' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Material sensorial de apriete SPEKS',
    'Set de imanes y bolitas para apretar, separar y armar formas, que ofrece estímulo sensorial y trabajo de motricidad fina.',
    'ETE-I10832-EB',
    'Entregá el set cuando detectes necesidad de regulación. El estudiante aprieta, separa y vuelve a unir las piezas con las manos mientras escucha o espera. Establecé de antemano cuándo y cómo se usa.',
    'Acordá con el estudiante que es un apoyo de regulación, no un juguete. Supervisá su uso por las piezas pequeñas. Guardalo entre usos y tené el segundo set como repuesto o para otro estudiante.',
    'La manipulación repetida y el apriete brindan estímulo táctil y propioceptivo que canaliza la inquietud, a la vez que ejercitan la motricidad fina y la coordinación de los dedos.',
    'Habilita una vía discreta para regularse y mantener las manos ocupadas sin interrumpir la dinámica de la clase.',
    'Estudiantes que necesitan regulación sensorial, descarga táctil o que se benefician de ejercitar la motricidad fina.',
    'Observar si se muestra más regulado, sostiene la atención y usa el material sin que se vuelva una distracción.',
    'En momentos de espera, transiciones o cuando el estudiante muestra signos de ansiedad, inquietud o necesidad de manipular algo con las manos.',
    2,
    'ETE-I10832-EB',
    'sensorial',
    NULL,
    TRUE,
    116
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();

INSERT INTO devices (organization_id, ramp_id, name, description, qr_code, how_to_use,
    recommendations, rationale, classroom_benefit, needs_description, evaluation_criteria,
    useful_when, quantity, product_code, product_family, stage, is_active, sort_order)
VALUES (
    (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión'),
    (SELECT id FROM ramps WHERE name = 'Organización' AND organization_id = (SELECT id FROM organizations WHERE name = 'Escuela Demo Inclusión')),
    'Organizador de tareas personalizable',
    'Panel con checklist deslizante para secuenciar las tareas del día y anticipar qué viene después.',
    'ETE-I10823-EB',
    'Armá la secuencia de tareas en el panel al inicio de la actividad. El estudiante desliza o marca cada paso a medida que lo completa y mira cuál sigue. Revisá el panel juntos al cerrar la jornada o la actividad.',
    'Empezá con pocas tareas e incorporá más a medida que el estudiante gana autonomía. Usá imágenes o palabras según su nivel de lectura. Ubicalo visible y al alcance. No lo uses como control sino como apoyo a la organización.',
    'Externalizar la secuencia de tareas reduce la carga sobre la memoria de trabajo y la ansiedad ante lo que viene, y el acto de marcar lo completado refuerza la sensación de avance y la función ejecutiva de planificación.',
    'Permite que el estudiante sepa qué tiene que hacer y en qué orden, avanzando con menor intervención constante del docente.',
    'Estudiantes que se desorganizan con varias consignas, pierden el hilo de la secuencia o se angustian ante lo imprevisto.',
    'Observar si sigue la secuencia de tareas con mayor autonomía, completa más pasos y pregunta menos qué tiene que hacer.',
    'Al inicio de actividades con varios pasos, en jornadas con varias materias o cuando el estudiante se bloquea sin saber por dónde seguir.',
    1,
    'ETE-I10823-EB',
    NULL,
    NULL,
    TRUE,
    117
)
ON CONFLICT (organization_id, product_code) WHERE product_code IS NOT NULL
DO UPDATE SET name = EXCLUDED.name, ramp_id = EXCLUDED.ramp_id, description = EXCLUDED.description,
    how_to_use = EXCLUDED.how_to_use, recommendations = EXCLUDED.recommendations, rationale = EXCLUDED.rationale,
    classroom_benefit = EXCLUDED.classroom_benefit, needs_description = EXCLUDED.needs_description,
    evaluation_criteria = EXCLUDED.evaluation_criteria, useful_when = EXCLUDED.useful_when,
    quantity = EXCLUDED.quantity, product_family = EXCLUDED.product_family, stage = EXCLUDED.stage,
    is_active = TRUE, updated_at = NOW();
