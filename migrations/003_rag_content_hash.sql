ALTER TABLE rag_documents
    ADD COLUMN IF NOT EXISTS content_hash TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS rag_documents_content_hash_idx
    ON rag_documents(content_hash);
