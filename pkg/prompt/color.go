package prompt

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
)

func Red(text string) string {
	return color.FgRed.Render(text)
}
func NewError(text string) error {
	return errors.New(color.FgRed.Render(text))
}
func Error(text string) {
	fmt.Println(Red(text))
}

func Green(text string) string {
	return color.Success.Render(text)
}
func Success(text string) {
	fmt.Println(Green(text))
}

func Bold(text string) string {
	return color.Bold.Render(text)
}
func Info(text string) {
	fmt.Println(Bold(text))
}

func Warning(text string) {
	color.Warn.Println(text)
}
