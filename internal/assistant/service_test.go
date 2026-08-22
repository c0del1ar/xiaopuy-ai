package assistant

import (
	"testing"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
	"github.com/c0del1ar/xiaopuy-ai/internal/rag"
)

func TestLastUserMessage(t *testing.T) {
	history := []ai.Message{
		{Role: "user", Content: "first"},
		{Role: "assistant", Content: "reply"},
		{Role: "user", Content: "latest"},
	}
	if got := lastUserMessage(history); got != "latest" {
		t.Fatalf("got %q, want latest", got)
	}
}

func TestLastUserMessageMissing(t *testing.T) {
	if got := lastUserMessage([]ai.Message{{Role: "assistant", Content: "hello"}}); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
}

func filterResultsByScore(results []rag.Result, threshold float32) []rag.Chunk {
	chunks := make([]rag.Chunk, 0, len(results))
	for _, result := range results {
		if result.Score >= threshold {
			chunks = append(chunks, result.Chunk)
		}
	}
	return chunks
}

func TestFilterResultsByScore(t *testing.T) {
	results := []rag.Result{
		{Chunk: rag.Chunk{ID: "high"}, Score: 0.91},
		{Chunk: rag.Chunk{ID: "edge"}, Score: 0.75},
		{Chunk: rag.Chunk{ID: "low"}, Score: 0.42},
	}
	chunks := filterResultsByScore(results, 0.75)
	if len(chunks) != 2 { t.Fatalf("got %d chunks, want 2", len(chunks)) }
	if chunks[0].ID != "high" || chunks[1].ID != "edge" { t.Fatalf("unexpected chunks: %+v", chunks) }
}

func TestFilterResultsByScoreNoMatch(t *testing.T) {
	results := []rag.Result{{Chunk: rag.Chunk{ID: "low"}, Score: 0.42}}
	if got := filterResultsByScore(results, 0.75); len(got) != 0 { t.Fatalf("got %d chunks, want 0", len(got)) }
}
