package cmd

import (
	"errors"
	"io"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
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

type inputTextErrorMock struct{}

func (inputTextErrorMock) Text(name string, required bool, helper ...string) (string, error) {
	return "", errors.New("error on input text")
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

func (inputIntMock) Int(name string, helper ...string) (int64, error) {
	return 0, nil
}

type inputIntErrorMock struct{}

func (inputIntErrorMock) Int(name string, helper ...string) (int64, error) {
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

type repoListerAdderCustomMock struct {
	list func() (formula.Repos, error)
	add  func(d formula.Repo) error
}

func (a repoListerAdderCustomMock) List() (formula.Repos, error) {
	return a.list()
}

func (a repoListerAdderCustomMock) Add(d formula.Repo) error {
	return a.add(d)
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

func (repoPrioritySetterMock) SetPriority(name formula.RepoName, priority int) error {
	return nil
}

type repoPrioritySetterCustomMock struct {
	setPriority func(name formula.RepoName, priority int) error
}

func (m repoPrioritySetterCustomMock) SetPriority(name formula.RepoName, priority int) error {
	return m.setPriority(name, priority)
}

type credSetterMock struct{}

func (credSetterMock) Set(d credential.Detail) error {
	return nil
}

type credSettingsMock struct {
	error
}

func (s credSettingsMock) ReadCredentialsFields(path string) (credential.Fields, error) {
	return credential.Fields{}, nil
}

func (s credSettingsMock) ReadCredentialsValue(path string) ([]credential.ListCredData, error) {
	return []credential.ListCredData{}, nil
}

func (s credSettingsMock) WriteDefaultCredentialsFields(path string) error {
	return nil
}

func (s credSettingsMock) WriteCredentialsFields(fields credential.Fields, path string) error {
	return nil
}

func (s credSettingsMock) ProviderPath() string {
	return ""
}

func (s credSettingsMock) CredentialsPath () string {
	return ""
}

type runnerMock struct {
	error error
}

func (r runnerMock) Run(def formula.Definition, inputType api.TermInputType, local bool) error {
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

type GitRepositoryMock struct {
	zipball   func(info github.RepoInfo, version string) (io.ReadCloser, error)
	tags      func(info github.RepoInfo) (github.Tags, error)
	latestTag func(info github.RepoInfo) (github.Tag, error)
}

func (m GitRepositoryMock) Zipball(info github.RepoInfo, version string) (io.ReadCloser, error) {
	return m.zipball(info, version)
}

func (m GitRepositoryMock) Tags(info github.RepoInfo) (github.Tags, error) {
	return m.tags(info)
}

func (m GitRepositoryMock) LatestTag(info github.RepoInfo) (github.Tag, error) {
	return m.latestTag(info)
}

type TutorialSetterMock struct{}

func (TutorialSetterMock) Set(tutorial string) (rtutorial.TutorialHolder, error) {
	return rtutorial.TutorialHolder{}, nil
}

type TutorialFinderMock struct{}

func (TutorialFinderMock) Find() (rtutorial.TutorialHolder, error) {
	return rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial}, nil
}

type TutorialFindSetterMock struct{}

func (TutorialFindSetterMock) Find() (rtutorial.TutorialHolder, error) {
	f := TutorialFinderMock{}
	return f.Find()
}

func (TutorialFindSetterMock) Set(tutorial string) (rtutorial.TutorialHolder, error) {
	s := TutorialSetterMock{}
	return s.Set(tutorial)
}

var (
	defaultRepoAdderMock = repoListerAdderCustomMock{
		add: func(d formula.Repo) error {
			return nil
		},
		list: func() (formula.Repos, error) {
			return formula.Repos{}, nil
		},
	}

	defaultGitRepositoryMock = GitRepositoryMock{
		latestTag: func(info github.RepoInfo) (github.Tag, error) {
			return github.Tag{}, nil
		},
		tags: func(info github.RepoInfo) (github.Tags, error) {
			return github.Tags{}, nil
		},
		zipball: func(info github.RepoInfo, version string) (io.ReadCloser, error) {
			return nil, nil
		},
	}
)
