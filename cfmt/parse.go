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

// Sprintf returns the formatted string with colour tags resolved when appropriate.
// Behaviour mirrors Printf but returns the string instead of printing it.
func Sprintf(s string, v ...any) string {
	// If the string contains colour tags like %red, they conflict with fmt verbs.
	// Replace colour tags with placeholder tokens, run fmt.Sprintf, then restore tags
	// and run the colour parser.
	escaped := strings.Builder{}
	tokens := map[string]string{}
	tokIdx := 0
	for i := 0; i < len(s); {
		if s[i] == '%' && i+1 < len(s) && isLetter(s[i+1]) {
			j := i + 1
			for j < len(s) && isLetter(s[j]) {
				j++
			}
			// Only treat as a colour tag if the tag length is > 1 (avoid catching format verbs like %s)
			if j-(i+1) > 1 {
				tag := s[i:j]
				token := fmt.Sprintf("\x00COL%03d\x00", tokIdx)
				tokens[token] = tag
				escaped.WriteString(token)
				tokIdx++
				i = j
				continue
			}
		}
		escaped.WriteByte(s[i])
		i++

	}

	formatted := fmt.Sprintf(escaped.String(), v...)

	// Restore tokens to their original %tag form
	for token, tag := range tokens {
		formatted = strings.ReplaceAll(formatted, token, tag)
	}

	buf := strings.Builder{}
	colour(&buf, formatted)
	return buf.String()
}

func isLetter(b byte) bool {
	c := rune(b)
	return unicode.IsLetter(c)
}

func stripTags(f string) string {
	var b strings.Builder
	for i := 0; i < len(f); {
		if f[i] == '%' && i+1 < len(f) && unicode.IsLetter(rune(f[i+1])) {
			j := i + 1
			for j < len(f) && unicode.IsLetter(rune(f[j])) {
				j++
			}
			i = j
			continue
		}
		b.WriteByte(f[i])
		i++
	}
	return b.String()
}

func colour(dst io.Writer, f string) {
	// If the user has disabled colour or we're not in a terminal, strip tags and write plain text.
	if !shouldColor() {
		dst.Write([]byte(stripTags(f)))
		return
	}

	// Process the input string character by character.
	for len(f) > 0 {
		c := f[0]
		if c == '%' {
			var key string
			key, f = parseKeyword(f)
			// Attempt to match the longest valid tag prefix within key.
			matched := ""
			remaining := ""
			for l := len(key); l > 0; l-- {
				cand := key[:l]
				if code := tagCode(cand); code != "" {
					matched = cand
					remaining = key[l:]
					// write code and set f to remaining+f
					dst.Write([]byte(code))
					f = remaining + f
					break
				}
			}

			if matched == "" {
				// unknown tag — write it verbatim
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
	// Parse tags consisting only of ASCII lowercase letters (a-z).
	// Returns the letters after the leading '%' and the rest of the string.
	if len(f) == 0 {
		return "", f
	}
	in := f[1:]
	i := 0
	for i < len(in) {
		c := in[i]
		if c < 'a' || c > 'z' {
			break
		}
		i++
	}
	key := in[:i]
	rest := in[i:]
	return key, rest
}

func shouldColor() bool {
	// 1. Check if user explicitly disabled it via NO_COLOR env var
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
