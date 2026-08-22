package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvDoesNotOverrideEnvironment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("XIAOPUY_TEST=value-from-file\nXIAOPUY_EXISTING=from-file\n"), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("XIAOPUY_EXISTING", "from-environment")
	if err := LoadDotEnv(path); err != nil {
		t.Fatal(err)
	}

	if got := os.Getenv("XIAOPUY_TEST"); got != "value-from-file" {
		t.Fatalf("XIAOPUY_TEST = %q, want %q", got, "value-from-file")
	}
	if got := os.Getenv("XIAOPUY_EXISTING"); got != "from-environment" {
		t.Fatalf("XIAOPUY_EXISTING = %q, want %q", got, "from-environment")
	}
}

func TestLoadDotEnvMissingFileIsNotAnError(t *testing.T) {
	if err := LoadDotEnv(filepath.Join(t.TempDir(), "missing.env")); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}
}
