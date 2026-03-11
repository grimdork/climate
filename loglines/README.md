# climate/loglines
Opinionated, simple logging utilities for small CLIs and services.

`loglines` provides three lightweight helpers for human-friendly output:

- `Msg(...)` — info-style lines to stdout with a timestamp
- `CMsg(...)` — colourised info lines using `cfmt` (timestamp grey, message cyan)
- `Err(...)` — error-style lines to stderr with a timestamp

Timestamps use the format: `<dayname> <mon> <day> HH:MM:SS.<microseconds> <year>`.

## Installation
```bash
go get github.com/grimdork/climate/loglines
```

## Usage
```go
package main

import (
	ll "github.com/grimdork/climate/loglines"
)

func main() {
	ll.Msg("Service started on port %d", 8080)
	ll.CMsg("Listening for requests on %s", ":8080")
	ll.Err("Failed to bind: %v", err)
}
```

## Notes
- `CMsg` uses `cfmt` under the hood and therefore respects the `NO_COLOR` environment variable — when `NO_COLOR` is set, the output falls back to plain text.
- Keep logging simple and human-readable; for machine-readable logs, consider writing JSON to a file or piping to a dedicated collector.
- The package is intentionally minimal and dependency-free (aside from internal climate packages) to stay TinyGo-friendly.
