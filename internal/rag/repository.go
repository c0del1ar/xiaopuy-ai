package rag

import "context"

// Repository stores indexed documents/chunks and retrieves the most relevant chunks.
type Repository interface {
	Upsert(ctx context.Context, document Document, chunks []Chunk, embeddings [][]float32) error
	Search(ctx context.Context, query []float32, limit int) ([]Result, error)
}
