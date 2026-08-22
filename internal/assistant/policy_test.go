package assistant

import "testing"

func TestRetrievalPolicy(t *testing.T) {
	p := RetrievalPolicy{MinQueryChars: 8, Limit: 7}
	if p.ShouldRetrieve("halo") {
		t.Fatal("short greeting should not trigger retrieval")
	}
	if !p.ShouldRetrieve("berapa harga website?") {
		t.Fatal("meaningful query should trigger retrieval")
	}
	if p.TopK() != 7 {
		t.Fatalf("TopK = %d, want 7", p.TopK())
	}
}

func TestRetrievalPolicyDefaults(t *testing.T) {
	p := RetrievalPolicy{}
	if p.ShouldRetrieve("1234567") {
		t.Fatal("query shorter than default threshold should not retrieve")
	}
	if !p.ShouldRetrieve("12345678") {
		t.Fatal("query at default threshold should retrieve")
	}
	if p.TopK() != 5 {
		t.Fatalf("default TopK = %d, want 5", p.TopK())
	}
}
