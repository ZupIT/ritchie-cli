package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)

type SurveyMultiline struct{}

func NewSurveyMultiline() SurveyMultiline {
	return SurveyMultiline{}
}

func (SurveyMultiline) MultiLineText(name string, required bool) (string, error) {

	var value string

	var validationQs []*survey.Question


	if required{
		validationQs = []*survey.Question{
			{
				Name:     "name",
				Prompt:   &survey.Multiline{Message: name},
				Validate: survey.Required,
			},
		}
	}else {
		validationQs = []*survey.Question{
			{
				Name:   "name",
				Prompt: &survey.Multiline{Message: name},
			},
		}
	}
	return value, survey.Ask(validationQs, &value)
}

