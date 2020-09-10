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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func Test_NewUpdateRepoCmd(t *testing.T) {
	type in struct {
		repo   formula.RepositoryListUpdater
		inList prompt.InputList
	}
	var tests = []struct {
		name       string
		in         in
		wantErr    bool
		inputStdin string
	}{
		{
			name: "success set formula run",
			in: in{
				repo: RepositoryListUpdaterCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Provider: "Github",
								Name:     "someRepo1",
								Version:  "1.0.0",
							},
						}, nil
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == "Select a repository to update: " {
							return "someRepo1", nil
						}
						if name == "Select your new version:" {
							return "1.0.0", nil
						}
						return "any", nil
					},
				},
			},
			wantErr:    false,
			inputStdin: "{\"name\": \"someRepo1\", \"version\": \"1.0.0\", \"url\": \"https://lala.com\", \"token\": \"\", priority:\"2\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := serverMock()
			defer server.Close()

			repoProviders := formula.NewRepoProviders()
			repoProviders.Add("Github", formula.Git{Repos: defaultGitRepositoryMock, NewRepoInfo: github.NewRepoInfo})

			newUpdateRepoPrompt := NewUpdateRepoCmd(server.Client(), tt.in.repo, repoProviders, inputTextMock{}, inputPasswordMock{}, inputURLMock{}, tt.in.inList, inputTrueMock{}, inputIntMock{})
			newUpdateRepoStdin := NewUpdateRepoCmd(server.Client(), tt.in.repo, repoProviders, inputTextMock{}, inputPasswordMock{}, inputURLMock{}, tt.in.inList, inputTrueMock{}, inputIntMock{})

			newUpdateRepoPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")
			newUpdateRepoStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

			newReader := strings.NewReader(tt.inputStdin)
			newUpdateRepoStdin.SetIn(newReader)

			if err := newUpdateRepoPrompt.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("new update repo type prompt command error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := newUpdateRepoStdin.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("new update repo type stdin command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func serverMock() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}
