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
		return fmt.Sprintf("%dB", int(num))
	}
	exp := int(math.Log(num) / math.Log(unit))
	pre := "kMGTPE"
	pre = pre[exp-1 : exp]
	if !si {
		pre = pre + "i"
	}
	r := n / uint64(math.Pow(unit, float64(exp)))
	return fmt.Sprintf("%d %sB", r, pre)
}
