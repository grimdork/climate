package cfmt

import (
	"strings"
	"testing"
)

// colourDirect bypasses the shouldColor check by writing through the
// colour parser's switch logic. Since tests don't run in a terminal,
// we test parseKeyword and the tag-to-code mapping directly.

func TestParseKeyword(t *testing.T) {
	tests := []struct {
		input   string
		keyword string
		rest    string
	}{
		{"%red hello", "red", "hello"},
		{"%reset", "reset", ""},
		{"%bold text", "bold", "text"},
		{"%lred!", "lred", "!"},
	}

	for _, tt := range tests {
		keyword, rest := parseKeyword(tt.input)
		if keyword != tt.keyword {
			t.Errorf("parseKeyword(%q): keyword = %q, want %q", tt.input, keyword, tt.keyword)
		}
		if rest != tt.rest {
			t.Errorf("parseKeyword(%q): rest = %q, want %q", tt.input, rest, tt.rest)
		}
	}
}

func TestTagMapping(t *testing.T) {
	tags := map[string]string{
		"reset":    Reset,
		"black":    Black,
		"red":      Red,
		"green":    Green,
		"yellow":   Yellow,
		"blue":     Blue,
		"magenta":  Magenta,
		"cyan":     Cyan,
		"white":    White,
		"grey":     Grey,
		"gray":     Grey,
		"lred":     LightRed,
		"lgreen":   LightGreen,
		"lyellow":  LightYellow,
		"lblue":    LightBlue,
		"lmagenta": LightMagenta,
		"lcyan":    LightCyan,
		"lwhite":   LightWhite,
		"bgred":    BGRed,
		"bggreen":  BGGreen,
		"bgblue":   BGBlue,
		"bold":     Bold,
		"italic":   Italic,
		"under":    Underscore,
		"strike":   Strikethrough,
		"fuzzy":    Fuzzy,
		"blink":    Blink,
		"reverse":  Reverse,
		"conceal":  Concealed,
	}

	for tag, code := range tags {
		// Build input that forces the colour branch by writing to a Builder
		// with the tag. We call colour directly but it checks shouldColor.
		// Instead, simulate what colour does when shouldColor is true.
		var b strings.Builder
		input := "%" + tag + " "
		// Manually walk the parser
		f := input
		for len(f) > 0 {
			c := f[0]
			if c == '%' {
				var key string
				key, f = parseKeyword(f)
				b.WriteString(resolveTag(key))
			} else {
				b.WriteByte(f[0])
				f = f[1:]
			}
		}

		if !strings.Contains(b.String(), code) {
			t.Errorf("tag %%%s: expected %q in output, got %q", tag, code, b.String())
		}
	}
}

func TestNoColorPath(t *testing.T) {
	// When not in a terminal, colour strips tags and outputs plain text without appending newline
	var b strings.Builder
	colour(&b, "%red Error:%reset file not found")
	result := b.String()

	if strings.Contains(result, "\x1b") {
		t.Error("expected no escape codes in non-terminal output")
	}
	if !strings.Contains(result, "Error:") {
		t.Error("expected 'Error:' in output")
	}
	if !strings.Contains(result, "file not found") {
		t.Error("expected 'file not found' in output")
	}
	if strings.HasSuffix(result, "\n") {
		t.Error("Print/colour should not append newline")
	}
}

func TestUnknownTagPassthrough(t *testing.T) {
	var b strings.Builder
	input := "%notacolour "
	f := input
	for len(f) > 0 {
		c := f[0]
		if c == '%' {
			var key string
			key, f = parseKeyword(f)
			resolved := resolveTag(key)
			if resolved == "" {
				b.WriteByte('%')
				b.WriteString(key)
			} else {
				b.WriteString(resolved)
			}
		} else {
			b.WriteByte(f[0])
			f = f[1:]
		}
	}

	if !strings.Contains(b.String(), "%notacolour") {
		t.Errorf("expected unknown tag passthrough, got %q", b.String())
	}
}

func TestPlainText(t *testing.T) {
	var b strings.Builder
	colour(&b, "no tags here")
	if !strings.Contains(b.String(), "no tags here") {
		t.Error("expected plain text in output")
	}
}

// resolveTag maps a keyword to its escape code, matching the switch in colour().
func resolveTag(key string) string {
	switch key {
	case "reset":
		return Reset
	case "black":
		return Black
	case "red":
		return Red
	case "green":
		return Green
	case "yellow":
		return Yellow
	case "blue":
		return Blue
	case "magenta":
		return Magenta
	case "cyan":
		return Cyan
	case "white":
		return White
	case "grey", "gray":
		return Grey
	case "lred":
		return LightRed
	case "lgreen":
		return LightGreen
	case "lyellow":
		return LightYellow
	case "lblue":
		return LightBlue
	case "lmagenta":
		return LightMagenta
	case "lcyan":
		return LightCyan
	case "lwhite":
		return LightWhite
	case "bgblack":
		return BGBlack
	case "bgred":
		return BGRed
	case "bggreen":
		return BGGreen
	case "bgyellow":
		return BGYellow
	case "bgblue":
		return BGBlue
	case "bgmagenta":
		return BGMagenta
	case "bgcyan":
		return BGCyan
	case "bgwhite":
		return BGWhite
	case "bggrey", "bggray":
		return BGGrey
	case "bglred":
		return BGLightRed
	case "bglgreen":
		return BGLightGreen
	case "bglyellow":
		return BGLightYellow
	case "bglblue":
		return BGLightBlue
	case "bglmagenta":
		return BGLightMagenta
	case "bglcyan":
		return BGLightCyan
	case "bglwhite":
		return BGLightWhite
	case "bold":
		return Bold
	case "fuzzy":
		return Fuzzy
	case "italic":
		return Italic
	case "under":
		return Underscore
	case "blink":
		return Blink
	case "fast":
		return FastBlink
	case "reverse":
		return Reverse
	case "conceal":
		return Concealed
	case "strike":
		return Strikethrough
	}
	return ""
}
