package router9

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
	Model   string
}

func New(baseURL, apiKey, model string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		Model:   model,
		HTTP:    &http.Client{},
	}
}

type request struct {
	Model    string        `json:"model,omitempty"`
	Messages []ai.Message `json:"messages"`
}

type response struct {
	Choices []struct {
		Message ai.Message `json:"message"`
	} `json:"choices"`
	Model string `json:"model,omitempty"`
}

// Chat implements the OpenAI-compatible chat contract expected from 9router.
func (c *Client) Chat(ctx context.Context, input ai.ChatRequest) (ai.ChatResponse, error) {
	model := input.Model
	if model == "" {
		model = c.Model
	}

	body, err := json.Marshal(request{Model: model, Messages: input.Messages})
	if err != nil {
		return ai.ChatResponse{}, fmt.Errorf("marshal 9router request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return ai.ChatResponse{}, fmt.Errorf("create 9router request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	res, err := c.HTTP.Do(req)
	if err != nil {
		return ai.ChatResponse{}, fmt.Errorf("call 9router: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return ai.ChatResponse{}, fmt.Errorf("9router returned HTTP %d", res.StatusCode)
	}

	var decoded response
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return ai.ChatResponse{}, fmt.Errorf("decode 9router response: %w", err)
	}
	if len(decoded.Choices) == 0 {
		return ai.ChatResponse{}, fmt.Errorf("9router returned no choices")
	}

	return ai.ChatResponse{
		Content: decoded.Choices[0].Message.Content,
		Model:   decoded.Model,
	}, nil
}
