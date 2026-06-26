-- align_categorias.sql
-- Normaliza los nombres de ramps a las 5 categorías canónicas del producto
-- (Flujo del docente, Epic ALZ-246). En la UI se muestran como "necesidad";
-- internamente son ramps. NO es schema: solo renombra (rename), no rompe FKs
-- de devices.ramp_id. Idempotente: re-correr matchea 0 filas.
--
-- Las 5 categorías: Atención y foco · Organización · Lectoescritura ·
-- Regulación sensorial y motriz · Tecnología.
--
-- Estado actual (prod): "Regulación sensorial y motriz", "Atención y foco" y
-- "Organización" ya coinciden. Solo cambian estas dos:
UPDATE ramps SET name = 'Lectoescritura', updated_at = NOW()
WHERE name = 'Acceso a la lectura, la escritura y la producción';

UPDATE ramps SET name = 'Tecnología', updated_at = NOW()
WHERE name = 'Acceso tecnológico adaptado';
