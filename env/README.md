# climate/env
Lightweight environment variable helpers with fallbacks and typed parsing.

`env` provides simple helpers to fetch environment variables with sensible defaults and common parsing conveniences (underscores in numbers, hex ints, localized decimals).

## Installation
```bash
go get github.com/grimdork/climate/env
```

## API
- Get(key, alt string) string
  - Returns the environment value if the variable is set (returns empty string if set to empty). If the variable is not set, returns alt.
- GetBool(key string, alt bool) bool
  - If the variable is not set, returns alt. If set, returns true only for explicit truthy values (true/yes/on/1/t); everything else (including empty, quoted non-truthy values, or invalid tokens) is false.
- GetInt(key string, alt int64) int64
  - If the variable is not set, returns alt. If set, attempts to parse as int64 (supports +/-, underscores, and 0x/0X hex). On parse failure or overflow, returns 0.
- GetFloat(key string, alt float64) float64
  - If the variable is not set, returns alt. If set, attempts to parse as float64 (supports underscores and treats a single comma as decimal separator when no dot is present). On parse failure or NaN/Inf, returns 0.

## Examples

Basic fetching:

```go
import "github.com/grimdork/climate/env"

port := env.Get("PORT", "8080")
```

Typed helpers:

```go
// Bool with fallback; when set to an invalid value, result is false
debug := env.GetBool("DEBUG", false)

// Int with fallback (accepts hex, underscores); invalid -> 0
workers := env.GetInt("WORKERS", 4)

// Float with fallback (accepts comma as decimal when no dot present); invalid -> 0
threshold := env.GetFloat("THRESHOLD", 0.75)
```

Integration with ini (fallback pattern):

```go
// Pass the INI-provided value as the alt/fallback so env overrides it when set.
val := env.Get("SOME_KEY", iniValue)
num := env.GetInt("NUM_KEY", iniNum)
```

## Notes
- The typed helpers distinguish unset (returns alt) from set-but-invalid (returns false/0). This makes them suitable for directly embedding in numeric or boolean calls without extra error handling.
- Unit tests for the typed helpers live in `env/vars_test.go` and exercise common formats and edge cases.
