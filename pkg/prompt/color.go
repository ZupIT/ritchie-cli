package prompt

import "fmt"

const (
	Teal   = "\033[1;36m%s\033[0m"
	Yellow = "\033[1;33m%s\033[0m"
	Red    = "\033[1;31m%s\033[0m"
	Green  = "\033[1;32m%s\033[0m"
)

func Error(text string) string {
	return fmt.Sprintf(Red,text)
}

func Warning(text string) string {
	return fmt.Sprintf(Yellow,text)
}

func Success(text string) string {
	return fmt.Sprintf(Green,text)
}

func Info(text string) string {
	return fmt.Sprintf(Teal,text)
}

