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

	gitRepo := new(mocks.GitRepositoryMock)
	gitRepo.On("Zipball", mock.Anything, mock.Anything).Return(nil, nil)
	gitRepo.On("Tags", mock.Anything).Return(git.Tags{git.Tag{Name: "1.0.0"}}, nil)
	gitRepo.On("LatestTag", mock.Anything).Return(git.Tag{}, nil)

	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: gitRepo, NewRepoInfo: github.NewRepoInfo})
	repoProviders.Add("GitLab", formula.Git{Repos: gitRepositoryWithoutTagsMock, NewRepoInfo: github.NewRepoInfo})
	repoProviders.Add("Bitbucket", formula.Git{Repos: gitRepositoryErrorsMock, NewRepoInfo: github.NewRepoInfo})

	repoList := repoProviders.List()

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
		name   string
		args   []string
		fields fields
		want   error
	}{
		{
			name: "run with success",
			args: []string{},
		},
		{
			name: "run with success when user add a new commons",
			args: []string{},
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: formula.RepoCommonsName}}, errList: nil},
				InputTextValidator: returnWithStringErr{formula.RepoCommonsName.String(), nil},
			},
		},
		{
			name: "run with success when user add a repo exitent",
			args: []string{},
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: "name-repo"}}, errList: nil},
				InputTextValidator: returnWithStringErr{"name-repo", nil},
			},
		},
		{
			name: "return nil when user add a new commons incorrectly",
			args: []string{},
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: formula.RepoCommonsName}}, errList: nil},
				InputTextValidator: returnWithStringErr{formula.RepoCommonsName.String(), nil},
				InputBool:          returnOfInputBool{false, nil},
			},
		},
		{
			name: "return nil success when user add a repo existent incorrectly",
			args: []string{},
			fields: fields{
				repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{formula.Repo{Provider: "Github", Name: "name-repo"}}, errList: nil},
				InputTextValidator: returnWithStringErr{"name-repo", nil},
				InputBool:          returnOfInputBool{false, nil},
			},
		},
		{
			name: "return error when len of tags is 0",
			args: []string{},
			fields: fields{
				InputList: []returnOfInputList{{response: "GitLab", err: nil}},
			},
			want: errors.New("please, generate a release to add your repository"),
		},
		{
			name: "return error when Repos.Tag fail",
			args: []string{},
			fields: fields{
				InputList: []returnOfInputList{{response: "Bitbucket", err: nil}},
			},
			want: errors.New("tag error"),
		},
		{
			name: "fail when repo.Add return err",
			args: []string{},
			fields: fields{
				repo:      returnOfRepoListerAdder{errAdd: someError, reposList: formula.Repos{}, errList: nil},
				InputList: []returnOfInputList{{response: "Github", err: nil}},
			},
			want: someError,
		},
		{
			name: "input bool error",
			args: []string{},
			fields: fields{
				InputBool: returnOfInputBool{false, someError},
			},
			want: someError,
		},
		{
			name: "input password error",
			args: []string{},
			fields: fields{
				InputPassword: returnWithStringErr{"", someError},
			},
			want: someError,
		},
		{
			name: "input list select provider return error",
			args: []string{},
			fields: fields{
				InputList: []returnOfInputList{{response: "item", err: someError}},
			},
			want: someError,
		},
		{
			name: "input list select version return error",
			args: []string{},
			fields: fields{
				InputList: []returnOfInputList{
					{question: "Select a tag version:", response: "", err: someError},
					{question: "", response: "Github", err: nil},
				},
			},
			want: someError,
		},
		{
			name: "input text error",
			args: []string{},
			fields: fields{
				InputTextValidator: returnWithStringErr{"", ErrFormulaCmdNotBeEmpty},
			},
			want: ErrFormulaCmdNotBeEmpty,
		},
		{
			name: "input url error",
			args: []string{},
			fields: fields{
				InputURL: returnWithStringErr{"http://localhost/mocked", someError},
			},
			want: someError,
		},
		{
			name: "tutorial status enabled",
			args: []string{},
			fields: fields{
				tutorialStatus: returnWithStringErr{"enabled", nil},
			},
		},
		{
			name: "return error when tutorial.Find fail",
			args: []string{},
			fields: fields{
				tutorialStatus: returnWithStringErr{"", someError},
			},
			want: someError,
		},
		{
			name: "fail when repo.List return err",
			args: []string{},
			fields: fields{
				repo: returnOfRepoListerAdder{errAdd: nil, reposList: nil, errList: someError},
			},
			want: someError,
		},
		{
			name: "run with success when input is stdin",
			args: []string{},
			fields: fields{
				stdin: `{"provider": "github", "name": "repo-name", "version": "0.0.0", "url": "https://url.com/repo", "token,omitempty": "", "priority": 5, "isLocal": false}\n`,
			},
		},
		{
			name: "run with success when input is stdin and version is not informed",
			args: []string{},
			fields: fields{
				stdin:           `{"provider": "github", "name": "repo-name", "version": "", "url": "https://url.com/repo", "token,omitempty": "", "priority": 5, "isLocal": false}\n`,
				detailLatestTag: "1.0.0",
			},
		},
		{
			name: "return error when user add a repo existent",
			args: []string{},
			fields: fields{
				repo:     returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{*repoTest}, errList: nil},
				InputURL: returnWithStringErr{repoTest.Url, nil},
				InputList: []returnOfInputList{
					{question: "Select a tag version:", response: "1.0.0", err: nil},
					{question: "", response: "Github", err: nil},
				},
			},
		},
		{
			name:   "fail flags with empty provider",
			args:   []string{"--provider="},
			fields: fields{},
			want:   errors.New(missingFlagText(providerFlagName)),
		},
		{
			name:   "fail flags with wrong provider",
			args:   []string{"--provider=github"},
			fields: fields{},
			want:   errors.New("please select a provider from " + strings.Join(repoList, ", ")),
		},
		{
			name:   "fail flags with empty name",
			args:   []string{"--provider=Github"},
			fields: fields{},
			want:   errors.New(missingFlagText(nameFlagName)),
		},
		{
			name:   "fail flags with empty repo url",
			args:   []string{"--provider=Github", "--name=my-repo"},
			fields: fields{},
			want:   errors.New(missingFlagText(repoUrlFlagName)),
		},
		{
			name:   "fail flags with empty tag",
			args:   []string{"--provider=Github", "--name=my-repo", "--repoUrl=github.com"},
			fields: fields{},
			want:   errors.New(missingFlagText(tagFlagName)),
		},
		{
			name:   "success flags",
			args:   []string{"--provider=Github", "--name=my-repo", "--repoUrl=github.com", "--tag=1.0.0"},
			fields: fields{},
		},
	}
	checkerManager := tree.NewChecker(treeMock{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := getFields(tt.fields)

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
			tutorialFindMock.On("Find").Return(rtutorial.TutorialHolder{Current: fields.tutorialStatus.string}, fields.tutorialStatus.error)
			repoListerAdderMock := new(mocks.RepoManager)
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

			cmd.SetArgs(tt.args)
			got := cmd.Execute()

			assert.Equal(t, tt.want, got)
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
	tutorialStatus     returnWithStringErr
}

func getFields(testFields fields) fields {
	fieldsNil := fields{
		repo:               returnOfRepoListerAdder{},
		InputTextValidator: returnWithStringErr{},
		InputPassword:      returnWithStringErr{},
		InputURL:           returnWithStringErr{},
		InputBool:          returnOfInputBool{},
		InputList:          []returnOfInputList{},
		stdin:              "",
		detailLatestTag:    "",
		tutorialStatus:     returnWithStringErr{},
	}

	fields := fields{
		repo:               returnOfRepoListerAdder{errAdd: nil, reposList: formula.Repos{}, errList: nil},
		InputTextValidator: returnWithStringErr{"mocked text", nil},
		InputPassword:      returnWithStringErr{"s3cr3t", nil},
		InputURL:           returnWithStringErr{"http://localhost/mocked", nil},
		InputBool:          returnOfInputBool{true, nil},
		InputList:          []returnOfInputList{{response: "Github", err: nil}},
		tutorialStatus:     returnWithStringErr{"disabled", nil},
		stdin:              "",
		detailLatestTag:    "",
	}

	if testFields.repo.reposList.Len() != fieldsNil.repo.reposList.Len() || testFields.repo.errAdd != fieldsNil.repo.errAdd || testFields.repo.errList != fieldsNil.repo.errList {
		fields.repo = testFields.repo
	}

	if testFields.InputTextValidator != fieldsNil.InputTextValidator {
		fields.InputTextValidator = testFields.InputTextValidator
	}

	if testFields.InputPassword != fieldsNil.InputPassword {
		fields.InputPassword = testFields.InputPassword
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

	return fields
}
