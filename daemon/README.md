# climate/daemon
Simple utilities for long-running processes and graceful shutdowns.

`daemon` provides helpers for programs that need to run in the background or stay active until interrupted.

## Installation
```bash
go get github.com/grimdork/climate/daemon
```

## Graceful shutdown with BreakChannel
Returns a channel that blocks until the program receives SIGINT or SIGTERM. Also cleans up the `^C` output from the terminal.

```go
package main

import (
	"fmt"
	"github.com/grimdork/climate/daemon"
)

func main() {
	fmt.Println("Server starting... Press Ctrl+C to stop.")

	// Start listeners, workers, etc.

	<-daemon.BreakChannel()

	fmt.Println("Shutting down gracefully...")
	// Cleanup here
}
```

## Privilege dropping with DegradeToUser
Drop from root to a specified user. Useful for services that bind to privileged ports before switching to a less privileged account.

```go
err := daemon.DegradeToUser("www-data")
if err != nil {
	// Either not running as root (daemon.ErrorNotRoot) or user lookup failed
	log.Fatal(err)
}
```

Sets the effective UID and primary GID of the specified user. Returns `daemon.ErrorNotRoot` if the process is not running as root.
