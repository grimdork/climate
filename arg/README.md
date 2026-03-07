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
	"github.com/grimdork/climate/arg"
)

func main() {
	opt := arg.New("greet")
	// Enable -h and --help automatically
	opt.SetDefaultHelp(true)
	// Set a non-required boolean that defaults to false.
	opt.SetOption(arg.GroupDefault, "v", "verbose", "Show more details in output.", false, false, arg.VarBool, nil)
	// Fields are: Display name in help, description, default, required, type set to string
	err := opt.SetPositional("name", "Name to greet.", "tester", true, arg.VarString)
	// That error should be handled, but we skip this for brevity
	err = opt.Parse(os.Args[1:])
	if err != nil {
		if err == arg.ErrNoArgs {
			// We got no args. Help text has already been displayed for us.
			return
		}

		// Handle different errors here
	}

	if opt.GetBool("verbose") {
		println("Debug: Starting greeting...")
	}
	fmt.Printf("Hello, %s!\n", opt.GetString("name"))
}
```

## Features
- Short & long flags: -v or --verbose.
- Parse from environment.
- Commands: Support for nested sub-commands.
- TinyGo Optimised: Should actually compile and run.

## Error handling
The Parse method returns specific error constants to help you control flow:
- arg.ErrNoArgs: Triggered when no arguments are provided (and help is shown, if enabled).
- arg.ErrRunCommand: Returned when a subcommand was successfully matched and executed.

Otherwise, the error returned is either an internal error or the returned error from a subcommand.

## Examples

### Simple flags and help
```go
	opt := arg.New("cmd")
	// Define a boolean flag.
	err := opt.SetOption(arg.GroupDefault, "v", "version", "Display the version and exit.", false, false, arg.VarBool, nil)
	// Ideally, handle errors here
	// Define a string option.
	err = opt.SetOption(arg.GroupDefault, "c", "config", "Specify an alternative configuration file.", false, false, arg.VarString, nil)
	// Define a short option with no long alternative.
	err = opt.SetOption(arg.GroupDefault, "x", "", "Extra flags.", false, false, arg.VarString, nil)

	// Parse, handle results
	err = opt.Parse(os.Args[1:])
	if err != nil {
		if err == arg.ErrNoArgs {
			return
		}

		// If it reaches this point, check the error and handle accordingly.
	}

	// Check for the short flag name
	if opt.GetBool("v") { ... }

	// You can also get an option by its long name
	cfgname := opt.GetString("config")

	// Optional method: Get the raw option for whichever dark purpose
	x := opt.GetOption("x")
```

### Positional arguments
```go
	opt := arg.New("pwhasher")
	opt.SetDefaultHelp(true)
	// First argument must be a string
	opt.SetPositional("PASSWORD", "The password to hash", "", true, arg.VarString)
	// Second also a string
	opt.SetPositional("HASH", "Hash to apply.", "", true, arg.VarString)
	// Third is parsed as a number
	opt.SetPositional("ROUNDS", "Number of rounds to apply.", 50000, false, arg.VarInt)
	err := opt.Parse(os.Args[1:])
	if err != nil || len(os.Args) < 3 {
			opt.PrintHelp()
			os.Exit(1)
	}

	// It's simple to get them
	s := GeneratePasswordHash(opt.GetPosString("PASSWORD"), opt.GetPosString("HASH"), opt.GetPosInt("ROUNDS"))
	fmt.Println(s)
```

### Subcommand
```go
func main() {
	opt := arg.New("testserver")
	opt.SetDefaultHelp(true)
	// The commands "serve" and "srv" will both work for this.
	opt.SetCommand("serve", "Start the server.", arg.GroupDefault, serve, []string{"srv"})
	// "config" is the command, and "cfg" its only alias.
	opt.SetCommand("config", "Configure the server.", arg.GroupDefault, config, []string{"cfg"})
	opt.AddOption("v", "verbose", "Verbose logging", false)
	// Parse and run commands
	err := opt.Parse(os.Args)
	if err != nil {
		if err == arg.ErrNoArgs {
			return
		}

		// If this is reached, either the function "serve" or "config" were called.
		if err == arg.ErrRunCommand {
			return
		}

		fmt.Printf("Error parsing arguments: %s\n", err.Error())
		os.Exit(1)
	}
}

// The args passed are what remains from the previous parsing. This can be nested deeper for sub-sub-commands.
func serve(args *arg.Options) error {
	// You can further define options here and reparse.
}
```
