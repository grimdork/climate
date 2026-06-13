//go:build tinygo

package prompter

import (
	"os"
)

func readPassword() ([]byte, error) {
	// TinyGo does not support terminal echo suppression.
	// Read a line from stdin without hiding input.
	var buf []byte
	tmp := make([]byte, 4096)
	for {
		n, err := os.Stdin.Read(tmp[:])
		if err != nil {
			return nil, err
		}
		buf = append(buf, tmp[:n]...)
		if n > 0 && tmp[n-1] == '\n' {
			break
		}
	}
	if len(buf) > 0 && buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}
	if len(buf) > 0 && buf[len(buf)-1] == '\r' {
		buf = buf[:len(buf)-1]
	}
	return buf, nil
}
