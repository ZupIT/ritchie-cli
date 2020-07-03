package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)


type InputMultiline interface {
	MultiLineText(name string, required bool) (string, error)
}


type surveyMultiline struct{}

func NewSurveyMultiline() surveyMultiline {
	return surveyMultiline{}
}

func (surveyMultiline) MultiLineText(name string, required bool) (string, error) {

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

