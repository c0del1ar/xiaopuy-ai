package ai

import "testing"

func TestPolicy(t *testing.T) {
	policy := Policy{}

	tests := []struct {
		name       string
		input      PolicyInput
		want       ReplyDecision
	}{
		{"empty", PolicyInput{ClientMode: true}, DoNotReply},
		{"normal client", PolicyInput{Message: "Halo, apakah bisa membuat website?", ClientMode: true}, AllowReply},
		{"owner request", PolicyInput{Message: "Saya ingin bicara langsung dengan Arya", ClientMode: true}, EscalateOwner},
		{"payment", PolicyInput{Message: "Kirim nomor rekening untuk transfer", ClientMode: true}, NeedContext},
		{"private", PolicyInput{Message: "Sayang, kamu di mana?", ClientMode: false}, AllowReply},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := policy.Decide(tt.input); got != tt.want {
				t.Fatalf("Decide() = %q, want %q", got, tt.want)
			}
		})
	}
}
