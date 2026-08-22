package ingest

import (
	"context"
	"testing"
)

func TestServiceCrawlAppliesConfiguration(t *testing.T) {
	fetcher := &mapFetcher{pages: map[string]string{
		"https://aryakun.id/": `<html><body><h1>Home page content sufficiently long for service tests.</h1><a href="/about">About</a></body></html>`,
		"https://aryakun.id/about": `<html><body><h1>About page content sufficiently long for service tests.</h1></body></html>`,
	}}
	crawler := &Crawler{Fetcher: fetcher}
	service := &Service{Crawler: crawler, Config: Config{
		MaxPages: 1, MaxDepth: 2, MaxConcurrency: 1,
		MaxBytes: 4096, MaxChars: 1000, ChunkOverlap: 100,
		AllowedDomains: []string{"aryakun.id"},
	}}

	result, err := service.Crawl(context.Background(), CrawlRequest{SeedURL: "https://aryakun.id/"})
	if err != nil { t.Fatalf("Crawl() error = %v", err) }
	if result.Visited != 1 { t.Fatalf("Visited = %d, want 1", result.Visited) }
	if crawler.Config.MaxPages != 1 || crawler.Config.MaxDepth != 2 || crawler.Config.MaxConcurrency != 1 {
		t.Fatalf("crawler limits were not applied: %+v", crawler.Config)
	}
}

func TestServiceRequiresAllowedDomains(t *testing.T) {
	service := &Service{Crawler: &Crawler{Fetcher: &mapFetcher{pages: map[string]string{}}}}
	_, err := service.Crawl(context.Background(), CrawlRequest{SeedURL: "https://aryakun.id/"})
	if err == nil { t.Fatal("Crawl() error = nil, want allowed-domain configuration error") }
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MaxPages != 50 || cfg.MaxDepth != 3 || cfg.MaxConcurrency != 3 || cfg.MaxChars != 1600 || cfg.ChunkOverlap != 200 {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}
