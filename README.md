# climate
A collection of packages for command line apps in Go, with minimal dependencies.

## Notes
Only Unix-like systems are supported. Windows support is not planned until I need it.

## Packages included

### arg
This is a simple option parser with no dependencies outside the Go stdlib. The goal is to have the typical, most common use cases supported, including nested tool commands, and to be compilable with TinyGo.

### paths
Helper to get typical directory paths for configuration data. Basically makes dot-directories in $HOME (most systems) for config directories in $HOME/Library/Application Support (macOS).
