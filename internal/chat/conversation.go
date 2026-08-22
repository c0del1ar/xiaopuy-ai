package chat

import "time"

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type Message struct {
	ID        string    `json:"id,omitempty"`
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Conversation struct {
	ID         string    `json:"id"`
	OwnerID    string    `json:"owner_id,omitempty"`
	ContactID  string    `json:"contact_id,omitempty"`
	ClientMode bool      `json:"client_mode"`
	Messages   []Message `json:"messages"`
}

func (c *Conversation) Add(role Role, content string) {
	c.Messages = append(c.Messages, Message{
		Role:      role,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	})
}
