package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/c0del1ar/xiaopuy-ai/internal/rag"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRAGRepositorySemanticSearch(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" { t.Skip("TEST_DATABASE_URL is not set") }
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second); defer cancel()
	pool, err := pgxpool.New(ctx, dsn); if err != nil { t.Fatalf("create pool: %v", err) }; defer pool.Close()
	if err := pool.Ping(ctx); err != nil { t.Fatalf("ping database: %v", err) }
	repo := NewRAGRepository(pool)
	prefix := "integration-vector-"
	doc := rag.Document{ID: prefix+"doc", Source: "test", URL: "https://example.test", Title: "Integration test", Type: "test", Trust: "high", Content: "RAG integration document", UpdatedAt: time.Now().UTC()}
	chunks := []rag.Chunk{
		{ID: prefix+"a", DocumentID: prefix+"doc", Index: 0, Content: "Laravel website development services", Metadata: map[string]string{"source":"test"}},
		{ID: prefix+"b", DocumentID: prefix+"doc", Index: 1, Content: "Logo and visual identity design", Metadata: map[string]string{"source":"test"}},
		{ID: prefix+"c", DocumentID: prefix+"doc", Index: 2, Content: "Linux server administration and hosting", Metadata: map[string]string{"source":"test"}},
	}
	embeddings := [][]float32{{1,0,0,0},{0,1,0,0},{0,0,1,0}}
	_, err = pool.Exec(ctx, `DELETE FROM rag_documents WHERE id = $1`, doc.ID); if err != nil { t.Fatalf("cleanup: %v", err) }
	if err := repo.Upsert(ctx, doc, chunks, embeddings); err != nil { t.Fatalf("upsert: %v", err) }
	defer pool.Exec(context.Background(), `DELETE FROM rag_documents WHERE id = $1`, doc.ID)
	results, err := repo.Search(ctx, []float32{0.98,0.02,0,0}, 2); if err != nil { t.Fatalf("search: %v", err) }
	if len(results) != 2 { t.Fatalf("result count = %d, want 2", len(results)) }
	if results[0].Chunk.ID != prefix+"a" { t.Fatalf("top result = %q, want %q", results[0].Chunk.ID, prefix+"a") }
	if results[0].Score <= results[1].Score { t.Fatalf("scores are not ordered: first=%v second=%v", results[0].Score, results[1].Score) }
}
