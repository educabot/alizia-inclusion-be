-- 000032_create_rag_corpus.down.sql
-- Revierte 000032. La extensión vector se deja (puede usarla otro esquema).
DROP TABLE IF EXISTS rag_chunks;
DROP TABLE IF EXISTS rag_resources;
