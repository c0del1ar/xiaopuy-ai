package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/c0del1ar/xiaopuy-ai/internal/chat"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Get(ctx context.Context, id string) (chat.Conversation, error) {
	var conversation chat.Conversation

	err := r.pool.QueryRow(ctx, `
		SELECT id, owner_id, contact_id, client_mode
		FROM conversations
		WHERE id = $1`, id).Scan(
		&conversation.ID,
		&conversation.OwnerID,
		&conversation.ContactID,
		&conversation.ClientMode,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return chat.Conversation{}, chat.ErrConversationNotFound
	}
	if err != nil {
		return chat.Conversation{}, fmt.Errorf("get conversation: %w", err)
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, role, content, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC, id ASC`, id)
	if err != nil {
		return chat.Conversation{}, fmt.Errorf("get messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message chat.Message
		if err := rows.Scan(&message.ID, &message.Role, &message.Content, &message.CreatedAt); err != nil {
			return chat.Conversation{}, fmt.Errorf("scan message: %w", err)
		}
		conversation.Messages = append(conversation.Messages, message)
	}
	if err := rows.Err(); err != nil {
		return chat.Conversation{}, fmt.Errorf("iterate messages: %w", err)
	}

	return conversation, nil
}

func (r *Repository) Save(ctx context.Context, conversation chat.Conversation) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // harmless after Commit

	_, err = tx.Exec(ctx, `
		INSERT INTO conversations (id, owner_id, contact_id, client_mode)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			owner_id = EXCLUDED.owner_id,
			contact_id = EXCLUDED.contact_id,
			client_mode = EXCLUDED.client_mode,
			updated_at = NOW()`,
		conversation.ID,
		conversation.OwnerID,
		conversation.ContactID,
		conversation.ClientMode,
	)
	if err != nil {
		return fmt.Errorf("upsert conversation: %w", err)
	}

	// Rebuilding the message set keeps Save semantics simple while the domain is
	// still small. Later we can introduce append-only message persistence.
	if _, err := tx.Exec(ctx, `DELETE FROM messages WHERE conversation_id = $1`, conversation.ID); err != nil {
		return fmt.Errorf("replace messages: %w", err)
	}

	for index, message := range conversation.Messages {
		id := message.ID
		if id == "" {
			id, err = newMessageID(index)
			if err != nil {
				return fmt.Errorf("generate message id: %w", err)
			}
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO messages (id, conversation_id, role, content, created_at)
			VALUES ($1, $2, $3, $4, $5)`,
			id,
			conversation.ID,
			message.Role,
			message.Content,
			message.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert message: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit conversation: %w", err)
	}
	return nil
}

func newMessageID(index int) (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("msg_%d_%s", index, hex.EncodeToString(buf)), nil
}
