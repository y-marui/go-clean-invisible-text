package mutate

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/y-marui/go-clean-invisible-text/internal/cleaner"
)

func writeFile(t *testing.T, dir, name string, content []byte, perm os.FileMode) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, perm); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func TestFile_CleanWrite(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "dirty.txt", []byte("a\u00A0b"), 0o644)

	res, err := File(path, Options{})
	if err != nil {
		t.Fatalf("File: %v", err)
	}
	if !res.Changed {
		t.Fatalf("Changed = false, want true")
	}
	if len(res.Findings) != 1 {
		t.Fatalf("Findings = %v, want exactly 1", res.Findings)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "a b" {
		t.Errorf("file content = %q, want %q", got, "a b")
	}
}

func TestFile_NoOpWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "clean.txt", []byte("already clean\n"), 0o644)

	before, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}

	// Ensure any rewrite would be detectable via mtime even on coarse
	// filesystem timestamp resolution.
	time.Sleep(10 * time.Millisecond)

	res, err := File(path, Options{})
	if err != nil {
		t.Fatalf("File: %v", err)
	}
	if res.Changed {
		t.Fatalf("Changed = true, want false")
	}
	if len(res.Findings) != 0 {
		t.Errorf("Findings = %v, want none", res.Findings)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("mtime changed: before=%v after=%v", before.ModTime(), after.ModTime())
	}
}

func TestFile_PreservesPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are not meaningful on windows")
	}

	dir := t.TempDir()
	path := writeFile(t, dir, "perm.txt", []byte("a\u00A0b"), 0o640)

	if _, err := File(path, Options{}); err != nil {
		t.Fatalf("File: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Errorf("mode = %v, want %v", info.Mode().Perm(), os.FileMode(0o640))
	}
}

func TestFile_PreservesCRLFAndTrailingNewline(t *testing.T) {
	dir := t.TempDir()
	original := "line one\r\nline\u00A0two\r\n"
	path := writeFile(t, dir, "crlf.txt", []byte(original), 0o644)

	if _, err := File(path, Options{}); err != nil {
		t.Fatalf("File: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	want := "line one\r\nline two\r\n"
	if string(got) != want {
		t.Errorf("file content = %q, want %q", got, want)
	}
}

func TestFile_SymlinkRejectedByDefault(t *testing.T) {
	dir := t.TempDir()
	target := writeFile(t, dir, "target.txt", []byte("a\u00A0b"), 0o644)
	link := filepath.Join(dir, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	_, err := File(link, Options{})
	if !errors.Is(err, ErrSymlink) {
		t.Fatalf("err = %v, want ErrSymlink", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "a\u00A0b" {
		t.Errorf("target was modified: %q", got)
	}
}

func TestFile_SymlinkAllowed(t *testing.T) {
	dir := t.TempDir()
	target := writeFile(t, dir, "target.txt", []byte("a\u00A0b"), 0o644)
	link := filepath.Join(dir, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	res, err := File(link, Options{AllowSymlink: true})
	if err != nil {
		t.Fatalf("File: %v", err)
	}
	if !res.Changed {
		t.Fatalf("Changed = false, want true")
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "a b" {
		t.Errorf("target content = %q, want %q", got, "a b")
	}
}

func TestFile_BinaryRejected(t *testing.T) {
	dir := t.TempDir()
	content := []byte("a\x00b")
	path := writeFile(t, dir, "binary.dat", content, 0o644)

	_, err := File(path, Options{})
	if !errors.Is(err, ErrBinary) {
		t.Fatalf("err = %v, want ErrBinary", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("file was modified: %q", got)
	}
}

func TestFile_InvalidUTF8Rejected(t *testing.T) {
	dir := t.TempDir()
	content := []byte{0xff, 0xfe, 0x41, 0x42}
	path := writeFile(t, dir, "invalid.txt", content, 0o644)

	_, err := File(path, Options{})
	if !errors.Is(err, cleaner.ErrInvalidUTF8) {
		t.Fatalf("err = %v, want cleaner.ErrInvalidUTF8", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("file was modified: %q", got)
	}
}

func TestFile_NoStrayTempFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "dirty.txt", []byte("a\u00A0b"), 0o644)

	if _, err := File(filepath.Join(dir, "dirty.txt"), Options{}); err != nil {
		t.Fatalf("File: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "dirty.txt" {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("directory entries = %v, want only [dirty.txt]", names)
	}
}
