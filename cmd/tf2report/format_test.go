package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormat_defaultIsMarkdown(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "basic.json")}})

	if !strings.HasPrefix(res.Stdout, "# Terraform Plan Summary") {
		t.Errorf("default output is not Markdown:\n%s", res.Stdout)
	}
}

func TestFormat_namesAndAliases(t *testing.T) {
	plan := fixturePath(t, "basic.json")
	run := func(format string) cliResult {
		return runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", format}})
	}

	markdown, text := run("markdown"), run("text")

	if got := run("md"); got.Stdout != markdown.Stdout || got.Code != 0 {
		t.Errorf("-f md differs from -f markdown")
	}
	if got := run("txt"); got.Stdout != text.Stdout || got.Code != 0 {
		t.Errorf("-f txt differs from -f text")
	}
	if strings.Contains(text.Stdout, "|") || !strings.HasPrefix(text.Stdout, "TERRAFORM PLAN SUMMARY") {
		t.Errorf("-f text is not plain text:\n%s", text.Stdout)
	}

	jsonRes := run("json")
	var report map[string]any
	if err := json.Unmarshal([]byte(jsonRes.Stdout), &report); err != nil {
		t.Fatalf("-f json output is not valid JSON: %v\n%s", err, jsonRes.Stdout)
	}
}

func TestFormat_sameCountsInEveryFormat(t *testing.T) {
	plan := fixturePath(t, "basic.json")

	markdown := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "markdown"}}).Stdout
	for _, want := range []string{"**Total Changes:** 4", "| Add | 1 |", "| Change | 1 |", "| Destroy | 1 |", "| Replace | 1 |"} {
		if !strings.Contains(markdown, want) {
			t.Errorf("markdown lacks %q", want)
		}
	}

	text := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "text"}}).Stdout
	for _, want := range []string{"Total Changes: 4", "Add:     1", "Change:  1", "Destroy: 1", "Replace: 1"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q", want)
		}
	}

	var report struct {
		Summary map[string]int `json:"summary"`
	}
	if err := json.Unmarshal([]byte(runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "json"}}).Stdout), &report); err != nil {
		t.Fatalf("json: %v", err)
	}
	want := map[string]int{"total_changes": 4, "to_add": 1, "to_change": 1, "to_destroy": 1, "to_replace": 1}
	for key, n := range want {
		if report.Summary[key] != n {
			t.Errorf("json summary[%s] = %d, want %d", key, report.Summary[key], n)
		}
	}
}

func TestFormat_unsupported(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "basic.json"), "-f", "xml"}})

	if res.Code != 1 {
		t.Errorf("exit code = %d, want 1", res.Code)
	}
	if res.Stdout != "" {
		t.Errorf("stdout = %q, want empty", res.Stdout)
	}
	if !strings.Contains(res.Stderr, "invalid format: xml (supported: markdown, text, json, svg)") {
		t.Errorf("stderr = %q", res.Stderr)
	}
}

func TestFormat_neverPrintsAttributeValues(t *testing.T) {
	markers := []string{"SYNTHETIC-SECRET-PASSWORD", "SYNTHETIC-SECRET-TOKEN", "SYNTHETIC-SECRET-OLD"}

	for _, format := range []string{"markdown", "text", "json"} {
		t.Run(format, func(t *testing.T) {
			res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "sensitive.json"), "-f", format}})

			if res.Code != 0 {
				t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
			}
			for _, marker := range markers {
				if strings.Contains(res.Stdout, marker) || strings.Contains(res.Stderr, marker) {
					t.Errorf("output contains attribute value %q", marker)
				}
			}
		})
	}
}
