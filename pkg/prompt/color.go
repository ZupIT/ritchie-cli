package prompt

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
)

func NewError(text string) error {
	return errors.New(color.FgRed.Render(text))
}

func Error(text string) {
	fmt.Println(color.FgRed.Render(text))
}

func Warning(text string) {
	color.Warn.Println(text)
}

func Success(text string) {
	color.Success.Println(text)
}

func Info(text string) {
	color.Bold.Println(text)
}

