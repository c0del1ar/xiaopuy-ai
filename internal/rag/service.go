package rag

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/c0del1ar/xiaopuy-ai/internal/embedding"
)

type Service struct {
	Embeddings embedding.Provider
	Repository Repository
	Dimension  int
}

func (s *Service) Index(ctx context.Context, document Document, maxChars, overlap int) error {
	if s == nil || s.Embeddings == nil || s.Repository == nil { return fmt.Errorf("RAG service is not configured") }
	if document.ID == "" { return fmt.Errorf("document ID cannot be empty") }
	if document.Content == "" { return fmt.Errorf("document has no content") }
	if document.ContentHash == "" { document.ContentHash = ContentHash(document.Content) }

	currentHash, found, err := s.Repository.GetDocumentHash(ctx, document.ID)
	if err != nil { return fmt.Errorf("check document hash: %w", err) }
	if found && currentHash == document.ContentHash { return nil }

	chunks := ChunkText(document.Content, maxChars, overlap)
	if len(chunks) == 0 { return fmt.Errorf("document has no content") }
	embeddings := make([][]float32, 0, len(chunks))
	items := make([]Chunk, 0, len(chunks))
	for i, content := range chunks {
		vector, err := s.Embeddings.Embed(ctx, content)
		if err != nil { return fmt.Errorf("embed chunk %d: %w", i, err) }
		if err := embedding.ValidateDimension(vector, s.Dimension); err != nil {
			return fmt.Errorf("validate embedding chunk %d: %w", i, err)
		}
		embeddings = append(embeddings, vector)
		items = append(items, Chunk{ID: contentID(document.ID, i), DocumentID: document.ID, Index: i, Content: content, Metadata: map[string]string{"source": document.Source, "url": document.URL, "title": document.Title, "type": document.Type, "trust": document.Trust}})
	}
	return s.Repository.Upsert(ctx, document, items, embeddings)
}

func (s *Service) Retrieve(ctx context.Context, query string, limit int) ([]Result, error) {
	if s == nil || s.Embeddings == nil || s.Repository == nil { return nil, fmt.Errorf("RAG service is not configured") }
	if query == "" { return nil, fmt.Errorf("query cannot be empty") }
	if limit <= 0 { limit = 5 }
	vector, err := s.Embeddings.Embed(ctx, query)
	if err != nil { return nil, fmt.Errorf("embed query: %w", err) }
	if err := embedding.ValidateDimension(vector, s.Dimension); err != nil { return nil, fmt.Errorf("validate query embedding: %w", err) }
	return s.Repository.Search(ctx, vector, limit)
}

func ContentHash(content string) string { sum := sha256.Sum256([]byte(content)); return hex.EncodeToString(sum[:]) }
func contentID(documentID string, index int) string { return fmt.Sprintf("%s:%d", documentID, index) }
