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
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddRepoCmd(t *testing.T) {
	someError := errors.New("some error")

	gitRepo := new(mocks.GitRepositoryMock)
	gitRepo.On("Zipball", mock.Anything, mock.Anything).Return(nil, nil)
	gitRepo.On("Tags", mock.Anything).Return(git.Tags{git.Tag{Name: "1.0.0"}}, nil)
	gitRepo.On("LatestTag", mock.Anything).Return(git.Tag{}, nil)

	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: gitRepo, NewRepoInfo: github.NewRepoInfo})

	repoTest := &formula.Repo{
		Provider: "Github",
		Name:     "someRepo1",
		Version:  "1.0.0",
		Url:      "https://github.com/owner/repo",
		Token:    "token",
		Priority: 2,
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

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Run with success",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: false,
		},
		{
			name: "Run with success when user add a new commons",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: formula.RepoCommonsName}}, errList: nil},
				InputTextValidator: returnWithStringErr{formula.RepoCommonsName.String(), nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: false,
		},
		{
			name: "Run with success when user add a repo exitent",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: "name-repo"}}, errList: nil},
				InputTextValidator: returnWithStringErr{"name-repo", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: false,
		},
		{
			name: "Return nil when user add a new commons incorrectly",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: formula.RepoCommonsName}}, errList: nil},
				InputTextValidator: returnWithStringErr{formula.RepoCommonsName.String(), nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{false, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: false,
		},
		{
			name: "Return nil success when user add a repo exitent incorrectly",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: "name-repo"}}, errList: nil},
				InputTextValidator: returnWithStringErr{"name-repo", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{false, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: false,
		},
		{
			name: "input bool error",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{false, someError},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: true,
		},
		{
			name: "input password error",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"", someError},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: true,
		},
		{
			name: "input list select provider return error",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "item", err: someError}},
				tutorialStatus:     "disabled",
			},
			wantErr: true,
		},
		{
			name: "input list select version return error",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList: []returnOfInputList{
					{question: "Select a tag version:", response: "", err: someError},
					{question: "", response: "Github", err: nil},
				},
				tutorialStatus: "disabled",
			},
			wantErr: true,
		},
		{
			name: "input text error",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", someError},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "item", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: true,
		},
		{
			name: "input url error",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", someError},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "item", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: true,
		},
		{
			name: "Tutorial status enabled",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				tutorialStatus:     "enabled",
			},
			wantErr: false,
		},
		{
			name: "Fail when repo.List return err",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: nil, errList: someError},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				tutorialStatus:     "disabled",
			},
			wantErr: true,
		},
		{
			name: "Run with success when input is stdin",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				stdin:              `{"provider": "github", "name": "repo-name", "version": "0.0.0", "url": "https://url.com/repo", "token,omitempty": "", "priority": 5, "isLocal": false}\n`,
				tutorialStatus:     "disabled",
			},
			wantErr: false,
		},
		{
			name: "Run with success when input is stdin and version is not informed",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
				stdin:              `{"provider": "github", "name": "repo-name", "version": "", "url": "https://url.com/repo", "token,omitempty": "", "priority": 5, "isLocal": false}\n`,
				detailLatestTag:    "1.0.0",
				tutorialStatus:     "disabled",
			},
			wantErr: false,
		},
		{
			name: "Return error when user add a repo existent",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{*repoTest}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{repoTest.Url, nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList: []returnOfInputList{
					{question: "Select a tag version:", response: "1.0.0", err: nil},
					{question: "", response: "Github", err: nil},
				},
				tutorialStatus: "disabled",
			},
			wantErr: false,
		},
	}
	checkerManager := tree.NewChecker(treeMock{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields

			detailMock := new(mocks.DetailManagerMock)
			detailMock.On("LatestTag", mock.Anything).Return(fields.detailLatestTag)
			inputURLMock := new(mocks.InputURLMock)
			inputURLMock.On("URL", mock.Anything, mock.Anything).Return(fields.InputURL.string, fields.InputURL.error)
			inputBoolMock := new(mocks.InputBoolMock)
			inputBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(fields.InputBool.bool, fields.InputBool.error)
			inputListMock := addInputList(fields.InputList)
			inputIntMock := new(mocks.InputIntMock)
			inputIntMock.On("Int", mock.Anything, mock.Anything).Return(int64(0), nil)
			inputPasswordMock := new(mocks.InputPasswordMock)
			inputPasswordMock.On("Password", mock.Anything, mock.Anything).Return(fields.InputPassword.string, fields.InputPassword.error)
			inputTextValidatorMock := new(mocks.InputTextValidatorMock)
			inputTextValidatorMock.On("Text", mock.Anything, mock.Anything).Return(fields.InputTextValidator.string, fields.InputTextValidator.error)
			tutorialFindMock := new(mocks.TutorialFindSetterMock)
			tutorialFindMock.On("Find").Return(rtutorial.TutorialHolder{Current: fields.tutorialStatus}, nil)
			repoListerAdderMock := new(mocks.RepoListerAdderMock)
			repoListerAdderMock.On("Add", mock.Anything).Return(fields.repo.errAdd)
			repoListerAdderMock.On("List").Return(fields.repo.reposList, fields.repo.errList)

			cmd := NewAddRepoCmd(
				repoListerAdderMock,
				repoProviders,
				inputTextValidatorMock,
				inputPasswordMock,
				inputURLMock,
				inputListMock,
				inputBoolMock,
				inputIntMock,
				tutorialFindMock,
				checkerManager,
				detailMock,
			)

			if fields.stdin != "" {
				newReader := strings.NewReader(fields.stdin)
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

type returnWithStringErr struct {
	string
	error
}

type returnOfInputBool struct {
	bool
	error
}

type returnOfInputList struct {
	question, response string
	err                error
}

type returnOfRepoListerAdder struct {
	errAdd, errList error
	reposList       formula.Repos
}

type fields struct {
	repo               returnOfRepoListerAdder
	InputTextValidator returnWithStringErr
	InputPassword      returnWithStringErr
	InputURL           returnWithStringErr
	InputList          []returnOfInputList
	InputBool          returnOfInputBool
	stdin              string
	detailLatestTag    string
	tutorialStatus     string
}

func defaultFields(testFields fields) fields {
	defaultFields := fields{
		repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
		InputTextValidator: returnWithStringErr{"mocked text", nil},
		InputPassword:      returnWithStringErr{"s3cr3t", nil},
		InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
		InputBool:          returnOfInputBool{true, nil},
		InputList:          []returnOfInputList{{response: "Github", err: nil}},
		tutorialStatus:     "disabled",
	}

	if testFields.repo != returnOfRepoListerAdder{} {
		defaultFields.repo = testFields.repo
	}

	if testFields.InputTextValidator != returnWithStringErr{} {
		defaultFields.InputTextValidator = testFields.InputTextValidator
	}

	if testFields.InputPassword != returnWithStringErr{} {
		defaultFields.InputPassword = testFields.InputPassword
	}

	if testFields.InputURL != returnWithStringErr{} {
		defaultFields.InputURL = testFields.InputURL
	}

	if testFields.InputBool != returnOfInputBool{} {
		defaultFields.InputBool = testFields.InputBool
	}

	if testFields.InputList != []returnOfInputList{} {
		defaultFields.InputList = testFields.InputList
	}

	if testFields.tutorialStatus != "" {
		defaultFields.tutorialStatus = testFields.tutorialStatus
	}

	return defaultFields
}
