package ingest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultFetchTimeout = 15 * time.Second
	DefaultMaxBodySize  = 2 << 20 // 2 MiB
	DefaultUserAgent    = "XiaopuyAI/1.0 (+https://aryakun.id)"
)

type Fetcher struct {
	Client    *http.Client
	Timeout   time.Duration
	MaxBytes  int64
	UserAgent string
}

type Response struct {
	URL         string
	StatusCode  int
	ContentType string
	Body        []byte
}

func (f *Fetcher) Fetch(ctx context.Context, url string) (Response, error) {
	if url == "" { return Response{}, fmt.Errorf("URL cannot be empty") }
	client := f.Client
	if client == nil { client = &http.Client{} }
	timeout := f.Timeout
	if timeout <= 0 { timeout = DefaultFetchTimeout }
	maxBytes := f.MaxBytes
	if maxBytes <= 0 { maxBytes = DefaultMaxBodySize }
	userAgent := f.UserAgent
	if userAgent == "" { userAgent = DefaultUserAgent }

	requestCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil { return Response{}, fmt.Errorf("create request: %w", err) }
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml;q=0.9,text/plain;q=0.8")

	resp, err := client.Do(req)
	if err != nil { return Response{}, fmt.Errorf("fetch %s: %w", url, err) }
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return Response{}, fmt.Errorf("fetch %s: unexpected HTTP status %d", url, resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !isSupportedContentType(contentType) { return Response{}, fmt.Errorf("fetch %s: unsupported content type %q", url, contentType) }

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil { return Response{}, fmt.Errorf("read %s: %w", url, err) }
	if int64(len(body)) > maxBytes { return Response{}, fmt.Errorf("fetch %s: response exceeds %d bytes", url, maxBytes) }
	return Response{URL: resp.Request.URL.String(), StatusCode: resp.StatusCode, ContentType: contentType, Body: body}, nil
}

func isSupportedContentType(contentType string) bool {
	for i := 0; i < len(contentType); i++ { if contentType[i] == ';' { contentType = contentType[:i]; break } }
	switch contentType {
	case "text/html", "application/xhtml+xml", "text/plain": return true
	default: return false
	}
}
