package main

import (
	"strings"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name  string
		opts  func(t *testing.T) cliOpts
		cause string
	}{
		{
			name:  "plan file not found",
			opts:  func(t *testing.T) cliOpts { return cliOpts{Args: []string{"-p", "nope.json"}} },
			cause: "plan file not found: nope.json",
		},
		{
			name: "plan is not valid JSON",
			opts: func(t *testing.T) cliOpts {
				return cliOpts{Args: []string{"-p", fixturePath(t, "invalid.json")}}
			},
			cause: "failed to parse plan:",
		},
		{
			name: "unsupported format",
			opts: func(t *testing.T) cliOpts {
				return cliOpts{Args: []string{"-p", fixturePath(t, "basic.json"), "-f", "xml"}}
			},
			cause: "invalid format: xml",
		},
		{
			name:  "explicit config file is missing",
			opts:  func(t *testing.T) cliOpts { return cliOpts{Args: []string{"--config", "missing.yaml"}} },
			cause: "failed to load config",
		},
		{
			name:  "no plan at the default path",
			opts:  func(t *testing.T) cliOpts { return cliOpts{} },
			cause: "plan file not found: terraform.tfplan.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := runCLI(t, tt.opts(t))

			if res.Code != 1 {
				t.Errorf("exit code = %d, want 1", res.Code)
			}
			if res.Stdout != "" {
				t.Errorf("stdout = %q, want empty (no partial report)", res.Stdout)
			}
			if !strings.Contains(res.Stderr, tt.cause) {
				t.Errorf("stderr = %q, want it to contain %q", res.Stderr, tt.cause)
			}
		})
	}
}

// PROVISIONAL: pins Known Gap "When the tool fails, the error message is printed twice on standard error,
// followed by the usage text" (specs/001-baseline-cli-behavior/spec.md). Update deliberately when fixed.
func TestKnownGap_ErrorPrintedTwiceWithUsage(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", "nope.json"}})

	cause := "plan file not found: nope.json"
	if n := strings.Count(res.Stderr, cause); n != 2 {
		t.Errorf("the cause appears %d times on stderr, want 2:\n%s", n, res.Stderr)
	}
	if !strings.Contains(res.Stderr, "Usage:") {
		t.Errorf("stderr lacks the usage text:\n%s", res.Stderr)
	}
}

// PROVISIONAL: pins Known Gap "A plan in the exact shape Terraform exports is rejected"
// (specs/001-baseline-cli-behavior/spec.md). Update deliberately when the gap is fixed.
func TestKnownGap_RealPlanShapeRejected(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "real-shape.json")}})

	if res.Code != 1 {
		t.Errorf("exit code = %d, want 1 (the real plan shape is rejected)", res.Code)
	}
	if !strings.Contains(res.Stderr, "failed to parse plan:") {
		t.Errorf("stderr = %q, want a parse failure", res.Stderr)
	}
}
