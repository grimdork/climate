package human

import (
	"fmt"
	"math"
)

// Float formats a floating-point number with the specified precision, optionally using SI or
// IEC (binary) prefixes. When si==true, uses 1000-based prefixes (k, M, G...). When
// si==false, uses 1024-based IEC prefixes (Ki, Mi, Gi...).
func Float(f float64, prec int, si bool) string {
	unit := 1000.0
	var prefixes []string
	if si {
		unit = 1000.0
		prefixes = []string{"k", "M", "G", "T", "P", "E"}
	} else {
		unit = 1024.0
		prefixes = []string{"Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}
	}

	if f < unit {
		// Print without any prefix
		return fmt.Sprintf("%.*f", prec, f)
	}

	exp := int(math.Log(f) / math.Log(unit))
	if exp > len(prefixes) {
		exp = len(prefixes)
	}

	pre := prefixes[exp-1]
	val := f / math.Pow(unit, float64(exp))

	// If the value is a whole number and precision is 0 or 1, print without trailing .0
	if prec <= 0 {
		return fmt.Sprintf("%.0f %s", val, pre)
	}

	// When precision > 0, show requested decimals. If the fractional part is zero, and
	// precision == 1, still show a single decimal to match previous behaviour for UInt.
	return fmt.Sprintf("%.*f %s", prec, val, pre)
}
