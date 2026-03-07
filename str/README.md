# climate/str
A type-aware, recursive extension of strings.Builder.

`str.Stringer` provides a fluent API for building complex strings from mixed data types. Unlike the standard library, it can recursively handle slices and maps, automatically formatting them based on your configuration.

## Installation
```bash
go get [github.com/grimdork/climate/str](https://github.com/grimdork/climate/str)
```

## Features
- Varargs support: Write multiple strings or mixed types in a single call.
- Recursive formatting: Deeply nested slices and maps are automatically traversed.
- Fluent configuration: Chain settings for delimiters (commas, equals signs).
- Zero external dependencies: Pure Go standard library (using reflect for type discovery).

## Usage Examples

### Mixed type writing (WriteI)
WriteI (Write Interface) accepts any number of arguments of different types.

```go
s := str.NewStringer()
s.WriteI("ID: ", 101, " | Active: ", true, " | Score: ", 98.6)

fmt.Println(s.String())
// Output: ID: 101 | Active: true | Score: 98.6
```

### Slices with custom delimiters
Enable SetSliceComma to automatically separate slice elements.

```go
s := str.NewStringer()
s.SetSliceComma(true).SetComma(';')

tags := []string{"go", "cli", "minimal"}
s.WriteI(tags)

fmt.Println(s.String())
// Output: go;cli;minimal
```

### Map Serialization
Maps are written as key=value. You can customize both the key-value joiner and the pair separator.

```go
s := str.NewStringer()
s.SetMapComma(true).SetEquals(':').SetComma('|')

metadata := map[string]int{"cpu": 4, "ram": 16}
s.WriteI(metadata)

fmt.Println(s.String())
// Output: cpu:4|ram:16 (NOTE: Map order in Go is random)
```

### Recursive depth
Since writeX is recursive, it can handle complex nested structures like slices of maps:

```go
s := str.NewStringer().SetSliceComma(true)
data := []map[string]string{
	{"name": "App"},
	{"version": "1.0"},
}
s.WriteI("Config: ", data)

fmt.Println(s.String())
// Output: Config: name=App,version=1.0
```

### Configuration Methods
| Method | Description | Default |
| :---- | :---- | :---- |
|SetSliceComma(bool) | Enable/disable commas between slice elements. | false |
|SetMapComma(bool) | Enable/disable commas between map pairs. | false |
|SetComma(byte) | Set the separator symbol (used for slices/maps). | , |
|SetEquals(byte) | Set the symbol joining keys and values. | = |

### Implementation notes
- Integers: Supports int and int64.
- Floats: Interpreted as float64 with precision optimized for readability.
- Performance: While it uses reflect to handle any types, it maintains the underlying performance benefits of strings.Builder for memory allocation.
