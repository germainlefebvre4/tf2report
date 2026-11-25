package report

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

var update = flag.Bool("update", false, "rewrite golden files")

// loadPlan parses a shared fixture and builds its summary.
func loadPlan(t *testing.T, name string) (*terraform.Plan, *terraform.Summary) {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "plans", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	plan, err := terraform.NewParser().Parse(data)
	if err != nil {
		t.Fatalf("parse fixture %s: %v", name, err)
	}
	return plan, terraform.NewSummary(plan.ResourceChanges)
}

// render runs a formatter into a buffer.
func render(t *testing.T, f Formatter, plan *terraform.Plan, summary *terraform.Summary) string {
	t.Helper()

	var buf bytes.Buffer
	if err := f.Format(plan, summary, &buf); err != nil {
		t.Fatalf("format: %v", err)
	}
	return buf.String()
}

// assertGolden compares got with testdata/<name>, rewriting it only with -update.
func assertGolden(t *testing.T, name, got string) {
	t.Helper()

	path := filepath.Join("testdata", name)
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatalf("create testdata: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden %s: %v", name, err)
		}
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s (run with -update to create it): %v", name, err)
	}
	if got != string(want) {
		t.Errorf("output differs from golden %s\n--- got ---\n%s\n--- want ---\n%s", name, got, want)
	}
}

// assertDeterministic renders twice and requires identical bytes (FR-014).
func assertDeterministic(t *testing.T, f Formatter, plan *terraform.Plan, summary *terraform.Summary) string {
	t.Helper()

	first := render(t, f, plan, summary)
	second := render(t, f, plan, summary)
	if first != second {
		t.Errorf("output is not deterministic across two renders")
	}
	return first
}
