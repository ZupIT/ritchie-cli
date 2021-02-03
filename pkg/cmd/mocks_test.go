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
	"encoding/json"
	"errors"
	"io"

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/git"

	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

type inputTextMock struct{}

func (inputTextMock) Text(name string, required bool, helper ...string) (string, error) {
	return "mocked text", nil
}

func (inputTextMock) TextWithValidate(name string, validate func(interface{}) error, helper ...string) (string, error) {
	return "mocked text", nil
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

type inputIntMock struct{}

func (inputIntMock) Int(name string, helper ...string) (int64, error) {
	return 0, nil
}

type inputIntErrorMock struct{}

func (inputIntErrorMock) Int(name string, helper ...string) (int64, error) {
	return 0, errors.New("some error")
}

type inputPasswordMock struct{}

func (inputPasswordMock) Password(label string, helper ...string) (string, error) {
	return "s3cr3t", nil
}

type inputPasswordErrorMock struct{}

func (inputPasswordErrorMock) Password(label string, helper ...string) (string, error) {
	return "", errors.New("password error")
}

type autocompleteGenMock struct{}

func (autocompleteGenMock) Generate(s autocomplete.ShellName, cmd *cobra.Command) (string, error) {
	return "autocomplete", nil
}

type inputTrueMock struct{}

func (inputTrueMock) Bool(name string, items []string, helper ...string) (bool, error) {
	return true, nil
}

type inputFalseMock struct{}

func (inputFalseMock) Bool(name string, items []string, helper ...string) (bool, error) {
	return false, nil
}

type inputBoolErrorMock struct{}

func (inputBoolErrorMock) Bool(name string, items []string, helper ...string) (bool, error) {
	return false, errors.New("error on boolean list")
}

type inputListMock struct{}

func (inputListMock) List(name string, items []string, helper ...string) (string, error) {
	return "item-mocked", nil
}

type inputListCustomMock struct {
	list func(name string, items []string) (string, error)
}

func (m inputListCustomMock) List(name string, items []string, helper ...string) (string, error) {
	return m.list(name, items)
}

type inputListErrorMock struct{}

func (inputListErrorMock) List(name string, items []string, helper ...string) (string, error) {
	return "item-mocked", errors.New("some error")
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

type envSetterMock struct{}

func (envSetterMock) Set(_ string) (env.Holder, error) {
	return env.Holder{}, nil
}

type envFinderCustomMock struct {
	find func() (env.Holder, error)
}

func (e envFinderCustomMock) Find() (env.Holder, error) {
	return e.find()
}

type envFinderMock struct{}

func (envFinderMock) Find() (env.Holder, error) {
	return env.Holder{}, nil
}

type envFindSetterMock struct{}

func (envFindSetterMock) Find() (env.Holder, error) {
	f := envFinderMock{}
	return f.Find()
}

func (envFindSetterMock) Set(env string) (env.Holder, error) {
	s := envSetterMock{}
	return s.Set(env)
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

func (s credSettingsMock) ReadCredentialsValueInEnv(path string, env string) ([]credential.ListCredData, error) {
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
	ReadCredentialsValueInEnvMock     func(path string, env string) ([]credential.ListCredData, error)
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

func (cscm credSettingsCustomMock) ReadCredentialsValueInEnv(path string, env string) ([]credential.ListCredData, error) {
	return cscm.ReadCredentialsValueInEnvMock(path, env)
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

func (t treeMock) Tree() (map[formula.RepoName]formula.Tree, error) {
	if t.value != "" {
		return map[formula.RepoName]formula.Tree{formula.RepoName(t.value): t.tree}, t.error
	}
	return map[formula.RepoName]formula.Tree{"test": t.tree}, t.error
}

func (t treeMock) MergedTree(bool) formula.Tree {
	return t.tree
}

type treeGeneratorMock struct {
	err error
}

func (t treeGeneratorMock) Generate(path string) (formula.Tree, error) {
	return formula.Tree{
		Commands: api.Commands{
			"root_group": {
				Parent: "root",
				Usage:  "group",
				Help:   "group for add",
			},
			"root_group_verb": {
				Parent:  "root_group",
				Usage:   "verb",
				Help:    "verb for add",
				Formula: true,
			},
		},
	}, t.err
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

type WorkspaceAddListerCustomMock struct {
	add  func(workspace formula.Workspace) error
	list func() (formula.Workspaces, error)
}

func (w WorkspaceAddListerCustomMock) Add(workspace formula.Workspace) error {
	return w.add(workspace)
}

func (w WorkspaceAddListerCustomMock) List() (formula.Workspaces, error) {
	return w.list()
}

var (
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

	gitRepositoryWithoutTagsMock = GitRepositoryMock{
		latestTag: func(info git.RepoInfo) (git.Tag, error) {
			return git.Tag{}, nil
		},
		tags: func(info git.RepoInfo) (git.Tags, error) {
			return git.Tags{}, nil
		},
		zipball: func(info git.RepoInfo, version string) (io.ReadCloser, error) {
			return nil, nil
		},
	}

	gitRepositoryErrorsMock = GitRepositoryMock{
		latestTag: func(info git.RepoInfo) (git.Tag, error) {
			return git.Tag{}, errors.New("latest tag error")
		},
		tags: func(info git.RepoInfo) (git.Tags, error) {
			return git.Tags{}, errors.New("tag error")
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

func createJSONEntry(v interface{}) string {
	s, _ := json.Marshal(v)
	return string(s)
}
