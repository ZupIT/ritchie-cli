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
		repo               returnOfRepoListerAdder
		InputTextValidator returnWithStringErr
		InputPassword      returnWithStringErr
		InputURL           returnWithStringErr
		InputList          []returnOfInputList
		InputBool          returnOfInputBool
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
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "Github", err: nil}},
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
			},
			wantErr: true,
		},
		{
			name: "input list error",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
				InputTextValidator: returnWithStringErr{"mocked text", nil},
				InputPassword:      returnWithStringErr{"s3cr3t", nil},
				InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
				InputBool:          returnOfInputBool{true, nil},
				InputList:          []returnOfInputList{{response: "item", err: someError}},
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
			},
			wantErr: true,
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
				stdin:              "{\"provider\": \"github\", \"name\": \"repo-name\", \"version\": \"0.0.0\", \"url\": \"https://url.com/repo\", \"token,omitempty\": \"\", \"priority\": 5, \"isLocal\": false}\n",
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
				stdin:              "{\"provider\": \"github\", \"name\": \"repo-name\", \"version\": \"\", \"url\": \"https://url.com/repo\", \"token,omitempty\": \"\", \"priority\": 5, \"isLocal\": false}\n",
				detailLatestTag:    "1.0.0",
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
			tutorialFindMock.On("Find").Return(rtutorial.TutorialHolder{Current: "disabled"}, nil)
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
