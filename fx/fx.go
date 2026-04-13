package fx

import (
	"encoding"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Options control rendering behaviour.
type Options struct {
	// SortMaps sorts supported map keys alphabetically/numerically before rendering.
	SortMaps bool
	// DisableColour suppresses ANSI output for built-in colour/style tokens.
	// NO_COLOR also disables colour automatically.
	DisableColour bool
}

// TokenFunc renders a named token. The modifier is the optional text after a colon,
// such as the "utc" in {tzoffset:utc}. The boolean reports whether the modifier
// was accepted.
type TokenFunc func(modifier string) (string, bool)

type stringer interface {
	String() string
}

type builderState struct {
	strings.Builder
	last byte
	has  bool
}

func (b *builderState) writeByte(ch byte) {
	b.Builder.WriteByte(ch)
	b.last = ch
	b.has = true
}

func (b *builderState) writeString(s string) {
	if s == "" {
		return
	}
	b.Builder.WriteString(s)
	b.last = s[len(s)-1]
	b.has = true
}

var (
	ansiPattern = regexp.MustCompile("\x1b\\[[0-9;]*m")

	colourCodes = map[string]int{
		"black":   0,
		"red":     1,
		"green":   2,
		"yellow":  3,
		"blue":    4,
		"magenta": 5,
		"cyan":    6,
		"white":   7,
	}

	aliasMu sync.RWMutex
	aliases = map[string]string{
		"danger":  "red bold",
		"warning": "yellow bold",
		"info":    "blue",
		"success": "green",
		"muted":   "dim",
	}

	tokenMu sync.RWMutex
	tokens  = map[string]TokenFunc{
		"date": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return time.Now().UTC().Format(time.RFC1123), true
		},
		"tzdate": func(modifier string) (string, bool) {
			t := time.Now()
			if modifier == "" {
				return t.Local().Format(time.RFC1123), true
			}
			if modifier == "utc" {
				return t.UTC().Format(time.RFC1123), true
			}
			return "", false
		},
		"time": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return time.Now().Format("15:04:05"), true
		},
		"stamp": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return time.Now().Format(time.RFC3339Nano), true
		},
		"year": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return strconv.Itoa(time.Now().Year()), true
		},
		"day": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return pad2(time.Now().Day()), true
		},
		"month": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return time.Now().Month().String(), true
		},
		"monthnum": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return pad2(int(time.Now().Month())), true
		},
		"dow": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return time.Now().Weekday().String(), true
		},
		"dayofweek": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return time.Now().Weekday().String(), true
		},
		"hour": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return pad2(time.Now().Hour()), true
		},
		"min": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return pad2(time.Now().Minute()), true
		},
		"sec": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			return pad2(time.Now().Second()), true
		},
		"tzoffset": func(modifier string) (string, bool) {
			t := time.Now()
			var off int
			switch modifier {
			case "":
				_, off = t.Local().Zone()
			case "utc":
				_, off = t.UTC().Zone()
			default:
				return "", false
			}
			sign := "+"
			if off < 0 {
				sign = "-"
				off = -off
			}
			h := off / 3600
			m := (off % 3600) / 60
			return sign + pad2(h) + ":" + pad2(m), true
		},
		"tzsecs": func(modifier string) (string, bool) {
			t := time.Now()
			var off int
			switch modifier {
			case "":
				_, off = t.Local().Zone()
			case "utc":
				_, off = t.UTC().Zone()
			default:
				return "", false
			}
			return strconv.Itoa(off), true
		},
		"tzname": func(modifier string) (string, bool) {
			t := time.Now()
			switch modifier {
			case "":
				name, _ := t.Local().Zone()
				return name, true
			case "utc":
				name, _ := t.UTC().Zone()
				return name, true
			default:
				return "", false
			}
		},
		"logstamp": func(modifier string) (string, bool) {
			if modifier != "" {
				return "", false
			}
			t := time.Now()
			wd := t.Weekday().String()
			mo := t.Month().String()
			if len(wd) >= 3 {
				wd = wd[:3]
			}
			if len(mo) >= 3 {
				mo = mo[:3]
			}
			return wd + " " + mo + " " + pad2(t.Day()) + " " + pad2(t.Hour()) + ":" + pad2(t.Minute()) + ":" + pad2(t.Second()) + "." + padN(t.Nanosecond()/1000, 6) + " " + strconv.Itoa(t.Year()), true
		},
	}
)

