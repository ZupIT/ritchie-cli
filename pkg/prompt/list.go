package prompt

import "github.com/manifoldco/promptui"

type InputList interface {
	ListI(name string, items []string) (string, error)
}

type inputList struct{}

func NewInputList() inputList {
	return inputList{}
}

// ListI show a prompt with options and parse to string.
func (inputList) ListI(name string, items []string) (string, error) {
	prompt := promptui.Select{
		Items:     items,
		Templates: defaultSelectTemplate(name),
	}
	_, result, err := prompt.Run()
	return result, err
}
