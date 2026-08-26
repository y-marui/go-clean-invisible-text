package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunExplain_DetailedOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dirty.txt")
	if err := os.WriteFile(path, []byte("a\u00A0b"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"explain", path}, nil, &stdout, &stderr)
	if code != exitFindings {
		t.Errorf("exit = %d, want %d", code, exitFindings)
	}
	out := stderr.String()
	for _, want := range []string{path + ":1:2:", "U+00A0", "NO-BREAK SPACE", "[nbsp]", `replace with " "`} {
		if !strings.Contains(out, want) {
			t.Errorf("stderr = %q, want it to contain %q", out, want)
		}
	}
}

func TestRunExplain_WarnAction(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "warn.txt")
	if err := os.WriteFile(path, []byte("a\u2007b"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"explain", path}, nil, &stdout, &stderr)
	if code != exitFindings {
		t.Errorf("exit = %d, want %d", code, exitFindings)
	}
	if !strings.Contains(stderr.String(), "-> warn") {
		t.Errorf("stderr = %q, want it to show the warn action", stderr.String())
	}
}

func TestRunExplain_CleanFileNoOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "clean.txt")
	if err := os.WriteFile(path, []byte("plain\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"explain", path}, nil, &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d", code, exitOK)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunExplain_NoFiles(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"explain"}, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
}
