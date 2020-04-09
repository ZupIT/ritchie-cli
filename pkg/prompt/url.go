package prompt

import (
	"github.com/ZupIT/ritchie-cli/pkg/validator"
	"github.com/manifoldco/promptui"
)

type InputURL interface {
	URL(name, defaultValue string) (string, error)
}

type inputURL struct{}

func NewInputURL() inputURL {
	return inputURL{}
}

// URL show a prompt and parse to string.
func (inputURL) URL(name, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Default:   defaultValue,
		Validate:  validator.IsValidURL,
		Templates: defaultTemplate(),
	}

	return prompt.Run()
}
