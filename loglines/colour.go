package loglines

import (
	"github.com/grimdork/climate/cfmt"
	"github.com/grimdork/climate/str"
)

// CMsg prints a colourised message to stdout. Timestamp is grey; message text is cyan by default.
// It respects colour tags in the format string and supports fmt verbs by using cfmt.Sprintf.
func CMsg(f string, v ...any) {
	b := str.NewStringer()
	b.WriteStrings(NowString(), ": ")
	msg := cfmt.Sprintf(f, v...)
	// Use %grey for timestamp and %cyan for message body, and reset at the end.
	cfmt.Println("%grey" + b.String() + "%reset %cyan" + msg + "%reset")
}