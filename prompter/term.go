//go:build !tinygo

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

	// Read in chunks to support arbitrary-length passwords
	var buf []byte
	tmp := make([]byte, 4096)
	for {
		n, err := os.Stdin.Read(tmp[:])
		if err != nil {
			return nil, err
		}
		// Append to buffer
		buf = append(buf, tmp[:n]...)
		// Stop at newline (terminal input ends with newline)
		if n > 0 && tmp[n-1] == '\n' {
			break
		}
	}

	// Strip trailing newline
	if len(buf) > 0 && buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}
	if len(buf) > 0 && buf[len(buf)-1] == '\r' {
		buf = buf[:len(buf)-1]
	}

	return buf, nil
}
