package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
)

var (
	boolOpts = map[string]bool{"yes": true, "no": false, "true": true, "false": false}
)

type InputBool interface {
	Bool(name string, items []string) (bool, error)
}

type inputBool struct{}

func NewInputBool() inputBool {
	return inputBool{}
}

type surveyBool struct{}

func NewSurveyBool() surveyBool {
	return surveyBool{}
}

func (surveyBool) Bool(name string, items []string) (bool, error) {

	choice := ""
	prompt := &survey.Select{
		Message: name,
		Options: items,
	}
	if err := survey.AskOne(prompt, &choice); err != nil {
		return false, err
	}

	return boolOpts[choice], nil
}

// Bool show a prompt with options and parse to bool.
func (inputBool) Bool(name string, items []string) (bool, error) {
	prompt := promptui.Select{
		Items:     items,
		Pointer:   promptui.PipeCursor,
		Templates: defaultSelectTemplate(name),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	b := boolOpts[result]
	return b, err
}
