package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func iconsPath(t *testing.T, name string) string {
	t.Helper()

	p, err := filepath.Abs(filepath.Join("..", "..", "testdata", "icons", name))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCLI_IconsDir(t *testing.T) {
	plan := fixturePath(t, "diagram/aws-network.json")

	t.Run("valid pack", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "svg", "--icons-dir", iconsPath(t, "valid")}})
		if res.Code != 0 || res.Stderr != "" || !strings.Contains(res.Stdout, `id="icon-aws-ec2"`) {
			t.Errorf("exit %d stderr %q", res.Code, res.Stderr)
		}
	})

	t.Run("missing directory", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "svg", "--icons-dir", "/nonexistent"}})
		if res.Code == 0 || res.Stdout != "" || !strings.Contains(res.Stderr, "/nonexistent") {
			t.Errorf("exit %d stdout %q stderr %q", res.Code, res.Stdout, res.Stderr)
		}
	})

	t.Run("ignored for json", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "json", "--icons-dir", "/nonexistent"}})
		if res.Code != 0 || res.Stderr != "" || !strings.HasPrefix(strings.TrimSpace(res.Stdout), "{") {
			t.Errorf("exit %d stderr %q", res.Code, res.Stderr)
		}
	})

	t.Run("ignored for json verbose", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "json", "-v", "--icons-dir", "/nonexistent"}})
		if res.Code != 0 || strings.Count(res.Stderr, "--icons-dir is ignored") != 1 {
			t.Errorf("exit %d stderr %q", res.Code, res.Stderr)
		}
	})

	t.Run("hostile pack verbose lists skipped files", func(t *testing.T) {
		res := runCLI(t, cliOpts{Args: []string{"-p", plan, "-f", "svg", "-v", "--icons-dir", iconsPath(t, "hostile")}})
		if res.Code != 0 {
			t.Fatalf("exit %d: %s", res.Code, res.Stderr)
		}
		for _, f := range []string{"script.svg", "onload.svg", "foreign.svg", "doctype.svg", "malformed.svg"} {
			if !strings.Contains(res.Stderr, "skipped icon "+f) {
				t.Errorf("stderr does not list %s", f)
			}
		}
		assertSVG(t, res.Stdout)
	})
}
