# climate/arg
A minimalist command-line argument and option parser.

`arg` provides a simple way to handle flags, commands, and positional arguments with zero dependencies. It is fully compatible with TinyGo.

## Installation
```bash
go get github.com/grimdork/climate/arg
```

## Quick Start

```go
package main

import (
	"fmt"
	"os"
	"github.com/grimdork/climate/arg"
)

func main() {
	opt := arg.New("greet", "Say hello to someone")
	opt.SetDefaultHelp(true)
	opt.SetOption(arg.GroupDefault, "v", "verbose", "Show more details.", false, false, arg.VarBool, nil)
	opt.SetPositional("name", "Name to greet.", "world", true, arg.VarString)

	err := opt.Parse(os.Args[1:])
	if err != nil {
		if err == arg.ErrNoArgs {
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if opt.GetBool("verbose") {
		fmt.Println("Debug: Starting greeting...")
	}
	fmt.Printf("Hello, %s!\n", opt.GetPosString("name"))
}
```

## Features
- Short & long flags: `-v` or `--verbose`
- Combined short flags: `-abc` (last one can take a value)
- Subcommands with aliases
- Positional arguments (including slice types)
- Choice restriction on options
- Environment variable parsing
- Bash completion generation
- TinyGo compatible

## Options

### SetOption
```go
// SetOption(group, short, long, help, default, required, type, choices)
opt.SetOption(arg.GroupDefault, "p", "port", "Server port.", 8080, false, arg.VarInt, nil)
opt.SetOption(arg.GroupDefault, "m", "mode", "Run mode.", "dev", false, arg.VarString, []any{"dev", "prod", "test"})
```

### SetFlag (shorthand for bool options)
```go
opt.SetFlag(arg.GroupDefault, "v", "verbose", "Enable verbose output.")
```

### Retrieving values
```go
opt.GetBool("v")       // or opt.GetBool("verbose")
opt.GetString("mode")
opt.GetInt("port")
opt.GetFloat("ratio")
```

Short and long names both work for lookup.

## Positional arguments

```go
opt.SetPositional("PASSWORD", "The password to hash.", "", true, arg.VarString)
opt.SetPositional("ROUNDS", "Hash rounds.", 50000, false, arg.VarInt)

// Slice positional — swallows all remaining args
opt.SetPositional("FILES", "Input files.", nil, false, arg.VarStringSlice)
```

Retrieve with `GetPosString`, `GetPosInt`, `GetPosFloat`, `GetPosBool`, and their slice variants.

## Subcommands

```go
func main() {
	opt := arg.New("server", "Run the server and handle requests")
	opt.SetDefaultHelp(true)
	opt.SetFlag(arg.GroupDefault, "v", "verbose", "Verbose logging.")

	cmd := opt.SetCommand("serve", "Start the server.", arg.GroupDefault, serve, []string{"srv"})
	cmd.Options.SetOption(arg.GroupDefault, "p", "port", "Listen port.", 8080, false, arg.VarInt, nil)

	opt.SetCommand("config", "Show configuration.", arg.GroupDefault, config, []string{"cfg"})

	err := opt.Parse(os.Args[1:])
	if err != nil {
		if err == arg.ErrNoArgs || err == arg.ErrRunCommand {
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// Receives remaining args after the command name
func serve(opts *arg.Options) error {
	port := opts.GetInt("port")
	fmt.Printf("Listening on :%d\n", port)
	return nil
}

func config(opts *arg.Options) error {
	fmt.Println("Current config...")
	return nil
}
```

## Environment variable parsing

Map long option names to environment variables:
```go
// Without prefix: --host reads from HOST
opt.ParseEnvironment("", ",")

// With prefix: --host reads from MYAPP_HOST
opt.ParseEnvironment("myapp", ",")
```

Slice types use the delimiter (default `,`) to split values.

## Convenience methods

### HelpOrFail
Parses args, prints help on no args, exits on error:
```go
opt.HelpOrFail() // uses os.Args[1:]
```

### Bash completions
Generate a completion script for your tool:
```go
script, err := opt.Completions()
```

## Option groups
Organise options into named groups for cleaner help output:

### Tests
This project includes unit tests for the `arg` package that cover option parsing, choice validation, positional arguments, completions, and more. Run them with:

```bash
go test ./arg -v
```

Contributions adding edge-case tests for long option parsing and environment variable behaviour are welcome.
```go
opt.AddGroup("Database")
opt.SetOption("Database", "H", "db-host", "Database host.", "localhost", false, arg.VarString, nil)
opt.SetOption("Database", "P", "db-port", "Database port.", 5432, false, arg.VarInt, nil)
```

## Error handling
- `arg.ErrNoArgs` — no arguments provided (help printed if enabled)
- `arg.ErrRunCommand` — a subcommand was matched and executed
- `arg.ErrIllegalChoice` — value not in the allowed choices list
- `arg.ErrMissingRequired` — a required option was not provided
- `arg.ErrMissingParam` — an option that expects a value didn't get one
- `arg.ErrUnknownOption` — unrecognised flag encountered

### Examples: parsing failures
Below are brief examples showing how parsing errors surface. The library returns `error` values; callers should inspect them and handle logging, help display, or program exit as appropriate.

1) Missing required option
```go
opt := arg.New("app", "Example application")
opt.SetOption(arg.GroupDefault, "u", "user", "User name", "", true, arg.VarString, nil)
err := opt.Parse([]string{"app"})
if err != nil {
    if errors.Is(err, arg.ErrMissingRequired) {
        fmt.Fprintln(os.Stderr, "Required option missing: --user")
        opt.PrintHelp()
        os.Exit(2)
    }
}
```

2) Illegal choice value
```go
opt := arg.New("app", "Example application")
opt.SetOption(arg.GroupDefault, "m", "mode", "Run mode", "dev", false, arg.VarString, []any{"dev","prod"})
err := opt.Parse([]string{"app", "--mode", "test"})
if err != nil {
    if errors.Is(err, arg.ErrIllegalChoice) {
        fmt.Fprintf(os.Stderr, "Invalid choice for --mode: %v\n", err)
        os.Exit(2)
    }
}
```

3) Unknown option (example of how unknown flags are reported)
```go
opt := arg.New("app", "Example application")
err := opt.Parse([]string{"app", "--nope"})
if err != nil {
    if errors.Is(err, arg.ErrUnknownOption) {
        fmt.Fprintf(os.Stderr, "Unknown option: %v\n", err)
        opt.PrintHelp()
        os.Exit(2)
    }
}
```

### Completion failures
Generation of shell completion scripts can fail only due to I/O issues or internal state problems. When using `Completions()`, check the returned error and surface it.

```go
script, err := opt.Completions()
if err != nil {
    // I/O or state error generating the completion script
    fmt.Fprintf(os.Stderr, "Failed to generate completions: %v\n", err)
    os.Exit(2)
}
fmt.Println(script)
```

## Variable types
| Constant | Go type | Getter |
| :--- | :--- | :--- |
| `VarBool` | `bool` | `GetBool` / `GetPosBool` |
| `VarInt` | `int` | `GetInt` / `GetPosInt` |
| `VarIntSlice` | `[]int` | `GetPosIntSlice` |
| `VarFloat` | `float64` | `GetFloat` / `GetPosFloat` |
| `VarFloatSlice` | `[]float64` | `GetPosFloatSlice` |
| `VarString` | `string` | `GetString` / `GetPosString` |
| `VarStringSlice` | `[]string` | `GetStringSlice` / `GetPosStringSlice` |
