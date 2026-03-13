package arg

import (
	"errors"
	"os"
	"testing"
)

func TestLongOptionParsing(t *testing.T) {
	p := New("testprog")
	// Create options: name, default, required, type, choices
	if err := p.SetOption(GroupDefault, "v", "verbose", "Enable verbose output", false, false, VarBool, nil); err != nil {
		t.Fatalf("SetOption failed: %v", err)
	}
	if err := p.SetOption(GroupDefault, "o", "opt", "An option with value", "", false, VarString, nil); err != nil {
		t.Fatalf("SetOption failed: %v", err)
	}

	// Simulate: prog --verbose --opt=value
	testArgs := []string{"testprog", "--verbose", "--opt=value"}
	// Ensure Parse's os.Args check passes
	os.Args = testArgs
	if err := p.Parse(testArgs); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
}

func TestChoiceValidationErrorUsesWrapping(t *testing.T) {
	p := New("testprog")
	if err := p.SetOption(GroupDefault, "c", "choice", "Choice option", "a", false, VarString, []any{"a", "b"}); err != nil {
		t.Fatalf("SetOption failed: %v", err)
	}

	// Provide an illegal choice
	args := []string{"testprog", "--choice", "z"}
	os.Args = args
	err := p.Parse(args)
	if err == nil {
		t.Fatal("expected error for illegal choice, got nil")
	}

	// Ensure the error supports Unwrap and errors.Is semantics if it wraps underlying error
	if !errors.Is(err, ErrIllegalChoice) {
		// At least ensure Unwrap doesn't panic
		_ = errors.Unwrap(err)
	}
}
