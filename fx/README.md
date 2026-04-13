# climate/fx
Small `{}`-style formatting with ANSI and time tokens.

`fx` is a tiny, `fmt`-free formatter for CLI output. It replaces `{}` placeholders with values, supports colour and style tokens such as `{red}` and `{bold}`, and can insert date and time values with tokens like `{time}` and `{stamp}`.

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
	_ = s
}
```

### Colours and styles
```go
fx.Println("{red}Error:{@} {}", "file not found")
fx.Println("{bold}{yellow}Warning:{@} something happened")
fx.Println("{success}Done{@}")
```

Tokens are case-insensitive. Use `{@}` or `{reset}` to clear formatting.

### Log-safe output
`Log` and `Logln` strip ANSI escape codes before writing the result.

```go
fx.Logln("{magenta}Coloured in the terminal, plain in logs{@}")
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
- Other slices and arrays are joined with `, ` and no brackets.
- Pointers and interfaces are dereferenced; `nil` becomes `<nil>`.
- Unknown tokens are left unchanged.

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

## Escaping and custom delimiters
Double a delimiter to emit it literally:

```go
fx.Sprint("{{red}}") // "{red}"
fx.Sprint("{{}}")    // "{}"
```

Use `SprintWithDelims` if you want a different single-byte delimiter pair:

```go
fx.SprintWithDelims("<", ">", "<red>Hello<@> {}", "world")
```

## ANSI stripping
```go
s := fx.Sprint("{green}OK{@}")
plain := fx.StripANSI(s)
```

## Zero dependencies
Uses only the Go standard library.
