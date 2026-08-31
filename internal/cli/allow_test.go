package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCheck_AllowFlag_PermitsSingleOccurrence(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "icons.txt", []byte("ab"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", "--allow", "codepoint=U+E000;reason=Nerd Font icon glyph", path}, nil, &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}
}

func TestRunCheck_AllowFlag_JSONReportsAllowAction(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "icons.txt", []byte("ab"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", "--json", "--allow", "codepoint=U+E000;reason=Nerd Font icon glyph", path}, nil, &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d", code, exitOK)
	}
	var reports []fileReport
	if err := json.Unmarshal(stdout.Bytes(), &reports); err != nil {
		t.Fatalf("json.Unmarshal: %v (stdout = %s)", err, stdout.String())
	}
	if len(reports) != 1 || len(reports[0].Findings) != 1 {
		t.Fatalf("reports = %+v, want a single report with a single finding", reports)
	}
	f := reports[0].Findings[0]
	if f.Action != "allow" {
		t.Errorf("Action = %q, want %q", f.Action, "allow")
	}
	if f.Reason != "Nerd Font icon glyph" {
		t.Errorf("Reason = %q, want the rule's reason", f.Reason)
	}
}

func TestRunCheck_AllowFlag_RunPastMaxRunStillFails(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "icons.txt", []byte("ab"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", "--allow", "codepoint=U+E000;reason=icons", path}, nil, &stdout, &stderr)
	if code != exitFindings {
		t.Errorf("exit = %d, want %d (a run of 2 exceeds the default MaxRun of 1)", code, exitFindings)
	}
}

func TestRunCheck_AllowFlag_PathScoped(t *testing.T) {
	dir := t.TempDir()
	mdPath := writeTestFile(t, dir, "notes.md", []byte("ab"))
	txtPath := writeTestFile(t, dir, "notes.txt", []byte("ab"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", "--json", "--allow", "codepoint=U+E000;reason=icons;paths=*.md", mdPath, txtPath}, nil, &stdout, &stderr)
	// notes.txt isn't covered by the *.md-scoped rule, so it still fails
	// check even though notes.md (same code point, same content) does not.
	if code != exitFindings {
		t.Errorf("exit = %d, want %d (notes.txt is not covered by the *.md-scoped rule)", code, exitFindings)
	}
	var reports []fileReport
	if err := json.Unmarshal(stdout.Bytes(), &reports); err != nil {
		t.Fatalf("json.Unmarshal: %v (stdout = %s)", err, stdout.String())
	}
	byPath := map[string]fileReport{}
	for _, r := range reports {
		byPath[r.Path] = r
	}
	if len(byPath[mdPath].Findings) != 1 || byPath[mdPath].Findings[0].Action != "allow" {
		t.Errorf("notes.md findings = %+v, want a single allow finding", byPath[mdPath].Findings)
	}
	if len(byPath[txtPath].Findings) != 1 || byPath[txtPath].Findings[0].Action != "warn" {
		t.Errorf("notes.txt findings = %+v, want a single warn finding (rule doesn't cover it)", byPath[txtPath].Findings)
	}
}

func TestRunCheck_AllowFile(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "icons.txt", []byte("ab"))
	allowPath := writeTestFile(t, dir, "allow.json", []byte(`[{"codepoint": "U+E000", "reason": "icons"}]`))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", "--allow-file", allowPath, path}, nil, &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}
}

func TestRunCheck_DefaultAllowFileAutoLoaded(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "icons.txt", []byte("ab"))
	writeTestFile(t, dir, ".clean-invisible-text-allow.json", []byte(`[{"codepoint": "U+E000", "reason": "icons"}]`))

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("Chdir back: %v", err)
		}
	}()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", filepath.Base(path)}, nil, &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}
}

func TestRunCheck_AllowFlag_InvalidRule(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "clean.txt", []byte("plain\n"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"check", "--allow", "codepoint=U+E000", path}, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d (missing required reason)", code, exitError)
	}
}

func TestRunFix_AllowFlag_PreservesAllowedCodepoint(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "icons.txt", []byte("ab c"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"fix", "--allow", "codepoint=U+E000;reason=icons", path}, nil, &stdout, &stderr)
	// NBSP is still Block-classified and gets replaced, so the file changes
	// and fix must still report exitFindings for that reason.
	if code != exitFindings {
		t.Errorf("exit = %d, want %d (stderr: %s)", code, exitFindings, stderr.String())
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "ab c" {
		t.Errorf("file = %q, want the allowed code point preserved and NBSP replaced", got)
	}
}

func TestRunClean_AllowFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("ab")
	code := Run([]string{"clean", "--allow", "codepoint=U+E000;reason=icons"}, stdin, &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}
	if stdout.String() != "ab" {
		t.Errorf("stdout = %q, want the allowed code point preserved", stdout.String())
	}
}

func TestRunClean_AllowFlag_PathScopedRuleNeverApplies(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("ab")
	code := Run([]string{"clean", "--allow", "codepoint=U+E000;reason=icons;paths=*.md"}, stdin, &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d", code, exitOK)
	}
	if stdout.String() != "ab" {
		t.Errorf("stdout = %q, want the code point removed (clean has no path for a path-scoped rule to match)", stdout.String())
	}
}
