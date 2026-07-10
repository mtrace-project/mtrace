package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

// CreateTempYAMLFile creates a temporary .mt.yaml file with the given content.
// It automatically cleans up the file when the test finishes.
func CreateTempYAMLFile(t *testing.T, dir string, name string, content string) string {
	t.Helper()

	filePath := filepath.Join(dir, name)
	err := os.WriteFile(filePath, []byte(content), 0o644) // nolint:gosec
	if err != nil {
		t.Fatalf("failed to create temp yaml file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(filePath)
	})

	return filePath
}
