package cfmt

import (
	"strings"
	"testing"
)

func TestSprintfBasicFormatting(t *testing.T) {
	out := Sprintf("Hello %s", "world")
	if out != "Hello world" {
		t.Fatalf("expected 'Hello world', got %q", out)
	}
}

func TestSprintfWithUnknownTagPassthrough(t *testing.T) {
	out := Sprintf("%notatag %s", "x")
	// Unknown tags should either pass through (in terminal mode) or be stripped (in non-terminal/NO_COLOR).
	if !strings.Contains(out, "x") {
		t.Fatalf("expected formatted value 'x' in output, got %q", out)
	}
	if strings.Contains(out, "%!s") {
		t.Fatalf("unexpected fmt artefact in output: %q", out)
	}
}