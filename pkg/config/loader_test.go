package config

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const fullConfig = `terraform_plan_path: plans/my.json
output_format: json
verbosity: debug
filters:
  resource_types:
    - aws_instance
  actions:
    - create
    - delete
`

func TestLoader_Load_defaults(t *testing.T) {
	isolate(t)

	cfg, err := NewLoader().Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := NewConfig()
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("Load() = %+v, want defaults %+v", cfg, want)
	}
}

func TestLoader_Load_explicitFile(t *testing.T) {
	workDir, _ := isolate(t)
	path := filepath.Join(workDir, "custom.yaml")
	writeFile(t, path, fullConfig)

	cfg, err := NewLoader().Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.TerraformPlanPath != "plans/my.json" || cfg.OutputFormat != "json" || cfg.Verbosity != "debug" {
		t.Errorf("scalar settings not applied: %+v", cfg)
	}
	if !reflect.DeepEqual(cfg.Filters.ResourceTypes, []string{"aws_instance"}) {
		t.Errorf("ResourceTypes = %v", cfg.Filters.ResourceTypes)
	}
	if !reflect.DeepEqual(cfg.Filters.Actions, []string{"create", "delete"}) {
		t.Errorf("Actions = %v", cfg.Filters.Actions)
	}
}

func TestLoader_Load_errors(t *testing.T) {
	t.Run("missing explicit file", func(t *testing.T) {
		workDir, _ := isolate(t)

		_, err := NewLoader().Load(filepath.Join(workDir, "missing.yaml"))
		if err == nil || !strings.Contains(err.Error(), "failed to open config file") {
			t.Errorf("error = %v, want one containing %q", err, "failed to open config file")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		workDir, _ := isolate(t)
		path := filepath.Join(workDir, "bad.yaml")
		writeFile(t, path, "output_format: [unclosed\n")

		_, err := NewLoader().Load(path)
		if err == nil || !strings.Contains(err.Error(), "failed to read config") {
			t.Errorf("error = %v, want one containing %q", err, "failed to read config")
		}
	})
}

func TestLoader_Load_searchPaths(t *testing.T) {
	t.Run("tf2report.yaml in the working directory", func(t *testing.T) {
		workDir, _ := isolate(t)
		writeFile(t, filepath.Join(workDir, "tf2report.yaml"), "output_format: text\n")

		cfg, err := NewLoader().Load("")
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if cfg.OutputFormat != "text" {
			t.Errorf("OutputFormat = %q, want text", cfg.OutputFormat)
		}
	})

	t.Run("tf2report.yaml under HOME/.config/tf2report", func(t *testing.T) {
		_, homeDir := isolate(t)
		writeFile(t, filepath.Join(homeDir, ".config", "tf2report", "tf2report.yaml"), "output_format: json\n")

		cfg, err := NewLoader().Load("")
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if cfg.OutputFormat != "json" {
			t.Errorf("OutputFormat = %q, want json", cfg.OutputFormat)
		}
	})
}

func TestLoader_Load_environment(t *testing.T) {
	t.Run("variables provide settings", func(t *testing.T) {
		isolate(t)
		t.Setenv("TF2REPORT_OUTPUT_FORMAT", "text")
		t.Setenv("TF2REPORT_VERBOSITY", "debug")
		t.Setenv("TF2REPORT_TERRAFORM_PLAN_PATH", "from-env.json")

		cfg, err := NewLoader().Load("")
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if cfg.OutputFormat != "text" || cfg.Verbosity != "debug" || cfg.TerraformPlanPath != "from-env.json" {
			t.Errorf("environment not applied: %+v", cfg)
		}
	})

	t.Run("environment wins over the file", func(t *testing.T) {
		workDir, _ := isolate(t)
		writeFile(t, filepath.Join(workDir, "tf2report.yaml"), "output_format: text\n")
		t.Setenv("TF2REPORT_OUTPUT_FORMAT", "json")

		cfg, err := NewLoader().Load("")
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if cfg.OutputFormat != "json" {
			t.Errorf("OutputFormat = %q, want json (environment over file)", cfg.OutputFormat)
		}
	})
}

// PROVISIONAL: pins Known Gap "The filter settings can be given through the environment only with a
// dotted variable name" (specs/001-baseline-cli-behavior/spec.md). Update deliberately when fixed.
func TestKnownGap_NestedFilterEnvVariable(t *testing.T) {
	t.Run("dotted name is honored", func(t *testing.T) {
		isolate(t)
		t.Setenv("TF2REPORT_FILTERS.ACTIONS", "delete")

		cfg, err := NewLoader().Load("")
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if !reflect.DeepEqual(cfg.Filters.Actions, []string{"delete"}) {
			t.Errorf("Actions = %v, want [delete]", cfg.Filters.Actions)
		}
	})

	t.Run("underscore name is ignored", func(t *testing.T) {
		isolate(t)
		t.Setenv("TF2REPORT_FILTERS_ACTIONS", "delete")

		cfg, err := NewLoader().Load("")
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if len(cfg.Filters.Actions) != 0 {
			t.Errorf("Actions = %v, want none (the underscore form is ignored)", cfg.Filters.Actions)
		}
	})
}

// PROVISIONAL: pins Known Gap "A file named tf2report (no extension) in the working directory is read
// as the configuration file" (specs/001-baseline-cli-behavior/spec.md). Update deliberately when fixed.
func TestKnownGap_ExtensionlessConfigFile(t *testing.T) {
	workDir, _ := isolate(t)
	writeFile(t, filepath.Join(workDir, "tf2report"), "output_format: text\n")

	cfg, err := NewLoader().Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("OutputFormat = %q, want text (the extensionless file is read as config)", cfg.OutputFormat)
	}
}
