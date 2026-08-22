package rag

import "time"

// Document is an authoritative knowledge source that can be indexed for retrieval.
type Document struct {
	ID          string
	Source      string
	URL         string
	Title       string
	Type        string
	Trust       string
	Content     string
	ContentHash string
	UpdatedAt   time.Time
}

// Chunk is the retrieval unit derived from a document.
type Chunk struct {
	ID         string
	DocumentID string
	Index      int
	Content    string
	Metadata   map[string]string
}
