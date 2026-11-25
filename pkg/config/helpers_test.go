package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// isolate moves the test into a fresh working directory and HOME, with no TF2REPORT_* variable.
func isolate(t *testing.T) (workDir, homeDir string) {
	t.Helper()

	workDir = t.TempDir()
	homeDir = t.TempDir()
	t.Chdir(workDir)
	t.Setenv("HOME", homeDir)

	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if strings.HasPrefix(name, "TF2REPORT_") {
			t.Setenv(name, "")
			os.Unsetenv(name)
		}
	}
	return workDir, homeDir
}

// writeFile writes content to path, creating parent directories.
func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
