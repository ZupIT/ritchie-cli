package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
)

type InputText interface {
	Text(name string, required bool, helper ...string) (string, error)
	TextWithValidate(name string, validate func(interface{}) error, helper ...string) (string, error)
}

type inputText struct{}

type surveyText struct{}

func NewInputText() inputText {
	return inputText{}
}

func NewSurveyText() surveyText {
	return surveyText{}
}

// Text show a prompt and parse to string.
func (inputText) Text(name string, required bool) (string, error) {
	var prompt promptui.Prompt

	if required {
		prompt = promptui.Prompt{
			Label:     name,
			Pointer:   promptui.PipeCursor,
			Validate:  validateEmptyInput,
			Templates: defaultTemplate(),
		}
	} else {
		prompt = promptui.Prompt{
			Label:     name,
			Pointer:   promptui.PipeCursor,
			Templates: defaultTemplate(),
		}
	}

	return prompt.Run()
}

func (inputText) TextWithValidate(name string, validate func(string) error, helper ...string) (string, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Pointer:   promptui.PipeCursor,
		Validate:  validate,
		Templates: defaultTemplate(),
	}

	return prompt.Run()
}

func (surveyText) Text(name string, required bool, helper ...string) (string, error) {

	var value string

	validationQs := []*survey.Question{
		{
			Name: "name",
		},
	}

	if required {
		validationQs[0].Validate = survey.Required
	}

	if len(helper) > 0 {
		validationQs[0].Prompt = &survey.Input{Message: name, Help: helper[0]}
	} else {
		validationQs[0].Prompt = &survey.Input{Message: name}
	}

	return value, survey.Ask(validationQs, &value)
}

func (surveyText) TextWithValidate(name string, validate func(interface{}) error, helper ...string) (string, error) {
	var value string
	validationQs := []*survey.Question{
		{
			Name:     "name",
			Validate: validate,
		},
	}

	if len(helper) > 0 {
		validationQs[0].Prompt = &survey.Input{Message: name, Help: helper[0]}
	} else {
		validationQs[0].Prompt = &survey.Input{Message: name}
	}

	return value, survey.Ask(validationQs, &value)
}
