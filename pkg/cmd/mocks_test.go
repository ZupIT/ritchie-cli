package cmd

import (
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/security/otp"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/spf13/cobra"
)

type inputTextMock struct{}

func (inputTextMock) Text(name string, required bool, helper ...string) (string, error) {
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

type inputURLMock struct{}

func (inputURLMock) URL(name, defaultValue string) (string, error) {
	return "http://localhost/mocked", nil
}

type inputIntMock struct{}

func (inputIntMock) Int(name string) (int64, error) {
	return 0, nil
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

func (loginManagerMock) Login(security.User) error {
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

type singleCredSettingsMock struct {}

func (s singleCredSettingsMock) WriteDefaultCredentials(path string) error {
	return nil
}

func (s singleCredSettingsMock) ReadCredentials(path string) (credential.Fields, error) {
	return nil, nil
}

func (s singleCredSettingsMock) WriteCredentials(fields credential.Fields, path string) error {
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

type passphraseManagerMock struct{}

func (passphraseManagerMock) Save(security.Passphrase) error {
	return nil
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

type loginManagerCustomMock struct {
	login func(security.User) error
}

func (m loginManagerCustomMock) Login(user security.User) error {
	return m.login(user)
}

type repoLoaderCustomMock struct {
	load func() error
}

func (m repoLoaderCustomMock) Load() error {
	return m.load()
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

type otpResolverMock struct{}

func (m otpResolverMock) RequestOtp(url, organization string) (otp.Response, error) {
	return otp.Response{Otp: true}, nil
}

type otpResolverCustomMock struct {
	requestOtp func(url, organization string) (otp.Response, error)
}

func (m otpResolverCustomMock) RequestOtp(url, organization string) (otp.Response, error) {
	return m.requestOtp(url, organization)
}
