package main

import (
	"fmt"
)

// tinygo-smoke: minimal program that avoids filesystem and OS features.
// Keeps dependencies tiny and TinyGo-friendly.
func main() {
	fmt.Println("tinygo smoke: ok")
}
