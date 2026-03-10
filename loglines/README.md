# climate/loglines
Opinionaed, simple logging.

`loglines`has two functions for general log lines and errors, plus a function to make a very specifically formatted timestamp. The timestamp is in the format `<dayname> <mon> <day> HH:MM:SS.<nanoseconds> <year>` and is printed to stdout or stderr, depending on the function used.

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
	ll.Msg("This is a log message.")
	ll.Err("This is an error message.")
}
```
