package loglines

import (
	"fmt"
	"os"

	"github.com/grimdork/climate/str"
)

// Err prints formatted messages to stderr, starting with a nicely formatted timestamp.
func Err(f string, v ...any) {
	b := str.NewStringer()
	b.WriteStrings(
		NowString(),
		":",
		fmt.Sprintf(f, v...),
	)
	fmt.Fprintln(os.Stderr, b.String())
}
