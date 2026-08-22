package ingest

import (
	"context"
	"testing"

	"github.com/c0del1ar/xiaopuy-ai/internal/rag"
)

type mapFetcher struct {
	pages map[string]string
	calls []string
}

func (f *mapFetcher) Fetch(_ context.Context, raw string) (Response, error) {
	f.calls = append(f.calls, raw)
	return Response{URL: raw, StatusCode: 200, ContentType: "text/html", Body: []byte(f.pages[raw])}, nil
}

type captureIndexer struct{ docs []rag.Document }

func (i *captureIndexer) Index(_ context.Context, doc rag.Document, _, _ int) error {
	i.docs = append(i.docs, doc)
	return nil
}

func TestCrawlerRespectsDomainDepthAndPageLimit(t *testing.T) {
	pages := map[string]string{
		"https://aryakun.id/": `<html><body><h1>Home page content for crawler tests.</h1><a href="/about">About</a><a href="https://example.com/out">External</a></body></html>`,
		"https://aryakun.id/about": `<html><body><h1>About page content for crawler tests.</h1><a href="/team">Team</a></body></html>`,
		"https://aryakun.id/team": `<html><body><h1>Team page content for crawler tests.</h1></body></html>`,
	}
	fetcher := &mapFetcher{pages: pages}
	indexer := &captureIndexer{}
	c := &Crawler{Fetcher: fetcher, Indexer: indexer, Config: CrawlerConfig{MaxPages: 10, MaxDepth: 1, AllowedDomains: []string{"aryakun.id"}}}
	result, err := c.Crawl(context.Background(), "https://aryakun.id/")
	if err != nil {
		t.Fatalf("Crawl() error=%v", err)
	}
	if result.Visited != 2 {
		t.Fatalf("Visited=%d want 2", result.Visited)
	}
	if len(result.Documents) != 2 {
		t.Fatalf("Documents=%d want 2", len(result.Documents))
	}
	if len(indexer.docs) != 2 {
		t.Fatalf("indexed documents=%d want 2", len(indexer.docs))
	}
	for _, call := range fetcher.calls {
		if call == "https://example.com/out" {
			t.Fatal("crawler followed external domain")
		}
	}
}

func TestCrawlerDeduplicatesURLs(t *testing.T) {
	pages := map[string]string{
		"https://aryakun.id/": `<html><body><h1>Home page content enough for crawler.</h1><a href="/about">A</a><a href="/about#section">B</a></body></html>`,
		"https://aryakun.id/about": `<html><body><h1>About page content enough for crawler.</h1></body></html>`,
	}
	fetcher := &mapFetcher{pages: pages}
	c := &Crawler{Fetcher: fetcher, Config: CrawlerConfig{MaxPages: 10, MaxDepth: 2, AllowedDomains: []string{"aryakun.id"}}}
	result, err := c.Crawl(context.Background(), "https://aryakun.id/")
	if err != nil {
		t.Fatalf("Crawl() error=%v", err)
	}
	if result.Visited != 2 {
		t.Fatalf("Visited=%d want 2", result.Visited)
	}
}

func TestNormalizeURL(t *testing.T) {
	got, err := normalizeURL("HTTPS://Aryakun.ID/path#section")
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://aryakun.id/path" {
		t.Fatalf("got %q", got)
	}
}
