package loglines

import (
	"fmt"

	"github.com/grimdork/climate/fx"
	"github.com/grimdork/climate/str"
)

// CMsg prints a colourised message to stdout. Timestamp is grey; message text is cyan by default.
// The format string supports Go fmt verbs (%s, %d, etc.) but not colour tags — use fx tags when needed.
func CMsg(f string, v ...any) {
	b := str.NewStringer()
	b.WriteStrings(NowString(), ": ")
	msg := fmt.Sprintf(f, v...)
	fx.Println("{grey}" + b.String() + "{@} {cyan}" + msg + "{@}")
}
