package ai

import "context"

// Message is a single turn in an AI conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the provider-neutral request sent to an LLM gateway.
type ChatRequest struct {
	Model       string    `json:"model,omitempty"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
}

// ChatResponse is the provider-neutral model response.
type ChatResponse struct {
	Content string `json:"content"`
	Model   string `json:"model,omitempty"`
}

// Provider isolates the AI core from a concrete model gateway such as 9router.
type Provider interface {
	Chat(ctx context.Context, request ChatRequest) (ChatResponse, error)
}
