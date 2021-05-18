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
	"fmt"
	"io" 
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/stretchr/testify/assert"
)

const (
	questionAVersion  = "Select your new version for \"someRepo1\":"
	questionAVersion2 = "Select your new version for \"someRepo2\":"
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
		LatestVersion: "2.0.0",
	}

	repoTest2 := &formula.Repo{
		Provider: "Github",
		Name:     "someRepo2",
		Version:  "1.0.0",
		Url:      "https://github.com/owner/repo",
		Token:    "token",
		Priority: 1,
		IsLocal:  true,
	}

	type in struct {
		repo   formula.RepositoryListUpdater
		inList prompt.InputList
		Repos  git.Repositories
	}
	var tests = []struct {
		name       string
		in         in
		wantErr    error
		inputStdin string
		inputFlag  []string
	}{
		{
			name: "success case update someRepo1",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest2, *repoTest}, nil
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
			wantErr:    nil,
			inputStdin: createJSONEntry(repoTest),
		},
		{
			name: "success case update ALL",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest, *repoTest2}, nil
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectARepo {
							return "ALL", nil
						}
						if name == questionAVersion {
							return "1.0.0", nil
						}
						if name == questionAVersion2 {
							return "1.0.0", nil
						}
						return "any", nil
					},
				},
				Repos: defaultGitRepositoryMock,
			},
			wantErr:    nil,
			inputStdin: "",
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
			wantErr:    someError,
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
			wantErr:    someError,
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
			wantErr:    someError,
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
			wantErr:    someError,
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
			wantErr:    someError,
			inputStdin: createJSONEntry(repoTest),
		},
		{
			name: "success with flags",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest2, *repoTest}, nil
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == repoName {
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
			wantErr:   nil,
			inputFlag: []string{"--name=someRepo1", "--version=1.0.0"},
		},
		{
			name:       "success flags with 'latest' version",
			in:         in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest2, *repoTest}, nil
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == repoName {
							return "someRepo1", nil
						}
						if name == questionAVersion {
							return "2.0.0", nil
						}
						return "any", nil
					},
				},
				Repos:  defaultGitRepositoryMock,
			},
			wantErr:    nil,
			inputFlag:  []string{"--name=someRepo1", "--version=latest"},
		},
		{
			name: "fail with flags, flag name empty",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, errors.New(missingFlagText(repoName))
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListMock{},
				Repos:  defaultGitRepositoryMock,
			},
			wantErr:   errors.New(missingFlagText(repoName)),
			inputFlag: []string{"--name=", "--version=1.0.0"},
		},
		{
			name: "fail with flags, missing value to version",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return nil, errors.New(missingFlagText(repoVersion))
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						return nil
					},
				},
				inList: inputListMock{},
				Repos:  defaultGitRepositoryMock,
			},
			wantErr:   errors.New(missingFlagText(repoVersion)),
			inputFlag: []string{"--name=someRepo1", "--version="},
		},
		{
			name: "fail with flags, invalid version value",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{*repoTest}, nil
					},
					update: func(name formula.RepoName, version formula.RepoVersion) error {
						var errorMsg = fmt.Sprintf("The version %q of repository %q was not found.\n", version, name)
						return errors.New(errorMsg)
					},
				},
				inList: inputListMock{},
				Repos:  GitRepositoryMock{
					tags: func(info git.RepoInfo) (git.Tags, error) {
						return git.Tags{
							git.Tag{Name: "1.0.0"},
						}, nil
					},
				},
			},
			wantErr:   errors.New(fmt.Sprintf("The version %q of repository %q was not found.\n", "3.0.0", repoTest.Name)),
			inputFlag: []string{"--name=someRepo1", "--version=3.0.0"},
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
			newUpdateRepoFlag := NewUpdateRepoCmd(server.Client(), tt.in.repo, repoProviders, inputTextMock{}, inputPasswordMock{}, inputURLMock{}, tt.in.inList, inputTrueMock{}, inputIntMock{})

			newUpdateRepoPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")
			newUpdateRepoStdin.PersistentFlags().Bool("stdin", true, "input by stdin")
			newUpdateRepoFlag.PersistentFlags().Bool("stdin", false, "input by stdin")

			newReader := strings.NewReader(tt.inputStdin)
			newUpdateRepoStdin.SetIn(newReader)
			newUpdateRepoStdin.SetArgs([]string{})

			newUpdateRepoPrompt.SetArgs([]string{})
			newUpdateRepoFlag.SetArgs(tt.inputFlag)

			if len(tt.inputFlag) != 0 {
				newUpdateRepoPrompt.SetArgs(tt.inputFlag)
			}

			itsTestCaseWithPrompt := len(tt.inputFlag) == 0
			if out := newUpdateRepoPrompt.Execute(); out != tt.wantErr && itsTestCaseWithPrompt {
				t.Errorf("Prompt command error = %v, wantErr %v", out, tt.wantErr)
			}

			itsTestCaseWithStdin := tt.inputStdin != ""
			if out := newUpdateRepoStdin.Execute(); out != tt.wantErr && itsTestCaseWithStdin {
				t.Errorf("Stdin command error = %v, wantErr %v", out, tt.wantErr)
			}

			// itsTestCaseWithFlag := len(tt.inputFlag) != 0
			// if out := newUpdateRepoFlag.Execute(); (out == tt.wantErr) && itsTestCaseWithFlag {
			// 	fmt.Println(out)
			// 	fmt.Println(tt.wantErr)
			// 	t.Errorf("Flag command error = %q, wantErr %q", out, tt.wantErr)
			// }

			flagOut := newUpdateRepoFlag.Execute()
			if len(tt.inputFlag) != 0 {
				assert.Equal(t, tt.wantErr, flagOut)
			}

		})
	}
}

func serverMock() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}
