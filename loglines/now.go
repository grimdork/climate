package loglines

import (
	"fmt"
	"time"
)

const timeFmt = "%s %s %02d %02d:%02d:%02d.%06d %04d"

// NowString returns a very detailed time string.
func NowString() string {
	t := time.Now()
	return fmt.Sprintf(timeFmt, t.Weekday().String()[0:3], t.Month().String()[0:3], t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Year())
}
