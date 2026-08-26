package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func TestRunCheck_CleanFile(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "clean.txt", []byte("plain text\n"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", path}, nil, &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
}

func TestRunCheck_BlockFinding(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "dirty.txt", []byte("a\u00A0b"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", path}, nil, &stdout, &stderr)
	if code != exitFindings {
		t.Errorf("exit = %d, want %d", code, exitFindings)
	}
	if !strings.Contains(stderr.String(), path) || !strings.Contains(stderr.String(), "1 finding") {
		t.Errorf("stderr = %q, want it to mention the path and finding count", stderr.String())
	}

	// check must not modify the file.
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "a\u00A0b" {
		t.Errorf("file was modified: %q", got)
	}
}

func TestRunCheck_WarnFinding(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "warn.txt", []byte("a\u2007b"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", path}, nil, &stdout, &stderr)
	if code != exitFindings {
		t.Errorf("exit = %d, want %d", code, exitFindings)
	}
	if !strings.Contains(stderr.String(), "1 finding") {
		t.Errorf("stderr = %q, want it to mention 1 finding", stderr.String())
	}
}

func TestRunCheck_InvalidUTF8(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "invalid.txt", []byte{0xff, 0xfe, 0x41})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", path}, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
	if !strings.Contains(stderr.String(), path) {
		t.Errorf("stderr = %q, want it to mention the path", stderr.String())
	}
}

func TestRunCheck_Binary(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "binary.dat", []byte("a\x00b"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", path}, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
}

func TestRunCheck_MultipleFilesWorstExit(t *testing.T) {
	dir := t.TempDir()
	clean := writeTestFile(t, dir, "clean.txt", []byte("plain\n"))
	dirty := writeTestFile(t, dir, "dirty.txt", []byte("a\u00A0b"))
	invalid := writeTestFile(t, dir, "invalid.txt", []byte{0xff, 0xfe})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", clean, dirty, invalid}, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d (worst of ok/findings/error)", code, exitError)
	}
}

func TestRunCheck_JSON(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "dirty.txt", []byte("a\u00A0b"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", "--json", path}, nil, &stdout, &stderr)
	if code != exitFindings {
		t.Errorf("exit = %d, want %d", code, exitFindings)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty in --json mode", stderr.String())
	}
	if !strings.Contains(stdout.String(), "\"path\"") || !strings.Contains(stdout.String(), path) {
		t.Errorf("stdout = %q, want JSON containing the path", stdout.String())
	}
}

func TestRunCheck_NoFiles(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"check"}, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
}
