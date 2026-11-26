package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// exampleSources maps each file in examples/ to the fixture it is a copy of.
var exampleSources = map[string]string{
	"sample-plan.json":                   "multi-provider.json",
	"legacy-format-0.1.json":             "legacy-format.json",
	"replace-create-before-destroy.json": "terraform/create-before-destroy.json",
	"replace-forced.json":                "terraform/forced-replace.json",
	"moved-resources.json":               "terraform/moved.json",
	"imported-resource.json":             "terraform/import.json",
	"no-changes.json":                    "terraform/no-changes.json",
	"empty-config.json":                  "terraform/empty-config.json",
}

func examplePath(t *testing.T, name string) string {
	t.Helper()

	path, err := filepath.Abs(filepath.Join("..", "..", "examples", name))
	if err != nil {
		t.Fatalf("resolve example %s: %v", name, err)
	}
	return path
}

// Every example plan must be listed here, be a faithful copy of its fixture, and be readable.
func TestExamples(t *testing.T) {
	files, err := filepath.Glob(examplePath(t, "*.json"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(files) != len(exampleSources) {
		t.Errorf("examples/ holds %d plans, exampleSources lists %d", len(files), len(exampleSources))
	}

	for _, f := range files {
		name := filepath.Base(f)
		source, ok := exampleSources[name]
		if !ok {
			t.Errorf("%s is not listed in exampleSources", name)
			continue
		}

		t.Run(name, func(t *testing.T) {
			got, err := os.ReadFile(f)
			if err != nil {
				t.Fatalf("read example: %v", err)
			}
			if want := fixtureContent(t, source); string(got) != want {
				t.Errorf("example differs from testdata/plans/%s; copy it again", source)
			}

			for _, format := range []string{"markdown", "text", "json"} {
				res := runCLI(t, cliOpts{Args: []string{"-p", f, "-f", format}})
				if res.Code != 0 || strings.TrimSpace(res.Stdout) == "" {
					t.Errorf("format %s: exit code %d, stdout %q, stderr %q", format, res.Code, res.Stdout, res.Stderr)
				}
			}
		})
	}
}
