package chat

import (
	"context"
	"testing"

	"github.com/c0del1ar/xiaopuy-ai/internal/ai"
)

type fakeProvider struct {
	response ai.ChatResponse
}

func (p fakeProvider) Chat(_ context.Context, _ ai.ChatRequest) (ai.ChatResponse, error) {
	return p.response, nil
}

func TestServiceReplyPersistsConversation(t *testing.T) {
	repo := NewMemoryRepository()
	conversation := Conversation{ID: "conversation-1", ClientMode: true}
	if err := repo.Save(context.Background(), conversation); err != nil {
		t.Fatal(err)
	}

	agent := &ai.Agent{
		Provider: fakeProvider{response: ai.ChatResponse{Content: "Halo, ada yang bisa saya bantu?"}},
		Persona:  ai.DefaultPersona(),
	}
	service := &Service{Agent: agent, Policy: ai.Policy{}, Repository: repo}

	result, saved, err := service.ReplyToConversation(context.Background(), conversation.ID, "Halo")
	if err != nil {
		t.Fatal(err)
	}
	if result.Decision != ai.AllowReply {
		t.Fatalf("decision = %q, want %q", result.Decision, ai.AllowReply)
	}
	if result.Response == "" {
		t.Fatal("expected a response")
	}
	if len(saved.Messages) != 2 {
		t.Fatalf("saved messages = %d, want 2", len(saved.Messages))
	}

	loaded, err := repo.Get(context.Background(), conversation.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Messages) != 2 {
		t.Fatalf("persisted messages = %d, want 2", len(loaded.Messages))
	}
}

func TestServiceDoesNotPersistBlockedReply(t *testing.T) {
	repo := NewMemoryRepository()
	conversation := Conversation{ID: "conversation-2", ClientMode: true}
	if err := repo.Save(context.Background(), conversation); err != nil {
		t.Fatal(err)
	}

	agent := &ai.Agent{
		Provider: fakeProvider{response: ai.ChatResponse{Content: "should not be called"}},
		Persona:  ai.DefaultPersona(),
	}
	service := &Service{Agent: agent, Policy: ai.Policy{}, Repository: repo}

	result, _, err := service.ReplyToConversation(context.Background(), conversation.ID, "Saya ingin bicara langsung dengan Arya")
	if err != nil {
		t.Fatal(err)
	}
	if result.Decision != ai.EscalateOwner {
		t.Fatalf("decision = %q, want %q", result.Decision, ai.EscalateOwner)
	}

	loaded, err := repo.Get(context.Background(), conversation.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Messages) != 0 {
		t.Fatalf("persisted messages = %d, want 0", len(loaded.Messages))
	}
}
