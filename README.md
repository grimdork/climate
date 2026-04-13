# climate
Minimalist toolkit for Go command-line applications.

[![Go Report Card](https://goreportcard.com/badge/github.com/grimdork/climate)](https://goreportcard.com/report/github.com/grimdork/climate)
[![Go Reference](https://pkg.go.dev/badge/github.com/grimdork/climate.svg)](https://pkg.go.dev/github.com/grimdork/climate)
![License](https://img.shields.io/github/license/grimdork/climate)

`climate` is a collection of modular packages designed to help you build fast, small, and robust command-line tools. Every package depends **only on the Go standard library**, making it the perfect companion for **TinyGo** projects and high-performance system utilities.

> [!IMPORTANT]
> Currently optimised for Unix-like systems (Linux, macOS, BSD). Windows support is not a priority.

---

## Package overview

| Package | Description |
| :--- | :--- |
| **[`arg`](./arg)** | No-dependency option parser with subcommands and positionals. |
| **[`cfmt`](./cfmt)** | `printf` wrapper with easy `%colour` formatting. |
| **[`fx`](./fx)** | Small `{}`-style formatter with ANSI and time tokens. |
| **[`daemon`](./daemon)** | Lifecycle utilities for background services and signal handling. |
| **[`prompter`](./prompter)** | Interactive user prompts with support for masked secrets. |
| **[`paths`](./paths)** | Resolves standard config paths (XDG on Linux, Library on macOS). |
| **[`str`](./str)** | Type-aware `strings.Builder` extension for complex data. |
| **[`tab`](./tab)** | Clean, elastic tabbed columns for terminal tables. |
| **[`human`](./human)** | Byte-to-human scaling (1024/1000) for readable metrics. |
| **[`env`](./env)** | Simple environment variable fetching with default fallbacks. |
| **[`ini`](./ini)** | Dependency-free INI parser with schema and env overrides. |
---

## Quick Start

### Installation
```bash
go get github.com/grimdork/climate
```

Or just type \<package>.\<function> and let `goimports` add it and `go mod tidy` download it like a sane person would.

## Argument Parsing (arg)
Handle flags, subcommands, and positional arguments without the overhead of heavy frameworks.

```go
package main

import (
	"fmt"
	"os"
	"github.com/grimdork/climate/arg"
)

func main() {
	opt := arg.New("mytool", "Short one-line description of mytool")
	opt.SetDefaultHelp(true)
	// Set a boolean flag (-v, --verbose)
	opt.SetOption(arg.GroupDefault, "v", "verbose", "Enable verbose output", false, false, arg.VarBool, nil)

	// Create a subcommand and add a positional to its Options
	cmd := opt.SetCommand("greet", "Greet a user", arg.GroupDefault, nil, nil)
	cmd.Options.SetPositional("name", "The name to greet", "world", true, arg.VarString)

	err := opt.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing args:", err)
	}
}
```

## Colour formatting (cfmt)
Use ANSI colours with Printf.

```go
import "github.com/grimdork/climate/cfmt"

cfmt.Printf("%red Error:%reset file not found: '%s'\n", fn)
cfmt.Print("%green Success!%reset Process completed.\n")
```

## Lightweight formatting (fx)
Use simple `{}` placeholders plus optional colour, style, and time tokens.

```go
import "github.com/grimdork/climate/fx"

fx.Println("{green}Status:{@} {}", "ready")
fx.Println("[{time}] Starting up")
```

## Signal Handling (daemon)
Cleanly handle Ctrl+C or SIGTERM in servers or background tasks. Also erases the ugly `^C` from the output.

```go
import "github.com/grimdork/climate/daemon"

// Block until termination signal
<-daemon.BreakChannel()
fmt.Println("Shutting down gracefully...")
```

## Interactive input (prompter)
```go
import "github.com/grimdork/climate/prompter"

// Ask with a default value and no masking
user := prompter.Ask("Username", "admin", false)

// Ask with masking enabled for passwords
pass := prompter.Ask("Password", "", true)
```

## XDG & macOS path handling (paths)
```go
import "github.com/grimdork/climate/paths"

// Returns ~/.config/myapp (Linux) or ~/Library/Application Support/myapp (macOS)
// Automatically creates the directory if it doesn't exist.
// If XDG path is set on Linux, it replaces ".config".
dir, err := paths.GetConfigDir("myapp")
```

## Enhanced string builder (str)
```go
import "github.com/grimdork/climate/str"

s := str.NewStringer().SetSliceComma(true)
tags := []string{"go", "cli", "minimal"}

// Recursive, type-aware writing
s.WriteI("Tags: ", tags, " | Count: ", len(tags))
// Output: Tags: go,cli,minimal | Count: 3
```

## Elastic tabulation (tab)
```go
import "github.com/grimdork/climate/tab"

input := "Name Role Status\nAlice Admin Active\nBob User Offline"
// Formats columns with consistent padding
output, _ := tab.Tabulate(input, false)
println(output)
```

## Readable sizes with rounding (human)
```go
import "github.com/grimdork/climate/human"

// Convert bytes to KiB/MiB/GiB
fmt.Println(human.UInt(1572864, false)) // "1.5 MiB"

// Convert to SI (k/M/G) units
fmt.Println(human.UInt(1500000, true))  // "1.5 MB"
```

## Environment helper (env)
```go
import "github.com/grimdork/climate/env"

// Simple string fallback: returns the env value if set, otherwise the alt.
host := env.Get("DB_HOST", "localhost")

// Typed helpers: return the supplied alt when the env var is unset or malformed.
// Safe to inline with existing config values (e.g. INI defaults).
debug := env.GetBool("DEBUG", false)
port := env.GetInt("PORT", 5432)
timeout := env.GetFloat("REQUEST_TIMEOUT", 1.5)
```

## INI subpackage (ini)

The ini package provides a small, dependency‑free INI parser/serializer with optional environment overrides and schema support.

Key features:
- Parse and Marshal INI documents in memory (Unmarshal/Marshal). Unmarshal("") returns an error on empty input.
- Declare expected types per section/key with DeclareType(section, key, type) and then call Parse(s) to honour those types (use section=="" for top-level properties).
- Environment overrides: enable env-first lookups with SetEnvFirst(true). Use ForceUpper() to consult UPPER_CASE env keys.
- Primary getters: GetString/GetBool/GetInt/GetFloat. When env-first is enabled, these consult the environment (using the climate/env helpers) and fall back to the INI value when unset or invalid.
- Save/Load helpers exist and use the host filesystem; on embedded or TinyGo targets without a filesystem, stubbed or build-tagbed implementations are recommended.

Minimal example:
```go
import (
    "github.com/grimdork/climate/ini"
)

// Parse into an INI with declared types
cfg, _ := ini.Unmarshal("port=8080\n[server]\nmax=1024")
// or declare types before parsing for deterministic typing
cfg2, _ := ini.New()
cfg2.DeclareType("", "port", ini.Int)
cfg2.Parse("port=0x1F\n") // parses port as hex int because we declared Int

// Environment overrides
cfg.SetEnvFirst(true)
cfg.ForceUpper() // makes env lookups use upper-case keys
port := cfg.GetInt("", "port")
```

Notes:
- The ini parser preserves integer-vs-float typing by trying integer parsing first, then float.
- Marshal produces a stable, sorted representation (properties then sections).
- For TinyGo/embedded use, consider supplying an EnvLookup or build-tagbed stubs for env and filesystem operations so the parsing/formatting code remains usable without an OS.


## Core Philosophy
- Standard library only: Don't pull in other 3rd party dependencies.
- TinyGo compatibility: Prioritise low-allocation and small-binary footprints. No use of reflect; all type handling is done via type switches and generics.
- Unix focus: Made for the systems I actually use every day.

Minimum Go version
- This project supports Go 1.25 as the minimum toolchain to preserve TinyGo compatibility. CI tests multiple Go versions (1.25 and 1.26) to catch regressions, but the module remains compatible with Go 1.25.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
