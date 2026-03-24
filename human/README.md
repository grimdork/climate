# climate/human
**Convert large numbers into friendly, readable strings.**

`human` takes raw numbers and formats them into the most appropriate unit (B, KB, MB, etc.). It supports both the **Binary (1024)** and **SI (1000)** standards.

## Installation
```bash
go get github.com/grimdork/climate/human
```

## Usage
### Binary Formatting (Default)
The library preserves historical/legacy behaviour for some binary prefixes. It uses base-1024 and returns units in a legacy-friendly form that matches existing expectations in the codebase and tests:

- For kilo (1024) the output uses a lowercase `k` together with the `iB` suffix (e.g. `"1 kiB"`).
- For larger binary prefixes it uses the more familiar `MiB`, `GiB`, `TiB`, etc.

```Go
package main

import (
	"fmt"
	"github.com/grimdork/climate/human"
)

func main() {
	size := uint64(1572864)
	fmt.Println(human.UInt(size, false)) // Output: 1.5 MiB
	fmt.Println(human.UInt(1024, false)) // Output: 1 kiB (legacy kilo casing)
}
```

### SI Formatting (Decimal)
Standardised by storage manufacturers and used more recently in operating systems. It uses base-1000 and produces SI-style prefixes (`k`, `M`, `G`...).

```Go
size := uint64(1500000)
fmt.Println(human.UInt(size, true)) // Output: 1.5 MB
```

### Floating-point formatting
The Float function formats floating-point values using either SI (1000) or IEC/binary (1024) prefixes. Note the subtle differences between `Float` and `UInt`:

- Float(si=false) uses IEC-style prefixes like `Ki`, `Mi`, `Gi` (capital `K`) and returns the value followed by the prefix (no trailing `B` by default): e.g. `"1.5 Ki"`.
- UInt preserves legacy unit casing for kilo (returns `kiB`) and appends a `B`/`iB` suffix.

```go
fmt.Println(human.Float(1536, 1, false)) // -> "1.5 Ki"
fmt.Println(human.Float(1500, 1, true))  // -> "1.5 k"
```

## Features
- Auto-scaling: Automatically chooses the best unit from Bytes to Exabytes.
- Smart precision: Shows decimals only when necessary for a cleaner look.
- TinyGo compatible: Uses standard `math` functions; no heavy dependencies.
- Zero dependencies: Pure Go standard library.

## Why use `human`?
Raw bytes are hard to read at a glance. `human` turns 124155123 into 118.4 MiB, making your CLI tool's output much more accessible to users.
