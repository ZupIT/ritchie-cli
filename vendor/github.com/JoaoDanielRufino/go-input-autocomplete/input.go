package input_autocomplete

import (
	"fmt"
	"runtime"
)

type Input struct {
	cursor      *Cursor
	fixedText   string
	currentText string
}

func NewInput(fixedText string) *Input {
	return &Input{
		cursor:      NewCursor(),
		fixedText:   fixedText,
		currentText: "",
	}
}

func (i *Input) canDeleteChar() bool {
	return i.cursor.GetPosition() >= 1
}

func (i *Input) AddChar(char rune) {
	pos := i.cursor.GetPosition()
	c := string(char)

	if pos == len(i.currentText) {
		i.currentText += c
		fmt.Print(c)
		i.cursor.IncrementPosition()
	} else {
		aux := len(i.currentText) - pos
		i.currentText = i.currentText[:pos] + c + i.currentText[pos:]
		i.cursor.SetPosition(len(i.currentText))
		i.Print()
		i.cursor.MoveLeftNPos(aux)
	}
}

func (i *Input) RemoveChar() {
	if i.canDeleteChar() {
		pos := i.cursor.GetPosition()
		aux := len(i.currentText) - pos
		i.currentText = i.currentText[:pos-1] + i.currentText[pos:]
		i.cursor.SetPosition(len(i.currentText))
		i.Print()
		i.cursor.MoveLeftNPos(aux)
	}
}

func (i *Input) MoveCursorLeft() {
	i.cursor.MoveLeft()
}

func (i *Input) MoveCursorRight() {
	if i.cursor.GetPosition() < len(i.currentText) {
		i.cursor.MoveRight()
	}
}

func (i *Input) Autocomplete() {
	if i.currentText == "" {
		return
	}
	autocompletedText := Autocomplete(i.currentText)
	i.currentText = autocompletedText
	i.cursor.SetPosition(len(i.currentText))
	i.Print()
}

func (i *Input) RemoveLastSlashIfNeeded() {
	os := runtime.GOOS
	size := len(i.currentText)
	var slash byte

	switch os {
	case "linux", "darwin":
		slash = '/'
	case "windows":
		slash = '\\'
	}

	if i.currentText[size-1] == slash {
		i.currentText = i.currentText[:size-1]
	}
}

func (i *Input) Print() {
	fmt.Print("\033[G\033[K")
	fmt.Print(i.fixedText + i.currentText)
}

func (i *Input) GetCurrentText() string {
	return i.currentText
}
