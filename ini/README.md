# ini — simple INI helper (typed values, env overrides, schema hints)

A compact INI parser/serializer with typed fields (bool, int, float, string),
convenience getters, environment-aware overrides, and small helpers for numeric parsing
(hex, localized decimals, underscores in numeric literals).

## Features

- Typed fields: `Bool`, `Int`, `Float`, `String`
- Deterministic `Marshal()` output (sorted) and `Save(filename, tabbed)`
- Environment overrides with `ForceUpper()` support; enable env-first lookup with `SetEnvFirst(true)` so getters consult the environment before falling back to INI values
- Schema hints: `DeclareType(section, key, type)` + `Parse(s)` so you can declare expected types before parsing
- Numeric parsing enhancements: hex integers (`0x...`), localized decimal comma (`"3,14" -> 3.14`), underscores (`"1_000"`)
- `Unmarshal("")` returns an error (empty input is rejected)

## Install

Import in your Go code:

```go
import "github.com/grimdork/climate/ini"
```

## Quick examples

### 1) Simple `Unmarshal` + getters

```go
data := `
 i=42
 f=3.14
 name=example
 enabled=true
`

cfg, err := ini.Unmarshal(data)
if err != nil {
    // handle error
}
fmt.Println(cfg.GetInt("", "i"))      // 42
fmt.Println(cfg.GetFloat("", "f"))    // 3.14
fmt.Println(cfg.GetString("", "name")) // example
fmt.Println(cfg.GetBool("", "enabled"))
```

Notes: top-level properties use `section == ""` in the `Get*` helpers.

### 2) Pre-declare types (schema hints) then `Parse`

Use this when you need deterministic typing (e.g. hex numbers or to treat `1` as a boolean):

```go
cfg, _ := ini.New()
// Top-level key "port" should be parsed as an integer (accepts 0x... too)
cfg.DeclareType("", "port", ini.Int)
// Section "db" key "timeout" should be a float
cfg.DeclareType("db", "timeout", ini.Float)

data := `
port=0x1F
[db]
 timeout=3,5
`

if err := cfg.Parse(data); err != nil {
    // handle
}
fmt.Println(cfg.GetInt("", "port"))        // 31 (hex parsed)
fmt.Println(cfg.GetFloat("db", "timeout")) // 3.5 (comma treated as decimal)
```

### 3) Environment overrides and `ForceUpper()`

Enable environment-first lookup and call `ForceUpper()` to prefer upper-case environment variable names (common on UNIX):

```go
cfg, _ := ini.Unmarshal("debug=false")
cfg.SetEnvFirst(true)
cfg.ForceUpper()
// If DEBUG=true is set in the environment, GetBool will reflect that
fmt.Println(cfg.GetBool("", "debug"))
```

### 4) `Marshal` and `Save`

```go
cfg, _ := ini.Unmarshal("a=1\n[sec]\nb=2")
out := cfg.Marshal()
fmt.Println(out)

// Save to a file; if tabbed=true each field line is prefixed with a tab
_ = cfg.Save("config.ini", false)
```

## Numeric parsing rules

- Parser order (when not explicitly declared): boolean words (true/false/yes/no/on/off), integer, float.
- Integers may contain underscores which are ignored (`1_000` → `1000`).
- Integer literals that start with `0x`/`0X` are parsed as hex.
- If a numeric token contains a comma but no dot (e.g. `3,14`) the comma is treated as the decimal separator.
- Use `DeclareType` to force a particular type for a key.

## Environment helpers and interaction

The `ini` package leverages `github.com/grimdork/climate/env` for environment-based overrides. When `SetEnvFirst(true)` is enabled, `Get*` helpers consult environment variables first and then fall back to parsed INI values. The env helpers accept a fallback (second) argument so callers can pass the INI-provided value directly:

```go
// Example: let the environment override the INI value
cfg.SetEnvFirst(true)
val := env.Get("SOME_KEY", cfg.GetString("", "some_key"))
num := env.GetInt("NUM_KEY", cfg.GetInt("", "num_key"))
```

This keeps the behavior consistent and centralises parsing rules (underscores, hex, localized decimals) in the `env` package.

## Notes and recommendations

- Use `DeclareType` + `Parse` when you need deterministic typing. `Unmarshal` is a convenience wrapper that returns a fresh `*INI`.
- If you want different locale handling, we can add a configurable normalizer function on the `INI` instance (e.g. `SetNormalizer(func(string) string)`).
- Tests in the repo demonstrate numeric handling and race-free concurrent access. Run with:

```sh
go test -race ./...
```
