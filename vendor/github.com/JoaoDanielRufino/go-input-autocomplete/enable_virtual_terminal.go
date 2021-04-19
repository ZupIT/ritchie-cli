// +build windows

package input_autocomplete

import (
	"os"

	"golang.org/x/sys/windows"
)

func EnableVirtalTerminalWindows() error {
	var originalMode uint32
	stdout := windows.Handle(os.Stdout.Fd())

	if err := windows.GetConsoleMode(stdout, &originalMode); err != nil {
		return err
	}

	return windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
