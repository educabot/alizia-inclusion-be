-- valija_chubut_images.sql
-- Puebla devices.image_url para la Valija Chubut (NO schema; la columna image_url
-- existe desde 000002_create_catalog).
--
-- Fuente de verdad: el binario lo sirve el propio BE embebido (package
-- src/app/web/static) bajo /images/devices/<product_code>.jpg. El nombre del
-- archivo ES el product_code, así que el mapeo device→imagen es determinístico.
-- El FE solo renderiza el image_url que manda el BE.
--
-- Imágenes: las reales del Excel/doc de Chubut (carpeta IMG-Inclusion, fuente de
-- verdad de Mercedes), validadas una por una contra el Sheet. Cobertura 24/24.
--
-- Idempotente: UPDATE por product_code. Scopeado por product_code y NO por org
-- a propósito: la imagen de un producto es la misma en cualquier org.
-- Lo corre scripts/dbmigrate (lib/pq); no usa meta-comandos de psql.

UPDATE devices SET image_url = '/images/devices/ETE-I10793-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10793-EB'; -- Teclado admouse con covertor acrílico
UPDATE devices SET image_url = '/images/devices/ETE-S10816-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-S10816-EB'; -- Auriculares con micrófono
UPDATE devices SET image_url = '/images/devices/ETE-I10817-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10817-EB'; -- Mouse trackball (AbleNet BIGtrack)
UPDATE devices SET image_url = '/images/devices/ETE-I10794-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10794-EB'; -- Mouse Admouse
UPDATE devices SET image_url = '/images/devices/ETE-I10795-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10795-EB'; -- Pulsador botón USB
UPDATE devices SET image_url = '/images/devices/ETE-I10818-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10818-EB'; -- Banda elástica Bouncyband (Elástico para silla)
UPDATE devices SET image_url = '/images/devices/ETE-I10819-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10819-EB'; -- Patas de silla x4 (Bouncyband wobble feet)
UPDATE devices SET image_url = '/images/devices/ETE-I10820-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10820-EB'; -- Auriculares de cancelación auditivo
UPDATE devices SET image_url = '/images/devices/ETE-I10821-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10821-EB'; -- Pen reader
UPDATE devices SET image_url = '/images/devices/ETE-I10824-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10824-EB'; -- Soporte para lápiz - etapa 4 (teardrop)
UPDATE devices SET image_url = '/images/devices/ETE-I10825-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10825-EB'; -- Soporte para lápiz - etapa 3 (crossover)
UPDATE devices SET image_url = '/images/devices/ETE-I10826-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10826-EB'; -- Soporte para lápiz - etapa 2 (alado)
UPDATE devices SET image_url = '/images/devices/ETE-I10827-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10827-EB'; -- Soporte para lápiz - etapa 1 (garra 3 dedos)
UPDATE devices SET image_url = '/images/devices/ETE-I10828-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10828-EB'; -- Pesas para lápices
UPDATE devices SET image_url = '/images/devices/ETE-I10829-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10829-EB'; -- Ayuda para la lectura - tamaño ajustable
UPDATE devices SET image_url = '/images/devices/ETE-I10830-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10830-EB'; -- Ayuda para la lectura - reglas guiada
UPDATE devices SET image_url = '/images/devices/ETE-I10831-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10831-EB'; -- Ayuda para la lectura - regla transparente con renglón
UPDATE devices SET image_url = '/images/devices/ETE-I10832-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10832-EB'; -- Material sensorial de apriete SPEKS
UPDATE devices SET image_url = '/images/devices/ETE-I10823-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10823-EB'; -- Organizador de tareas personalizable
UPDATE devices SET image_url = '/images/devices/ETE-I10833-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10833-EB'; -- Sacapuntas para zurdos (Maped)
UPDATE devices SET image_url = '/images/devices/ETE-I10834-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10834-EB'; -- Lapicera para zurdos (ergonómica)
UPDATE devices SET image_url = '/images/devices/ETE-I10835-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10835-EB'; -- Tijera para zurdos (Maped)
UPDATE devices SET image_url = '/images/devices/ETE-I10836-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10836-EB'; -- Tijera adaptada - etapa 1 (loop de mano)
UPDATE devices SET image_url = '/images/devices/ETE-I10837-EB.jpg', updated_at = NOW() WHERE product_code = 'ETE-I10837-EB'; -- Tijera adaptada - etapa 2 (mesa c/ventosa)
