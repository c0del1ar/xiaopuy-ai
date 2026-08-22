package chat

import (
	"context"
	"fmt"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
)

type Service struct {
	Agent *ai.Agent
}

func (s *Service) Reply(ctx context.Context, conversation *Conversation, userText string) (string, error) {
	if s == nil || s.Agent == nil {
		return "", fmt.Errorf("chat service has no AI agent")
	}
	if userText == "" {
		return "", fmt.Errorf("message cannot be empty")
	}

	conversation.Add(RoleUser, userText)

	history := make([]ai.Message, 0, len(conversation.Messages))
	for _, message := range conversation.Messages {
		history = append(history, ai.Message{
			Role:    string(message.Role),
			Content: message.Content,
		})
	}

	response, err := s.Agent.Reply(ctx, history, conversation.ClientMode)
	if err != nil {
		return "", err
	}

	conversation.Add(RoleAssistant, response.Content)
	return response.Content, nil
}
