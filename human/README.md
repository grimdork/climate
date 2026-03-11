# climate/human
**Convert large numbers into friendly, readable strings.**

`human` takes raw numbers and formats them into the most appropriate unit (B, KB, MB, etc.). It supports both the **Binary (1024)** and **SI (1000)** standards.

## Installation
```bash
go get github.com/grimdork/climate/human
```

## Usage
### Binary Formatting (Default)
The old way, still preferred by some. It uses base-1024 and units like KiB, MiB, etc.

```Go
package main

import (
	"fmt"
	"github.com/grimdork/climate/human"
)

func main() {
	size := uint64(1572864)
	fmt.Println(human.UInt(size, false)) // Output: 1.5 MiB
}
```

### SI Formatting (Decimal)
Standardised by storage manufacturers and used more recently in operating systems. It uses base-1000.

```Go
size := uint64(1500000)
fmt.Println(human.UInt(size, true)) // Output: 1.5 MB
```

## Features
- Auto-scaling: Automatically chooses the best unit from Bytes to Exabytes.
- Smart precision: Shows decimals only when necessary for a cleaner look.
- TinyGo compatible: Uses standard `math` functions; no heavy dependencies.
- Zero dependencies: Pure Go standard library.

## Why use `human`?
Raw bytes are hard to read at a glance. `human` turns 124155123 into 118.4 MiB, making your CLI tool's output much more accessible to users.
