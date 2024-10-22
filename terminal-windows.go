//go:build windows
// +build windows

package main

import (
	"golang.org/x/sys/windows"
)

func getWindowSize() (*windowSize, error) {
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Stdout, &info); err != nil {
		return nil, err
	}

	ws := &windowSize{
		Row: uint16(info.Window.Bottom - info.Window.Top + 1),
		Col: uint16(info.Window.Right - info.Window.Left + 1),
	}
	return ws, nil
}
