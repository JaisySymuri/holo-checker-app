//go:build windows

package utility

import "syscall"

func consoleAttached() bool {
	h, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	return err == nil && h != 0 && h != syscall.InvalidHandle
}
