package daemon

import "errors"

// ErrNotRoot is returned when the user doesn't have root privileges.
var ErrNotRoot = errors.New("not running as root, so can't drop privileges to specified user")
