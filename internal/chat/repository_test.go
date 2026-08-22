package chat

import (
	"context"
	"testing"
)

func TestMemoryRepositoryRoundTrip(t *testing.T) {
	repo := NewMemoryRepository()
	conversation := Conversation{ID: "conv-1", ContactID: "client-1", ClientMode: true}
	conversation.Add(RoleUser, "Halo")
	conversation.Add(RoleAssistant, "Halo, ada yang bisa saya bantu?")

	if err := repo.Save(context.Background(), conversation); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := repo.Get(context.Background(), conversation.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(got.Messages) != 2 {
		t.Fatalf("Get() returned %d messages, want 2", len(got.Messages))
	}
	if got.Messages[0].Content != "Halo" || got.Messages[1].Role != RoleAssistant {
		t.Fatalf("conversation contents were not preserved: %+v", got.Messages)
	}
}

func TestMemoryRepositoryMissingConversation(t *testing.T) {
	repo := NewMemoryRepository()

	_, err := repo.Get(context.Background(), "missing")
	if err != ErrConversationNotFound {
		t.Fatalf("Get() error = %v, want ErrConversationNotFound", err)
	}
}
