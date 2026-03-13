package paths

import (
	"path/filepath"
)

// Paths holds OS-specific paths for data.
type Paths struct {
	// AppName is used to build the config paths.
	AppName string
	// UserBase is the base path for user apps.
	UserBase string
	// ServerBase is the base path for server apps.
	ServerBase string
}

// New returns a Paths struct.
func New(appname string) (*Paths, error) {
	cp := &Paths{
		AppName: appname,
	}

	return cp, cp.setup()
}

func (cp *Paths) setup() error {
	err := cp.SetBase("")
	if err != nil {
		return err
	}

	err = cp.SetServerBase("")
	return err
}

// SetBase to something different than the default.
func (cp *Paths) SetBase(s string) error {
	var err error
	if s == "" {
		cp.UserBase, err = basePath()
		if err != nil {
			return err
		}

		cp.UserBase = filepath.Join(cp.UserBase, cp.AppName)
	} else {
		cp.UserBase = filepath.Join(s, cp.AppName)
	}

	return nil
}

// SetServerBase to something different than the default.
func (cp *Paths) SetServerBase(s string) error {
	var err error
	if s == "" {
		cp.ServerBase, err = baseServerPath()
		if err != nil {
			return err
		}

		cp.ServerBase = filepath.Join(cp.ServerBase, cp.AppName)
	} else {
		cp.ServerBase = filepath.Join(s, cp.AppName)
	}

	return nil
}