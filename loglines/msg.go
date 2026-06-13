package loglines

import (
	"fmt"
	"io"
	"os"

	"github.com/grimdork/climate/str"
)

// write outputs a formatted message to w with a timestamp prefix.
func write(w io.Writer, f string, v ...any) {
	b := str.NewStringer()
	b.WriteStrings(
		NowString(),
		": ",
		fmt.Sprintf(f, v...),
	)
	fmt.Fprintln(w, b.String())
}

// Msg prints formatted messages to stdout, starting with a nicely formatted timestamp.
func Msg(f string, v ...any) {
	write(os.Stdout, f, v...)
}
