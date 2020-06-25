package prompt

import (
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
)

type InputInt interface {
	Int(name string) (int64, error)
}

type inputInt struct{}

type surveyInt struct{}

func NewInputInt() inputInt {
	return inputInt{}
}

func NewSurveyInt() surveyInt {
	return surveyInt{}
}

//Int show a prompt and parse to int.
func (inputInt) Int(name string) (int64, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Pointer: promptui.PipeCursor,
		Validate:  validateIntIn,
		Templates: defaultTemplate(),
	}

	promptResult, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	parseInt, err := strconv.ParseInt(promptResult, 0, 64)
	if err != nil {
		return 0, err
	}
	return parseInt, nil
}

func (surveyInt) Int(name string) (int64, error) {

	var value string

	validationQs := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: name},
			Validate: validateSurveyIntIn,
		},
	}
	if err := survey.Ask(validationQs, &value); err != nil {
		return 0, err
	}

	parseInt, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return 0, err
	}
	return parseInt, nil
}

