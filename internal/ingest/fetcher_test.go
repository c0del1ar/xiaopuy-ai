package ingest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetcherFetchHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != DefaultUserAgent { t.Errorf("User-Agent = %q", r.Header.Get("User-Agent")) }
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<html>hello</html>")
	}))
	defer server.Close()

	got, err := (&Fetcher{}).Fetch(context.Background(), server.URL)
	if err != nil { t.Fatalf("Fetch() error = %v", err) }
	if got.StatusCode != http.StatusOK { t.Fatalf("status = %d, want 200", got.StatusCode) }
	if string(got.Body) != "<html>hello</html>" { t.Fatalf("body = %q", got.Body) }
}

func TestFetcherRejectsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	if _, err := (&Fetcher{}).Fetch(context.Background(), server.URL); err == nil || !strings.Contains(err.Error(), "status 404") { t.Fatalf("error = %v, want HTTP 404", err) }
}

func TestFetcherRejectsUnsupportedContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "application/pdf"); fmt.Fprint(w, "pdf") }))
	defer server.Close()
	if _, err := (&Fetcher{}).Fetch(context.Background(), server.URL); err == nil || !strings.Contains(err.Error(), "unsupported content type") { t.Fatalf("error = %v", err) }
}

func TestFetcherEnforcesBodyLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "text/html"); fmt.Fprint(w, "1234567890") }))
	defer server.Close()
	_, err := (&Fetcher{MaxBytes: 5}).Fetch(context.Background(), server.URL)
	if err == nil || !strings.Contains(err.Error(), "exceeds 5 bytes") { t.Fatalf("error = %v", err) }
}

func TestFetcherHonorsContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err := (&Fetcher{Timeout: time.Second}).Fetch(ctx, server.URL)
	if err == nil { t.Fatal("Fetch() error = nil, want cancellation error") }
}
