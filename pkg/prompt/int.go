package prompt

import (
	"errors"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

var ErrInvalidNumber = errors.New("invalid number")

type SurveyInt struct{}

func NewSurveyInt() SurveyInt {
	return SurveyInt{}
}

func (SurveyInt) Int(name string) (int64, error) {

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

func validateSurveyIntIn(input interface{}) error {
	if _, err := strconv.ParseInt(input.(string), 0, 64); err != nil {
		return ErrInvalidNumber
	}
	return nil
}
