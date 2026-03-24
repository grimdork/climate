package human

import (
	"fmt"
	"math"
)

// Float formats a floating-point number with the specified precision and returns the
// value using the same unit prefixes and casing as UInt so callers can rely on
// identical output for human-readable sizes.
//
// When si==true, uses 1000-based prefixes (k, M, G...) and a trailing "B" suffix
// (e.g. "1.5 MB"). When si==false, uses legacy binary/IEC behaviour with a
// lowercase kilo together with the "iB" suffix for 1024 ("kiB") and uppercase
// letters for larger prefixes ("MiB", "GiB"), matching UInt.
func Float(f float64, prec int, si bool) string {
	unit := 1024.0
	if si {
		unit = 1000.0
	}

	// Under the unit size, show the raw number with a B suffix and respect whole numbers
	if f < unit {
		if f == math.Floor(f) {
			return fmt.Sprintf("%.0f B", f)
		}
		return fmt.Sprintf("%.*f B", prec, f)
	}

	exp := int(math.Log(f) / math.Log(unit))
	if exp > 6 {
		exp = 6
	}

	var prefixes []string
	if si {
		prefixes = []string{"k", "M", "G", "T", "P", "E"}
	} else {
		// Preserve legacy kilo casing for binary ("k" + "iB") and use uppercase for larger
		// prefixes to maintain backward compatibility with existing output.
		prefixes = []string{"k", "M", "G", "T", "P", "E"}
	}

	pre := prefixes[exp-1]
	suffix := "B"
	if !si {
		suffix = "iB"
	}

	val := f / math.Pow(unit, float64(exp))

	// Match UInt behaviour: if value is whole, don't show decimals
	if val == math.Floor(val) {
		return fmt.Sprintf("%.0f %s%s", val, pre, suffix)
	}

	// Otherwise honor the requested precision
	return fmt.Sprintf("%.*f %s%s", prec, val, pre, suffix)
}
