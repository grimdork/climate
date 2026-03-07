# climate/env
Lightweight environment variable fetching with fallbacks.

`env` provides a clean, single-function interface for retrieving environment variables. It eliminates the boilerplate of checking for empty strings and manually assigning default values.

## Installation
```bash
go get [github.com/grimdork/climate/env](https://github.com/grimdork/climate/env)
```

## Why use env?
Standard Go os.Getenv returns an empty string if a variable is not set. This often leads to repetitive code:

```go
// The "verbose" way
port := os.Getenv("PORT")
if port == "" {
	port = "8080"
}
```

With climate/env, this becomes a one-liner:

```go
port := env.Get("PORT", "8080")
```

## Usage

### Basic fetching
Retrieve a value or use a supplied alternative if the environment variable is unset or empty.

```go
package main

import (
	"fmt"
	"github.com/grimdork/climate/env"
)

func main() {
	dbUser := env.Get("DB_USER", "postgres")
	mode := env.Get("APP_MODE", "development")
	fmt.Printf("Connecting as %s in %s mode...\n", dbUser, mode)
}
```

### Integration with CLI flags
A common pattern is to use env to set the default value of a CLI flag, allowing users to configure the app via both environment variables and command-line arguments.

```go
// Use the environment variable as the default for an 'arg' option
p.AddOption("p", "port", "Server port", env.Get("PORT", "8080"))
```

## Best Practices
Environment variables are the preferred way to handle secrets (like API keys) and environment-specific settings (like database URLs). Using env.Get ensures that your application has sensible defaults for development while being fully configurable in production environments like Kubernetes or Docker.