// Sprint returns a formatted string using {} placeholders.
// Each {} is replaced by the next argument's string representation.
// If there are more args than placeholders they are appended separated by spaces.
// If there are fewer args than placeholders the remaining {} are left as-is.
func Sprint(format string, args ...any) string {
	return Render(format, args...)
}

// Render formats a string using the default options.
func Render(format string, args ...any) string {
	return RenderWithOptions(Options{}, format, args...)
}

// RenderPlain formats a string and strips ANSI escape sequences from the result.
func RenderPlain(format string, args ...any) string {
	return RenderPlainWithOptions(Options{}, format, args...)
}

// RenderWithOptions formats a string using the provided options.
func RenderWithOptions(opts Options, format string, args ...any) string {
	out := renderWithChars('{', '}', opts, format, args...)
	if !colourEnabled(opts) {
		return StripANSI(out)
	}
	return out
}

// RenderPlainWithOptions formats a string using the provided options and strips ANSI.
func RenderPlainWithOptions(opts Options, format string, args ...any) string {
	return StripANSI(renderWithChars('{', '}', opts, format, args...))
}

// SprintWithDelims is like Sprint but uses the provided open/close single-character
// delimiters instead of { and } (e.g. '<' and '>'). If open or close are not
// single bytes the call falls back to Sprint.
func SprintWithDelims(open, close string, format string, args ...any) string {
	if len(open) != 1 || len(close) != 1 {
		return Sprint(format, args...)
	}
	out := renderWithChars(open[0], close[0], Options{}, format, args...)
	if !colourEnabled(Options{}) {
		return StripANSI(out)
	}
	return out
}

// Fprint writes the rendered output to w with no automatic newline.
func Fprint(w io.Writer, format string, args ...any) {
	_, _ = io.WriteString(w, Render(format, args...))
}

// Fprintln writes the rendered output to w and ensures a trailing newline.
func Fprintln(w io.Writer, format string, args ...any) {
	out := Render(format, args...)
	if len(out) == 0 || out[len(out)-1] != '\n' {
		out += "\n"
	}
	_, _ = io.WriteString(w, out)
}

// Flog writes the rendered output to w with ANSI sequences stripped.
func Flog(w io.Writer, format string, args ...any) {
	_, _ = io.WriteString(w, RenderPlain(format, args...))
}

// Flogln writes the rendered output to w with ANSI sequences stripped and ensures a trailing newline.
func Flogln(w io.Writer, format string, args ...any) {
	out := RenderPlain(format, args...)
	if len(out) == 0 || out[len(out)-1] != '\n' {
		out += "\n"
	}
	_, _ = io.WriteString(w, out)
}

// Print writes the formatted output to stdout (no automatic newline).
func Print(format string, args ...any) {
	Fprint(os.Stdout, format, args...)
}

// Println writes the formatted output to stdout and ensures a trailing newline.
// If the formatted text already ends with a newline, it is not doubled.
func Println(format string, args ...any) {
	Fprintln(os.Stdout, format, args...)
}

// Log writes the formatted output but with ANSI sequences stripped (useful for logs).
func Log(format string, args ...any) {
	Flog(os.Stdout, format, args...)
}

// Logln like Log but ensures a trailing newline.
func Logln(format string, args ...any) {
	Flogln(os.Stdout, format, args...)
}

// AddAlias registers or updates an alias. Name is case-insensitive.
func AddAlias(name, expansion string) {
	aliasMu.Lock()
	aliases[strings.ToLower(name)] = expansion
	aliasMu.Unlock()
}

// DeleteAlias removes a registered alias. Name is case-insensitive.
func DeleteAlias(name string) {
	aliasMu.Lock()
	delete(aliases, strings.ToLower(name))
	aliasMu.Unlock()
}

// AddToken registers or updates a token renderer. Name is case-insensitive.
func AddToken(name string, fn TokenFunc) {
	tokenMu.Lock()
	tokens[strings.ToLower(name)] = fn
	tokenMu.Unlock()
}

// DeleteToken removes a registered token renderer. Name is case-insensitive.
func DeleteToken(name string) {
	tokenMu.Lock()
	delete(tokens, strings.ToLower(name))
	tokenMu.Unlock()
}

func getAlias(name string) (string, bool) {
	aliasMu.RLock()
	exp, ok := aliases[strings.ToLower(name)]
	aliasMu.RUnlock()
	return exp, ok
}

