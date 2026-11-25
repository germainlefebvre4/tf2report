package main

import (
	"strings"
	"testing"
)

func TestDiagnostics_verbose(t *testing.T) {
	plan := fixturePath(t, "basic.json")
	quiet := runCLI(t, cliOpts{Args: []string{"-p", plan}})

	t.Run("quiet by default", func(t *testing.T) {
		if quiet.Stderr != "" {
			t.Errorf("stderr = %q, want empty", quiet.Stderr)
		}
	})

	t.Run("progress goes to stderr and stdout is unchanged", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-v"}})

		if !strings.Contains(res.Stderr, "Parsed plan with 4 resource changes") {
			t.Errorf("stderr = %q", res.Stderr)
		}
		if res.Stdout != quiet.Stdout {
			t.Errorf("-v changed the report on stdout")
		}
	})

	t.Run("type filter progress", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-v", "-t", "aws_instance"}})

		if !strings.Contains(res.Stderr, "Filtered by type: 2 resources remain") {
			t.Errorf("stderr = %q", res.Stderr)
		}
	})

	t.Run("action filter progress", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-v", "-a", "create"}})

		if !strings.Contains(res.Stderr, "Filtered by action: 1 resources remain") {
			t.Errorf("stderr = %q", res.Stderr)
		}
	})
}

func TestDiagnostics_help(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"--help"}})

	if res.Code != 0 {
		t.Errorf("exit code = %d, want 0", res.Code)
	}
	if !strings.Contains(res.Stdout, "Usage:") {
		t.Errorf("stdout lacks the usage text:\n%s", res.Stdout)
	}
	if res.Stderr != "" {
		t.Errorf("stderr = %q, want empty", res.Stderr)
	}
}

func TestDiagnostics_version(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"--version"}})

	if res.Code != 0 {
		t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
	}
	want := "tf2report\nVersion: test-1.2.3\nCommit: abc1234\nBuild time: 2026-01-01T00:00:00Z\n"
	if res.Stdout != want {
		t.Errorf("stdout = %q, want %q", res.Stdout, want)
	}
}
