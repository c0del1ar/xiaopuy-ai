package ingest

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/c0del1ar/xiaopuy-ai/internal/rag"
	"golang.org/x/net/html"
)

type Fetcher interface { Fetch(context.Context, string) (string, error) }
type Indexer interface { Index(context.Context, rag.Document, int, int) error }

type CrawlerConfig struct {
	MaxPages       int
	MaxDepth       int
	MaxConcurrency int
	MaxChars       int
	ChunkOverlap   int
	AllowedDomains []string
}

type CrawlResult struct {
	Documents []rag.Document
	Visited   int
	Skipped   int
}

type Crawler struct { Fetcher Fetcher; Indexer Indexer; Config CrawlerConfig }

func (c *Crawler) Crawl(ctx context.Context, seed string) (CrawlResult, error) {
	if c == nil || c.Fetcher == nil { return CrawlResult{}, fmt.Errorf("crawler fetcher is not configured") }
	if seed == "" { return CrawlResult{}, fmt.Errorf("seed URL cannot be empty") }
	cfg := c.Config
	if cfg.MaxPages <= 0 { cfg.MaxPages = 50 }
	if cfg.MaxDepth < 0 { cfg.MaxDepth = 0 }
	if cfg.MaxConcurrency <= 0 { cfg.MaxConcurrency = 2 }
	seedURL, err := normalizeURL(seed); if err != nil { return CrawlResult{}, err }
	if !allowed(seedURL, cfg.AllowedDomains) { return CrawlResult{}, fmt.Errorf("seed URL is outside allowed domains") }

	type item struct { raw string; depth int }
	queue := []item{{seedURL, 0}}
	visited := map[string]bool{}
	var result CrawlResult
	sem := make(chan struct{}, cfg.MaxConcurrency)
	var mu sync.Mutex

	for len(queue) > 0 && result.Visited < cfg.MaxPages {
		select { case <-ctx.Done(): return result, ctx.Err(); default: }
		current := queue[0]; queue = queue[1:]
		u, err := normalizeURL(current.raw); if err != nil { result.Skipped++; continue }
		if visited[u] || !allowed(u, cfg.AllowedDomains) { result.Skipped++; continue }
		visited[u] = true
		sem <- struct{}{}
		raw, err := c.Fetcher.Fetch(ctx, u)
		<-sem
		result.Visited++
		if err != nil { result.Skipped++; continue }
		page, err := ParseHTML(raw, u); if err != nil { result.Skipped++; continue }
		docID := "web:" + u
		doc := page.Document(docID, u)
		result.Documents = append(result.Documents, doc)
		if c.Indexer != nil {
			if err := c.Indexer.Index(ctx, doc, cfg.MaxChars, cfg.ChunkOverlap); err != nil { return result, fmt.Errorf("index %s: %w", u, err) }
		}
		if current.depth >= cfg.MaxDepth { continue }
		links := extractLinks(raw, u)
		mu.Lock()
		for _, link := range links { if !visited[link] && allowed(link, cfg.AllowedDomains) { queue = append(queue, item{link, current.depth+1}) } }
		mu.Unlock()
	}
	return result, nil
}

func normalizeURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw)); if err != nil { return "", fmt.Errorf("invalid URL %q: %w", raw, err) }
	if u.Scheme != "http" && u.Scheme != "https" { return "", fmt.Errorf("unsupported URL scheme %q", u.Scheme) }
	u.Fragment = ""
	u.Host = strings.ToLower(u.Host)
	if u.Path == "" { u.Path = "/" }
	return u.String(), nil
}

func allowed(raw string, domains []string) bool {
	if len(domains) == 0 { return false }
	u, err := url.Parse(raw); if err != nil { return false }
	host := strings.ToLower(strings.TrimSuffix(u.Hostname(), "."))
	for _, d := range domains { d = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(d, "."))); if d != "" && (host == d || strings.HasSuffix(host, "."+d)) { return true } }
	return false
}

func extractLinks(raw, base string) []string {
	doc, err := html.Parse(strings.NewReader(raw)); if err != nil { return nil }
	baseURL, err := url.Parse(base); if err != nil { return nil }
	seen := map[string]bool{}; var links []string
	var walk func(*html.Node)
	walk = func(n *html.Node) { if n.Type == html.ElementNode && strings.EqualFold(n.Data,"a") { href := attr(n,"href"); if href != "" { if u,err:=url.Parse(href); err==nil && (u.Scheme=="" || u.Scheme=="http" || u.Scheme=="https") { resolved:=baseURL.ResolveReference(u); resolved.Fragment=""; if normalized,err:=normalizeURL(resolved.String()); err==nil && !seen[normalized] {seen[normalized]=true; links=append(links,normalized)} } } }; for child:=n.FirstChild; child!=nil; child=child.NextSibling {walk(child)} }
	walk(doc)
	return links
}
