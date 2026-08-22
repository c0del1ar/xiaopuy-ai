package chat

import (
	"context"
	"fmt"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
)

type Service struct {
	Agent  *ai.Agent
	Policy ai.Policy
}

type ReplyResult struct {
	Decision ai.ReplyDecision
	Response string
}

func (s *Service) Reply(ctx context.Context, conversation *Conversation, userText string) (ReplyResult, error) {
	if s == nil || s.Agent == nil {
		return ReplyResult{}, fmt.Errorf("chat service has no AI agent")
	}
	if userText == "" {
		return ReplyResult{}, fmt.Errorf("message cannot be empty")
	}

	decision := s.Policy.Decide(ai.PolicyInput{
		Message:    userText,
		ClientMode: conversation.ClientMode,
	})
	if decision != ai.AllowReply {
		return ReplyResult{Decision: decision}, nil
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
		return ReplyResult{}, err
	}

	conversation.Add(RoleAssistant, response.Content)
	return ReplyResult{
		Decision: ai.AllowReply,
		Response: response.Content,
	}, nil
}
