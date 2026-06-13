package loglines

import "os"

// Err prints formatted messages to stderr, starting with a nicely formatted timestamp.
func Err(f string, v ...any) {
	write(os.Stderr, f, v...)
}
