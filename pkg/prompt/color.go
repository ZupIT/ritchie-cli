package prompt

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
)

func Red(text string) string {
	return color.FgRed.Render(text)
}

func PrintRed(text string) {
	fmt.Println(Red(text))
}

func Error(text string) error {
	return errors.New(Red(text))
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