func getToken(name string) (TokenFunc, bool) {
	tokenMu.RLock()
	fn, ok := tokens[strings.ToLower(name)]
	tokenMu.RUnlock()
	return fn, ok
}

// renderWithChars is the underlying implementation that allows different delimiter bytes.
func renderWithChars(open, close byte, opts Options, format string, args ...any) string {
	var b builderState
	argIndex := 0
	for i := 0; i < len(format); {
		ch := format[i]
		if ch != open && ch != close {
			b.writeByte(ch)
			i++
			continue
		}
		if ch == open {
			if i+1 < len(format) && format[i+1] == open {
				b.writeByte(open)
				i += 2
				continue
			}
			if i+1 < len(format) && format[i+1] == close {
				if argIndex < len(args) {
					b.writeString(sprintValue(args[argIndex], opts))
					argIndex++
				} else {
					b.writeByte(open)
					b.writeByte(close)
				}
				i += 2
				continue
			}
			j := i + 1
			for j < len(format) && format[j] != close {
				j++
			}
			if j >= len(format) {
				b.writeByte(open)
				i++
				continue
			}
			tok := strings.ToLower(format[i+1 : j])
			if writeToken(&b, tok, 0) {
				i = j + 1
				continue
			}
			b.writeByte(open)
			b.writeString(format[i+1 : j])
			b.writeByte(close)
			i = j + 1
			continue
		}
		if i+1 < len(format) && format[i+1] == close {
			b.writeByte(close)
			i += 2
			continue
		}
		b.writeByte(close)
		i++
	}
	for j := argIndex; j < len(args); j++ {
		if b.has && b.last != ' ' {
			b.writeByte(' ')
		}
		b.writeString(sprintValue(args[j], opts))
	}
	return b.String()
}

func writeToken(b *builderState, tok string, depth int) bool {
	if depth > 8 {
		return false
	}
	if exp, ok := getAlias(tok); ok {
		parts := strings.Fields(exp)
		for _, part := range parts {
			if writeToken(b, strings.ToLower(part), depth+1) {
				continue
			}
		}
		return true
	}
	name, modifier := splitToken(tok)
	if fn, ok := getToken(name); ok {
		if v, ok := fn(modifier); ok {
			b.writeString(v)
			return true
		}
	}
	if esc, ok := tokenEscape(name); ok {
		b.writeString(esc)
		return true
	}
	return false
}

func splitToken(tok string) (string, string) {
	if idx := strings.Index(tok, ":"); idx != -1 {
		return tok[:idx], tok[idx+1:]
	}
	return tok, ""
}

