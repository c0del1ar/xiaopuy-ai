package embedding

import "fmt"

// ValidateDimension ensures a provider returns vectors with the configured size.
// The dimension is a deployment property because pgvector indexes require a fixed dimension.
func ValidateDimension(vector []float32, expected int) error {
	if len(vector) == 0 {
		return fmt.Errorf("embedding vector is empty")
	}
	if expected > 0 && len(vector) != expected {
		return fmt.Errorf("embedding dimension = %d, want %d", len(vector), expected)
	}
	return nil
}
