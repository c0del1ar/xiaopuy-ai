package rag

import "testing"

func TestChunkText(t *testing.T) {
	text := "First paragraph with useful information.\n\nSecond paragraph with more information.\n\nThird paragraph."
	chunks := ChunkText(text, 55, 8)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	for i, chunk := range chunks {
		if len(chunk) > 55 {
			t.Fatalf("chunk %d has length %d > 55", i, len(chunk))
		}
	}
}

func TestChunkTextEmpty(t *testing.T) {
	if got := ChunkText("   ", 100, 10); got != nil {
		t.Fatalf("expected nil for empty text, got %#v", got)
	}
}

func TestChunkTextLongParagraph(t *testing.T) {
	text := "one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty"
	chunks := ChunkText(text, 30, 5)
	if len(chunks) < 3 {
		t.Fatalf("expected long paragraph to split, got %d chunks", len(chunks))
	}
	for i, chunk := range chunks {
		if len(chunk) > 30 {
			t.Fatalf("chunk %d has length %d > 30", i, len(chunk))
		}
	}
}
