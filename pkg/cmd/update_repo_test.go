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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func TestUpdateRepoRun(t *testing.T) {
	someError := errors.New("some error")

	repoTest := &formula.Repo{
		Provider: "Github",
		Name:     "someRepo1",
		Version:  "1.0.0",
		Url:      "https://github.com/owner/repo",
		Token:    "token",
		Priority: 2,
	}

	type in struct {
		repo   formula.RepositoryListUpdater
		inList prompt.InputList
		Repos  git.Repositories
	}
	var tests = []struct {
		name       string
		in         in
		wantErr    bool
		inputStdin string
	}{
		{
			name: "success case",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest}, nil
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectARepo {
							return "someRepo1", nil
						}
						if name == questionAVersion {
							return "1.0.0", nil
						}
						return "any", nil
					},
				},
				Repos: defaultGitRepositoryMock,
			},
			wantErr:    false,
			inputStdin: createJSONEntry(repoTest),
		},
		{
			name: "fails when repo list returns an error",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, someError
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListMock{},
				Repos:  defaultGitRepositoryMock,
			},
			wantErr:    true,
			inputStdin: "",
		},
		{
			name: "fails when question about select repo returns an error",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest}, nil
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectARepo {
							return "", someError
						}
						return "any", nil
					},
				},
				Repos: defaultGitRepositoryMock,
			},
			wantErr:    true,
			inputStdin: "",
		},
		{
			name: "fails when repos tags returns an error",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest}, nil
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectARepo {
							return "someRepo1", nil
						}
						if name == questionAVersion {
							return "1.0.0", nil
						}
						return "any", nil
					},
				},
				Repos: GitRepositoryMock{
					latestTag: func(info git.RepoInfo) (git.Tag, error) {
						return git.Tag{}, nil
					},
					tags: func(info git.RepoInfo) (git.Tags, error) {
						return git.Tags{}, someError
					},
					zipball: func(info git.RepoInfo, version string) (io.ReadCloser, error) {
						return nil, nil
					},
				},
			},
			wantErr:    true,
			inputStdin: "",
		},
		{
			name: "fails when question about select version returns an error",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest}, nil
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectARepo {
							return "someRepo1", nil
						}
						if name == questionAVersion {
							return "", someError
						}
						return "any", nil
					},
				},
				Repos: defaultGitRepositoryMock,
			},
			wantErr:    true,
			inputStdin: "",
		},
		{
			name: "fails when repo update returns an error",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest}, nil
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return someError
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectARepo {
							return "someRepo1", nil
						}
						if name == questionAVersion {
							return "1.0.0", nil
						}
						return "any", nil
					},
				},
				Repos: defaultGitRepositoryMock,
			},
			wantErr:    true,
			inputStdin: createJSONEntry(repoTest),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := serverMock()
			defer server.Close()

			repoProviders := formula.NewRepoProviders()
			repoProviders.Add("Github", formula.Git{Repos: tt.in.Repos, NewRepoInfo: github.NewRepoInfo})

			newUpdateRepoPrompt := NewUpdateRepoCmd(server.Client(), tt.in.repo, repoProviders, inputTextMock{}, inputPasswordMock{}, inputURLMock{}, tt.in.inList, inputTrueMock{}, inputIntMock{})
			newUpdateRepoStdin := NewUpdateRepoCmd(server.Client(), tt.in.repo, repoProviders, inputTextMock{}, inputPasswordMock{}, inputURLMock{}, tt.in.inList, inputTrueMock{}, inputIntMock{})

			newUpdateRepoPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")
			newUpdateRepoStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

			newReader := strings.NewReader(tt.inputStdin)
			newUpdateRepoStdin.SetIn(newReader)

			if err := newUpdateRepoPrompt.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Prompt command error = %v, wantErr %v", err, tt.wantErr)
			}

			itsTestCaseWithStdin := tt.inputStdin != ""
			if err := newUpdateRepoStdin.Execute(); (err != nil) != tt.wantErr && itsTestCaseWithStdin {
				t.Errorf("Stdin command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func serverMock() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}
