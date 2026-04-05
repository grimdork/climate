//go:build windows

package paths

import (
	"os"
	"path/filepath"
)

func basePath() (string, error) {
	appData := os.Getenv("AppData")
	if appData != "" {
		return appData, nil
	}

	return filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming"), nil
}

func baseServerPath() (string, error) {
	programData := os.Getenv("ProgramData")
	if programData != "" {
		return programData, nil
	}

	return filepath.Join(os.Getenv("SystemDrive"), "ProgramData"), nil
}
