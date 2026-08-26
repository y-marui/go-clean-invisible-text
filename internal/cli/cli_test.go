package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_NoCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(nil, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Errorf("stderr = %q, want usage text", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"bogus"}, nil, &stdout, &stderr)
	if code != exitError {
		t.Errorf("exit = %d, want %d", code, exitError)
	}
	if !strings.Contains(stderr.String(), "bogus") {
		t.Errorf("stderr = %q, want it to name the unknown command", stderr.String())
	}
}

func TestRun_Help(t *testing.T) {
	for _, flag := range []string{"-h", "--help", "help"} {
		t.Run(flag, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run([]string{flag}, nil, &stdout, &stderr)
			if code != exitOK {
				t.Errorf("exit = %d, want %d", code, exitOK)
			}
			if !strings.Contains(stdout.String(), "Usage:") {
				t.Errorf("stdout = %q, want usage text", stdout.String())
			}
			if stderr.Len() != 0 {
				t.Errorf("stderr = %q, want empty", stderr.String())
			}
		})
	}
}
