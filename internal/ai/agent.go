package ai

import "context"

// Agent builds the model context and delegates generation to a provider.
type Agent struct {
	Provider Provider
	Persona  Persona
	Context  ContextBuilder
}

func (a *Agent) Reply(ctx context.Context, history []Message, clientMode bool) (ChatResponse, error) {
	base := append([]Message(nil), history...)
	var messages []Message
	var err error
	if a.Context != nil {
		messages, err = a.Context.Build(base, clientMode)
		if err != nil {
			return ChatResponse{}, err
		}
	} else {
		messages = append([]Message{{
			Role:    "system",
			Content: a.Persona.SystemPrompt(clientMode),
		}}, base...)
	}

	return a.Provider.Chat(ctx, ChatRequest{Messages: messages})
}
