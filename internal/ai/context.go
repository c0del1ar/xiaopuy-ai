package ai

// ContextBuilder is the future boundary for memory and RAG. Keeping it separate
// prevents those concerns from leaking into providers and channel adapters.
type ContextBuilder interface {
	Build(base []Message, clientMode bool) ([]Message, error)
}

type PassthroughContext struct{}

func (PassthroughContext) Build(base []Message, _ bool) ([]Message, error) {
	return base, nil
}
