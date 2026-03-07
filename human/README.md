# climate/human
**Convert large numbers into friendly, readable strings.**

`human` takes raw numbers and formats them into the most appropriate unit (B, KB, MB, etc.). It supports both the **Binary (1024)** and **SI (1000)** standards.

## Installation
```bash
go get [github.com/grimdork/climate/human](https://github.com/grimdork/climate/human)
Usage
1. Binary Formatting (Default)

Used by most operating systems (Windows, Linux) to represent file sizes. It uses the iB suffix (KiB, MiB) to denote base-1024.

Go
package main

import (
    "fmt"
    "[github.com/grimdork/climate/human](https://github.com/grimdork/climate/human)"
)

func main() {
    size := uint64(1572864)
    fmt.Println(human.UInt(size, false)) // Output: 1.5 MiB
}
2. SI Formatting (Decimal)

Standardized by storage manufacturers and used by macOS. It uses base-1000.

Go
size := uint64(1500000)
fmt.Println(human.UInt(size, true)) // Output: 1.5 MB
Features
Auto-Scaling: Automatically chooses the best unit from Bytes to Exabytes.

Smart Precision: Shows decimals only when necessary for a cleaner look.

TinyGo Compatible: Uses standard math functions; no heavy dependencies.

Zero Dependencies: Pure Go standard library.

Why use human?
Raw bytes are hard to read at a glance. human turns 124155123 into 118.4 MiB, making your CLI tool's output much more accessible to users.
