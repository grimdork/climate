# climate/cfmt
Simple ANSI color formatting using printf-style templates.

`cfmt` allows you to add color and style to your terminal output using a `%keyword` syntax. It handles the ANSI escape codes for you, keeping your Go code free of messy control characters.

## Installation
```bash
go get github.com/grimdork/climate/cfmt
```

## Usage
Color and style tags use a `%` prefix followed by a keyword. The keyword ends at the first non-letter character. Always include `%reset` after colored output to prevent color bleeding into the rest of the terminal.

**Note:** Both `Print` and `Printf` automatically append a newline.

### Basic Example
```go
package main

import "github.com/grimdork/climate/cfmt"

func main() {
	// Simple colored output
	cfmt.Printf("%red Error:%reset file not found: '%s'", fn)
	cfmt.Print("%green Success!%reset Configuration loaded")

	// Combine styles
	cfmt.Printf("%bold %yellow Warning:%reset This action is permanent.")
}
```

### Printf with format verbs
Color tags and standard Go format verbs work together:
```go
cfmt.Printf("%cyan Status:%reset %s (%d items)", status, count)
```

### NO_COLOR support
`cfmt` respects the `NO_COLOR` environment variable. If set, or if stdout is not a terminal, all color tags are stripped and output is plain text. You can check this yourself with `cfmt.IsTerminal()`.

### Available Tags
| Category | Tags |
| :--- | :--- |
| Reset | `%reset` |
| Text | `%black`, `%red`, `%green`, `%yellow`, `%blue`, `%magenta`, `%cyan`, `%white` |
| Light text | `%grey`/`%gray`, `%lred`, `%lgreen`, `%lyellow`, `%lblue`, `%lmagenta`, `%lcyan`, `%lwhite` |
| Background | `%bgblack`, `%bgred`, `%bggreen`, `%bgyellow`, `%bgblue`, `%bgmagenta`, `%bgcyan`, `%bgwhite` |
| Light bg | `%bggrey`/`%bggray`, `%bglred`, `%bglgreen`, `%bglyellow`, `%bglblue`, `%bglmagenta`, `%bglcyan`, `%bglwhite` |
| Styles | `%bold`, `%fuzzy`, `%italic`, `%under`, `%blink`, `%fast`, `%reverse`, `%conceal`, `%strike` |
