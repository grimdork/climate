package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// BreakChannel generates a goroutine which waits for ctrl-C to be pressed, and returns a channel to wait on.
// Using it can be as simple as:
//  <-daemon.BreakChannel()
func BreakChannel() chan bool {
	quit := make(chan bool)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		fmt.Print("\b\b") // Remove the CTRL-C symbol from stdout.
		quit <- true
	}()

	return quit
}
