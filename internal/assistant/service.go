package assistant

import (
	"context"
	"fmt"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
	"github.com/c0del1ar/xiaopuy-ai/internal/rag"
)

type Service struct {
	Agent    *ai.Agent
	RAG      *rag.Service
	Context  rag.ContextBuilder
	Policy   RetrievalPolicy
	MinScore float32
}

func (s *Service) Reply(ctx context.Context, history []ai.Message, clientMode bool) (ai.ChatResponse, error) {
	if s == nil || s.Agent == nil { return ai.ChatResponse{}, fmt.Errorf("assistant service is not configured") }
	if s.RAG == nil { return s.Agent.Reply(ctx, history, clientMode) }
	query := lastUserMessage(history)
	if query == "" { return ai.ChatResponse{}, fmt.Errorf("conversation has no user message") }

	if !s.Policy.ShouldRetrieve(query) {
		return s.Agent.Reply(ctx, history, clientMode)
	}

	results, err := s.RAG.Retrieve(ctx, query, s.Policy.TopK())
	if err != nil { return ai.ChatResponse{}, fmt.Errorf("retrieve knowledge: %w", err) }

	threshold := s.MinScore
	if threshold <= 0 { threshold = 0.75 }
	chunks := make([]rag.Chunk, 0, len(results))
	for _, result := range results {
		if result.Score >= threshold { chunks = append(chunks, result.Chunk) }
	}
	if len(chunks) == 0 { return s.Agent.Reply(ctx, history, clientMode) }

	knowledge := s.Context.Build(chunks)
	agent := *s.Agent
	agent.Context = ai.KnowledgeContext{Persona: agent.Persona, Knowledge: knowledge}
	return agent.Reply(ctx, history, clientMode)
}

func lastUserMessage(history []ai.Message) string {
	for i := len(history)-1; i >= 0; i-- { if history[i].Role == "user" { return history[i].Content } }
	return ""
}
