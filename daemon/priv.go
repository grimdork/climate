// Drop privileges.
package daemon

import (
	"os/user"
	"strconv"
	"syscall"
)

// DegradeToUser drops down to a specific user and its primary group if run by root.
func DegradeToUser(uname string) error {
	uid := syscall.Geteuid()
	if uid == 0 {
		u, err := user.Lookup(uname)
		if err != nil {
			return err
		}

		uid, err := strconv.Atoi(u.Uid)
		if err != nil {
			return err
		}

		gid, err := strconv.Atoi(u.Gid)
		if err != nil {
			return err
		}

		err = syscall.Setgid(gid)
		if err != nil {
			return err
		}

		err = syscall.Setreuid(-1, uid)
		if err != nil {
			return err
		}
	} else {
		return ErrorNotRoot
	}

	return nil
}
