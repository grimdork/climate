package fx

import (
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Sprint returns a formatted string using {} placeholders.
// Each {} is replaced by the next argument's string representation.
// If there are more args than placeholders they are appended separated by spaces.
// If there are fewer args than placeholders the remaining {} are left as-is.
func Sprint(format string, args ...any) string {
	return sprintWithChars('{', '}', format, args...)
}

// SprintWithDelims is like Sprint but uses the provided open/close single-character
// delimiters instead of { and } (e.g. '<' and '>'). If open or close are not
// single bytes the call falls back to Sprint.
func SprintWithDelims(open, close string, format string, args ...any) string {
	if len(open) != 1 || len(close) != 1 {
		return Sprint(format, args...)
	}
	return sprintWithChars(open[0], close[0], format, args...)
}

// underlying implementation that allows different delimiter bytes
func sprintWithChars(open, close byte, format string, args ...any) string {
	var b strings.Builder
	argIndex := 0
	for i := 0; i < len(format); {
		ch := format[i]
		if ch != open && ch != close {
			b.WriteByte(ch)
			i++
			continue
		}
		// handle doubled escapes like {{ or }} -> single brace
		if ch == open {
			if i+1 < len(format) && format[i+1] == open {
				b.WriteByte(open)
				i += 2
				continue
			}
			// placeholder {}
			if i+1 < len(format) && format[i+1] == close {
				if argIndex < len(args) {
					b.WriteString(sprintValue(args[argIndex]))
					argIndex++
				} else {
					b.WriteByte(open)
					b.WriteByte(close)
				}
				i += 2
				continue
			}
			// token like {red}
			j := i + 1
			for j < len(format) && format[j] != close {
				j++
			}
			if j >= len(format) {
				// no closing delimiter; treat as literal
				b.WriteByte(open)
				i++
				continue
			}
			token := strings.ToLower(format[i+1 : j])
			// check aliases first
			if exp, ok := Aliases[token]; ok {
				// expand alias into multiple tokens; each part may be an escape or a value token
				parts := strings.Fields(exp)
				for _, p := range parts {
					if writeToken(&b, p) {
						continue
					}
					// if neither escape nor value, ignore
				}
				i = j + 1
				continue
			}
			// first try value tokens (date/time)
			if writeToken(&b, token) {
				i = j + 1
				continue
			}
			// next try ANSI escapes
			if esc := tokenEscape(token); esc != "" {
				b.WriteString(esc)
				i = j + 1
				continue
			}
			// unknown token: write literal
			b.WriteByte(open)
			b.WriteString(format[i+1 : j])
			b.WriteByte(close)
			i = j + 1
			continue
		}
		// ch == close
		if i+1 < len(format) && format[i+1] == close {
			b.WriteByte(close)
			i += 2
			continue
		}
		// unmatched close -> write it
		b.WriteByte(close)
		i++
	}
	// append remaining args if any
	for j := argIndex; j < len(args); j++ {
		if b.Len() > 0 && b.String()[b.Len()-1] != ' ' {
			b.WriteByte(' ')
		}
		b.WriteString(sprintValue(args[j]))
	}
	return b.String()
}

// Print writes the formatted output to stdout (no automatic newline).
func Print(format string, args ...any) {
	os.Stdout.WriteString(Sprint(format, args...))
}

// Println writes the formatted output to stdout and ensures a trailing newline.
// If the formatted text already ends with a newline, it is not doubled.
func Println(format string, args ...any) {
	out := Sprint(format, args...)
	if len(out) == 0 || out[len(out)-1] != '\n' {
		out += "\n"
	}
	os.Stdout.WriteString(out)
}

// sprintValue converts a single value to a string without using fmt.
func sprintValue(v any) string {
	if v == nil {
		return "<nil>"
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rv.String()
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Slice, reflect.Array:
		// special-case []byte -> string
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			if rv.Kind() == reflect.Slice {
				return string(rv.Bytes())
			}
			// array of bytes - build manually
			var bs []byte
			for i := 0; i < rv.Len(); i++ {
				bs = append(bs, byte(rv.Index(i).Uint()))
			}
			return string(bs)
		}
		parts := make([]string, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			parts = append(parts, sprintValue(rv.Index(i).Interface()))
		}
		return strings.Join(parts, ", ")
	case reflect.Ptr:
		if rv.IsNil() {
			return "<nil>"
		}
		return sprintValue(rv.Elem().Interface())
	case reflect.Interface:
		if rv.IsNil() {
			return "<nil>"
		}
		return sprintValue(rv.Elem().Interface())
	default:
		// fallback: type name in angle brackets
		return "<" + rv.Type().String() + ">"
	}
}

// Aliases map holds named alias expansions. Use AddAlias to register friendly names
// that expand to space-separated tokens (e.g. "danger" -> "red bold").
var Aliases = map[string]string{
	"danger":  "red bold",
	"warning": "yellow bold",
	"info":    "blue",
	"success": "green",
	"muted":   "dim",
}

// AddAlias registers or updates an alias. Name is case-insensitive.
func AddAlias(name, expansion string) {
	Aliases[strings.ToLower(name)] = expansion
}

