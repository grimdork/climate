package loglines

import (
	"fmt"
	"os"

	"github.com/grimdork/str"
)

// Err prints formatted messages to stderr, starting with a nicely formatted timestamp.
func Err(f string, v ...interface{}) {
	b := str.NewStringer()
	b.WriteStrings(
		NowString(),
		":",
		fmt.Sprintf(f, v...),
	)
	fmt.Fprintln(os.Stderr, b.String())
}
