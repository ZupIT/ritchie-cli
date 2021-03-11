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
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

func TestListRepoRunFunc(t *testing.T) {
	someError := errors.New("some error")
	type in struct {
		RepositoryLister formula.RepositoryLister
		Tutorial         rtutorial.Finder
	}
	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name:     "someRepo1",
								Provider: "Github",
								Url:      "https://github.com/owner/repo",
								Token:    "token",
							},
						}, nil
					},
				},
				Tutorial: TutorialFinderMockReturnDisabled{},
			},
			wantErr: false,
		},
		{
			name: "Run with success with more than 1 repo",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name:     "someRepo1",
								Provider: "Github",
								Url:      "https://github.com/owner/repo1",
								Token:    "token",
							},
							{
								Name:     "someRepo2",
								Provider: "Github",
								Url:      "https://github.com/owner/repo2",
								Token:    "token",
							},
						}, nil
					},
				},
				Tutorial: TutorialFinderMockReturnDisabled{},
			},
			wantErr: false,
		},
		{
			name: "Run with success when tutorial enabled",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name:     "someRepo1",
								Provider: "Github",
								Url:      "https://github.com/owner/repo",
								Token:    "token",
							},
						}, nil
					},
				},
				Tutorial: TutorialFinderMock{},
			},
			wantErr: false,
		},
		{
			name: "Return err when list fail",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return nil, someError
					},
				},
				Tutorial: TutorialFinderMockReturnDisabled{},
			},
			wantErr: true,
		},
		{
			name: "Return err when find tutorial fail",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name:     "someRepo1",
								Provider: "Github",
								Url:      "https://github.com/owner/repo",
								Token:    "token",
							},
						}, nil
					},
				},
				Tutorial: TutorialFindSetterCustomMock{
					find: func() (rtutorial.TutorialHolder, error) {
						return rtutorial.TutorialHolder{}, someError
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lr := NewListRepoCmd(tt.in.RepositoryLister, tt.in.Tutorial)
			lr.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := lr.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("setCredentialCmd_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type RepositoryListerCustomMock struct {
	list func() (formula.Repos, error)
}

func (m RepositoryListerCustomMock) List() (formula.Repos, error) {
	return m.list()
}
