package prompt

import "github.com/manifoldco/promptui"

type InputList interface {
	List(name string, items []string) (string, error)
}

type inputList struct{}

func NewInputList() inputList {
	return inputList{}
}

// List show a prompt with options and parse to string.
func (inputList) List(name string, items []string) (string, error) {
	prompt := promptui.Select{
		Items:     items,
		Pointer: promptui.PipeCursor,
		Templates: defaultSelectTemplate(name),
	}
	_, result, err := prompt.Run()
	return result, err
}
