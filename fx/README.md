# climate/fx
Small `{}`-style formatting with ANSI and time tokens.

`fx` is a tiny formatter for CLI output. It replaces `{}` placeholders with values, supports colour and style tokens such as `{red}` and `{bold}`, can insert date and time values with tokens like `{time}` and `{stamp}`, and includes plain-text rendering helpers for logs and file output.

## Installation
```bash
go get github.com/grimdork/climate/fx
```

## Usage

### Basic formatting
```go
package main

import "github.com/grimdork/climate/fx"

func main() {
	fx.Println("Hello, {}!", "Boss")
	fx.Print("Count: {}", 42)

	s := fx.Sprint("Joined: {} and {}", "first", "second")
	plain := fx.RenderPlain("{green}OK{@} {}", "done")
	_ = s
	_ = plain
}
```

### Colours and styles
```go
fx.Println("{red}Error:{@} {}", "file not found")
fx.Println("{bold}{yellow}Warning:{@} something happened")
fx.Println("{success}Done{@}")
```

Tokens are case-insensitive. Use `{@}` or `{reset}` to clear formatting.

| Category | Tokens |
|---|---|---|
| Foreground | `{black}`, `{red}`, `{green}`, `{yellow}`, `{blue}`, `{magenta}`, `{cyan}`, `{white}`, `{grey}`/`{gray}` |
| Bright foreground | `{brightblack}`/`{brightgrey}`/`{brightgray}`, `{brightred}`, …, `{brightwhite}` |
| Background | `{bgblack}`, `{bgred}`, …, `{bgwhite}`, `{bggrey}`/`{bggray}` |
| Bright background | `{bgbrightblack}`/`{bgbrightgrey}`/`{bgbrightgray}`, `{bgbrightred}`, …, `{bgbrightwhite}` |
| Styles | `{bold}`, `{dim}`, `{italic}`, `{underline}`/`{under}`, `{blink}`, `{fast}`, `{inverse}`/`{invert}`, `{hidden}`/`{conceal}`/`{concealed}`, `{strike}`/`{strikethrough}` |
| Reset | `{@}`, `{reset}` |

### Log-safe output
`RenderPlain`, `Log`, and `Logln` strip ANSI escape codes from the result.

```go
fx.Logln("{magenta}Coloured in the terminal, plain in logs{@}")
plain := fx.RenderPlain("{red}Error:{@} {}", "boom")
```

If you want to write to your own destination, use the writer variants:

```go
fx.Fprint(w, "{green}Status:{@} {}", "ready")
fx.Flogln(w, "{magenta}Plain in logs{@}")
```

### Time tokens
```go
fx.Println("[{time}] Started")
fx.Println("Zone: {tzname} ({tzoffset})")
fx.Println("Stamp: {stamp}")
```

Available value tokens include:
- `{date}` and `{tzdate}`
- `{time}` and `{stamp}`
- `{logstamp}`
- `{year}`, `{month}`, `{monthnum}`, `{day}`
- `{dow}` / `{dayofweek}`
- `{hour}`, `{min}`, `{sec}`
- `{tzoffset}`, `{tzsecs}`, `{tzname}`

Some timezone tokens accept a `:utc` modifier, such as `{tzoffset:utc}`.

## Notes
- `{}` is replaced by the next argument.
- Extra arguments are appended with spaces.
- Missing arguments leave `{}` unchanged.
- `[]byte` is emitted as a string.
- Common slices are joined with `, ` and no brackets.
- `error` values and types with a `String() string` method use their own string output.
- Common string-key and int-key maps are rendered as `{key=value}` pairs.
- `complex64` and `complex128` are supported.
- Unsupported values render as `<unsupported>` unless they provide their own string or text output.
- Unknown tokens are left unchanged.
- Setting `NO_COLOR` suppresses ANSI output automatically.

## Format specs
Use `{:spec}` to format an argument with Go's `fmt` formatting verbs:

```go
fx.Println("Hex: {:x}", 255)              // "ff"
fx.Println("Bin: {:b}", 42)               // "101010"
fx.Println("{:.2f}%", 98.5)               // "98.50%"
fx.Println("{:10.3f}", 3.14)              // "     3.140"
fx.Println("{:09d}", 42)                  // "000000042"
fx.Println("{:+d}", 42)                   // "+42"
fx.Println("{:10s}", "hello")             // "     hello"
fx.Println("{:-10s}", "hello")            // "hello     "
fx.Println("Type: {:T}", 42)              // "int"
fx.Println("Repr: {:#v}", &x)             // Go syntax
```

The spec after `:` is passed directly as the format string to `fmt.Sprintf`, so any valid Go `%` verb works. Commonly used specs:

| Spec | Example arg | Output |
|---|---|---|
| `{:x}` | `255` | `"ff"` |
| `{:X}` | `255` | `"FF"` |
| `{:b}` | `42` | `"101010"` |
| `{:o}` | `42` | `"52"` |
| `{:.2f}` | `3.14159` | `"3.14"` |
| `{:10.3f}` | `3.14` | `"     3.140"` |
| `{:09d}` | `42` | `"000000042"` |
| `{:+d}` | `42` | `"+42"` |
| `{:q}` | `"hello"` | `"\"hello\""` |
| `{:c}` | `65` | `"A"` |
| `{:t}` | `true` | `"true"` |
| `{:p}` | pointer | `"0xc000..."` |
| `{:+v}` | struct | `"{A:1}"` |
| `{:#v}` | any | Go-syntax repr |
| `{:T}` | any | Go type name |

Plain `{}` uses the type-switch and is preferred when no special formatting is needed.

## Aliases
Predefined aliases:
- `danger` → `red bold`
- `warning` → `yellow bold`
- `info` → `blue`
- `success` → `green`
- `muted` → `dim`

You can add your own:

```go
fx.AddAlias("note", "magenta italic")
fx.Println("{note}Remember this{@}")
```

Aliases are safe to add and remove at runtime.

## Custom tokens
Register your own tokens with `AddToken`:

```go
fx.AddToken("app", func(mod string) (string, bool) {
	if mod == "short" {
		return "climate", true
	}
	if mod == "" {
		return "github.com/grimdork/climate", true
	}
	return "", false
})

fx.Println("Running {app:short}")
```

## Escaping and custom delimiters
Double a delimiter to emit it literally:

```go
fx.Sprint("{{red}}") // "{red}"
fx.Sprint("{{}}")    // "{}"
```

Use `SprintWithDelims` for a different single-byte delimiter pair:

```go
fx.SprintWithDelims("<", ">", "<red>Hello<@> <>", "world")
```

To also pass options, use `SprintWithDelimsOptions`:

```go
fx.SprintWithDelimsOptions(fx.Options{DisableColour: true}, "<", ">", "<red>Hello<@>")
```

## Terminal detection
`IsTerminal()` reports whether stdout is a character device (terminal). Colour output is automatically suppressed when stdout is a pipe or file, or when `NO_COLOR` is set.

```go
if fx.IsTerminal() {
    // Terminal-specific output
}
```

## Render options
Use `RenderWithOptions` when you want to change behaviour for a single call:

```go
out := fx.RenderWithOptions(fx.Options{SortMaps: true}, "{}", map[string]int{
	"b": 2,
	"a": 1,
})
// out == "{a=1, b=2}"
```

Custom delimiters also support options via `SprintWithDelimsOptions`:

```go
```

## ANSI stripping
```go
s := fx.Sprint("{green}OK{@}")
plain := fx.StripANSI(s)
```

## Zero dependencies
Uses only the Go standard library.
