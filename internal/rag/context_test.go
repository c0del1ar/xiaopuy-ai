package rag

import "testing"

func TestContextBuilder(t *testing.T) {
	builder := ContextBuilder{MaxChars: 200}
	chunks := []Chunk{
		{Content: "Laravel website development", Metadata: map[string]string{"source": "aryakun.id", "title": "Services"}},
		{Content: "Logo design services", Metadata: map[string]string{"source": "aryakun.id", "title": "Design"}},
	}

	got := builder.Build(chunks)
	if got == "" {
		t.Fatal("expected context")
	}
	if len(got) > 200 {
		t.Fatalf("context length = %d, want <= 200", len(got))
	}
	if got == "Laravel website development" {
		t.Fatal("expected attributed context")
	}
}

func TestContextBuilderEmpty(t *testing.T) {
	if got := (ContextBuilder{}).Build(nil); got != "" {
		t.Fatalf("expected empty context, got %q", got)
	}
}
