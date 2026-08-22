package presence

import "testing"

func TestStore(t *testing.T) {
	store := NewStore(Offline)
	if got := store.Get(); got != Offline {
		t.Fatalf("Get() = %q, want %q", got, Offline)
	}
	if !store.Set(Online) {
		t.Fatal("Set(Online) returned false")
	}
	if got := store.Get(); got != Online {
		t.Fatalf("Get() = %q, want %q", got, Online)
	}
	if store.Set(Status("invalid")) {
		t.Fatal("Set(invalid) returned true")
	}
	if got := store.Get(); got != Online {
		t.Fatalf("invalid Set changed status to %q", got)
	}
}

func TestMayAutoReply(t *testing.T) {
	for _, status := range []Status{Offline, AIOnly} {
		if !MayAutoReply(status) {
			t.Errorf("MayAutoReply(%q) = false, want true", status)
		}
	}
	for _, status := range []Status{Online, Busy, Away} {
		if MayAutoReply(status) {
			t.Errorf("MayAutoReply(%q) = true, want false", status)
		}
	}
}