// sprintValue converts a single value to a string without using reflection.
func sprintValue(v any, opts Options) string {
	if v == nil {
		return "<nil>"
	}

	switch x := v.(type) {
	case error:
		return x.Error()
	case stringer:
		return x.String()
	case encoding.TextMarshaler:
		if text, err := x.MarshalText(); err == nil {
			return string(text)
		}
	case string:
		return x
	case bool:
		return strconv.FormatBool(x)
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case uintptr:
		return strconv.FormatUint(uint64(x), 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	case []byte:
		return string(x)
	case []string:
		return formatSlice(x, opts)
	case []bool:
		return formatSlice(x, opts)
	case []int:
		return formatSlice(x, opts)
	case []int8:
		return formatSlice(x, opts)
	case []int16:
		return formatSlice(x, opts)
	case []int32:
		return formatSlice(x, opts)
	case []int64:
		return formatSlice(x, opts)
	case []uint:
		return formatSlice(x, opts)
	case []uint16:
		return formatSlice(x, opts)
	case []uint32:
		return formatSlice(x, opts)
	case []uint64:
		return formatSlice(x, opts)
	case []uintptr:
		return formatSlice(x, opts)
	case []float32:
		return formatSlice(x, opts)
	case []float64:
		return formatSlice(x, opts)
	case []any:
		return formatSlice(x, opts)
	case map[string]string:
		return formatStringMap(x, opts)
	case map[string]bool:
		return formatStringMap(x, opts)
	case map[string]int:
		return formatStringMap(x, opts)
	case map[string]int8:
		return formatStringMap(x, opts)
	case map[string]int16:
		return formatStringMap(x, opts)
	case map[string]int32:
		return formatStringMap(x, opts)
	case map[string]int64:
		return formatStringMap(x, opts)
	case map[string]uint:
		return formatStringMap(x, opts)
	case map[string]uint8:
		return formatStringMap(x, opts)
	case map[string]uint16:
		return formatStringMap(x, opts)
	case map[string]uint32:
		return formatStringMap(x, opts)
	case map[string]uint64:
		return formatStringMap(x, opts)
	case map[string]float32:
		return formatStringMap(x, opts)
	case map[string]float64:
		return formatStringMap(x, opts)
	case map[string]any:
		return formatStringMap(x, opts)
	case map[int]string:
		return formatIntMap(x, opts)
	case map[int]bool:
		return formatIntMap(x, opts)
	case map[int]int:
		return formatIntMap(x, opts)
	case map[int]int8:
		return formatIntMap(x, opts)
	case map[int]int16:
		return formatIntMap(x, opts)
	case map[int]int32:
		return formatIntMap(x, opts)
	case map[int]int64:
		return formatIntMap(x, opts)
	case map[int]uint:
		return formatIntMap(x, opts)
	case map[int]uint8:
		return formatIntMap(x, opts)
	case map[int]uint16:
		return formatIntMap(x, opts)
	case map[int]uint32:
		return formatIntMap(x, opts)
	case map[int]uint64:
		return formatIntMap(x, opts)
	case map[int]float32:
		return formatIntMap(x, opts)
	case map[int]float64:
		return formatIntMap(x, opts)
	case map[int]any:
		return formatIntMap(x, opts)
	}

	return "<unsupported>"
}

func formatSlice[T any](vals []T, opts Options) string {
	parts := make([]string, 0, len(vals))
	for _, v := range vals {
		parts = append(parts, sprintValue(any(v), opts))
	}
	return strings.Join(parts, ", ")
}

func formatStringMap[T any](m map[string]T, opts Options) string {
	parts := make([]string, 0, len(m))
	if opts.SortMaps {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			parts = append(parts, k+"="+sprintValue(any(m[k]), opts))
		}
	} else {
		for k, v := range m {
			parts = append(parts, k+"="+sprintValue(any(v), opts))
		}
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func formatIntMap[T any](m map[int]T, opts Options) string {
	parts := make([]string, 0, len(m))
	if opts.SortMaps {
		keys := make([]int, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			parts = append(parts, strconv.Itoa(k)+"="+sprintValue(any(m[k]), opts))
		}
	} else {
		for k, v := range m {
			parts = append(parts, strconv.Itoa(k)+"="+sprintValue(any(v), opts))
		}
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// tokenEscape returns the ANSI escape sequence for a known token (like "red", "bgwhite", "bold", "@" (reset)).
// If the token is unknown it returns false.
func tokenEscape(tok string) (string, bool) {
	if tok == "@" || tok == "reset" {
		return "\x1b[0m", true
	}
	switch tok {
	case "bold":
		return "\x1b[1m", true
	case "dim":
		return "\x1b[2m", true
	case "italic":
		return "\x1b[3m", true
	case "underline":
		return "\x1b[4m", true
	case "inverse", "invert":
		return "\x1b[7m", true
	case "hidden":
		return "\x1b[8m", true
	}

	if strings.HasPrefix(tok, "bg") {
		c := strings.TrimSpace(strings.TrimPrefix(tok, "bg"))
		isBright := false
		if strings.HasPrefix(c, "bright") {
			isBright = true
			c = strings.TrimPrefix(c, "bright")
		}
		if code, ok := colourCodes[c]; ok {
			if isBright {
				return "\x1b[" + strconv.Itoa(100+code) + "m", true
			}
			return "\x1b[" + strconv.Itoa(40+code) + "m", true
		}
		return "", false
	}

	c := strings.TrimSpace(tok)
	isBright := false
	if strings.HasPrefix(c, "bright") {
		isBright = true
		c = strings.TrimPrefix(c, "bright")
	}
	if code, ok := colourCodes[c]; ok {
		if isBright {
			return "\x1b[" + strconv.Itoa(90+code) + "m", true
		}
		return "\x1b[" + strconv.Itoa(30+code) + "m", true
	}
	return "", false
}

func colourEnabled(opts Options) bool {
	if opts.DisableColour {
		return false
	}
	return os.Getenv("NO_COLOR") == ""
}

// StripANSI removes common ANSI SGR escape sequences (like \x1b[31m) from s.
func StripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
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
