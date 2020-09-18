/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"errors"
	"io"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/git"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
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

type inputTextValidatorCustomMock struct {
	text func(name string, validate func(interface{}) error, helper ...string) (string, error)
}

func (i inputTextValidatorCustomMock) Text(name string, validate func(interface{}) error, helper ...string) (string, error) {
	return i.text(name, validate)
}

type inputTextValidatorErrorMock struct{}

func (inputTextValidatorErrorMock) Text(name string, validate func(interface{}) error, helper ...string) (string, error) {
	return "mocked text", errors.New("error on input text")
}

type inputTextErrorMock struct{}

func (inputTextErrorMock) Text(name string, required bool, helper ...string) (string, error) {
	return "", errors.New("error on input text")
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

type inputURLErrorMock struct{}

func (inputURLErrorMock) URL(name, defaultValue string) (string, error) {
	return "http://localhost/mocked", errors.New("error on input url")
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

type inputPasswordErrorMock struct{}

func (inputPasswordErrorMock) Password(label string) (string, error) {
	return "", errors.New("password error")
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

type inputBoolErrorMock struct{}

func (inputBoolErrorMock) Bool(name string, items []string) (bool, error) {
	return false, errors.New("error on boolean list")
}

type inputListMock struct{}

func (inputListMock) List(name string, items []string) (string, error) {
	return "item-mocked", nil
}

type inputListCustomMock struct {
	list func(name string, items []string) (string, error)
}

func (m inputListCustomMock) List(name string, items []string) (string, error) {
	return m.list(name, items)
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

func (workspaceForm) Delete(workspace formula.Workspace) error {
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

func (s credSettingsMock) CredentialsPath() string {
	return ""
}

type credSettingsCustomMock struct {
	ReadCredentialsValueMock          func(path string) ([]credential.ListCredData, error)
	ReadCredentialsFieldsMock         func(path string) (credential.Fields, error)
	WriteDefaultCredentialsFieldsMock func(path string) error
	WriteCredentialsFieldsMock        func(fields credential.Fields, path string) error
	ProviderPathMock                  func() string
	CredentialsPathMock               func() string
}

func (cscm credSettingsCustomMock) ReadCredentialsFields(path string) (credential.Fields, error) {
	return cscm.ReadCredentialsFieldsMock(path)
}

func (cscm credSettingsCustomMock) ReadCredentialsValue(path string) ([]credential.ListCredData, error) {
	return cscm.ReadCredentialsValueMock(path)
}

func (cscm credSettingsCustomMock) WriteDefaultCredentialsFields(path string) error {
	return nil
}

func (cscm credSettingsCustomMock) WriteCredentialsFields(fields credential.Fields, path string) error {
	return nil
}

func (cscm credSettingsCustomMock) ProviderPath() string {
	return ""
}

func (cscm credSettingsCustomMock) CredentialsPath() string {
	return ""
}

type treeMock struct {
	tree  formula.Tree
	error error
	value string
}

func (t treeMock) Tree() (map[string]formula.Tree, error) {
	if t.value != "" {
		return map[string]formula.Tree{t.value: t.tree}, t.error
	}
	return map[string]formula.Tree{"test": t.tree}, t.error
}

func (t treeMock) MergedTree(bool) formula.Tree {
	return t.tree
}

type treeGeneratorMock struct {
}

func (t treeGeneratorMock) Generate(path string) (formula.Tree, error) {
	return formula.Tree{
		Commands: api.Commands{
			{
				Id:     "root_group",
				Parent: "root",
				Usage:  "group",
				Help:   "group for add",
			},
			{
				Id:      "root_group_verb",
				Parent:  "root_group",
				Usage:   "verb",
				Help:    "verb for add",
				Formula: true,
			},
		},
	}, nil
}

type GitRepositoryMock struct {
	zipball   func(info git.RepoInfo, version string) (io.ReadCloser, error)
	tags      func(info git.RepoInfo) (git.Tags, error)
	latestTag func(info git.RepoInfo) (git.Tag, error)
}

func (m GitRepositoryMock) Zipball(info git.RepoInfo, version string) (io.ReadCloser, error) {
	return m.zipball(info, version)
}

func (m GitRepositoryMock) Tags(info git.RepoInfo) (git.Tags, error) {
	return m.tags(info)
}

func (m GitRepositoryMock) LatestTag(info git.RepoInfo) (git.Tag, error) {
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

type TutorialFindSetterCustomMock struct {
	find func() (rtutorial.TutorialHolder, error)
	set  func(tutorial string) (rtutorial.TutorialHolder, error)
}

func (t TutorialFindSetterCustomMock) Find() (rtutorial.TutorialHolder, error) {
	return t.find()
}

func (t TutorialFindSetterCustomMock) Set(tutorial string) (rtutorial.TutorialHolder, error) {
	return t.set(tutorial)
}

type TutorialFinderMockReturnDisabled struct{}

func (TutorialFinderMockReturnDisabled) Find() (rtutorial.TutorialHolder, error) {
	return rtutorial.TutorialHolder{Current: "disabled"}, nil
}

type DirManagerCustomMock struct {
	exists func(dir string) bool
	list   func(dir string, hiddenDir bool) ([]string, error)
	isDir  func(dir string) bool
	create func(dir string) error
}

func (d DirManagerCustomMock) Exists(dir string) bool {
	return d.exists(dir)
}

func (d DirManagerCustomMock) List(dir string, hiddenDir bool) ([]string, error) {
	return d.list(dir, hiddenDir)
}

func (d DirManagerCustomMock) IsDir(dir string) bool {
	return d.isDir(dir)
}

func (d DirManagerCustomMock) Create(dir string) error {
	return d.create(dir)
}

type LocalBuilderMock struct {
	build func(workspacePath, formulaPath string) error
}

func (l LocalBuilderMock) Build(workspacePath, formulaPath string) error {
	return l.build(workspacePath, formulaPath)
}

type WatcherMock struct {
	watch func(workspacePath, formulaPath string)
}

func (w WatcherMock) Watch(workspacePath, formulaPath string) {
	w.watch(workspacePath, formulaPath)
}

type WorkspaceAddListValidatorCustomMock struct {
	add      func(workspace formula.Workspace) error
	list     func() (formula.Workspaces, error)
	validate func(workspace formula.Workspace) error
}

func (w WorkspaceAddListValidatorCustomMock) Add(workspace formula.Workspace) error {
	return w.add(workspace)
}

func (w WorkspaceAddListValidatorCustomMock) List() (formula.Workspaces, error) {
	return w.list()
}

func (w WorkspaceAddListValidatorCustomMock) Validate(workspace formula.Workspace) error {
	return w.validate(workspace)
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
		latestTag: func(info git.RepoInfo) (git.Tag, error) {
			return git.Tag{}, nil
		},
		tags: func(info git.RepoInfo) (git.Tags, error) {
			return git.Tags{git.Tag{Name: "1.0.0"}}, nil
		},
		zipball: func(info git.RepoInfo, version string) (io.ReadCloser, error) {
			return nil, nil
		},
	}
)

type FormulaExecutorMock struct {
	err error
}

func (f FormulaExecutorMock) Execute(exe formula.ExecuteData) error {
	return f.err
}

type ConfigRunnerMock struct {
	runType   formula.RunnerType
	createErr error
	findErr   error
}

func (c ConfigRunnerMock) Create(runType formula.RunnerType) error {
	return c.createErr
}

func (c ConfigRunnerMock) Find() (formula.RunnerType, error) {
	return c.runType, c.findErr
}

type RepositoryListUpdaterCustomMock struct {
	list   func() (formula.Repos, error)
	update func(name formula.RepoName, version formula.RepoVersion) error
}

func (m RepositoryListUpdaterCustomMock) List() (formula.Repos, error) {
	return m.list()
}

func (m RepositoryListUpdaterCustomMock) Update(name formula.RepoName, version formula.RepoVersion) error {
	return m.update(name, version)
}
