-- 000024_create_rag_corpus.up.sql
-- Versiona el corpus RAG (rag_resources / rag_chunks) que YA EXISTE en producción y
-- alimenta la búsqueda híbrida (tool search_content_hibrido + endpoint /search-content/hybrid).
-- Idempotente (IF NOT EXISTS): no-op en prod, crea el esquema en dev/local.
-- El esquema CANÓNICO vive en prod; este archivo lo documenta para version-control.
-- Si la introspección de prod difiere (tipos/índices), reconciliar este DDL.

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS rag_resources (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rag_chunks (
    id          BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL REFERENCES rag_resources(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL DEFAULT 0,
    page_start  INTEGER NOT NULL DEFAULT 0,
    page_end    INTEGER NOT NULL DEFAULT 0,
    content     TEXT NOT NULL DEFAULT '',
    summary     TEXT,
    concepts    TEXT[] NOT NULL DEFAULT '{}',
    embedding   vector(1536), -- dim = config.EmbeddingDim (text-embedding-3-small)
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rag_chunks_resource_id ON rag_chunks (resource_id);

-- Índice ANN para el operador de distancia coseno (<=>) usado por vector_hits.
CREATE INDEX IF NOT EXISTS idx_rag_chunks_embedding
    ON rag_chunks USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
