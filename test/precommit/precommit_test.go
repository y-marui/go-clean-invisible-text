// Package precommit_test is an end-to-end fixture for the pre-commit hook
// contract described in docs/integrations/pre-commit.md: check and fix run
// against every staged file in one process, and fix must fail its own run
// when it changes a file so the diff is reviewed before commit.
package precommit_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const zwsp = "​"

func TestPreCommitHooks(t *testing.T) {
	if _, err := exec.LookPath("pre-commit"); err != nil {
		t.Skip("pre-commit not installed; skipping end-to-end fixture")
	}

	bin := buildBinary(t)
	repo := initRepo(t)
	writeConfig(t, repo, bin)

	fileA := filepath.Join(repo, "a.txt")
	fileB := filepath.Join(repo, "b.txt")
	original := "hello" + zwsp + "world\n"
	mustWriteFile(t, fileA, original)
	mustWriteFile(t, fileB, original)

	out, err := runHook(t, repo, "clean-invisible-text-check", fileA, fileB)
	if err == nil {
		t.Fatalf("check: want failure on findings, got success\n%s", out)
	}
	if !strings.Contains(out, "a.txt") || !strings.Contains(out, "b.txt") {
		t.Fatalf("check: want both files reported by one process, got:\n%s", out)
	}
	assertFileContent(t, fileA, original)
	assertFileContent(t, fileB, original)

	out, err = runHook(t, repo, "clean-invisible-text-fix", fileA, fileB)
	if err == nil {
		t.Fatalf("fix: want failure after changing files, got success\n%s", out)
	}
	want := "helloworld\n"
	assertFileContent(t, fileA, want)
	assertFileContent(t, fileB, want)

	if out, err := runHook(t, repo, "clean-invisible-text-fix", fileA, fileB); err != nil {
		t.Fatalf("fix: want success once files are already clean, got failure\n%s", out)
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "clean-invisible-text")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin, "github.com/y-marui/go-clean-invisible-text/cmd/clean-invisible-text")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return bin
}

func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	return dir
}

// writeConfig defines the check and fix hooks as repo: local so the fixture
// exercises the entry/types contract that .pre-commit-hooks.yaml ships,
// without pre-commit needing network access to clone this repository.
func writeConfig(t *testing.T, repo, bin string) {
	t.Helper()
	config := fmt.Sprintf(`repos:
  - repo: local
    hooks:
      - id: clean-invisible-text-check
        name: clean-invisible-text check
        language: system
        entry: %s
        args: [check]
        types: [text]
      - id: clean-invisible-text-fix
        name: clean-invisible-text fix
        language: system
        entry: %s
        args: [fix]
        types: [text]
`, bin, bin)
	mustWriteFile(t, filepath.Join(repo, ".pre-commit-config.yaml"), config)
}

func runHook(t *testing.T, repo, hookID string, files ...string) (string, error) {
	t.Helper()
	args := append([]string{"run", hookID, "--files"}, files...)
	cmd := exec.Command("pre-commit", args...)
	cmd.Dir = repo
	cmd.Env = append(os.Environ(), "PRE_COMMIT_HOME="+t.TempDir())
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("%s: got %q, want %q", path, got, want)
	}
}
