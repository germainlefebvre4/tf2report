package main

import (
	"encoding/xml"
	"strings"
	"testing"
)

func assertSVG(t *testing.T, doc string) {
	t.Helper()

	dec := xml.NewDecoder(strings.NewReader(doc))
	for {
		if _, err := dec.Token(); err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatalf("stdout is not well-formed XML: %v", err)
		}
	}
	if !strings.HasPrefix(doc, "<svg ") {
		t.Errorf("stdout does not start with <svg: %.40q", doc)
	}
}

func TestCLI_SVG(t *testing.T) {
	plan := fixturePath(t, "diagram/aws-network.json")
	res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "svg"}})
	if res.Code != 0 || res.Stderr != "" {
		t.Fatalf("exit %d, stderr %q", res.Code, res.Stderr)
	}
	assertSVG(t, res.Stdout)

	again := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "svg"}})
	if again.Stdout != res.Stdout {
		t.Errorf("two runs differ")
	}
}

func TestCLI_SVG_Selection(t *testing.T) {
	plan := fixturePath(t, "diagram/aws-network.json")
	tests := []struct {
		name string
		opts cliOpts
	}{
		{"flag", cliOpts{Args: []string{"-p", plan, "-f", "svg"}}},
		{"environment", cliOpts{Args: []string{"-p", plan}, Env: []string{"TF2REPORT_OUTPUT_FORMAT=svg"}}},
		{"config file", cliOpts{Args: []string{"-p", plan, "--config", "c.yaml"}, Files: map[string]string{"c.yaml": "output_format: svg\n"}}},
		{"flag beats environment", cliOpts{Args: []string{"-p", plan, "-f", "svg"}, Env: []string{"TF2REPORT_OUTPUT_FORMAT=text"}}},
		{"environment beats file", cliOpts{Args: []string{"-p", plan, "--config", "c.yaml"}, Env: []string{"TF2REPORT_OUTPUT_FORMAT=svg"},
			Files: map[string]string{"c.yaml": "output_format: text\n"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := runCLI(t, tt.opts)
			if res.Code != 0 {
				t.Fatalf("exit %d: %s", res.Code, res.Stderr)
			}
			assertSVG(t, res.Stdout)
		})
	}

	t.Run("flag beats file", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "json", "--config", "c.yaml"},
			Files: map[string]string{"c.yaml": "output_format: svg\n"}})
		if strings.HasPrefix(res.Stdout, "<svg") {
			t.Errorf("config file overrode the flag")
		}
	})
}

func TestCLI_SVG_EmptySelection(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "terraform/empty-config.json"), "-f", "svg"}})
	if res.Code != 0 {
		t.Fatalf("exit %d: %s", res.Code, res.Stderr)
	}
	assertSVG(t, res.Stdout)
	if !strings.Contains(res.Stdout, "No resources to display") {
		t.Errorf("no empty message")
	}
}

func TestCLI_InvalidFormat(t *testing.T) {
	res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "basic.json"), "-f", "xml"}})
	if res.Code != 1 || res.Stdout != "" || !strings.Contains(res.Stderr, "invalid format: xml (supported: markdown, text, json, svg)") {
		t.Errorf("exit %d stdout %q stderr %q", res.Code, res.Stdout, res.Stderr)
	}
}

func TestCLI_ExistingFormatsUnchanged(t *testing.T) {
	for _, f := range []string{"markdown", "text", "json"} {
		res := runCLI(t, cliOpts{Args: []string{"-p", fixturePath(t, "basic.json"), "-f", f}})
		if res.Code != 0 || res.Stderr != "" || strings.Contains(res.Stdout, "<svg") || res.Stdout == "" {
			t.Errorf("%s: exit %d stderr %q", f, res.Code, res.Stderr)
		}
	}
}
