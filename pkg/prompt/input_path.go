package prompt

import input_autocomplete "github.com/JoaoDanielRufino/go-input-autocomplete"

type InputAutocomplete struct{}

func NewInputAutocomplete() InputAutocomplete {
	return InputAutocomplete{}
}

func (InputAutocomplete) Read(text string) (string, error) {
	return input_autocomplete.Read(Green("? ") + Bold(text))
}
