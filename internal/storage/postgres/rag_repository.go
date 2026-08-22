package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/c0del1ar/xiaopuy-ai/internal/rag"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RAGRepository struct {
	pool *pgxpool.Pool
}

func NewRAGRepository(pool *pgxpool.Pool) *RAGRepository {
	return &RAGRepository{pool: pool}
}

func (r *RAGRepository) Upsert(ctx context.Context, chunks []rag.Chunk, embeddings [][]float32) error {
	if len(chunks) != len(embeddings) {
		return fmt.Errorf("chunks and embeddings length mismatch: %d != %d", len(chunks), len(embeddings))
	}
	if len(chunks) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin RAG transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for i, chunk := range chunks {
		vector := pgVectorLiteral(embeddings[i])
		_, err := tx.Exec(ctx, `
			INSERT INTO rag_chunks (id, document_id, chunk_index, content, metadata, embedding)
			VALUES ($1, $2, $3, $4, $5::jsonb, $6::vector)
			ON CONFLICT (id) DO UPDATE SET
				document_id = EXCLUDED.document_id,
				chunk_index = EXCLUDED.chunk_index,
				content = EXCLUDED.content,
				metadata = EXCLUDED.metadata,
				embedding = EXCLUDED.embedding`,
			chunk.ID,
			chunk.DocumentID,
			chunk.Index,
			chunk.Content,
			metadataJSON(chunk.Metadata),
			vector,
		)
		if err != nil {
			return fmt.Errorf("upsert RAG chunk %q: %w", chunk.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit RAG transaction: %w", err)
	}
	return nil
}

func (r *RAGRepository) Search(ctx context.Context, query []float32, limit int) ([]rag.Chunk, error) {
	if len(query) == 0 {
		return nil, fmt.Errorf("query embedding is empty")
	}
	if limit <= 0 {
		limit = 5
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, document_id, chunk_index, content, metadata::text
		FROM rag_chunks
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> $1::vector
		LIMIT $2`, pgVectorLiteral(query), limit)
	if err != nil {
		return nil, fmt.Errorf("search RAG chunks: %w", err)
	}
	defer rows.Close()

	var results []rag.Chunk
	for rows.Next() {
		var chunk rag.Chunk
		var metadata string
		if err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.Index, &chunk.Content, &metadata); err != nil {
			return nil, fmt.Errorf("scan RAG chunk: %w", err)
		}
		chunk.Metadata = parseMetadata(metadata)
		results = append(results, chunk)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate RAG search: %w", err)
	}
	return results, nil
}

func pgVectorLiteral(vector []float32) string {
	parts := make([]string, len(vector))
	for i, value := range vector {
		parts[i] = fmt.Sprintf("%g", value)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func metadataJSON(metadata map[string]string) string {
	if len(metadata) == 0 {
		return "{}"
	}
	parts := make([]string, 0, len(metadata))
	for key, value := range metadata {
		parts = append(parts, fmt.Sprintf("%q:%q", key, value))
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func parseMetadata(value string) map[string]string {
	// Metadata is written by this package and is intentionally kept flat.
	result := make(map[string]string)
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "{")
	value = strings.TrimSuffix(value, "}")
	for _, pair := range strings.Split(value, ",") {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.Trim(parts[0], " \"")
		val := strings.Trim(parts[1], " \"")
		if key != "" {
			result[key] = val
		}
	}
	return result
}

var _ rag.Repository = (*RAGRepository)(nil)
var _ = pgxpool.Pool{}
