package env

import (
	"math"
	"os"
	"strconv"
	"strings"
)

// helper: strip surrounding single or double quotes
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// Get environment variable or alternative string.
// If the variable is set (even to an empty string) its value is returned.
// If the variable is not set at all, alt is returned.
func Get(key, alt string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return alt
}

// GetBool fetches an environment variable and parses it as a boolean.
// If the variable is not set, the provided alt value is returned.
// When the variable is set, truthiness rules apply: only explicit "true"-like
// values are considered true. If the variable is set but not a recognized
// truthy value, the alt fallback is returned (i.e. invalid -> alt).
// Accepted true values (case-insensitive): "true", "yes", "on", "1", "t".
func GetBool(key string, alt bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return alt
	}
	v = strings.TrimSpace(stripQuotes(v))
	if v == "" {
		return alt
	}
	v = strings.ToLower(v)
	switch v {
	case "true", "yes", "on", "1", "t":
		return true
	case "false", "no", "off", "0", "f":
		return false
	default:
		return alt
	}
}

// GetInt fetches an environment variable and parses it as an int64.
// If the variable is not set, alt is returned. If it is set but fails to
// parse (including overflow), alt is returned. Supports underscores, a
// leading sign (+/-), and 0x/0X hex prefixes.
func GetInt(key string, alt int64) int64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return alt
	}
	v = strings.TrimSpace(stripQuotes(v))
	if v == "" {
		return alt
	}
	// Remove underscores which are allowed for readability.
	v = strings.ReplaceAll(v, "_", "")
	// Parse with base 0 so 0x... is accepted as hex.
	x, err := strconv.ParseInt(v, 0, 64)
	if err != nil {
		return alt
	}
	return x
}

// GetFloat fetches an environment variable and parses it as a float64.
// If the variable is not set, alt is returned. If it is set but fails to
// parse (including overflow/NaN/Inf), alt is returned. Supports underscores
// and treats a single comma as decimal separator when no dot is present.
func GetFloat(key string, alt float64) float64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return alt
	}
	v = strings.TrimSpace(stripQuotes(v))
	if v == "" {
		return alt
	}
	// Remove underscores
	v = strings.ReplaceAll(v, "_", "")
	// If there is a comma but no dot, treat comma as decimal separator.
	if strings.Contains(v, ",") && !strings.Contains(v, ".") {
		v = strings.ReplaceAll(v, ",", ".")
	}
	x, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return alt
	}
	// Treat NaN/Inf as invalid -> alt
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return alt
	}
	return x
}
