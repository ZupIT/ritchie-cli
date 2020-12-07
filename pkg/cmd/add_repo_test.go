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

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddRepoCmd(t *testing.T) {
	someError := errors.New("some error")
	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: defaultGitRepositoryMock, NewRepoInfo: github.NewRepoInfo})

	repoTest := &formula.Repo{
		Provider: "Github",
		Name:     "someRepo1",
		Version:  "1.0.0",
		Url:      "https://github.com/owner/repo",
		Token:    "token",
		Priority: 2,
	}

	repoListerPopulated := new(mocks.RepoListerAdderMock)
	repoListerPopulated.On("Add", mock.Anything).Return(nil)
	repoListerPopulated.On("List").Return(formula.Repos{*repoTest}, nil)

	type returnOfInputBool struct {
		bool
		error
	}
	type returnOfInputList struct {
		question string
		response string
		err      error
	}

	addInputList := func(input []returnOfInputList) *mocks.InputListMock {
		inputListMock := new(mocks.InputListMock)

		for _, input := range input {
			question := input.question
			if input.question == "" {
				question = mock.Anything
			}
			inputListMock.On("List", question, mock.Anything, mock.Anything).Return(input.response, input.err)
		}
		return inputListMock
	}

	type fields struct {
		repo               formula.RepositoryAddLister
		repoProviders      formula.RepoProviders
		InputTextValidator prompt.InputTextValidator
		InputPassword      prompt.InputPassword
		InputURLText       string
		InputURLErr        error
		InputList          []returnOfInputList
		InputBool          returnOfInputBool
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
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURLText:       "http://localhost/mocked",
				InputURLErr:        nil,
				InputBool:          returnOfInputBool{true, nil},
				InputInt:           inputIntMock{},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
			},
			wantErr: false,
		},
		{
			name: "input bool error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURLText:       "http://localhost/mocked",
				InputURLErr:        nil,
				InputBool:          returnOfInputBool{false, someError},
				InputInt:           inputIntMock{},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
			},
			wantErr: true,
		},
		{
			name: "input password error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordErrorMock{},
				InputURLText:       "http://localhost/mocked",
				InputURLErr:        nil,
				InputBool:          returnOfInputBool{true, nil},
				InputInt:           inputIntMock{},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
			},
			wantErr: true,
		},
		{
			name: "input list error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURLText:       "http://localhost/mocked",
				InputURLErr:        nil,
				InputBool:          returnOfInputBool{true, nil},
				InputInt:           inputIntMock{},
				InputList:          []returnOfInputList{{response: "item", err: someError}},
			},
			wantErr: true,
		},
		{
			name: "input text error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				InputTextValidator: inputTextValidatorErrorMock{},
				InputPassword:      inputPasswordMock{},
				InputURLText:       "http://localhost/mocked",
				InputURLErr:        nil,
				InputBool:          returnOfInputBool{true, nil},
				InputInt:           inputIntMock{},
				InputList:          []returnOfInputList{{response: "item", err: nil}},
			},
			wantErr: true,
		},
		{
			name: "input url error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURLText:       "http://localhost/mocked",
				InputURLErr:        errors.New("error on input url"),
				InputBool:          returnOfInputBool{true, nil},
				InputInt:           inputIntMock{},
				InputList:          []returnOfInputList{{response: "item", err: nil}},
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
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURLText:       "http://localhost/mocked",
				InputURLErr:        nil,
				InputBool:          returnOfInputBool{true, nil},
				InputInt:           inputIntMock{},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
			},
			wantErr: true,
		},
		{
			name: "Run with success when input is stdin",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURLText:       "http://localhost/mocked",
				InputURLErr:        nil,
				InputBool:          returnOfInputBool{true, nil},
				InputInt:           inputIntMock{},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				stdin:              "{\"provider\": \"github\", \"name\": \"repo-name\", \"version\": \"0.0.0\", \"url\": \"https://url.com/repo\", \"token,omitempty\": \"\", \"priority\": 5, \"isLocal\": false}\n",
			},
			wantErr: false,
		},
		{
			name: "Run with success when input is stdin and version is not informed",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURLText:       "http://localhost/mocked",
				InputURLErr:        nil,
				InputBool:          returnOfInputBool{true, nil},
				InputInt:           inputIntMock{},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				stdin:              "{\"provider\": \"github\", \"name\": \"repo-name\", \"version\": \"\", \"url\": \"https://url.com/repo\", \"token,omitempty\": \"\", \"priority\": 5, \"isLocal\": false}\n",
				detailLatestTag:    "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Return error when user add a repo existent",
			fields: fields{
				repo:               repoListerPopulated,
				repoProviders:      repoProviders,
				InputTextValidator: inputTextValidatorMock{},
				InputPassword:      inputPasswordMock{},
				InputURLText:       repoTest.Url,
				InputURLErr:        nil,
				InputBool:          returnOfInputBool{true, nil},
				InputInt:           inputIntMock{},
				InputList: []returnOfInputList{
					{question: "Select a tag version:", response: "1.0.0", err: nil},
					{question: "", response: "Github", err: nil},
				},
			},
			wantErr: false,
		},
	}
	checkerManager := tree.NewChecker(treeMock{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detailMock := new(mocks.DetailManagerMock)
			detailMock.On("LatestTag", mock.Anything).Return(tt.fields.detailLatestTag)
			inputURLMock := new(mocks.InputURLMock)
			inputURLMock.On("URL", mock.Anything, mock.Anything).Return(tt.fields.InputURLText, tt.fields.InputURLErr)
			inputBoolMock := new(mocks.InputBoolMock)
			inputBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(tt.fields.InputBool.bool, tt.fields.InputBool.error)
			inputListMock := addInputList(tt.fields.InputList)

			cmd := NewAddRepoCmd(
				tt.fields.repo,
				tt.fields.repoProviders,
				tt.fields.InputTextValidator,
				tt.fields.InputPassword,
				inputURLMock,
				inputListMock,
				inputBoolMock,
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
