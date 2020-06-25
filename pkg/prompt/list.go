package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
)

type InputList interface {
	List(name string, items []string) (string, error)
}

type inputList struct{}

type surveyList struct{}

func NewInputList() inputList {
	return inputList{}
}

func NewSurveyList() surveyList {
	return surveyList{}
}

// List show a prompt with options and parse to string.
func (inputList) List(name string, items []string) (string, error) {
	prompt := promptui.Select{
		Items:     items,
		Pointer: promptui.PipeCursor,
		Templates: defaultSelectTemplate(name),
	}
	_, result, err := prompt.Run()
	return result, err
}

// List show a prompt with options and parse to string.
func (surveyList) List(name string, items []string) (string, error) {
	choice := ""
	prompt := &survey.Select{
		Message: name,
		Options: items,
	}
	if err := survey.AskOne(prompt, &choice); err != nil{
		return "", err
	}

	return choice, nil
}
