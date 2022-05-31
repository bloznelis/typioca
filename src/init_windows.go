package main

import (
	"runtime"

	"github.com/muesli/termenv"
)

func OsInit() {
	// enable colors for one guy who uses windows
	if runtime.GOOS == "windows" {
		mode, err := termenv.EnableWindowsANSIConsole()
		if err != nil {
			panic(err)
		}
		defer termenv.RestoreWindowsConsole(mode)
	}
}
