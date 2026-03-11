package loglines

import (
	"fmt"
	"github.com/grimdork/climate/cfmt"
	"github.com/grimdork/climate/str"
)

// CMsg prints a colourised message to stdout. Timestamp is grey; message text is cyan by default.
func CMsg(f string, v ...any) {
	b := str.NewStringer()
	b.WriteStrings(NowString(), ": ")
	msg := fmt.Sprintf(f, v...)
	// Use %grey for timestamp and %cyan for message body, and reset at the end.
	cfmt.Println("%grey" + b.String() + "%reset %cyan" + msg + "%reset")
}
