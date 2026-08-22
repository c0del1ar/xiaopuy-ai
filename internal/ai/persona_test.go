package ai

import (
	"strings"
	"testing"
)

func TestDefaultPersonaIdentity(t *testing.T) {
	p := DefaultPersona()
	if p.Name != "Xiao-Puy" { t.Fatalf("Name = %q, want Xiao-Puy", p.Name) }
	if p.PersonaName != "Pupuy" { t.Fatalf("PersonaName = %q, want Pupuy", p.PersonaName) }
	if p.OwnerName != "Aryakun" { t.Fatalf("OwnerName = %q, want Aryakun", p.OwnerName) }
	if p.Website != "aryakun.id" { t.Fatalf("Website = %q, want aryakun.id", p.Website) }

	for mode, prompt := range map[string]string{"private": p.SystemPrompt(false), "client": p.SystemPrompt(true)} {
		if !strings.Contains(prompt, "Xiao-Puy") { t.Errorf("%s prompt does not identify Xiao-Puy", mode) }
		if !strings.Contains(prompt, "Pupuy") { t.Errorf("%s prompt does not identify Pupuy", mode) }
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
	if private == client { t.Fatal("private and client persona prompts must differ") }
	if !strings.Contains(private, "personal AI assistant") { t.Fatal("private prompt missing personal-assistant identity") }
	if !strings.Contains(private, `"sayang"`) { t.Fatal("private prompt missing affectionate owner behavior") }
	if !strings.Contains(client, "client-facing assistant") { t.Fatal("client prompt missing client-facing identity") }
	if !strings.Contains(client, "Do not flirt with clients") { t.Fatal("client prompt missing client boundary") }
	if strings.Contains(client, `call Aryakun "sayang"`) { t.Fatal("client prompt leaked owner-only affectionate behavior") }
}
