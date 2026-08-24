package ingest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPHandlerCrawl(t *testing.T) {
	fetcher := &mapFetcher{pages: map[string]string{
		"https://aryakun.id/": "<html><body><h1>Website knowledge content that is long enough for the ingestion parser to retain as a document.</h1></body></html>",
	}}
	handler := &HTTPHandler{Service: &Service{
		Crawler: &Crawler{Fetcher: fetcher},
		Config:  Config{AllowedDomains: []string{"aryakun.id"}},
	}}

	req := httptest.NewRequest(http.MethodPost, "/v1/ingest/crawl", bytes.NewBufferString(`{"seed_url":"https://aryakun.id/"}`))
	res := httptest.NewRecorder()
	handler.Crawl(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", res.Code, http.StatusOK, res.Body.String())
	}
	if got := res.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if got := res.Body.String(); got != "{\"visited\":1,\"skipped\":0,\"documents\":1}\n" {
		t.Fatalf("body = %q", got)
	}
}

func TestHTTPHandlerCrawlRejectsInvalidJSON(t *testing.T) {
	handler := &HTTPHandler{Service: &Service{}}
	req := httptest.NewRequest(http.MethodPost, "/v1/ingest/crawl", bytes.NewBufferString(`{`))
	res := httptest.NewRecorder()
	handler.Crawl(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
	}
}
