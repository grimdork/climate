//go:build darwin

package paths

import (
	"os/user"
	"path/filepath"
)

func basePath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(u.HomeDir, "Library", "Application Support")
	return dir, nil
}

func baseServerPath() (string, error) {
	return basePath()
}
