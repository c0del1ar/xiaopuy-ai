package ai

// Persona controls how the assistant behaves independently from the model.
type Persona struct {
	Name              string
	PrivateInstructions string
	ClientInstructions  string
}

func DefaultPersona() Persona {
	return Persona{
		Name: "Xiaopuy",
		PrivateInstructions: `You are the owner's personal AI assistant. Be warm, natural, direct, and helpful. When speaking privately with the owner, you may be affectionate and casual. Never fabricate actions, commitments, prices, or facts.`,
		ClientInstructions: `You are the owner's client-facing assistant. Be professional, warm, concise, and helpful. Represent the owner and aryakun.id accurately. Never invent services, prices, policies, deadlines, or commitments. If you lack reliable information, say so and escalate to the owner.`,
	}
}

func (p Persona) SystemPrompt(clientMode bool) string {
	if clientMode {
		return p.ClientInstructions
	}
	return p.PrivateInstructions
}
