package fx_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/grimdork/climate/fx"
)

type stringy struct{}

func (stringy) String() string { return "stringy" }

type sample struct {
	Name  string
	Count int
}

func TestSprintPlaceholders(t *testing.T) {
	got := fx.Sprint("Hello, {}!", "Boss")
	if got != "Hello, Boss!" {
		t.Fatalf("expected placeholder replacement, got %q", got)
	}

	got = fx.Sprint("{}", "one", "two")
	if got != "one two" {
		t.Fatalf("expected extra args to be appended, got %q", got)
	}

	got = fx.Sprint("{} {}", "one")
	if got != "one {}" {
		t.Fatalf("expected missing args to leave placeholder intact, got %q", got)
	}
}

func TestStringerAndErrorSupport(t *testing.T) {
	got := fx.Sprint("{} {}", stringy{}, errors.New("bad"))
	if got != "stringy bad" {
		t.Fatalf("expected Stringer and error output, got %q", got)
	}
}

func TestRenderPlainAndNoColor(t *testing.T) {
	plain := fx.RenderPlain("{red}Error{@} {}", "boom")
	if plain != "Error boom" {
		t.Fatalf("expected ANSI-free output, got %q", plain)
	}

	old := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", old)
	os.Setenv("NO_COLOR", "1")

	got := fx.Render("{green}OK{@}")
	if got != "OK" {
		t.Fatalf("expected NO_COLOR to suppress ANSI, got %q", got)
	}
}

func TestWriterVariants(t *testing.T) {
	var out bytes.Buffer
	fx.Fprint(&out, "{}", "alpha")
	fx.Fprintln(&out, "{}", "beta")
	fx.Flog(&out, "{red}{}{@}", "gamma")
	fx.Flogln(&out, "{blue}{}{@}", "delta")

	got := out.String()
	expected := "alphabeta\ngammadelta\n"
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestMapFormattingAndSorting(t *testing.T) {
	m := map[string]int{"b": 2, "a": 1}
	got := fx.RenderWithOptions(fx.Options{SortMaps: true}, "{}", m)
	if got != "{a=1, b=2}" {
		t.Fatalf("expected sorted map output, got %q", got)
	}

	m2 := map[int]string{2: "two", 1: "one"}
	got = fx.RenderWithOptions(fx.Options{SortMaps: true}, "{}", m2)
	if got != "{1=one, 2=two}" {
		t.Fatalf("expected sorted int-key map output, got %q", got)
	}
}

func TestUnsupportedStructFallback(t *testing.T) {
	got := fx.Sprint("{}", sample{Name: "serf", Count: 3})
	if got != "<unsupported>" {
		t.Fatalf("expected unsupported fallback, got %q", got)
	}
}

func TestTokenRegistry(t *testing.T) {
	fx.AddToken("greet", func(modifier string) (string, bool) {
		if modifier == "boss" {
			return "hello, Boss", true
		}
		if modifier == "" {
			return "hello", true
		}
		return "", false
	})
	defer fx.DeleteToken("greet")

	got := fx.Sprint("{greet} / {greet:boss} / {greet:nope}")
	if got != "hello / hello, Boss / {greet:nope}" {
		t.Fatalf("expected custom token output, got %q", got)
	}
}

func TestAliasRegistryAndConcurrency(t *testing.T) {
	fx.AddAlias("note", "magenta italic")
	defer fx.DeleteAlias("note")

	got := fx.RenderPlain("{note}memo{@}")
	if got != "memo" {
		t.Fatalf("expected alias to resolve, got %q", got)
	}

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fx.AddAlias("temp", "blue")
			_ = fx.Render("{temp}x{@}")
			fx.DeleteAlias("temp")
		}()
	}
	wg.Wait()
}

func TestSprintWithDelims(t *testing.T) {
	got := fx.SprintWithDelims("<", ">", "<red>Hello<@> <>", "world")
	if !strings.Contains(got, "Hello") || !strings.Contains(got, "world") {
		t.Fatalf("expected custom delimiters to work, got %q", got)
	}

	plain := fx.StripANSI(got)
	if plain != "Hello world" {
		t.Fatalf("expected custom-delimiter output to strip to plain text, got %q", plain)
	}
}

func TestDisableColourOption(t *testing.T) {
	got := fx.RenderWithOptions(fx.Options{DisableColour: true}, "{red}warn{@}")
	if got != "warn" {
		t.Fatalf("expected DisableColour to suppress ANSI, got %q", got)
	}
}
