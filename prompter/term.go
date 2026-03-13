package prompter

import (
	"os"
	"syscall"
	"unsafe"
)

func readPassword() ([]byte, error) {
	fd := int(syscall.Stdin)

	// Get current terminal state
	var oldState syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	// Disable echo
	newState := oldState
	newState.Lflag &^= syscall.ECHO

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	// Restore on exit
	defer syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)

	var buf [256]byte
	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		return nil, err
	}

	// Strip trailing newline
	if n > 0 && buf[n-1] == '\n' {
		n--
	}
	if n > 0 && buf[n-1] == '\r' {
		n--
	}

	return buf[:n], nil
}