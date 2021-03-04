package input_autocomplete

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/eiannone/keyboard"
)

func keyboardListener(input *Input) error {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		switch key {
		case keyboard.KeyEnter:
			fmt.Println("")
			return nil
		case keyboard.KeyArrowLeft:
			input.MoveCursorLeft()
		case keyboard.KeyArrowRight:
			input.MoveCursorRight()
		case keyboard.KeyBackspace:
			input.RemoveChar()
		case keyboard.KeyBackspace2:
			input.RemoveChar()
		case keyboard.KeyTab:
			input.Autocomplete()
		case keyboard.KeyCtrlC:
			return errors.New("Aborted")

		default:
			input.AddChar(char)
		}
	}
}

func Read(text string) (string, error) {
	if err := keyboard.Open(); err != nil {
		return "", err
	}

	defer keyboard.Close()

	os := runtime.GOOS
	if os == "windows" {
		if err := EnableVirtalTerminalWindows(); err != nil {
			return "", err
		}
	}

	input := NewInput(text)

	input.Print()

	if err := keyboardListener(input); err != nil {
		return "", err
	}

	input.RemoveLastSlashIfNeeded()

	return input.GetCurrentText(), nil
}
