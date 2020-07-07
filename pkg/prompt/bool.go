package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)

var (
	boolOpts = map[string]bool{"yes": true, "no": false, "true": true, "false": false}
)

type SurveyBool struct{}

func NewSurveyBool() SurveyBool {
	return SurveyBool{}
}

func (SurveyBool) Bool(name string, items []string) (bool, error) {
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
