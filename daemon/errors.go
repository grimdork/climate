package daemon

import "errors"

// ErrorNotRoot is returned when the user doesn't have root privileges.
var ErrorNotRoot = errors.New("not running as root, so can't drop privileges to specified user")
