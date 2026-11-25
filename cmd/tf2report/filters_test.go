package main

import (
	"strings"
	"testing"
)

func TestFilters(t *testing.T) {
	plan := fixturePath(t, "basic.json")

	tests := []struct {
		name     string
		args     []string
		total    string
		contains []string
		excludes []string
	}{
		{
			name: "one type", args: []string{"-t", "aws_instance"}, total: "**Total Changes:** 2",
			contains: []string{"`web`", "`app`"}, excludes: []string{"aws_s3_bucket", "aws_security_group"},
		},
		{
			name: "repeated type flag", args: []string{"-t", "aws_s3_bucket", "-t", "aws_security_group"}, total: "**Total Changes:** 2",
			contains: []string{"aws_s3_bucket", "aws_security_group"}, excludes: []string{"aws_instance"},
		},
		{
			name: "comma separated types", args: []string{"-t", "aws_s3_bucket,aws_security_group"}, total: "**Total Changes:** 2",
			contains: []string{"aws_s3_bucket", "aws_security_group"}, excludes: []string{"aws_instance"},
		},
		{
			name: "one action", args: []string{"-a", "create"}, total: "**Total Changes:** 1",
			contains: []string{"aws_s3_bucket"}, excludes: []string{"aws_instance", "aws_security_group"},
		},
		{
			name: "two actions", args: []string{"-a", "delete", "-a", "replace"}, total: "**Total Changes:** 2",
			contains: []string{"aws_security_group", "`app`"}, excludes: []string{"aws_s3_bucket", "`web`"},
		},
		{
			name: "type and action must both match", args: []string{"-t", "aws_instance", "-a", "update"}, total: "**Total Changes:** 1",
			contains: []string{"`web`"}, excludes: []string{"`app`", "aws_s3_bucket"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := runCLI(t, cliOpts{Args: append([]string{"-p", plan}, tt.args...)})

			if res.Code != 0 {
				t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
			}
			if !strings.Contains(res.Stdout, tt.total) {
				t.Errorf("stdout lacks %q:\n%s", tt.total, res.Stdout)
			}
			for _, want := range tt.contains {
				if !strings.Contains(res.Stdout, want) {
					t.Errorf("stdout lacks %q:\n%s", want, res.Stdout)
				}
			}
			for _, unwanted := range tt.excludes {
				if strings.Contains(res.Stdout, unwanted) {
					t.Errorf("stdout contains %q:\n%s", unwanted, res.Stdout)
				}
			}
		})
	}
}

func TestFilters_matchingNothingIsNotAnError(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "basic.json"), "-t", "aws_lambda_function"}})

	if res.Code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr: %s", res.Code, res.Stderr)
	}
	if !strings.Contains(res.Stdout, noChanges) {
		t.Errorf("stdout lacks %q:\n%s", noChanges, res.Stdout)
	}
}

// PROVISIONAL: pins Known Gap "An unknown value for the action filter is accepted silently"
// (specs/001-baseline-cli-behavior/spec.md). Update deliberately when the gap is fixed.
func TestKnownGap_UnknownActionFilter(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "basic.json"), "-a", "bogus"}})

	if res.Code != 0 {
		t.Fatalf("exit code = %d, want 0 (the unknown action is accepted); stderr: %s", res.Code, res.Stderr)
	}
	if !strings.Contains(res.Stdout, noChanges) {
		t.Errorf("an unknown action must yield an empty report:\n%s", res.Stdout)
	}
}
