package prompt

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
)

// NewError returns new error with red message
func NewError(text string) error {
	return errors.New(Red(text))
}

// Red returns a red string
func Red(text string) string {
	return color.FgRed.Render(text)
}

// Error is a Println with red message
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

func Yellow(text string) string {
	return color.Warn.Render(text)
}
func Warning(text string) {
	color.Warn.Println(text)
}
