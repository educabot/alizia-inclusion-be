-- 000027_adaptation_ficha.down.sql
DROP INDEX IF EXISTS idx_adaptations_ramp_id;
ALTER TABLE adaptations DROP COLUMN IF EXISTS ramp_id;
ALTER TABLE adaptations DROP COLUMN IF EXISTS steps;
