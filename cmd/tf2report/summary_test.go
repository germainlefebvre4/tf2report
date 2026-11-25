package main

import (
	"strings"
	"testing"
)

const noChanges = "No changes. Infrastructure is up-to-date."

func TestSummary_basicPlan(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "basic.json")}})

	if res.Code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", res.Code, res.Stderr)
	}
	if res.Stderr != "" {
		t.Errorf("stderr = %q, want empty", res.Stderr)
	}
	for _, want := range []string{
		"# Terraform Plan Summary",
		"**Total Changes:** 4",
		"| Add | 1 |",
		"| Change | 1 |",
		"| Destroy | 1 |",
		"| Replace | 1 |",
	} {
		if !strings.Contains(res.Stdout, want) {
			t.Errorf("stdout does not contain %q:\n%s", want, res.Stdout)
		}
	}
}

func TestSummary_replaceIsCountedOnce(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "basic.json"), "-a", "replace"}})

	if res.Code != 0 {
		t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
	}
	if !strings.Contains(res.Stdout, "**Total Changes:** 1") || !strings.Contains(res.Stdout, "| Replace | 1 |") {
		t.Errorf("a replacement must count as one replace:\n%s", res.Stdout)
	}
	if strings.Contains(res.Stdout, "| Add |") || strings.Contains(res.Stdout, "| Destroy |") {
		t.Errorf("a replacement must not count as add or destroy:\n%s", res.Stdout)
	}
}

func TestSummary_ignoredEntries(t *testing.T) {
	tests := []struct {
		fixture string
	}{
		{"noop-and-read.json"},
		{"unknown-action.json"},
	}

	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, tt.fixture)}})

			if res.Code != 0 {
				t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
			}
			if !strings.Contains(res.Stdout, "**Total Changes:** 1") {
				t.Errorf("only the real change must be counted:\n%s", res.Stdout)
			}
		})
	}
}

func TestSummary_noChanges(t *testing.T) {
	for _, fixture := range []string{"no-changes.json", "empty-object.json"} {
		t.Run(fixture, func(t *testing.T) {
			res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, fixture)}})

			if res.Code != 0 {
				t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
			}
			if !strings.Contains(res.Stdout, noChanges) {
				t.Errorf("stdout does not state %q:\n%s", noChanges, res.Stdout)
			}
		})
	}
}

// Regular test (not a known gap): the plan path defaults to terraform.tfplan.json (FR-001, US1 scenario 5).
func TestSummary_DefaultPlanPath(t *testing.T) {
	res := runCLI(t, cliOpts{Files: map[string]string{
		"terraform.tfplan.json": fixtureContent(t, "basic.json"),
	}})

	if res.Code != 0 {
		t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
	}
	if !strings.Contains(res.Stdout, "**Total Changes:** 4") {
		t.Errorf("default plan path was not used:\n%s", res.Stdout)
	}
}
