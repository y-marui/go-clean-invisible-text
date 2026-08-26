package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunFix_RewritesDirtyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dirty.txt")
	if err := os.WriteFile(path, []byte("a\u00A0b"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"fix", path}, nil, &stdout, &stderr)
	if code != exitFindings {
		t.Errorf("exit = %d, want %d", code, exitFindings)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "a b" {
		t.Errorf("file content = %q, want %q", got, "a b")
	}
	if !strings.Contains(stderr.String(), path) {
		t.Errorf("stderr = %q, want it to mention the change", stderr.String())
	}
}

func TestRunFix_CleanFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "clean.txt")
	if err := os.WriteFile(path, []byte("plain\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"fix", path}, nil, &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d", code, exitOK)
	}
}

func TestRunFix_KeepWarnings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "warn.txt")
	if err := os.WriteFile(path, []byte("a\u2007b"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"fix", "--keep-warnings", path}, nil, &stdout, &stderr)
	// The Finding is still reported (exitFindings), even though the byte
	// itself survives in the file.
	if code != exitOK {
		t.Errorf("exit = %d, want %d (nothing removed, so Changed should be false)", code, exitOK)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "a\u2007b" {
		t.Errorf("file content = %q, want unchanged %q", got, "a\u2007b")
	}
}

func TestRunFix_SymlinkRejected(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	if err := os.WriteFile(target, []byte("a\u00A0b"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	link := filepath.Join(dir, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"fix", link}, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
	if !strings.Contains(stderr.String(), "symbolic link") {
		t.Errorf("stderr = %q, want it to mention the symlink rejection", stderr.String())
	}
}

func TestRunFix_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dirty.txt")
	if err := os.WriteFile(path, []byte("a\u00A0b"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"fix", "--json", path}, nil, &stdout, &stderr)
	if code != exitFindings {
		t.Errorf("exit = %d, want %d", code, exitFindings)
	}
	if !strings.Contains(stdout.String(), `"changed": true`) {
		t.Errorf("stdout = %q, want JSON with changed: true", stdout.String())
	}
}

func TestRunFix_NoFiles(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"fix"}, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
}
