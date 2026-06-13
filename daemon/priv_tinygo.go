//go:build tinygo

package daemon

// DegradeToUser is not supported in TinyGo.
func DegradeToUser(uname string) error {
	return ErrNotRoot
}
