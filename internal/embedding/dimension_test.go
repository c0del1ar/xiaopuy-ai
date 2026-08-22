package embedding

import "testing"

func TestValidateDimension(t *testing.T) {
	vector := make([]float32, 3)
	if err := ValidateDimension(vector, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateDimensionMismatch(t *testing.T) {
	vector := make([]float32, 3)
	if err := ValidateDimension(vector, 4); err == nil {
		t.Fatal("expected dimension mismatch")
	}
}

func TestValidateDimensionUnset(t *testing.T) {
	vector := make([]float32, 3)
	if err := ValidateDimension(vector, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
