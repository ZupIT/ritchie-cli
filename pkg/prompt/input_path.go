package prompt

import input_autocomplete "github.com/JoaoDanielRufino/go-input-autocomplete"

type InputPath struct{}

func NewInputPath() InputPath {
	return InputPath{}
}

func (InputPath) Read(text string) (string, error) {
	return input_autocomplete.Read(Green("? ") + Bold(text))
}
