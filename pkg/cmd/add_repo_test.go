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
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func Test_addRepoCmd_runPrompt(t *testing.T) {
	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: defaultGitRepositoryMock, NewRepoInfo: github.NewRepoInfo})

	type fields struct {
		repo               formula.RepositoryAddLister
		repoProviders      formula.RepoProviders
		InputTextValidator prompt.InputTextValidator
		InputPassword      prompt.InputPassword
		InputURL           prompt.InputURL
		InputList          prompt.InputList
		InputBool          prompt.InputBool
		InputInt           prompt.InputInt
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
			name: "input bool error",
			fields: fields{
				repo:               defaultRepoAdderMock,
				repoProviders:      repoProviders,
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
			name:
			"Fail when repo.List return err",
			fields: fields{
				repo: repoListerAdderCustomMock{
					list: func() (formula.Repos, error) {
						return nil, errors.New("some error")
					},
				},
				repoProviders:      repoProviders,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewAddRepoCmd(
				tt.fields.repo,
				tt.fields.repoProviders,
				tt.fields.InputTextValidator,
				tt.fields.InputPassword,
				tt.fields.InputURL,
				tt.fields.InputList,
				tt.fields.InputBool,
				tt.fields.InputInt,
				TutorialFinderMock{},
			)
			o.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := o.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("init_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
