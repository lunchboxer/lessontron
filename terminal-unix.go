//go:build !windows
// +build !windows

package main

import (
	"syscall"
	"unsafe"
)

func getWindowSize() (*windowSize, error) {
	ws := &windowSize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		return nil, errno
	}
	return ws, nil
}
