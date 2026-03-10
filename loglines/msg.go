package loglines

import (
	"fmt"
	"os"

	"github.com/grimdork/str"
)

// Msg prints formatted messages to stdout, starting with a nicely formatted timestamp.
func Msg(f string, v ...interface{}) {
	b := str.NewStringer()
	b.WriteStrings(
		NowString(),
		":",
		fmt.Sprintf(f, v...),
	)
	fmt.Fprintln(os.Stdout, b.String())
}
