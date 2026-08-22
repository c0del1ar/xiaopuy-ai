package router9

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmbeddingClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/embeddings" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"embedding":[0.1,0.2,0.3],"index":0}],"model":"test-embed"}`))
	}))
	defer server.Close()

	client := NewEmbedding(server.URL, "test-key", "test-embed")
	vector, err := client.Embed(context.Background(), "hello")
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}
	if len(vector) != 3 || vector[0] != 0.1 || vector[2] != 0.3 {
		t.Fatalf("unexpected vector: %#v", vector)
	}
}

func TestEmbeddingClientRejectsEmptyInput(t *testing.T) {
	client := NewEmbedding("http://example.test", "", "test-embed")
	if _, err := client.Embed(context.Background(), " "); err == nil {
		t.Fatal("expected empty input error")
	}
}
