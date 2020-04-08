package prompt

import (
	"errors"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

const (
	// PasswordType type
	PasswordType = "password"
)

var boolmap = map[string]bool{"yes": true, "no": false, "true": true, "false": false}

// String show a prompt and parse to string.
func String(name string, required bool) (string, error) {
	var prompt promptui.Prompt

	if required {
		prompt = promptui.Prompt{
			Label:     name,
			Validate:  validateEmptyInput,
			Templates: defaultTemplate(),
		}
	} else {
		prompt = promptui.Prompt{
			Label:     name,
			Templates: defaultTemplate(),
		}
	}

	return prompt.Run()
}

// String show a prompt and parse to string.
func Email(name string) (string, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Validate:  validator.IsValidEmail,
		Templates: defaultTemplate(),
	}

	return prompt.Run()
}

// String show a prompt and parse to string.
func URL(name, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Default:   defaultValue,
		Validate:  validator.IsValidURL,
		Templates: defaultTemplate(),
	}

	return prompt.Run()
}

// Integer show a prompt and parse to int.
func Integer(name string) (int64, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Validate:  validateIntegerNumberInput,
		Templates: defaultTemplate(),
	}

	promptResult, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	parseInt, _ := strconv.ParseInt(promptResult, 0, 64)
	return parseInt, nil
}

// List show a prompt with options and parse to string.
func List(name string, items []string) (string, error) {
	prompt := promptui.Select{
		Items:     items,
		Templates: defaultSelectTemplate(name),
	}
	_, result, err := prompt.Run()
	return result, err
}

// ListBool show a prompt with options and parse to bool.
func ListBool(name string, items []string) (bool, error) {
	prompt := promptui.Select{
		Items:     items,
		Templates: defaultSelectTemplate(name),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	b := boolmap[result]
	return b, err
}

// Password show a masked prompt and parse to string.
func Password(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:     label,
		Mask:      '*',
		Validate:  validateEmptyInput,
		Templates: defaultTemplate(),
	}

	return prompt.Run()
}

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

func validateIntegerNumberInput(input string) error {
	_, err := strconv.ParseInt(input, 0, 64)
	if err != nil {
		return errors.New("invalid number")
	}
	return nil
}
