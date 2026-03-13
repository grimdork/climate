package cfmt_test

import (
	"github.com/grimdork/climate/cfmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNO_COLOREnvDisablesEscapes(t *testing.T) {
	old := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", old)
	os.Setenv("NO_COLOR", "1")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe failed: %v", err)
	}
	defer r.Close()
	defer w.Close()

	// Capture stdout
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	// Use Printf which previously processed formatting; ensure no escape sequences appear
	cfmt.Printf("%red Error:%reset file not found: %s", "fn")
	w.Close()

	outBytes, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	out := string(outBytes)

	if strings.Contains(out, "\x1b[") {
		t.Errorf("expected no ANSI escape codes when NO_COLOR set, got: %q", out)
	}
	if !strings.Contains(out, "Error:") || !strings.Contains(out, "file not found") {
		t.Errorf("expected plain message text in output, got: %q", out)
	}
}