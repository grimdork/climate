# climate
A collection of packages for command line apps in Go, with minimal dependencies.

## Notes
Only Unix-like systems are supported. Windows support is not planned until I need it.

## Packages included

### arg
This is a simple option parser with no dependencies outside the Go stdlib. The goal is to have the typical, most common use cases supported, including nested tool commands, and to be compilable with TinyGo.

Currently supports:
- Short and long options
- Commands
- Positional arguments
- Catch-all positional arguments (last args become a slice)

### cfmt
Colour formatting, printf style.

### daemon
A few utility functions for daemon style programs (servers). The most important function is BreakChannel, which returns a channel that waits for Ctrl-C and returns. Use it like this:
	<-daemon.BreakChannel()

### env
At the moment, just the convenience function Get(), which returns the value of an environment variable or a supplied alternative.

### human
Humanised numbers (and possibly other things in the future).

### paths
Helper to get typical directory paths for configuration data. Basically makes dot-directories in $HOME (most systems) for config directories in $HOME/Library/Application Support (macOS).

### prompter
A simple tool to ask question in the terminal. Provide questions, and optionally default answers and whether the reply should be hidden.

### str
Currently just an extended strings.Builder which can write out all common variable types, maps and slices with the same function. Also satisfies the io.Writer interface.
