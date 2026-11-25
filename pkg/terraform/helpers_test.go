package terraform

import (
	"os"
	"path/filepath"
	"testing"
)

// fixture reads a shared synthetic plan from testdata/plans.
func fixture(t *testing.T, name string) []byte {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "plans", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}
