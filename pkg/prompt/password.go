package prompt

import "github.com/manifoldco/promptui"

const (
	// PasswordType type
	PasswordType = "password"
)

type InputPassword interface {
	Password(label string) (string, error)
}

type inputPassword struct{}

func NewInputPassword() inputPassword {
	return inputPassword{}
}

// Password show a masked prompt and parse to string.
func (inputPassword) Password(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:     label,
		Pointer: promptui.PipeCursor,
		Mask:      '*',
		Validate:  validateEmptyInput,
		Templates: defaultTemplate(),
	}

	return prompt.Run()
}
