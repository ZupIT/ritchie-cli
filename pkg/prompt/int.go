package prompt

import (
	"github.com/manifoldco/promptui"
	"strconv"
)

type InputInt interface {
	Int(name string) (int64, error)
}

type inputInt struct{}

func NewInputInt() inputInt {
	return inputInt{}
}

// Int show a prompt and parse to int.
func (inputInt) Int(name string) (int64, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Validate:  validateIntegerNumberInput,
		Templates: defaultTemplate(),
	}

	promptResult, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	parseInt, err := strconv.ParseInt(promptResult, 0, 64)
	if err != nil {
		return 0, err
	}
	return parseInt, nil
}
