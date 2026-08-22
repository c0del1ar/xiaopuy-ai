package router9

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// EmbeddingClient implements the OpenAI-compatible /v1/embeddings contract.
// It is separate from chat generation so the embedding model can be configured independently.
type EmbeddingClient struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
	Model   string
}

func NewEmbedding(baseURL, apiKey, model string) *EmbeddingClient {
	return &EmbeddingClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		Model:   model,
		HTTP:    &http.Client{},
	}
}

type embeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model,omitempty"`
}

func (c *EmbeddingClient) Embed(ctx context.Context, text string) ([]float32, error) {
	if c == nil {
		return nil, fmt.Errorf("embedding client is nil")
	}
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("embedding input cannot be empty")
	}
	if c.Model == "" {
		return nil, fmt.Errorf("embedding model is not configured")
	}

	body, err := json.Marshal(embeddingRequest{Model: c.Model, Input: text})
	if err != nil {
		return nil, fmt.Errorf("marshal embedding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create embedding request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call embedding provider: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("embedding provider returned HTTP %d", res.StatusCode)
	}

	var decoded embeddingResponse
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode embedding response: %w", err)
	}
	if len(decoded.Data) == 0 || len(decoded.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("embedding provider returned no vector")
	}
	return decoded.Data[0].Embedding, nil
}
