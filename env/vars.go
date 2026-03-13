package env

import "os"

// Get environment variable or alternative string.
func Get(key, alt string) string {
	s := os.Getenv(key)
	if s == "" {
		return alt
	}

	return s
}
