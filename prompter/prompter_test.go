package prompter

import (
	"bytes"
	"strings"
	"testing"
)

func TestAskPlaintext(t *testing.T) {
	input := strings.NewReader("Alice\n")
	pr := NewWithReader([]Question{
		{Question: "Name", Default: "anon"},
	}, input, nil, nil)

	err := pr.Ask()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pr.Answers[0] != "Alice" {
		t.Fatalf("expected 'Alice', got '%s'", pr.Answers[0])
	}
}

func TestAskDefault(t *testing.T) {
	input := strings.NewReader("\n")
	pr := NewWithReader([]Question{
		{Question: "Name", Default: "anon"},
	}, input, nil, nil)

	err := pr.Ask()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pr.Answers[0] != "anon" {
		t.Fatalf("expected default 'anon', got '%s'", pr.Answers[0])
	}
}

func TestAskSecret(t *testing.T) {
	fakePass := func() ([]byte, error) {
		return []byte("s3cret"), nil
	}

	input := strings.NewReader("")
	pr := NewWithReader([]Question{
		{Question: "Password", Secret: true, Default: ""},
	}, input, nil, fakePass)

	err := pr.Ask()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pr.Answers[0] != "s3cret" {
		t.Fatalf("expected 's3cret', got '%s'", pr.Answers[0])
	}
}

func TestAskMultiple(t *testing.T) {
	input := strings.NewReader("Bob\n")
	fakePass := func() ([]byte, error) {
		return []byte("hunter2"), nil
	}

	pr := NewWithReader([]Question{
		{Question: "User", Default: "admin"},
		{Question: "Pass", Secret: true},
	}, input, nil, fakePass)

	err := pr.Ask()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pr.Answers[0] != "Bob" {
		t.Fatalf("expected 'Bob', got '%s'", pr.Answers[0])
	}
	if pr.Answers[1] != "hunter2" {
		t.Fatalf("expected 'hunter2', got '%s'", pr.Answers[1])
	}
}

func TestPromptOutput(t *testing.T) {
	input := strings.NewReader("test\n")
	var output bytes.Buffer

	pr := NewWithReader([]Question{
		{Question: "Name", Default: "x"},
	}, input, &output, nil)

	err := pr.Ask()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output.String(), "Name [x]:") {
		t.Fatalf("expected prompt in output, got '%s'", output.String())
	}
}