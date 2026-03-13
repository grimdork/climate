//go:build aix || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris

package paths

import (
	"os"
	"os/user"
	"path/filepath"
)

func basePath() (string, error) {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		return xdg, nil
	}

	// Fall back to $HOME/.config
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(u.HomeDir, ".config"), nil
}

func baseServerPath() (string, error) {
	return "/etc", nil
}
