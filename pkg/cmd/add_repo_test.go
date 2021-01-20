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
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

func TestAddRepoCmd(t *testing.T) {
	someError := errors.New("some error")

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
			name:    "Run with success",
			wantErr: false,
		},
		{
			name: "Run with success when user add a new commons",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: formula.RepoCommonsName}}, errList: nil},
				InputTextValidator: returnWithStringErr{formula.RepoCommonsName.String(), nil},
			},
			wantErr: false,
		},
		{
			name: "Run with success when user add a repo exitent",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: "name-repo"}}, errList: nil},
				InputTextValidator: returnWithStringErr{"name-repo", nil},
			},
			wantErr: false,
		},
		{
			name: "Return nil when user add a new commons incorrectly",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: formula.RepoCommonsName}}, errList: nil},
				InputTextValidator: returnWithStringErr{formula.RepoCommonsName.String(), nil},
				InputBool:          returnOfInputBool{false, nil},
			},
			wantErr: false,
		},
		{
			name: "Return nil success when user add a repo exitent incorrectly",
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: "name-repo"}}, errList: nil},
				InputTextValidator: returnWithStringErr{"name-repo", nil},
				InputBool:          returnOfInputBool{false, nil},
			},
			wantErr: false,
		},
		{
			name: "return error when len of tags is 0",
			fields: fields{
				InputList: []returnOfInputList{{response: "GitLab", err: nil}},
			},
			wantErr: true,
		},
		{
			name: "return error when Repos.Tag fail",
			fields: fields{
				InputList: []returnOfInputList{{response: "Bitbucket", err: nil}},
			},
			wantErr: true,
		},
		{
			name: "Fail when repo.Add return err",
			fields: fields{
				repo:      returnOfRepoListerAdder{errAdd: someError, reposList: formula.Repos{}, errList: nil},
				InputList: []returnOfInputList{{response: "Github", err: nil}},
			},
			wantErr: true,
		},
		{
			name: "Fail when repo.Add return 401 err",
			fields: fields{
				repo:       returnOfRepoListerAdder{errAdd: someError, reposList: formula.Repos{}, errList: nil},
				InputList:  []returnOfInputList{{response: "Github", err: nil}},
				gitRepoTag: returnOfGitRepoTag{nil, errors.New("401 - Unauthorized")},
			},
			wantErr: true,
		},
		{
			name: "input bool error",
			fields: fields{
				InputBool: returnOfInputBool{false, someError},
			},
			wantErr: true,
		},
		{
			name: "input list select provider return error",
			fields: fields{
				InputList: []returnOfInputList{{response: "item", err: someError}},
			},
			wantErr: true,
		},
		{
			name: "input list select version return error",
			fields: fields{
				InputList: []returnOfInputList{
					{question: "Select a tag version:", response: "", err: someError},
					{question: "", response: "Github", err: nil},
				},
			},
			wantErr: true,
		},
		{
			name: "input text error",
			fields: fields{
				InputTextValidator: returnWithStringErr{"mocked text", someError},
			},
			wantErr: true,
		},
		{
			name: "input url error",
			fields: fields{
				InputURL: returnWithStringErr{"http://localhost/mocked", someError},
			},
			wantErr: true,
		},
		{
			name: "Tutorial status enabled",
			fields: fields{
				tutorialStatus: returnWithStringErr{"enabled", nil},
			},
			wantErr: false,
		},
		{
			name: "return error when tutorial.Find fail",
			fields: fields{
				tutorialStatus: returnWithStringErr{"", someError},
			},
			wantErr: true,
		},
		{
			name: "Fail when repo.List return err",
			fields: fields{
				repo: returnOfRepoListerAdder{errAdd: nil, reposList: nil, errList: someError},
			},
			wantErr: true,
		},
		{
			name: "Run with success when input is stdin",
			fields: fields{
				stdin: `{"provider": "github", "name": "repo-name", "version": "0.0.0", "url": "https://url.com/repo", "token,omitempty": "", "priority": 5, "isLocal": false}\n`,
			},
			wantErr: false,
		},
		{
			name: "Run with success when input is stdin and version is not informed",
			fields: fields{
				stdin:           `{"provider": "github", "name": "repo-name", "version": "", "url": "https://url.com/repo", "token,omitempty": "", "priority": 5, "isLocal": false}\n`,
				detailLatestTag: "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Return error when user add a repo existent",
			fields: fields{
				repo:     returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{*repoTest}, errList: nil},
				InputURL: returnWithStringErr{repoTest.Url, nil},
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
			fields := getFields(tt.fields)

			gitRepo := new(mocks.GitRepositoryMock)
			gitRepo.On("Zipball", mock.Anything, mock.Anything).Return(nil, nil)
			gitRepo.On("Tags", mock.Anything).Return(fields.gitRepoTag.Tags, fields.gitRepoTag.error)
			gitRepo.On("LatestTag", mock.Anything).Return(git.Tag{}, nil)

			repoProviders := formula.NewRepoProviders()
			repoProviders.Add("Github", formula.Git{Repos: gitRepo, NewRepoInfo: github.NewRepoInfo})
			repoProviders.Add("GitLab", formula.Git{Repos: gitRepositoryWithoutTagsMock, NewRepoInfo: github.NewRepoInfo})
			repoProviders.Add("Bitbucket", formula.Git{Repos: gitRepositoryErrorsMock, NewRepoInfo: github.NewRepoInfo})

			detailMock := new(mocks.DetailManagerMock)
			detailMock.On("LatestTag", mock.Anything).Return(fields.detailLatestTag)
			inputURLMock := new(mocks.InputURLMock)
			inputURLMock.On("URL", mock.Anything, mock.Anything).Return(fields.InputURL.string, fields.InputURL.error)
			inputBoolMock := new(mocks.InputBoolMock)
			inputBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(fields.InputBool.bool, fields.InputBool.error)
			inputListMock := addInputList(fields.InputList)
			inputIntMock := new(mocks.InputIntMock)
			inputIntMock.On("Int", mock.Anything, mock.Anything).Return(int64(0), nil)
			inputTextValidatorMock := new(mocks.InputTextValidatorMock)
			inputTextValidatorMock.On("Text", mock.Anything, mock.Anything).Return(fields.InputTextValidator.string, fields.InputTextValidator.error)
			tutorialFindMock := new(mocks.TutorialFindSetterMock)
			tutorialFindMock.On("Find").Return(rtutorial.TutorialHolder{Current: fields.tutorialStatus.string}, fields.tutorialStatus.error)
			repoListerAdderMock := new(mocks.RepoListerAdderMock)
			repoListerAdderMock.On("Add", mock.Anything).Return(fields.repo.errAdd)
			repoListerAdderMock.On("List").Return(fields.repo.reposList, fields.repo.errList)
			credResolverMock := new(mocks.CredResolverMock)
			credResolverMock.On("Resolve", mock.Anything).Return(fields.credResolver.string, fields.credResolver.error)

			cmd := NewAddRepoCmd(
				repoListerAdderMock,
				repoProviders,
				credResolverMock,
				inputTextValidatorMock,
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

type returnOfGitRepoTag struct {
	git.Tags
	error
}

type fields struct {
	repo               returnOfRepoListerAdder
	InputTextValidator returnWithStringErr
	InputURL           returnWithStringErr
	InputList          []returnOfInputList
	InputBool          returnOfInputBool
	stdin              string
	detailLatestTag    string
	tutorialStatus     returnWithStringErr
	credResolver       returnWithStringErr
	gitRepoTag         returnOfGitRepoTag
}

func getFields(testFields fields) fields {
	fieldsNil := fields{
		repo:               returnOfRepoListerAdder{},
		InputTextValidator: returnWithStringErr{},
		InputURL:           returnWithStringErr{},
		InputBool:          returnOfInputBool{},
		InputList:          []returnOfInputList{},
		stdin:              "",
		detailLatestTag:    "",
		tutorialStatus:     returnWithStringErr{},
		credResolver:       returnWithStringErr{},
		gitRepoTag:         returnOfGitRepoTag{},
	}

	fields := fields{
		repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
		InputTextValidator: returnWithStringErr{"mocked text", nil},
		InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
		InputBool:          returnOfInputBool{true, nil},
		InputList:          []returnOfInputList{{response: "Github", err: nil}},
		tutorialStatus:     returnWithStringErr{"disabled", nil},
		stdin:              "",
		detailLatestTag:    "",
		credResolver:       returnWithStringErr{"token", nil},
		gitRepoTag:         returnOfGitRepoTag{git.Tags{git.Tag{Name: "1.0.0"}}, nil},
	}

	if testFields.repo.reposList.Len() != fieldsNil.repo.reposList.Len() || testFields.repo.errAdd != fieldsNil.repo.errAdd || testFields.repo.errList != fieldsNil.repo.errList {
		fields.repo = testFields.repo
	}

	if testFields.InputTextValidator != fieldsNil.InputTextValidator {
		fields.InputTextValidator = testFields.InputTextValidator
	}

	if testFields.InputURL != fieldsNil.InputURL {
		fields.InputURL = testFields.InputURL
	}

	if testFields.InputBool != fieldsNil.InputBool {
		fields.InputBool = testFields.InputBool
	}

	if len(testFields.InputList) != 0 {
		fields.InputList = testFields.InputList
	}

	if testFields.tutorialStatus != fieldsNil.tutorialStatus {
		fields.tutorialStatus = testFields.tutorialStatus
	}

	if testFields.stdin != fieldsNil.stdin {
		fields.stdin = testFields.stdin
	}

	if testFields.detailLatestTag != fieldsNil.detailLatestTag {
		fields.detailLatestTag = testFields.detailLatestTag
	}

	if testFields.credResolver != fieldsNil.credResolver {
		fields.credResolver = testFields.credResolver
	}

	if len(testFields.gitRepoTag.Tags) != len(fieldsNil.gitRepoTag.Tags) || testFields.gitRepoTag.error != fieldsNil.gitRepoTag.error {
		fields.gitRepoTag = testFields.gitRepoTag
	}

	return fields
}
