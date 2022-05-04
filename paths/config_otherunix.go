//go:build aix || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris

package paths

import (
	"os/user"
	"path/filepath"
	"strings"
)

func basePath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return u.HomeDir, nil
}

func baseServerPath() (string, error) {
	return "/etc", nil
}

// Path to application-specific configuration directory.
func (cp *ConfigPath) Path() string {
	if cp.dirty {
		cp.Reset()
		cp.WriteString(filepath.Join(cp.base, "."+strings.ToLower(cp.name)))
		cp.dirty = false
	}

	return cp.String()
}
