-- Context Engine — Capa E: RAG de contenido pedagógico (keyword-first + embeddings fallback).
-- pgvector ya está disponible en el servidor (v0.8.2); acá se habilita la extensión.

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS pedagogical_content (
    id              BIGSERIAL PRIMARY KEY,
    parent_id       BIGINT REFERENCES pedagogical_content(id) ON DELETE SET NULL,  -- jerarquía tipo Notion
    type            VARCHAR(50),                       -- libro / paper / material / capitulo
    title           VARCHAR(512),
    status          VARCHAR(50) NOT NULL DEFAULT 'draft',
    keywords        TEXT[],                            -- temas + discapacidades (ej. "TEA")
    organization_id UUID REFERENCES organizations(id), -- NULL = global (Educabot)
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pedagogical_content_parent ON pedagogical_content(parent_id);
CREATE INDEX IF NOT EXISTS idx_pedagogical_content_keywords ON pedagogical_content USING GIN (keywords);

CREATE TABLE IF NOT EXISTS pedagogical_content_chunks (
    id         BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES pedagogical_content(id) ON DELETE CASCADE,
    chunk_text TEXT,
    tags       TEXT[],
    embedding  VECTOR,   -- dimensión sin fijar hasta confirmar el modelo de embeddings (Azure).
                         -- El índice vectorial (hnsw/ivfflat) requiere dim fija → migración aparte.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pedagogical_content_chunks_content ON pedagogical_content_chunks(content_id);
CREATE INDEX IF NOT EXISTS idx_pedagogical_content_chunks_tags ON pedagogical_content_chunks USING GIN (tags);
