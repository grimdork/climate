//go:build (darwin || freebsd || netbsd || openbsd || dragonfly) && !tinygo

package prompter

import "syscall"

const (
	ioctlReadTermios  = syscall.TIOCGETA
	ioctlWriteTermios = syscall.TIOCSETA
)
