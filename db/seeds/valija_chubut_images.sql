-- valija_chubut_images.sql
-- Puebla devices.image_url para la Valija Chubut (NO schema; la columna image_url
-- existe desde 000002_create_catalog).
--
-- Fuente de verdad: el binario lo sirve el propio BE embebido (package
-- src/app/web/static) bajo /images/devices/<product_code>.png. El nombre del
-- archivo ES el product_code, así que el mapeo device→imagen es determinístico.
-- El FE solo renderiza el image_url que manda el BE.
--
-- Imágenes: las reales del Excel/doc de Chubut (fuente de verdad de Mercedes)
-- para la mayoría; algunas legacy sembradas del FE hasta tener la real.
--
-- Idempotente: UPDATE por product_code. Scopeado por product_code y NO por org
-- a propósito: la imagen de un producto es la misma en cualquier org.
-- Lo corre scripts/dbmigrate (lib/pq); no usa meta-comandos de psql.

-- ============================================================
-- 22/24 devices activos con imagen.
-- ============================================================
UPDATE devices SET image_url = '/images/devices/ETE-S10816-EB.png', updated_at = NOW() WHERE product_code = 'ETE-S10816-EB'; -- Auriculares con micrófono (FE)
UPDATE devices SET image_url = '/images/devices/ETE-I10820-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10820-EB'; -- Auriculares de cancelación (FE)
UPDATE devices SET image_url = '/images/devices/ETE-I10817-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10817-EB'; -- Mouse trackball (AbleNet BIGtrack, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10794-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10794-EB'; -- Mouse Admouse (real)
UPDATE devices SET image_url = '/images/devices/ETE-I10795-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10795-EB'; -- Pulsador botón USB (real)
UPDATE devices SET image_url = '/images/devices/ETE-I10793-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10793-EB'; -- Teclado admouse (FE, revisar)
UPDATE devices SET image_url = '/images/devices/ETE-I10818-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10818-EB'; -- Elástico para silla (FE)
UPDATE devices SET image_url = '/images/devices/ETE-I10819-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10819-EB'; -- Patas de silla x4 (Bouncyband chair feet, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10821-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10821-EB'; -- Pen reader (FE)
UPDATE devices SET image_url = '/images/devices/ETE-I10823-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10823-EB'; -- Organizador de tareas (To-Do list, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10824-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10824-EB'; -- Soporte lápiz et.4 (teardrop, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10825-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10825-EB'; -- Soporte lápiz et.3 (crossover, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10826-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10826-EB'; -- Soporte lápiz et.2 (alado, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10827-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10827-EB'; -- Soporte lápiz et.1 (garra 3 dedos, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10829-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10829-EB'; -- Ayuda lectura tamaño ajustable (overlay, real, revisar)
UPDATE devices SET image_url = '/images/devices/ETE-I10830-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10830-EB'; -- Reglas de lectura guiada (FE, revisar)
UPDATE devices SET image_url = '/images/devices/ETE-I10831-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10831-EB'; -- Reglas lectura transparente (FE, revisar)
UPDATE devices SET image_url = '/images/devices/ETE-I10832-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10832-EB'; -- Material sensorial SPEKS (squishy, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10833-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10833-EB'; -- Sacapuntas para zurdos (Maped, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10835-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10835-EB'; -- Tijera para zurdos (Maped, real)
UPDATE devices SET image_url = '/images/devices/ETE-I10836-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10836-EB'; -- Tijera adaptada et.1 (FE, revisar: foto real es la de mano)
UPDATE devices SET image_url = '/images/devices/ETE-I10837-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10837-EB'; -- Tijera adaptada et.2 (mesa c/ventosa, real)

-- ============================================================
-- PENDIENTES (sin foto en ningún lado; ni FE ni doc de Chubut).
-- Re-exportar de la carpeta IMG-Inclusion (faltaron img_p1_9/10/11) o del Sheet,
-- dropear <product_code>.png en src/app/web/static/images/devices/ y descomentar.
-- ============================================================
-- UPDATE devices SET image_url = '/images/devices/ETE-I10828-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10828-EB'; -- Pesas para lápices
-- UPDATE devices SET image_url = '/images/devices/ETE-I10834-EB.png', updated_at = NOW() WHERE product_code = 'ETE-I10834-EB'; -- Lapicera para zurdos
