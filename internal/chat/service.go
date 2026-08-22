package chat

import (
	"context"
	"fmt"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
)

type Service struct {
	Agent      *ai.Agent
	Policy     ai.Policy
	Repository Repository
}

type ReplyResult struct {
	Decision ai.ReplyDecision
	Response string
}

func (s *Service) Reply(ctx context.Context, conversation *Conversation, userText string) (ReplyResult, error) {
	if s == nil || s.Agent == nil {
		return ReplyResult{}, fmt.Errorf("chat service has no AI agent")
	}
	if conversation == nil {
		return ReplyResult{}, fmt.Errorf("conversation cannot be nil")
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

	if s.Repository != nil {
		if err := s.Repository.Save(ctx, *conversation); err != nil {
			return ReplyResult{}, fmt.Errorf("persist conversation: %w", err)
		}
	}

	return ReplyResult{
		Decision: ai.AllowReply,
		Response: response.Content,
	}, nil
}

// ReplyToConversation loads the latest conversation state before generating a reply.
// This is the persistence-aware entry point used by channel adapters.
func (s *Service) ReplyToConversation(ctx context.Context, conversationID string, userText string) (ReplyResult, Conversation, error) {
	if s == nil || s.Repository == nil {
		return ReplyResult{}, Conversation{}, fmt.Errorf("chat service has no repository")
	}

	conversation, err := s.Repository.Get(ctx, conversationID)
	if err != nil {
		return ReplyResult{}, Conversation{}, fmt.Errorf("load conversation: %w", err)
	}

	result, err := s.Reply(ctx, &conversation, userText)
	if err != nil {
		return ReplyResult{}, Conversation{}, err
	}
	return result, conversation, nil
}
