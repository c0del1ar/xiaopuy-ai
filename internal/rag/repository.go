package rag

import "context"

// Repository stores indexed chunks and retrieves the most relevant chunks.
type Repository interface {
	Upsert(ctx context.Context, chunks []Chunk, embeddings [][]float32) error
	Search(ctx context.Context, query []float32, limit int) ([]Chunk, error)
}