// pad2 returns a zero-padded 2-digit representation of n.
func pad2(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

// padN returns n as a zero-padded decimal string with width digits.
func padN(n int, width int) string {
	s := strconv.Itoa(n)
	if len(s) >= width {
		return s
	}
	return strings.Repeat("0", width-len(s)) + s
}

// writeToken tries to write either a value token (date/time etc.) or an ANSI escape
// to the provided builder. It returns true if something was written.
func writeToken(b *strings.Builder, tok string) bool {
	if v, ok := tokenValue(tok); ok {
		b.WriteString(v)
		return true
	}
	if esc := tokenEscape(tok); esc != "" {
		b.WriteString(esc)
		return true
	}
	return false
}

// tokenValue returns a string value for tokens that insert dynamic values
// such as date/time. The boolean indicates whether the token was recognized.
// tok may include a modifier after a colon, e.g. "tzoffset:utc".
func tokenValue(tok string) (string, bool) {
	// split optional modifier
	name := tok
	mod := ""
	if idx := strings.Index(tok, ":"); idx != -1 {
		name = tok[:idx]
		mod = tok[idx+1:]
	}
	name = strings.ToLower(name)
	t := time.Now()
	switch name {
	case "date":
		return t.UTC().Format(time.RFC1123), true
	case "tzdate":
		return func() string {
			if mod == "utc" {
				return t.UTC().Format(time.RFC1123)
			}
			return t.Local().Format(time.RFC1123)
		}(), true
	case "time":
		return t.Format("15:04:05"), true
	case "stamp":
		return t.Format(time.RFC3339Nano), true
	case "year":
		return strconv.Itoa(t.Year()), true
	case "day":
		return pad2(t.Day()), true
	case "month":
		// month name
		return t.Month().String(), true
	case "monthnum":
		// numeric month, zero-padded
		return pad2(int(t.Month())), true
	case "dow", "dayofweek":
		return t.Weekday().String(), true
	case "hour":
		return pad2(t.Hour()), true
	case "min":
		return pad2(t.Minute()), true
	case "sec":
		return pad2(t.Second()), true
	case "tzoffset":
		// timezone offset like +01:00 or -05:30; mod="utc" forces UTC
		var off int
		if mod == "utc" {
			_, off = t.UTC().Zone()
		} else {
			_, off = t.Local().Zone()
		}
		sign := "+"
		if off < 0 {
			sign = "-"
			off = -off
		}
		h := off / 3600
		m := (off % 3600) / 60
		return sign + pad2(h) + ":" + pad2(m), true
	case "tzsecs":
		// numeric offset in seconds; mod="utc" forces UTC
		var off int
		if mod == "utc" {
			_, off = t.UTC().Zone()
		} else {
			_, off = t.Local().Zone()
		}
		return strconv.Itoa(off), true
	case "tzname":
		if mod == "utc" {
			name, _ := t.UTC().Zone()
			return name, true
		}
		name, _ := t.Local().Zone()
		return name, true
	case "logstamp":
		// Format similar to: "Mon Jan 02 15:04:05.123456 2006"
		wd := t.Weekday().String()
		mo := t.Month().String()
		day := pad2(t.Day())
		hour := pad2(t.Hour())
		min := pad2(t.Minute())
		sec := pad2(t.Second())
		micro := t.Nanosecond() / 1000
		microStr := padN(micro, 6)
		year := strconv.Itoa(t.Year())
		// take first 3 letters of weekday/month
		if len(wd) >= 3 {
			wd = wd[:3]
		}
		if len(mo) >= 3 {
			mo = mo[:3]
		}
		return wd + " " + mo + " " + day + " " + hour + ":" + min + ":" + sec + "." + microStr + " " + year, true
	}
	return "", false
}

// tokenEscape returns the ANSI escape sequence for a known token (like "red", "bgwhite", "bold", "@" (reset)).
// If the token is unknown it returns an empty string.
func tokenEscape(tok string) string {
	if tok == "@" || tok == "reset" {
		return "\x1b[0m"
	}
	switch tok {
	case "bold":
		return "\x1b[1m"
	case "dim":
		return "\x1b[2m"
	case "italic":
		return "\x1b[3m"
	case "underline":
		return "\x1b[4m"
	case "inverse", "invert":
		return "\x1b[7m"
	case "hidden":
		return "\x1b[8m"
	}

	// colours
	colours := map[string]int{
		"black":   0,
		"red":     1,
		"green":   2,
		"yellow":  3,
		"blue":    4,
		"magenta": 5,
		"cyan":    6,
		"white":   7,
	}

	// background? prefixed with "bg"
	if strings.HasPrefix(tok, "bg") {
		c := strings.TrimPrefix(tok, "bg")
		isBright := false
		if strings.HasPrefix(c, "bright") {
			isBright = true
			c = strings.TrimPrefix(c, "bright")
		}
		c = strings.TrimSpace(c)
		if code, ok := colours[c]; ok {
			if isBright {
				return "\x1b[" + strconv.Itoa(100+code) + "m"
			}
			return "\x1b[" + strconv.Itoa(40+code) + "m"
		}
		return ""
	}

	// foreground bright?
	isBright := false
	c := tok
	if strings.HasPrefix(c, "bright") {
		isBright = true
		c = strings.TrimPrefix(c, "bright")
	}
	c = strings.TrimSpace(c)
	if code, ok := colours[c]; ok {
		if isBright {
			return "\x1b[" + strconv.Itoa(90+code) + "m"
		}
		return "\x1b[" + strconv.Itoa(30+code) + "m"
	}
	return ""
}

// StripANSI removes common ANSI SGR escape sequences (like \x1b[31m) from s.
func StripANSI(s string) string {
	re := regexp.MustCompile("\x1b\\[[0-9;]*m")
	return re.ReplaceAllString(s, "")
}

// Log writes the formatted output but with ANSI sequences stripped (useful for logs).
func Log(format string, args ...any) {
	out := Sprint(format, args...)
	os.Stdout.WriteString(StripANSI(out))
}

// Logln like Log but ensures a trailing newline.
func Logln(format string, args ...any) {
	out := Sprint(format, args...)
	if len(out) == 0 || out[len(out)-1] != '\n' {
		out += "\n"
	}
	os.Stdout.WriteString(StripANSI(out))
}
