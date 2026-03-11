# climate/str
A type-aware extension of strings.Builder.

`str.Stringer` provides a fluent API for building complex strings from mixed data types. It handles slices and maps with configurable delimiters, using type switches for zero-dependency, TinyGo-friendly operation.

## Installation
```bash
go get github.com/grimdork/climate/str
```

## Usage

### Mixed type writing (WriteI)
`WriteI` accepts any number of arguments of different types:
```go
s := str.NewStringer()
s.WriteI("ID: ", 101, " | Active: ", true, " | Score: ", 98.6)

fmt.Println(s.String())
// Output: ID: 101 | Active: true | Score: 98.6
```

### Multiple strings (WriteStrings)
When you only have strings:
```go
s := str.NewStringer()
s.WriteStrings("hello", " ", "world")
// Output: hello world
```

### Slices with custom delimiters
Enable `SetSliceComma` to separate slice elements:
```go
s := str.NewStringer()
s.SetSliceComma(true).SetComma(';')

tags := []string{"go", "cli", "minimal"}
s.WriteI(tags)

fmt.Println(s.String())
// Output: go;cli;minimal
```

### Map serialization
Maps are written as key=value pairs. Customize the joiner and separator:
```go
s := str.NewStringer()
s.SetMapComma(true).SetEquals(':').SetComma('|')

metadata := map[string]int{"cpu": 4, "ram": 16}
s.WriteI(metadata)

fmt.Println(s.String())
// Output: cpu:4|ram:16 (map order is random)
```

### Supported types
| Category | Types |
| :--- | :--- |
| Scalars | `bool`, `string`, `int`, `int64`, `float64` |
| Slices | `[]bool`, `[]string`, `[]int`, `[]int64`, `[]float64`, `[]any` |
| Maps | `map[string]string`, `map[string]int`, `map[string]any`, `map[int]string`, `map[int]int`, `map[int]any` |

Unsupported types are silently skipped.

### Configuration methods
| Method | Description | Default |
| :--- | :--- | :--- |
| `SetSliceComma(bool)` | Enable/disable separator between slice elements | `false` |
| `SetMapComma(bool)` | Enable/disable separator between map pairs | `false` |
| `SetComma(byte)` | Set the separator character | `,` |
| `SetEquals(byte)` | Set the key-value joiner character | `=` |

All setters return `*Stringer` for chaining.
