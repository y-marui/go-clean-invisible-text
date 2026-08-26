package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/y-marui/go-clean-invisible-text/internal/cleaner"
)

func TestLocate(t *testing.T) {
	cases := []struct {
		name       string
		data       string
		offset     int
		wantLine   int
		wantColumn int
	}{
		{"start-of-text", "abc", 0, 1, 1},
		{"end-of-text", "abc", 3, 1, 4},
		{"second-line", "a\nb", 2, 2, 1},
		{"crlf-only-lf-counts", "a\r\nb", 3, 2, 1},
		{"multibyte-rune-before-offset", "\u00e9b", 2, 1, 2},
		{"third-line", "a\nb\nc", 4, 3, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			line, col := locate([]byte(tc.data), tc.offset)
			if line != tc.wantLine || col != tc.wantColumn {
				t.Errorf("locate(%q, %d) = (%d, %d), want (%d, %d)",
					tc.data, tc.offset, line, col, tc.wantLine, tc.wantColumn)
			}
		})
	}
}

func TestWorstExit(t *testing.T) {
	cases := []struct {
		a, b, want int
	}{
		{exitOK, exitOK, exitOK},
		{exitOK, exitFindings, exitFindings},
		{exitFindings, exitOK, exitFindings},
		{exitFindings, exitError, exitError},
		{exitError, exitFindings, exitError},
		{exitOK, exitError, exitError},
	}
	for _, tc := range cases {
		if got := worstExit(tc.a, tc.b); got != tc.want {
			t.Errorf("worstExit(%d, %d) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestBuildReportAndJSONRoundTrip(t *testing.T) {
	data := []byte("a\u00A0b")
	res, err := cleaner.Clean(data, cleaner.Options{})
	if err != nil {
		t.Fatalf("cleaner.Clean: %v", err)
	}
	report := buildReport("notes.txt", data, res.Findings, false)

	var buf bytes.Buffer
	if err := writeJSON(&buf, []fileReport{report}); err != nil {
		t.Fatalf("writeJSON: %v", err)
	}

	var decoded []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v (raw: %s)", err, buf.String())
	}
	if len(decoded) != 1 {
		t.Fatalf("decoded reports = %d, want 1", len(decoded))
	}
	obj := decoded[0]
	if obj["path"] != "notes.txt" {
		t.Errorf("path = %v, want notes.txt", obj["path"])
	}
	if obj["changed"] != false {
		t.Errorf("changed = %v, want false", obj["changed"])
	}
	if obj["error"] != nil {
		t.Errorf("error = %v, want nil", obj["error"])
	}
	findings, ok := obj["findings"].([]any)
	if !ok || len(findings) != 1 {
		t.Fatalf("findings = %v, want a single-element array", obj["findings"])
	}
	f := findings[0].(map[string]any)
	for _, field := range []string{"line", "column", "offset", "rune", "name", "category", "action", "replacement"} {
		if _, ok := f[field]; !ok {
			t.Errorf("finding missing field %q: %v", field, f)
		}
	}
	if f["rune"] != "U+00A0" {
		t.Errorf("rune = %v, want U+00A0", f["rune"])
	}
	if f["category"] != "nbsp" {
		t.Errorf("category = %v, want nbsp", f["category"])
	}
	if f["action"] != "replace" {
		t.Errorf("action = %v, want replace", f["action"])
	}
}

func TestErrorReport(t *testing.T) {
	r := errorReport("bad.bin", mustErr("boom"))
	if r.Error == nil || *r.Error != "boom" {
		t.Errorf("Error = %v, want \"boom\"", r.Error)
	}
	if len(r.Findings) != 0 {
		t.Errorf("Findings = %v, want empty", r.Findings)
	}
	if r.Changed {
		t.Errorf("Changed = true, want false")
	}
}

type simpleErr string

func (e simpleErr) Error() string { return string(e) }

func mustErr(msg string) error { return simpleErr(msg) }
