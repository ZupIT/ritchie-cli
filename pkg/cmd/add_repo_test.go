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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func TestAddRepoCmd(t *testing.T) {
	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: defaultGitRepositoryMock, NewRepoInfo: github.NewRepoInfo})

	type fields struct {
		repo               formula.RepositoryAddLister
		repoProviders      formula.RepoProviders
		repoDeleter        formula.RepositoryDeleter
		InputTextValidator prompt.InputTextValidator
		InputPassword      prompt.InputPassword
		InputURL           prompt.InputURL
		InputList          prompt.InputList
		InputBool          prompt.InputBool
		InputInt           prompt.InputInt
		stdin              string
		detailLatestTag    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Run with success",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				repoDeleter:        repositoryDeleterMock{},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fail when repo.Add return err",
			fields: fields{
				repo: repoListerAdderCustomMock{
					add: func(d formula.Repo) error {
						return errors.New("")
					},
					list: func() (formula.Repos, error) {
						return formula.Repos{}, nil
					},
				},
				repoProviders: repoProviders,
				repoDeleter: repositoryDeleterMock{
					deleteMock: func(repoName formula.RepoName) error {
						return nil
					},
				},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Fail when repo.Add and repoDeleter.Delete return err",
			fields: fields{
				repo: repoListerAdderCustomMock{
					add: func(d formula.Repo) error {
						return errors.New("")
					},
					list: func() (formula.Repos, error) {
						return formula.Repos{}, nil
					},
				},
				repoProviders: repoProviders,
				repoDeleter: repositoryDeleterMock{
					deleteMock: func(repoName formula.RepoName) error {
						return errors.New("")
					},
				},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "input bool error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				repoDeleter:        repositoryDeleterMock{},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputBoolErrorMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "input password error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				repoDeleter:        repositoryDeleterMock{},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordErrorMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "input list error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				repoDeleter:        repositoryDeleterMock{},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList:          inputListErrorMock{},
			},
			wantErr: true,
		},
		{
			name: "input text error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				repoDeleter:        repositoryDeleterMock{},
				InputTextValidator: inputTextValidatorErrorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList:          inputListMock{},
			},
			wantErr: true,
		},
		{
			name: "input text error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				repoDeleter:        repositoryDeleterMock{},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLErrorMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList:          inputListMock{},
			},
			wantErr: true,
		},
		{
			name: "Fail when repo.List return err",
			fields: fields{
				repo: repoListerAdderCustomMock{
					list: func() (formula.Repos, error) {
						return nil, errors.New("some error")
					},
				},
				repoProviders:      repoProviders,
				repoDeleter:        repositoryDeleterMock{},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Run with success when input is stdin",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				repoDeleter:        repositoryDeleterMock{},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
				stdin: "{\"provider\": \"github\", \"name\": \"repo-name\", \"version\": \"0.0.0\", \"url\": \"https://url.com/repo\", \"token,omitempty\": \"\", \"priority\": 5, \"isLocal\": false}\n",
			},
			wantErr: false,
		},
		{
			name: "Run with success when input is stdin and version is not informed",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				repoDeleter:        repositoryDeleterMock{},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
				stdin:           "{\"provider\": \"github\", \"name\": \"repo-name\", \"version\": \"\", \"url\": \"https://url.com/repo\", \"token,omitempty\": \"\", \"priority\": 5, \"isLocal\": false}\n",
				detailLatestTag: "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Fail when repo.Add return err Stdin",
			fields: fields{
				repo: repoListerAdderCustomMock{
					add: func(d formula.Repo) error {
						return errors.New("")
					},
				},
				repoProviders: repoProviders,
				repoDeleter: repositoryDeleterMock{
					deleteMock: func(repoName formula.RepoName) error {
						return nil
					},
				},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
				stdin:           "{\"provider\": \"github\", \"name\": \"repo-name\", \"version\": \"\", \"url\": \"https://url.com/repo\", \"token,omitempty\": \"\", \"priority\": 5, \"isLocal\": false}\n",
				detailLatestTag: "1.0.0",
			},
			wantErr: true,
		},
		{
			name: "Fail when repo.Add and repoDeleter.Delete return err stdin",
			fields: fields{
				repo: repoListerAdderCustomMock{
					add: func(d formula.Repo) error {
						return errors.New("error to add")
					},
				},
				repoProviders: repoProviders,
				repoDeleter: repositoryDeleterMock{
					deleteMock: func(repoName formula.RepoName) error {
						return errors.New("error to delete")
					},
				},
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURL:           inputURLMock{},
				InputBool:          inputTrueMock{},
				InputInt:           inputIntMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "Github", nil
					},
				},
				stdin:           "{\"provider\": \"github\", \"name\": \"repo-name\", \"version\": \"\", \"url\": \"https://url.com/repo\", \"token,omitempty\": \"\", \"priority\": 5, \"isLocal\": false}\n",
				detailLatestTag: "1.0.0",
			},
			wantErr: true,
		},
	}
	checkerManager := tree.NewChecker(treeMock{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detailMock := new(mocks.DetailManagerMock)
			detailMock.On("LatestTag", mock.Anything).Return(tt.fields.detailLatestTag)
			cmd := NewAddRepoCmd(
				tt.fields.repo,
				tt.fields.repoProviders,
				tt.fields.repoDeleter,
				tt.fields.InputTextValidator,
				tt.fields.InputPassword,
				tt.fields.InputURL,
				tt.fields.InputList,
				tt.fields.InputBool,
				tt.fields.InputInt,
				TutorialFinderMock{},
				checkerManager,
				detailMock,
			)

			if tt.fields.stdin != "" {
				newReader := strings.NewReader(tt.fields.stdin)
				cmd.SetIn(newReader)
				cmd.PersistentFlags().Bool("stdin", true, "input by stdin")
			} else {
				cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			}

			err := cmd.Execute()

			assert.Equal(t, tt.wantErr, (err != nil))
		})
	}
}
