package prompt

import "github.com/manifoldco/promptui"

type InputText interface {
	Text(name string, required bool) (string, error)
}

type inputText struct{}

func NewInputText() inputText {
	return inputText{}
}

// Text show a prompt and parse to string.
func (inputText) Text(name string, required bool) (string, error) {
	var prompt promptui.Prompt

	if required {
		prompt = promptui.Prompt{
			Label:     name,
			Validate:  validateEmptyInput,
			Templates: defaultTemplate(),
		}
	} else {
		prompt = promptui.Prompt{
			Label:     name,
			Templates: defaultTemplate(),
		}
	}

	return prompt.Run()
}
