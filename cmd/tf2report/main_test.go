package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

type cliOpts struct {
	Args      []string
	Files     map[string]string // written into the working directory
	HomeFiles map[string]string // written into HOME
	Env       []string          // extra KEY=VALUE entries
}

type cliResult struct {
	Stdout string
	Stderr string
	Code   int
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	dir, err := os.MkdirTemp("", "tf2report-bin-")
	if err != nil {
		fmt.Fprintln(os.Stderr, "create temp dir:", err)
		return 1
	}
	defer os.RemoveAll(dir)

	binaryPath = filepath.Join(dir, "tf2report")
	ldflags := "-X main.appVersion=test-1.2.3 -X main.buildCommit=abc1234 -X main.buildTime=2026-01-01T00:00:00Z"
	build := exec.Command("go", "build", "-ldflags", ldflags, "-o", binaryPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build binary: %v\n%s", err, out)
		return 1
	}

	return m.Run()
}

// runCLI runs the built binary in an isolated working directory and HOME.
func runCLI(t *testing.T, o cliOpts) cliResult {
	t.Helper()

	workDir := t.TempDir()
	homeDir := t.TempDir()
	writeFiles(t, workDir, o.Files)
	writeFiles(t, homeDir, o.HomeFiles)

	var env []string
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if strings.HasPrefix(name, "TF2REPORT_") || name == "HOME" {
			continue
		}
		env = append(env, kv)
	}
	env = append(env, "HOME="+homeDir)
	env = append(env, o.Env...)

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(binaryPath, o.Args...)
	cmd.Dir = workDir
	cmd.Env = env
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	code := 0
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("run binary: %v", err)
		}
		code = exitErr.ExitCode()
	}
	return cliResult{Stdout: stdout.String(), Stderr: stderr.String(), Code: code}
}

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()

	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", name, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
}

// fixturePath returns the absolute path of a shared plan fixture.
func fixturePath(t *testing.T, name string) string {
	t.Helper()

	path, err := filepath.Abs(filepath.Join("..", "..", "testdata", "plans", name))
	if err != nil {
		t.Fatalf("resolve fixture %s: %v", name, err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("fixture %s: %v", name, err)
	}
	return path
}

// fixtureContent returns the content of a shared plan fixture.
func fixtureContent(t *testing.T, name string) string {
	t.Helper()

	data, err := os.ReadFile(fixturePath(t, name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(data)
}
