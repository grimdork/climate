# climate/daemon
Simple utilities for long-running processes and graceful shutdowns.

`daemon` provides helpers for programs that need to run in the background or stay active until interrupted. Its primary feature is a clean, channel-based way to handle system interrupts (like `Ctrl+C`).

## Installation
```bash
go get github.com/grimdork/climate/daemon
```

## Core Features

### Graceful Shutdown with BreakChannel
Instead of manually setting up os/signal notify patterns, BreakChannel returns a channel that blocks until the program receives a termination signal (SIGINT or SIGTERM).

```go
package main

import (
	"fmt"
	"github.com/grimdork/climate/daemon"
)

func main() {
	println("Server starting... Press Ctrl+C to stop.")

	// Do your setup here (start database, listeners, etc.)

	// Block until Ctrl+C is pressed
	<-daemon.BreakChannel()

	println("\nShutting down gracefully...")
	// Perform cleanup here
}
```

### TODO
- Context version
