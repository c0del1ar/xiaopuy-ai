package rag

import (
	"context"
	"testing"

	"github.com/c0del1ar/xiaopuy-ai/internal/embedding"
)

type countingEmbeddingProvider struct {
	calls int
}

func (p *countingEmbeddingProvider) Embed(_ context.Context, text string) ([]float32, error) {
	p.calls++
	return []float32{float32(len(text))}, nil
}

var _ embedding.Provider = (*countingEmbeddingProvider)(nil)

type hashRepository struct {
	hash       string
	found      bool
	upsertCall int
}

func (r *hashRepository) GetDocumentHash(_ context.Context, _ string) (string, bool, error) {
	return r.hash, r.found, nil
}

func (r *hashRepository) Upsert(_ context.Context, _ Document, _ []Chunk, _ [][]float32) error {
	r.upsertCall++
	return nil
}

func (r *hashRepository) Search(_ context.Context, _ []float32, _ int) ([]Result, error) {
	return nil, nil
}

var _ Repository = (*hashRepository)(nil)

func TestIndexSkipsUnchangedDocument(t *testing.T) {
	content := "A stable document with enough content to index."
	repo := &hashRepository{hash: ContentHash(content), found: true}
	embeddings := &countingEmbeddingProvider{}
	service := &Service{Embeddings: embeddings, Repository: repo}

	err := service.Index(context.Background(), Document{ID: "doc-1", Content: content}, 20, 0)
	if err != nil { t.Fatalf("Index() error = %v", err) }
	if embeddings.calls != 0 { t.Fatalf("embedding calls = %d, want 0", embeddings.calls) }
	if repo.upsertCall != 0 { t.Fatalf("upsert calls = %d, want 0", repo.upsertCall) }
}

func TestIndexEmbedsChangedDocument(t *testing.T) {
	content := "A changed document with enough content to index."
	repo := &hashRepository{hash: ContentHash("old content"), found: true}
	embeddings := &countingEmbeddingProvider{}
	service := &Service{Embeddings: embeddings, Repository: repo}

	err := service.Index(context.Background(), Document{ID: "doc-1", Content: content}, 20, 0)
	if err != nil { t.Fatalf("Index() error = %v", err) }
	if embeddings.calls == 0 { t.Fatal("embedding provider was not called for changed document") }
	if repo.upsertCall != 1 { t.Fatalf("upsert calls = %d, want 1", repo.upsertCall) }
}
