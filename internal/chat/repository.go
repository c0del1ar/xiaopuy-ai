package chat

import (
	"context"
	"sync"
)

// Repository persists conversation state. The application depends on this
// interface so PostgreSQL can be introduced without coupling the chat domain
// to a database driver.
type Repository interface {
	Get(ctx context.Context, id string) (Conversation, error)
	Save(ctx context.Context, conversation Conversation) error
}

// MemoryRepository is useful for local development and unit tests.
type MemoryRepository struct {
	mu            sync.RWMutex
	conversations map[string]Conversation
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{conversations: make(map[string]Conversation)}
}

func (r *MemoryRepository) Get(_ context.Context, id string) (Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conversation, ok := r.conversations[id]
	if !ok {
		return Conversation{}, ErrConversationNotFound
	}
	return cloneConversation(conversation), nil
}

func (r *MemoryRepository) Save(_ context.Context, conversation Conversation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.conversations[conversation.ID] = cloneConversation(conversation)
	return nil
}

func cloneConversation(src Conversation) Conversation {
	dst := src
	dst.Messages = append([]Message(nil), src.Messages...)
	return dst
}
