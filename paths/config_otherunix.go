//go:build aix || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris

package paths

import (
	"os/user"
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
