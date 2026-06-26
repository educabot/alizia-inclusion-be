-- 000027_adaptation_ficha.up.sql
-- Ficha del recurso (Epic ALZ-246): paso a paso estructurado + categoría explícita.
-- steps: array jsonb de {orden, texto, checkbox} — la guía imprimible del recurso.
-- ramp_id: categoría/"necesidad" explícita (FK lógica a ramps); permite categorizar
--   recursos sin material físico (estrategia_aula) que no tienen device del cual inferirla.
-- Aditivo e idempotente: no quita columnas (notes/outcome quedan, deprecadas).
ALTER TABLE adaptations ADD COLUMN IF NOT EXISTS steps jsonb NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE adaptations ADD COLUMN IF NOT EXISTS ramp_id bigint;
CREATE INDEX IF NOT EXISTS idx_adaptations_ramp_id ON adaptations (ramp_id);
