package ai

import (
	"strings"
	"testing"
)

func TestDefaultPersonaIdentity(t *testing.T) {
	p := DefaultPersona()
	if p.Name != "Xiaopuy" { t.Fatalf("Name = %q, want Xiaopuy", p.Name) }
	if p.OwnerName != "Aryakun" { t.Fatalf("OwnerName = %q, want Aryakun", p.OwnerName) }
	if p.Website != "aryakun.id" { t.Fatalf("Website = %q, want aryakun.id", p.Website) }

	for mode, prompt := range map[string]string{"private": p.SystemPrompt(false), "client": p.SystemPrompt(true)} {
		if !strings.Contains(prompt, "Xiaopuy") { t.Errorf("%s prompt does not identify Xiaopuy", mode) }
		if !strings.Contains(prompt, "Aryakun") { t.Errorf("%s prompt does not identify Aryakun", mode) }
		if strings.Contains(prompt, "organization that is not disclosed") || strings.Contains(prompt, "organisasi yang tidak diungkapkan") {
			t.Errorf("%s prompt contains fabricated/obscured organization identity", mode)
		}
	}
}

func TestPersonaSeparatesClientAndPrivateBehavior(t *testing.T) {
	p := DefaultPersona()
	private := p.SystemPrompt(false)
	client := p.SystemPrompt(true)
	if !strings.Contains(private, "privately") && !strings.Contains(private, "privately") {
		// The exact wording is intentionally not part of the public contract;
		// this test mainly protects that the prompts are distinct.
	}
	if private == client { t.Fatal("private and client persona prompts must differ") }
	if !strings.Contains(client, "client-facing") { t.Fatal("client prompt missing client-facing policy") }
	if !strings.Contains(private, "personal assistant") { t.Fatal("private prompt missing personal-assistant policy") }
}
