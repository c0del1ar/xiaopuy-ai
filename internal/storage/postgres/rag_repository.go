package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/c0del1ar/xiaopuy-ai/internal/rag"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RAGRepository struct { pool *pgxpool.Pool }

func NewRAGRepository(pool *pgxpool.Pool) *RAGRepository { return &RAGRepository{pool: pool} }

func (r *RAGRepository) GetDocumentHash(ctx context.Context, documentID string) (string, bool, error) {
	var hash string
	err := r.pool.QueryRow(ctx, `SELECT content_hash FROM rag_documents WHERE id = $1`, documentID).Scan(&hash)
	if err != nil {
		if err.Error() == "no rows in result set" { return "", false, nil }
		return "", false, fmt.Errorf("get RAG document hash %q: %w", documentID, err)
	}
	return hash, true, nil
}

func (r *RAGRepository) Upsert(ctx context.Context, document rag.Document, chunks []rag.Chunk, embeddings [][]float32) error {
	if len(chunks) != len(embeddings) { return fmt.Errorf("chunks and embeddings length mismatch: %d != %d", len(chunks), len(embeddings)) }
	if len(chunks) == 0 { return nil }
	if document.ContentHash == "" { document.ContentHash = rag.ContentHash(document.Content) }

	tx, err := r.pool.Begin(ctx); if err != nil { return fmt.Errorf("begin RAG transaction: %w", err) }
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO rag_documents (id, source, url, title, type, trust, content, content_hash, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (id) DO UPDATE SET source=EXCLUDED.source, url=EXCLUDED.url, title=EXCLUDED.title, type=EXCLUDED.type, trust=EXCLUDED.trust, content=EXCLUDED.content, content_hash=EXCLUDED.content_hash, updated_at=EXCLUDED.updated_at`, document.ID, document.Source, document.URL, document.Title, document.Type, document.Trust, document.Content, document.ContentHash, document.UpdatedAt)
	if err != nil { return fmt.Errorf("upsert RAG document %q: %w", document.ID, err) }
	for i, chunk := range chunks {
		_, err := tx.Exec(ctx, `INSERT INTO rag_chunks (id, document_id, chunk_index, content, metadata, embedding) VALUES ($1,$2,$3,$4,$5::jsonb,$6::vector) ON CONFLICT (id) DO UPDATE SET document_id=EXCLUDED.document_id, chunk_index=EXCLUDED.chunk_index, content=EXCLUDED.content, metadata=EXCLUDED.metadata, embedding=EXCLUDED.embedding`, chunk.ID, chunk.DocumentID, chunk.Index, chunk.Content, metadataJSON(chunk.Metadata), pgVectorLiteral(embeddings[i]))
		if err != nil { return fmt.Errorf("upsert RAG chunk %q: %w", chunk.ID, err) }
	}
	if err := tx.Commit(ctx); err != nil { return fmt.Errorf("commit RAG transaction: %w", err) }
	return nil
}

func (r *RAGRepository) Search(ctx context.Context, query []float32, limit int) ([]rag.Result, error) {
	if len(query) == 0 { return nil, fmt.Errorf("query embedding is empty") }
	if limit <= 0 { limit = 5 }
	rows, err := r.pool.Query(ctx, `SELECT id, document_id, chunk_index, content, metadata::text, 1-(embedding <=> $1::vector) AS score FROM rag_chunks WHERE embedding IS NOT NULL ORDER BY embedding <=> $1::vector LIMIT $2`, pgVectorLiteral(query), limit)
	if err != nil { return nil, fmt.Errorf("search RAG chunks: %w", err) }
	defer rows.Close()
	var results []rag.Result
	for rows.Next() {
		var result rag.Result; var metadata string
		if err := rows.Scan(&result.Chunk.ID, &result.Chunk.DocumentID, &result.Chunk.Index, &result.Chunk.Content, &metadata, &result.Score); err != nil { return nil, fmt.Errorf("scan RAG chunk: %w", err) }
		result.Chunk.Metadata = parseMetadata(metadata); results = append(results, result)
	}
	if err := rows.Err(); err != nil { return nil, fmt.Errorf("iterate RAG search: %w", err) }
	return results, nil
}

func pgVectorLiteral(vector []float32) string { parts:=make([]string,len(vector)); for i,v:=range vector { parts[i]=fmt.Sprintf("%g",v) }; return "["+strings.Join(parts,",")+"]" }
func metadataJSON(metadata map[string]string) string { if len(metadata)==0{return "{}"}; parts:=make([]string,0,len(metadata)); for k,v:=range metadata {parts=append(parts,fmt.Sprintf("%q:%q",k,v))}; return "{"+strings.Join(parts,",")+"}" }
func parseMetadata(value string) map[string]string { result:=make(map[string]string); value=strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(value),"{"),"}"); for _,pair:=range strings.Split(value,","){parts:=strings.SplitN(pair,":",2); if len(parts)!=2{continue}; k:=strings.Trim(parts[0]," \""); v:=strings.Trim(parts[1]," \""); if k!=""{result[k]=v}}; return result }

var _ rag.Repository = (*RAGRepository)(nil)
