package cfmt

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// Print with colours, but no other formatting.
// Skips colours if the user has disabled them via the NO_COLOR envvar or we're not in a terminal.
func Print(s string) {
	colour(os.Stdout, s)
}

// Println prints with colours and appends a newline.
func Println(s string) {
	colour(os.Stdout, s)
	os.Stdout.Write([]byte("\n"))
}

// Printf with colours uses fmt.Printf().
// Skips colours if the user has disabled them via the NO_COLOR envvar or we're not in a terminal.
func Printf(s string, v ...any) {
	buf := strings.Builder{}
	colour(&buf, s)
	fmt.Printf(buf.String(), v...)
}

func colour(dst io.Writer, f string) {
	// If the user has disabled colour or we're not in a terminal, just write the string as-is.
	if !shouldColor() {
		dst.Write([]byte(f))
		return
	}

	// Process the input string character by character.
	for len(f) > 0 {
		c := f[0]
		if c == '%' {
			var key string
			key, f = parseKeyword(f)
			switch key {
			case "reset":
				dst.Write([]byte(Reset))

			// Text colour
			case "black":
				dst.Write([]byte(Black))
			case "red":
				dst.Write([]byte(Red))
			case "green":
				dst.Write([]byte(Green))
			case "yellow":
				dst.Write([]byte(Yellow))
			case "blue":
				dst.Write([]byte(Blue))
			case "magenta":
				dst.Write([]byte(Magenta))
			case "cyan":
				dst.Write([]byte(Cyan))
			case "white":
				dst.Write([]byte(White))

			case "grey", "gray":
				dst.Write([]byte(Grey))
			case "lred":
				dst.Write([]byte(LightRed))
			case "lgreen":
				dst.Write([]byte(LightGreen))
			case "lyellow":
				dst.Write([]byte(LightYellow))
			case "lblue":
				dst.Write([]byte(LightBlue))
			case "lmagenta":
				dst.Write([]byte(LightMagenta))
			case "lcyan":
				dst.Write([]byte(LightCyan))
			case "lwhite":
				dst.Write([]byte(LightWhite))

			// Background colour
			case "bgblack":
				dst.Write([]byte(BGBlack))
			case "bgred":
				dst.Write([]byte(BGRed))
			case "bggreen":
				dst.Write([]byte(BGGreen))
			case "bgyellow":
				dst.Write([]byte(BGYellow))
			case "bgblue":
				dst.Write([]byte(BGBlue))
			case "bgmagenta":
				dst.Write([]byte(BGMagenta))
			case "bgcyan":
				dst.Write([]byte(BGCyan))
			case "bgwhite":
				dst.Write([]byte(BGWhite))

			case "bggrey", "bggray":
				dst.Write([]byte(BGGrey))
			case "bglred":
				dst.Write([]byte(BGLightRed))
			case "bglgreen":
				dst.Write([]byte(BGLightGreen))
			case "bglyellow":
				dst.Write([]byte(BGLightYellow))
			case "bglblue":
				dst.Write([]byte(BGLightBlue))
			case "bglmagenta":
				dst.Write([]byte(BGLightMagenta))
			case "bglcyan":
				dst.Write([]byte(BGLightCyan))
			case "bglwhite":
				dst.Write([]byte(BGLightWhite))

			// Other styling
			case "bold":
				dst.Write([]byte(Bold))
			case "fuzzy":
				dst.Write([]byte(Fuzzy))
			case "italic":
				dst.Write([]byte(Italic))
			case "under":
				dst.Write([]byte(Underscore))
			case "blink":
				dst.Write([]byte(Blink))
			case "fast":
				dst.Write([]byte(FastBlink))
			case "reverse":
				dst.Write([]byte(Reverse))
			case "conceal":
				dst.Write([]byte(Concealed))
			case "strike":
				dst.Write([]byte(Strikethrough))
			default:
				dst.Write([]byte("%"))
				dst.Write([]byte(key))
			}
		} else {
			dst.Write([]byte{f[0]})
			f = f[1:]
		}
	}
}

// parseKeyword returns the parsed keyword and the rest of the input string.
func parseKeyword(f string) (string, string) {
	var b strings.Builder
	if len(f) == 0 {
		return "", f
	}

	b.WriteByte(f[0])
	in := f[1:]
	loop := true
	for len(in) > 0 && loop {
		if !unicode.IsLetter(rune(in[0])) {
			loop = false
			if len(in) > 1 && in[0] == ' ' {
				in = in[1:]
			}
		} else {
			b.WriteByte(in[0])
			in = in[1:]
		}
	}
	return b.String()[1:], in
}

func shouldColor() bool {
	// 1. Check if user explicitly disabled it via env var
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	// 2. Check if we are actually in a terminal
	return IsTerminal()
}

// IsTerminal returns true if the standard output is a terminal.
func IsTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	// Check if the bit for "Character Device" is set.
	// Terminals are character devices; pipes and files are not.
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
