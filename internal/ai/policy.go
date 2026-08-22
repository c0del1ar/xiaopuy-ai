package ai

import "strings"

// ReplyDecision describes what the assistant is allowed to do with an inbound message.
type ReplyDecision string

const (
	AllowReply    ReplyDecision = "allow_reply"
	NeedContext   ReplyDecision = "need_context"
	EscalateOwner ReplyDecision = "escalate_owner"
	DoNotReply    ReplyDecision = "do_not_reply"
)

type PolicyInput struct {
	Message    string
	ClientMode bool
}

// Policy is intentionally deterministic. The model should not decide whether it is
// allowed to act; the application decides first, then asks the model to generate text.
type Policy struct{}

func (Policy) Decide(input PolicyInput) ReplyDecision {
	if strings.TrimSpace(input.Message) == "" {
		return DoNotReply
	}

	if !input.ClientMode {
		return AllowReply
	}

	text := normalize(input.Message)

	// Requests that explicitly ask to talk to the owner should reach the owner.
	for _, phrase := range []string{
		"bicara dengan arya",
		"bicara langsung dengan arya",
		"hubungi arya",
		"chat arya",
		"langsung dengan arya",
	} {
		if contains(text, phrase) {
			return EscalateOwner
		}
	}

	// Sensitive account/payment requests require verified context before answering.
	for _, phrase := range []string{
		"password",
		"otp",
		"kode otp",
		"rekening",
		"transfer",
		"kartu kredit",
		"credit card",
	} {
		if contains(text, phrase) {
			return NeedContext
		}
	}

	return AllowReply
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func contains(s, part string) bool {
	return strings.Contains(s, part)
}
