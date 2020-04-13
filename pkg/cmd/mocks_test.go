package cmd

import (
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
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

type inputBoolMock struct{}

func (inputBoolMock) Bool(name string, items []string) (bool, error) {
	return true, nil
}

type inputListMock struct{}

func (inputListMock) List(name string, items []string) (string, error) {
	return "item-mocked", nil
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

type findRemoverMock struct{}

func (findRemoverMock) Find() (rcontext.ContextHolder, error) {
	return rcontext.ContextHolder{}, nil
}

func (findRemoverMock) Remove(ctx string) (rcontext.ContextHolder, error) {
	return rcontext.ContextHolder{}, nil
}

type repoDeleterMock struct{}

func (repoDeleterMock) Delete(name string) error {
	return nil
}

type repoListerMock struct{}

func (repoListerMock) List() ([]formula.Repository, error) {
	return []formula.Repository{}, nil
}
