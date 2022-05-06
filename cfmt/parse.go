package cfmt

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/grimdork/climate/str"
)

// Print with colours, but no other formatting.
func Print(s string) {
	colour(os.Stdout, s)
}

// Printf with colours uses fmt.Printf().
func Printf(s string, v ...interface{}) {
	buf := str.NewStringer()
	colour(buf, s)
	fmt.Printf(buf.String(), v...)
}

func colour(dst io.Writer, f string) {
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

			case "grey":
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

			case "bggrey":
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
	dst.Write([]byte("\n"))
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
