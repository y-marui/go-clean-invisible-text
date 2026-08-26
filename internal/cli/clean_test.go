package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunClean_StdinToStdout(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"clean"}, strings.NewReader("a\u00A0b"), &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d (clean always exits 0 on success)", code, exitOK)
	}
	if stdout.String() != "a b" {
		t.Errorf("stdout = %q, want %q", stdout.String(), "a b")
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunClean_NoFindingsAlsoExitOK(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"clean"}, strings.NewReader("plain text"), &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d", code, exitOK)
	}
	if stdout.String() != "plain text" {
		t.Errorf("stdout = %q, want unchanged input", stdout.String())
	}
}

func TestRunClean_InvalidUTF8(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"clean"}, bytes.NewReader([]byte{0xff, 0xfe, 0x41}), &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want empty on error", stdout.String())
	}
}

func TestRunClean_KeepWarnings(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"clean", "--keep-warnings"}, strings.NewReader("a\u2007b"), &stdout, &stderr)
	if code != exitOK {
		t.Errorf("exit = %d, want %d", code, exitOK)
	}
	if stdout.String() != "a\u2007b" {
		t.Errorf("stdout = %q, want unchanged input (Warn byte kept)", stdout.String())
	}
}

func TestRunClean_UnexpectedArgument(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"clean", "somefile.txt"}, strings.NewReader(""), &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
}
