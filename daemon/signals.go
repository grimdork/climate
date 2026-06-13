package daemon

import (
	"os"
	"os/signal"
	"syscall"
)

// BreakChannel generates a goroutine which waits for ctrl-C to be pressed, and returns a channel to wait on.
// Using it can be as simple as:
//
//	<-daemon.BreakChannel()
func BreakChannel() chan bool {
	quit := make(chan bool, 1)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		quit <- true
	}()

	return quit
}
