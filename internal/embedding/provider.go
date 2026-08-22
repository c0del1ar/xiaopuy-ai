package embedding

import "context"

// Provider converts text into vectors. Generation and embedding are intentionally
// separate so the RAG layer is not coupled to a particular model gateway.
type Provider interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}
