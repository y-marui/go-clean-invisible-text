// Package mutate safely rewrites a file on disk with the output of
// internal/cleaner, per the "File mutation" and "Inputs" sections of
// docs/specification.md.
package mutate

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/y-marui/go-clean-invisible-text/internal/cleaner"
)

// resolve returns the real path to read from and write to for path: path
// itself, unless path is a symbolic link, in which case it's the fully
// resolved link target. Writing must always target the resolved path — an
// atomic rename against a symlink path replaces the symlink itself rather
// than following it, which would silently detach the link from its target.
func resolve(path string, isSymlink bool) (string, error) {
	if !isSymlink {
		return path, nil
	}
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	return target, nil
}

// ErrSymlink is returned when path is a symbolic link and Options.AllowSymlink is false.
var ErrSymlink = errors.New("mutate: symbolic link mutation is rejected by default")

// ErrBinary is returned when the file content looks binary (contains a NUL byte).
var ErrBinary = errors.New("mutate: binary-looking input")

// Options configures File.
type Options struct {
	// AllowSymlink permits mutating a file reached through a symbolic link.
	// The zero value (false) matches the spec's reject-by-default behavior.
	AllowSymlink bool
	// KeepWarnings is forwarded to cleaner.Options.KeepWarnings.
	KeepWarnings bool
}

// Result is the outcome of File.
type Result struct {
	// Changed reports whether the file's content was rewritten.
	Changed bool
	// Findings is what cleaner.Clean found, whether or not the file was changed.
	Findings []cleaner.Finding
}

// File cleans the UTF-8 text file at path in place, per the character policy
// implemented by cleaner.Clean. The file is rewritten only when its cleaned
// content differs from what's on disk, using a same-directory temp file and
// an atomic rename so a failure never leaves partial output in path.
func File(path string, opts Options) (Result, error) {
	link, err := os.Lstat(path)
	if err != nil {
		return Result{}, fmt.Errorf("mutate: %s: %w", path, err)
	}
	isSymlink := link.Mode()&os.ModeSymlink != 0
	if isSymlink && !opts.AllowSymlink {
		return Result{}, fmt.Errorf("%w: %s", ErrSymlink, path)
	}

	realPath, err := resolve(path, isSymlink)
	if err != nil {
		return Result{}, fmt.Errorf("mutate: %s: %w", path, err)
	}

	info, err := os.Stat(realPath)
	if err != nil {
		return Result{}, fmt.Errorf("mutate: %s: %w", path, err)
	}

	data, err := os.ReadFile(realPath)
	if err != nil {
		return Result{}, fmt.Errorf("mutate: %s: %w", path, err)
	}

	if bytes.IndexByte(data, 0) != -1 {
		return Result{}, fmt.Errorf("%w: %s", ErrBinary, path)
	}

	cleaned, err := cleaner.Clean(data, cleaner.Options{KeepWarnings: opts.KeepWarnings})
	if err != nil {
		return Result{}, fmt.Errorf("mutate: %s: %w", path, err)
	}

	if bytes.Equal(cleaned.Cleaned, data) {
		return Result{Changed: false, Findings: cleaned.Findings}, nil
	}

	if err := replace(realPath, info.Mode().Perm(), cleaned.Cleaned); err != nil {
		return Result{}, fmt.Errorf("mutate: %s: %w", path, err)
	}

	return Result{Changed: true, Findings: cleaned.Findings}, nil
}

// replace writes data to path via a same-directory temp file, chmod'd to
// perm, fsync'd, and atomically renamed over path. On any error the temp
// file is removed and path is left untouched.
func replace(path string, perm os.FileMode, data []byte) (err error) {
	dir := filepath.Dir(path)

	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() {
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	if err = tmp.Chmod(perm); err != nil {
		tmp.Close()
		return err
	}
	if _, err = tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}
