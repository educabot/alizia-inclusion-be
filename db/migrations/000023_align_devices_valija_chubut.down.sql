-- 000023_align_devices_valija_chubut.down.sql
-- Revierte 000023: elimina el índice único y las 4 columnas agregadas a devices.
-- Idempotente: usa IF EXISTS para no fallar si ya fueron removidas.

DROP INDEX IF EXISTS uq_devices_org_product_code;

ALTER TABLE devices DROP COLUMN IF EXISTS is_active;
ALTER TABLE devices DROP COLUMN IF EXISTS stage;
ALTER TABLE devices DROP COLUMN IF EXISTS product_family;
ALTER TABLE devices DROP COLUMN IF EXISTS product_code;
