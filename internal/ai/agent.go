package ai

import "context"

// Agent builds the model context and delegates generation to a provider.
type Agent struct {
	Provider Provider
	Persona  Persona
}

func (a *Agent) Reply(ctx context.Context, history []Message, clientMode bool) (ChatResponse, error) {
	messages := make([]Message, 0, len(history)+1)
	messages = append(messages, Message{
		Role:    "system",
		Content: a.Persona.SystemPrompt(clientMode),
	})
	messages = append(messages, history...)

	return a.Provider.Chat(ctx, ChatRequest{Messages: messages})
}
