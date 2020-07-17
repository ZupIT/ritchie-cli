package cmd

import (
	"errors"

	"github.com/docker/docker/api/server"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type inputTextMock struct{}

func (inputTextMock) Text(name string, required bool, helper ...string) (string, error) {
	return "mocked text", nil
}

func (inputTextMock) TextWithValidate(name string, validate func(interface{}) error, helper ...string) (string, error) {
	return "mocked text", nil
}

type inputTextValidatorMock struct{}

func (inputTextValidatorMock) Text(name string, validate func(interface{}) error, helper ...string) (string, error) {
	return "mocked text", nil
}

type inputSecretMock struct{}

func (inputSecretMock) Text(name string, required bool, helper ...string) (string, error) {
	return "username=ritchie", nil
}

func (inputSecretMock) TextWithValidate(name string, validate func(interface{}) error, helper ...string) (string, error) {
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

type inputIntErrorMock struct{}

func (inputIntErrorMock) Int(name string) (int64, error) {
	return 0, errors.New("some error")
}

type inputPasswordMock struct{}

func (inputPasswordMock) Password(label string) (string, error) {
	return "s3cr3t", nil
}

type autocompleteGenMock struct{}

func (autocompleteGenMock) Generate(s autocomplete.ShellName, cmd *cobra.Command) (string, error) {
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

type inputListCustomMock struct {
	name string
}

func (m inputListCustomMock) List(name string, items []string) (string, error) {
	return m.name, nil
}

type inputListCredMock struct{}

func (inputListCredMock) List(name string, items []string) (string, error) {
	return "me", nil
}

type inputListErrorMock struct{}

func (inputListErrorMock) List(name string, items []string) (string, error) {
	return "item-mocked", errors.New("some error")
}

type repoAdder struct{}

func (a repoAdder) List() (formula.Repos, error) {
	return formula.Repos{}, nil
}

func (repoAdder) Add(d formula.Repo) error {
	return nil
}

type formCreator struct{}

func (formCreator) Create(cf formula.Create) error {
	return nil
}

func (formCreator) Build(workspacePath, formulaPath string) error {
	return nil
}

type workspaceForm struct{}

func (workspaceForm) Add(workspace formula.Workspace) error {
	return nil
}

func (workspaceForm) List() (formula.Workspaces, error) {
	return formula.Workspaces{}, nil
}

func (workspaceForm) Validate(workspace formula.Workspace) error {
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

type repoListerMock struct{}

func (repoListerMock) List() (formula.Repos, error) {
	return formula.Repos{}, nil
}

type repoListerNonEmptyMock struct{}

func (repoListerNonEmptyMock) List() (formula.Repos, error) {
	return formula.Repos{
		{
			Name:     "repoName",
			Priority: 0,
		},
	}, nil
}

type repoListerErrorMock struct{}

func (repoListerErrorMock) List() (formula.Repos, error) {
	return formula.Repos{}, errors.New("some error")
}

type repoPrioritySetterMock struct{}

func (repoPrioritySetterMock) SetPriority(repo string, priority int) error {
	return nil
}

type repoPrioritySetterCustomMock struct {
	setPriority func(repo string, priority int) error
}

func (m repoPrioritySetterCustomMock) SetPriority(repo string, priority int) error {
	return m.setPriority(repo, priority)
}

type repoLoaderMock struct{}

func (repoLoaderMock) Load() error {
	return nil
}

type repoUpdaterMock struct{}

func (repoUpdaterMock) Update() error {
	return nil
}

type credSetterMock struct{}

func (credSetterMock) Set(d credential.Detail) error {
	return nil
}

type credSettingsMock struct{}

type singleCredSettingsMock struct{}

func (s singleCredSettingsMock) WriteDefaultCredentialsFields(path string) error {
	return nil
}

func (s singleCredSettingsMock) ReadCredentialsFields(path string) (credential.Fields, error) {
	return nil, nil
}

func (s singleCredSettingsMock) WriteCredentialsFields(fields credential.Fields, path string) error {
	return nil
}

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

type runnerMock struct {
	error error
}

func (r runnerMock) Run(def formula.Definition, inputType api.TermInputType) error {
	return r.error
}

type treeMock struct {
	tree  formula.Tree
	error error
}

func (t treeMock) Tree() (map[string]formula.Tree, error) {
	return map[string]formula.Tree{"test": t.tree}, t.error
}

func (t treeMock) MergedTree(bool) formula.Tree {
	return t.tree
}

type findSetterServerMock struct{}

func (findSetterServerMock) Set(*server.Config) error {
	return nil
}

func (findSetterServerMock) Find() (server.Config, error) {
	return server.Config{}, nil
}

type findSetterServerCustomMock struct {
	set  func(*server.Config) error
	find func() (server.Config, error)
}

func (m findSetterServerCustomMock) Set(ctx *server.Config) error {
	return m.set(ctx)
}

func (m findSetterServerCustomMock) Find() (server.Config, error) {
	return m.find()
}

type inputBoolCustomMock struct {
	bool func(name string, items []string) (bool, error)
}

func (m inputBoolCustomMock) Bool(name string, items []string) (bool, error) {
	return m.bool(name, items)
}

type inputTextCustomMock struct {
	text             func(name string, required bool) (string, error)
	textWithValidate func(name string, validate func(interface{}) error) (string, error)
}

func (m inputTextCustomMock) Text(name string, required bool, helper ...string) (string, error) {
	return m.text(name, required)
}

func (m inputTextCustomMock) TextWithValidate(name string, validate func(interface{}) error, helper ...string) (string, error) {
	return m.textWithValidate(name, validate)
}

type inputURLCustomMock struct {
	url func(name, defaultValue string) (string, error)
}

func (m inputURLCustomMock) URL(name, defaultValue string) (string, error) {
	return m.url(name, defaultValue)
}

type inputPasswordCustomMock struct {
	password func(label string) (string, error)
}

func (m inputPasswordCustomMock) Password(label string) (string, error) {
	return m.password(label)
}

type InputMultilineMock struct{}

func (InputMultilineMock) MultiLineText(name string, required bool) (string, error) {
	return "username=ritchie", nil
}
