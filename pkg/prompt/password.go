package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)

const PasswordType = "password"

type SurveyPassword struct{}

func NewSurveyPassword() SurveyPassword {
	return SurveyPassword{}
}

func (SurveyPassword) Password(label string) (string, error) {
	password := ""
	prompt := &survey.Password{
		Message: label,
	}

	return password, survey.AskOne(prompt, &password)
}
