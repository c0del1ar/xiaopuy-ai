package ingest

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// Config contains the operational limits for website ingestion.
type Config struct {
	MaxPages       int
	MaxDepth       int
	MaxConcurrency int
	MaxBytes       int64
	MaxChars       int
	ChunkOverlap   int
	AllowedDomains []string
}

func DefaultConfig() Config {
	return Config{
		MaxPages:       50,
		MaxDepth:       3,
		MaxConcurrency: 3,
		MaxBytes:       DefaultMaxBodySize,
		MaxChars:       1600,
		ChunkOverlap:   200,
	}
}

// CrawlRequest describes one explicit ingestion run.
type CrawlRequest struct {
	SeedURL string
}

// Service owns the ingestion workflow. HTTP handlers, schedulers, and CLI code
// should call this service instead of constructing or controlling a Crawler.
type Service struct {
	Crawler *Crawler
	Config  Config
	mu      sync.Mutex
}

func (s *Service) Crawl(ctx context.Context, request CrawlRequest) (CrawlResult, error) {
	if s == nil || s.Crawler == nil {
		return CrawlResult{}, fmt.Errorf("ingestion service is not configured")
	}
	if strings.TrimSpace(request.SeedURL) == "" {
		return CrawlResult{}, fmt.Errorf("seed URL cannot be empty")
	}

	// A crawl mutates the visited/indexing state of the run. Serialize service
	// runs for now; parallelism inside a run is still bounded by crawler config.
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg := s.Config
	if cfg.MaxPages <= 0 || cfg.MaxDepth < 0 || cfg.MaxConcurrency <= 0 {
		defaults := DefaultConfig()
		if cfg.MaxPages <= 0 { cfg.MaxPages = defaults.MaxPages }
		if cfg.MaxDepth < 0 { cfg.MaxDepth = defaults.MaxDepth }
		if cfg.MaxConcurrency <= 0 { cfg.MaxConcurrency = defaults.MaxConcurrency }
	}
	if cfg.MaxBytes <= 0 { cfg.MaxBytes = DefaultMaxBodySize }

	// Keep fetcher limits and crawler limits in one service-level configuration.
	if fetcher, ok := s.Crawler.Fetcher.(*Fetcher); ok {
		if fetcher.MaxBytes <= 0 || fetcher.MaxBytes != cfg.MaxBytes {
			fetcher.MaxBytes = cfg.MaxBytes
		}
	}

	s.Crawler.Config.MaxPages = cfg.MaxPages
	s.Crawler.Config.MaxDepth = cfg.MaxDepth
	s.Crawler.Config.MaxConcurrency = cfg.MaxConcurrency
	s.Crawler.Config.MaxChars = cfg.MaxChars
	s.Crawler.Config.ChunkOverlap = cfg.ChunkOverlap
	s.Crawler.Config.AllowedDomains = append([]string(nil), cfg.AllowedDomains...)

	if len(s.Crawler.Config.AllowedDomains) == 0 {
		return CrawlResult{}, fmt.Errorf("allowed domains are not configured")
	}

	return s.Crawler.Crawl(ctx, request.SeedURL)
}
