# climate
Minimalist toolkit for Go command-line applications.

[![Go Report Card](https://goreportcard.com/badge/github.com/grimdork/climate)](https://goreportcard.com/report/github.com/grimdork/climate)
[![Go Reference](https://pkg.go.dev/badge/github.com/grimdork/climate.svg)](https://pkg.go.dev/github.com/grimdork/climate)
![License](https://img.shields.io/github/license/grimdork/climate)

`climate` is a collection of modular packages designed to help you build fast, small, and robust command-line tools. Every package depends **only on the Go standard library**, making it the perfect companion for **TinyGo** projects and high-performance system utilities.

> [!IMPORTANT]
> Currently optimized for Unix-like systems (Linux, macOS, BSD). Windows support is not a priority.

---

## Package overview

| Package | Description |
| :--- | :--- |
| **[`arg`](./arg)** | No-dependency option parser with subcommands and positionals. |
| **[`cfmt`](./cfmt)** | `printf` wrapper with easy `{{.Color}}` template formatting. |
| **[`daemon`](./daemon)** | Lifecycle utilities for background services and signal handling. |
| **[`prompter`](./prompter)** | Interactive user prompts with support for masked secrets. |
| **[`paths`](./paths)** | Resolves standard config paths (XDG on Linux, Library on macOS). |
| **[`str`](./str)** | Type-aware `strings.Builder` extension for complex data. |
| **[`tab`](./tab)** | Clean, elastic tabbed columns for terminal tables. |
| **[`human`](./human)** | Byte-to-human scaling (1024/1000) for readable metrics. |
| **[`env`](./env)** | Simple environment variable fetching with default fallbacks. |

---

## Quick Start

### Installation
```bash
go get [github.com/grimdork/climate](https://github.com/grimdork/climate)
```

Or just type \<package>.\<function> and let `goimports` add it and `go mod tidy` download it like a sane person would.

## Argument Parsing (arg)
Handle flags, subcommands, and positional arguments without the overhead of heavy frameworks.

```go
package main

import "github.com/grimdork/climate/arg"

func main() {
	p := arg.New("mytool")
	p.AddOption("v", "verbose", "Enable verbose output", false)
	cmd := p.AddCommand("greet", "Greet a user")
	cmd.AddPositional("name", "The name to greet", true)
	p.Parse()
}
```

## Color formatting (cfmt)
Use ANSI colours with Printf.

```go
import "github.com/grimdork/climate/cfmt"

cfmt.Printf("{{.Red}}Error:{{.Reset}} file not found: '%s'\n", fn)
cfmt.Print("{{.Green}}Success!{{.Reset}} Process completed.\n")
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

// Use environment variable or fall back to "localhost"
host := env.Get("DB_HOST", "localhost")
```

## Core Philosophy
- Standard library only: Don't pull in other 3rd party dependencies.
- TinyGo compatibility: Prioritise low-allocation and small-binary footprints. While reflect is used sparingly in some packages, the core logic is designed to be as "TinyGo-friendly" as possible.
- Unix focus: Made for the systems I actually use every day.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
