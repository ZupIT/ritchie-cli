package prompt

import "fmt"

const (
	Info    = "\033[1;36m%s\033[0m"
	Yellow  = "\033[1;33m%s\033[0m"
	Error   = "\033[1;31m%s\033[0m"
	Success = "\033[1;32m%s\033[0m"
)

func Red(text string) string {
	return fmt.Sprintf(Error,text)
}

func Yaellow(text string) string {
	return fmt.Sprintf(Error,text)
}

func Green(text string) string {
	return fmt.Sprintf(Error,text)
}

func Red(text string) string {
	return fmt.Sprintf(Error,text)
}
