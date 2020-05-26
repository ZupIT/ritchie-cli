package cmd

import (
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/security"
)

type inputTextMock struct{}

func (inputTextMock) Text(name string, required bool) (string, error) {
	return "mocked text", nil
}

type inputSecretMock struct{}

func (inputSecretMock) Text(name string, required bool) (string, error) {
	return "username=ritchie", nil
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

type inputTrueMock struct{}

func (inputTrueMock) Bool(name string, items []string) (bool, error) {
	return true, nil
}

type inputFalseMock struct{}

func (inputFalseMock) Bool(name string, items []string) (bool, error) {
	return false, nil
}

type inputListMock struct{}

func (inputListMock) List(name string, items []string) (string, error) {
	return "item-mocked", nil
}

type inputListCredMock struct{}

func (inputListCredMock) List(name string, items []string) (string, error) {
	return "me", nil
}

type repoAdder struct{}

func (a repoAdder) List() ([]formula.Repository, error) {
	return []formula.Repository{}, nil
}

func (repoAdder) Add(d formula.Repository) error {
	return nil
}

type repoCleaner struct{}

func (repoCleaner) Clean(name string) error {
	return nil
}

type formCreator struct{}

func (formCreator) Create(formulaCmd, lang string) (formula.CreateManager, error) {
	return formula.CreateManager{}, nil
}

type userManagerMock struct{}

func (userManagerMock) Create(u security.User) error {
	return nil
}

func (userManagerMock) Delete(u security.User) error {
	return nil
}

type ctxSetterMock struct{}

func (ctxSetterMock) Set(ctx string) (rcontext.ContextHolder, error) {
	return rcontext.ContextHolder{}, nil
}

type ctxFinderMock struct{}

func (ctxFinderMock) Find() (rcontext.ContextHolder, error) {
	return rcontext.ContextHolder{}, nil
}

type ctxFindRemoverMock struct{}

func (ctxFindRemoverMock) Find() (rcontext.ContextHolder, error) {
	f := ctxFinderMock{}
	return f.Find()
}

func (ctxFindRemoverMock) Remove(ctx string) (rcontext.ContextHolder, error) {
	return rcontext.ContextHolder{}, nil
}

type ctxFindSetterMock struct{}

func (ctxFindSetterMock) Find() (rcontext.ContextHolder, error) {
	f := ctxFinderMock{}
	return f.Find()
}

func (ctxFindSetterMock) Set(ctx string) (rcontext.ContextHolder, error) {
	s := ctxSetterMock{}
	return s.Set(ctx)
}

type repoDeleterMock struct{}

func (m repoDeleterMock) List() ([]formula.Repository, error) {
	return []formula.Repository{}, nil
}

func (repoDeleterMock) Delete(name string) error {
	return nil
}

type repoListerMock struct{}

func (repoListerMock) List() ([]formula.Repository, error) {
	return []formula.Repository{}, nil
}

type repoLoaderMock struct{}

func (repoLoaderMock) Load() error {
	return nil
}

type repoUpdaterMock struct{}

func (repoUpdaterMock) Update() error {
	return nil
}

type loginManagerMock struct{}

func (loginManagerMock) Login() error {
	return nil
}

type logoutManagerMock struct{}

func (logoutManagerMock) Logout() error {
	return nil
}

type credSetterMock struct{}

func (credSetterMock) Set(d credential.Detail) error {
	return nil
}

type credSettingsMock struct{}

func (credSettingsMock) Fields() (credential.Fields, error) {
	return credential.Fields{
		"github": []credential.Field{
			{
				Name: "username",
				Type: "text",
			},
			{
				Name: "token",
				Type: "password",
			},
		},
	}, nil
}
