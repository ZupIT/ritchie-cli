package prompt

import (
	"fmt"

	"github.com/gookit/color"
)

func Red(text string) string {
	return color.FgRed.Render(text)
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
