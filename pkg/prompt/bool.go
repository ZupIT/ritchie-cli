package prompt

import "github.com/manifoldco/promptui"

var (
	boolOpts = map[string]bool{"yes": true, "no": false, "true": true, "false": false}
)

type InputBool interface {
	Bool(name string, items []string) (bool, error)
}

type inputBool struct{}

func NewInputBool() inputBool {
	return inputBool{}
}

// Bool show a prompt with options and parse to bool.
func (inputBool) Bool(name string, items []string) (bool, error) {
	prompt := promptui.Select{
		Items:     items,
		Pointer: promptui.PipeCursor,
		Templates: defaultSelectTemplate(name),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	b := boolOpts[result]
	return b, err
}
