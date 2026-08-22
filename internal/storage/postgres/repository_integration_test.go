package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/c0del1ar/xiaopuy-ai/internal/chat"
)

func TestRepositoryRoundTrip(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo, err := New(ctx, dsn)
	if err != nil {
		t.Fatalf("create repository: %v", err)
	}
	defer repo.Close()

	conversation := chat.NewConversation("integration-test", "owner", "contact", true)
	conversation.Add(chat.RoleUser, "Hello")
	conversation.Add(chat.RoleAssistant, "Hi, how can I help?")

	if err := repo.Save(ctx, conversation); err != nil {
		t.Fatalf("save conversation: %v", err)
	}

	got, err := repo.Get(ctx, conversation.ID)
	if err != nil {
		t.Fatalf("get conversation: %v", err)
	}

	if got.ID != conversation.ID {
		t.Fatalf("ID = %q, want %q", got.ID, conversation.ID)
	}
	if got.ClientMode != conversation.ClientMode {
		t.Fatalf("ClientMode = %v, want %v", got.ClientMode, conversation.ClientMode)
	}
	if len(got.Messages) != 2 {
		t.Fatalf("message count = %d, want 2", len(got.Messages))
	}
	if got.Messages[0].Content != "Hello" || got.Messages[1].Content != "Hi, how can I help?" {
		t.Fatalf("messages were not persisted correctly: %+v", got.Messages)
	}
}
