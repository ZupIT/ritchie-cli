package cmd

import (
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/security"
)

type inputTextMock struct{}

func (inputTextMock) Text(name string, required bool) (string, error) {
	return "mocked text", nil
}

type inputURLMock struct{}

func (inputURLMock) URL(name, defaultValue string) (string, error) {
	return "http://localhost/mocked", nil
}

type inputIntMock struct{}

func (inputIntMock) Int(name string) (int64, error) {
	return 0, nil
}

type inputEmailMock struct{}

func (inputEmailMock) Email(name string) (string, error) {
	return "dennis@ritchie.io", nil
}

type inputPasswordMock struct{}

func (inputPasswordMock) Password(label string) (string, error) {
	return "s3cr3t", nil
}

type autocompleteGenMock struct{}

func (autocompleteGenMock) Generate(s autocomplete.ShellName) (string, error) {
	return "autocomplete", nil
}

type repoAdder struct{}

func (repoAdder) Add(d formula.Repository) error {
	return nil
}

type formCreator struct{}

func (formCreator) Create(formulaCmd string) error {
	return nil
}

type userManagerMock struct{}

func (userManagerMock) Create(u security.User) error {
	return nil
}

func (userManagerMock) Delete(u security.User) error {
	return nil
}
