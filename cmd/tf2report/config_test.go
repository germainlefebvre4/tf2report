package main

import (
	"strings"
	"testing"
)

// planFile is the working-directory name used by tests that rely on a configuration file.
const planFile = "plan.json"

func TestConfig_fileInWorkingDirectory(t *testing.T) {
	res := runCLI(t, cliOpts{Files: map[string]string{
		planFile: fixtureContent(t, "basic.json"),
		"tf2report.yaml": "terraform_plan_path: " + planFile + "\n" +
			"output_format: text\n" +
			"filters:\n  actions:\n    - delete\n",
	}})

	if res.Code != 0 {
		t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
	}
	if !strings.HasPrefix(res.Stdout, "TERRAFORM PLAN SUMMARY") {
		t.Errorf("configured format text was not applied:\n%s", res.Stdout)
	}
	if !strings.Contains(res.Stdout, "Total Changes: 1") || !strings.Contains(res.Stdout, "aws_security_group") {
		t.Errorf("configured filter was not applied:\n%s", res.Stdout)
	}
}

func TestConfig_explicitFile(t *testing.T) {
	res := runCLI(t, cliOpts{
		Args: []string{"--config", "custom.yaml"},
		Files: map[string]string{
			planFile:      fixtureContent(t, "basic.json"),
			"custom.yaml": "terraform_plan_path: " + planFile + "\noutput_format: json\n",
		},
	})

	if res.Code != 0 {
		t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
	}
	if !strings.HasPrefix(res.Stdout, "{") {
		t.Errorf("explicit config format json was not applied:\n%s", res.Stdout)
	}
}

func TestConfig_flagsOverrideFile(t *testing.T) {
	res := runCLI(t, cliOpts{
		Args: []string{"-f", "json", "-p", planFile},
		Files: map[string]string{
			planFile:         fixtureContent(t, "basic.json"),
			"tf2report.yaml": "terraform_plan_path: missing.json\noutput_format: text\n",
		},
	})

	if res.Code != 0 {
		t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
	}
	if !strings.HasPrefix(res.Stdout, "{") {
		t.Errorf("-f json must override output_format: text:\n%s", res.Stdout)
	}
}

func TestConfig_environmentVariable(t *testing.T) {
	res := runCLI(t, cliOpts{
		Args: []string{"-p", fixturePath(t, "basic.json")},
		Env:  []string{"TF2REPORT_OUTPUT_FORMAT=json"},
	})

	if res.Code != 0 {
		t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
	}
	if !strings.HasPrefix(res.Stdout, "{") {
		t.Errorf("TF2REPORT_OUTPUT_FORMAT=json was not applied:\n%s", res.Stdout)
	}
}

func TestConfig_fileUnderHome(t *testing.T) {
	res := runCLI(t, cliOpts{
		Args:      []string{"-p", fixturePath(t, "basic.json")},
		HomeFiles: map[string]string{".config/tf2report/tf2report.yaml": "output_format: text\n"},
	})

	if res.Code != 0 {
		t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
	}
	if !strings.HasPrefix(res.Stdout, "TERRAFORM PLAN SUMMARY") {
		t.Errorf("configuration under HOME was not applied:\n%s", res.Stdout)
	}
}

func TestConfig_defaultsWithoutAnyConfiguration(t *testing.T) {
	res := runCLI(t, cliOpts{Files: map[string]string{"terraform.tfplan.json": fixtureContent(t, "basic.json")}})

	if res.Code != 0 {
		t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
	}
	if !strings.HasPrefix(res.Stdout, "# Terraform Plan Summary") || !strings.Contains(res.Stdout, "**Total Changes:** 4") {
		t.Errorf("defaults (markdown, terraform.tfplan.json, no filter) were not applied:\n%s", res.Stdout)
	}
}

// PROVISIONAL: pins Known Gap "A file named tf2report (no extension) in the working directory is read
// as the configuration file" (specs/001-baseline-cli-behavior/spec.md). Update deliberately when fixed.
func TestKnownGap_ExtensionlessConfigFile(t *testing.T) {
	res := runCLI(t, cliOpts{
		Args:  []string{"-p", fixturePath(t, "basic.json")},
		Files: map[string]string{"tf2report": "hello world\n"},
	})

	if res.Code != 1 {
		t.Fatalf("exit code = %d, want 1 (the file is parsed as configuration)", res.Code)
	}
	if res.Stdout != "" {
		t.Errorf("stdout = %q, want empty", res.Stdout)
	}
	if !strings.Contains(res.Stderr, "failed to load config") {
		t.Errorf("stderr = %q, want a configuration loading failure", res.Stderr)
	}
}

// PROVISIONAL: pins Known Gap "The verbosity setting only has an effect through verbose mode"
// (specs/001-baseline-cli-behavior/spec.md). Update deliberately when the gap is fixed.
func TestKnownGap_VerbosityOnlyViaVerboseFlag(t *testing.T) {
	files := map[string]string{"tf2report.yaml": "verbosity: debug\n"}
	plan := fixturePath(t, "basic.json")

	t.Run("verbosity debug in the file prints nothing", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan}, Files: files})

		if res.Code != 0 {
			t.Fatalf("exit code = %d; stderr: %s", res.Code, res.Stderr)
		}
		if res.Stderr != "" {
			t.Errorf("stderr = %q, want empty (verbosity in the file has no effect)", res.Stderr)
		}
	})

	t.Run("the verbose flag prints progress", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-v"}})

		if !strings.Contains(res.Stderr, "Parsed plan with 4 resource changes") {
			t.Errorf("stderr = %q", res.Stderr)
		}
	})
}
