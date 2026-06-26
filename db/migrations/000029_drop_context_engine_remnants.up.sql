-- 000029_drop_context_engine_remnants.up.sql
-- Limpieza del "context engine" inconcluso. Dos tablas quedaron sin uso real:
--   - response_examples (000019): banco de few-shot, nunca cableado en código Go.
--   - integradora_assignments (000022): provider instanciado pero no inyectado en
--     ningún usecase (wiring muerto).
-- NO toca las columnas source_conversation_id/source_message_id/was_edited de
-- adaptations (también de 000019) que SÍ están en uso (trazabilidad de origen IA).
DROP TABLE IF EXISTS response_examples;
DROP TABLE IF EXISTS integradora_assignments;
