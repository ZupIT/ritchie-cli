package prompt

import (
	"errors"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

func defaultTemplate() *promptui.PromptTemplates {
	return &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | bold }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}
}

func defaultSelectTemplate(label string) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label: label,
	}
}

func validateEmptyInput(input string) error {
	if len(strings.TrimSpace(input)) < 1 {
		return errors.New("this input must not be empty")
	}
	return nil
}

func validateIntIn(input string) error {
	_, err := strconv.ParseInt(input, 0, 64)
	if err != nil {
		return errors.New("invalid number")
	}
	return nil
}

func validateSurveyIntIn(input interface{}) error {
	_, err := strconv.ParseInt(input.(string), 0, 64)
	if err != nil {
		return errors.New("invalid number")
	}
	return nil
}