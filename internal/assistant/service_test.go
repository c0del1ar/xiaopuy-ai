package assistant

import (
	"testing"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
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
