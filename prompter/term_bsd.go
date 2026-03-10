//go:build darwin || freebsd || netbsd || openbsd || dragonfly

package prompter

import "syscall"

const (
	ioctlReadTermios  = syscall.TIOCGETA
	ioctlWriteTermios = syscall.TIOCSETA
)
