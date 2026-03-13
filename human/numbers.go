package human

import (
	"fmt"
	"math"
)

// UInt formats an unsigned 64-bit number into bytes, kilobytes etc. strings, optionally with SI numbers.
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

	pre := "kMGTPE"[exp-1 : exp]
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