package human

import (
	"fmt"
	"math"
)

// UInt formats an unsigned 64-bit number into bytes, kilobytes etc. strings, optionally with SI numbers.
// When si==true, uses 1000-based prefixes (k, M, G...). When si==false, uses IEC 1024-based prefixes (Ki, Mi, Gi...).
func UInt(n uint64, si bool) string {
	num := float64(n)
	var unit float64 = 1024
	if si {
		unit = 1000
	}

	if num < unit {
		return fmt.Sprintf("%.0f B", num)
	}

	exp := int(math.Log(num) / math.Log(unit))
	// Safety check for prefix range
	if exp > 6 {
		exp = 6
	}

	var prefixes []string
	if si {
		prefixes = []string{"k", "M", "G", "T", "P", "E"}
	} else {
		// Historical behaviour uses a lowercase 'k' together with the 'iB' suffix
		// (e.g. "kiB") for the kilo prefix, and uppercase letters for larger
		// prefixes (MiB, GiB...). Preserve that to match existing expectations.
		prefixes = []string{"k", "M", "G", "T", "P", "E"}
	}

	pre := prefixes[exp-1]
	// Use a 'B' suffix for SI and 'iB' for binary (legacy style: "kiB", "MiB").
	suffix := "B"
	if !si {
		suffix = "iB"
	}

	// Calculate value with one decimal place for better readability
	val := num / math.Pow(unit, float64(exp))

	// If it's a whole number, don't show .0
	if val == math.Floor(val) {
		return fmt.Sprintf("%.0f %s%s", val, pre, suffix)
	}
	return fmt.Sprintf("%.1f %s%s", val, pre, suffix)
}
