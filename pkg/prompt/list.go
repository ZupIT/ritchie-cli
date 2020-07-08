package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)

type SurveyList struct{}

func NewSurveyList() SurveyList {
	return SurveyList{}
}

// List show a prompt with options and parse to string.
func (SurveyList) List(name string, items []string) (string, error) {
	choice := ""
	prompt := &survey.Select{
		Message: name,
		Options: items,
	}
	if err := survey.AskOne(prompt, &choice); err != nil {
		return "", err
	}

	return choice, nil
}
