package prompt

import (
	"github.com/ZupIT/ritchie-cli/pkg/validator"
	"github.com/manifoldco/promptui"
)

type InputEmail interface {
	Email(name string) (string, error)
}

type inputEmail struct{}

func NewInputEmail() inputEmail {
	return inputEmail{}
}

type surveyEmail struct{}

func NewSurveyEmail() surveyEmail {
	return surveyEmail{}
}

// Email show a prompt and parse the string to email.
func (inputEmail) Email(name string) (string, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Pointer: promptui.PipeCursor,
		Validate:  validator.IsValidEmail,
		Templates: defaultTemplate(),
	}

	return prompt.Run()
}

func (surveyEmail) Email(name string) (string, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Pointer: promptui.PipeCursor,
		Validate:  validator.IsValidEmail,
		Templates: defaultTemplate(),
	}

	return prompt.Run()
}
