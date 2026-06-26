-- valija_chubut_images.sql
-- Puebla devices.image_url para la Valija Chubut (NO schema; la columna image_url
-- existe desde 000002_create_catalog).
--
-- Fuente de verdad: el binario lo sirve el propio BE embebido (package
-- src/app/web/static) bajo /images/devices/<product_code>.<ext>. El nombre del
-- archivo ES el product_code, así que el mapeo device→imagen es determinístico
-- (sin match por nombre). El FE solo renderiza el image_url que manda el BE.
--
-- Idempotente: son UPDATE por product_code. Scopeado por product_code y NO por
-- organización a propósito: la imagen de un producto es la misma en cualquier org
-- (demo o prod), y product_code identifica al producto del catálogo Chubut.
-- Lo corre scripts/dbmigrate (lib/pq), por eso no usa meta-comandos de psql.

-- ============================================================
-- Devices con imagen disponible (sembradas desde el FE actual).
-- Extensión explícita porque conviven .png / .jpg / .jpeg.
-- ============================================================
UPDATE devices SET image_url = '/images/devices/ETE-S10816-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-S10816-EB'; -- Auriculares con micrófono
UPDATE devices SET image_url = '/images/devices/ETE-I10820-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10820-EB'; -- Auriculares de cancelación
UPDATE devices SET image_url = '/images/devices/ETE-I10817-EB.jpg',  updated_at = NOW() WHERE product_code = 'ETE-I10817-EB'; -- Mouse trackball
UPDATE devices SET image_url = '/images/devices/ETE-I10795-EB.jpeg', updated_at = NOW() WHERE product_code = 'ETE-I10795-EB'; -- Pulsador botón USB
UPDATE devices SET image_url = '/images/devices/ETE-I10821-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10821-EB'; -- Pen reader
UPDATE devices SET image_url = '/images/devices/ETE-I10818-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10818-EB'; -- Elástico para silla
UPDATE devices SET image_url = '/images/devices/ETE-I10793-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10793-EB'; -- Teclado admouse
UPDATE devices SET image_url = '/images/devices/ETE-I10794-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10794-EB'; -- Mouse Admouse
UPDATE devices SET image_url = '/images/devices/ETE-I10819-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10819-EB'; -- Patas de silla x4
UPDATE devices SET image_url = '/images/devices/ETE-I10829-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10829-EB'; -- Ayuda lectura tamaño ajustable
UPDATE devices SET image_url = '/images/devices/ETE-I10830-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10830-EB'; -- Reglas de lectura guiada
UPDATE devices SET image_url = '/images/devices/ETE-I10831-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10831-EB'; -- Reglas lectura transparente
UPDATE devices SET image_url = '/images/devices/ETE-I10832-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10832-EB'; -- Material sensorial SPEKS
UPDATE devices SET image_url = '/images/devices/ETE-I10835-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10835-EB'; -- Tijera para zurdos
UPDATE devices SET image_url = '/images/devices/ETE-I10836-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10836-EB'; -- Tijera adaptada etapa 1
UPDATE devices SET image_url = '/images/devices/ETE-I10837-EB.png',  updated_at = NOW() WHERE product_code = 'ETE-I10837-EB'; -- Tijera adaptada etapa 2

-- ============================================================
-- PENDIENTES (sin imagen en el FE; solo existen en el Excel de Chubut).
-- Cuando se exporten del Sheet, dropear el archivo en
-- src/app/web/static/images/devices/<product_code>.<ext> y descomentar acá.
-- ============================================================
-- UPDATE devices SET image_url = '/images/devices/ETE-I10823-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10823-EB'; -- Organizador de tareas
-- UPDATE devices SET image_url = '/images/devices/ETE-I10824-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10824-EB'; -- Soporte para lápiz 1 - etapa 4
-- UPDATE devices SET image_url = '/images/devices/ETE-I10825-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10825-EB'; -- Soporte para lápiz 2 - etapa 3
-- UPDATE devices SET image_url = '/images/devices/ETE-I10826-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10826-EB'; -- Soporte para lápiz 3 - etapa 2
-- UPDATE devices SET image_url = '/images/devices/ETE-I10827-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10827-EB'; -- Soporte para lápiz 4 - etapa 1
-- UPDATE devices SET image_url = '/images/devices/ETE-I10828-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10828-EB'; -- Pesas para lápices
-- UPDATE devices SET image_url = '/images/devices/ETE-I10833-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10833-EB'; -- Sacapuntas para zurdos
-- UPDATE devices SET image_url = '/images/devices/ETE-I10834-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10834-EB'; -- Lapicera para zurdos
