-- 000023_align_devices_valija_chubut.up.sql
-- Alinea la tabla devices con el catálogo de la Valija Chubut (SOLO schema, sin datos).
-- Agrega metadatos de producto (product_code, product_family, stage) y un flag de
-- activación (is_active), e impone unicidad de product_code por organización.
-- Idempotente: usa IF NOT EXISTS para poder re-correrse sin error.
-- Los seeds/backfill de datos viven aparte en db/seeds/valija_chubut_align.sql.

-- Código de producto ETE de la Valija (ej. ETE-I10827-EB); nullable para devices legacy.
ALTER TABLE devices ADD COLUMN IF NOT EXISTS product_code VARCHAR(50);

-- Familia/agrupación funcional del producto (ej. soporte_lapiz, ayuda_lectura, tijera).
ALTER TABLE devices ADD COLUMN IF NOT EXISTS product_family VARCHAR(80);

-- Etapa progresiva dentro de una familia (1..N); NULL para productos sin progresión.
ALTER TABLE devices ADD COLUMN IF NOT EXISTS stage SMALLINT;

-- Flag de catálogo activo: permite ocultar devices descontinuados sin borrarlos.
ALTER TABLE devices ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- Unicidad de product_code por organización (parcial: ignora filas sin product_code).
CREATE UNIQUE INDEX IF NOT EXISTS uq_devices_org_product_code
    ON devices (organization_id, product_code)
    WHERE product_code IS NOT NULL;
