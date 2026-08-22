package rag

import (
	"context"
	"fmt"

	"github.com/c0del1ar/xiaopuy-ai/internal/embedding"
)

type Service struct {
	Embeddings embedding.Provider
	Repository Repository
}

func (s *Service) Index(ctx context.Context, document Document, maxChars, overlap int) error {
	if s == nil || s.Embeddings == nil || s.Repository == nil {
		return fmt.Errorf("RAG service is not configured")
	}
	chunks := ChunkText(document.Content, maxChars, overlap)
	if len(chunks) == 0 {
		return fmt.Errorf("document has no content")
	}

	embeddings := make([][]float32, 0, len(chunks))
	items := make([]Chunk, 0, len(chunks))
	for i, content := range chunks {
		vector, err := s.Embeddings.Embed(ctx, content)
		if err != nil {
			return fmt.Errorf("embed chunk %d: %w", i, err)
		}
		embeddings = append(embeddings, vector)
		items = append(items, Chunk{
			ID:         fmt.Sprintf("%s:%d", document.ID, i),
			DocumentID: document.ID,
			Index:      i,
			Content:    content,
			Metadata: map[string]string{
				"source": document.Source,
				"url":    document.URL,
				"title":  document.Title,
				"type":   document.Type,
				"trust":  document.Trust,
			},
		})
	}

	return s.Repository.Upsert(ctx, items, embeddings)
}

func (s *Service) Retrieve(ctx context.Context, query string, limit int) ([]Chunk, error) {
	if s == nil || s.Embeddings == nil || s.Repository == nil {
		return nil, fmt.Errorf("RAG service is not configured")
	}
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}
	if limit <= 0 {
		limit = 5
	}
	vector, err := s.Embeddings.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}
	return s.Repository.Search(ctx, vector, limit)
}
