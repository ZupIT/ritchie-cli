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

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func TestNewInitCmd(t *testing.T) {
	cmd := NewInitCmd(
		defaultRepoAdderMock,
		defaultGitRepositoryMock,
		TutorialFinderMock{},
		inputListCustomMock{
			list: func(name string, items []string) (string, error) {
				if name == AddCommonsQuestion {
					return "yes", nil
				}
				return AcceptMetrics, nil
			},
		},
		sMocks.FileWriteReadExisterCustomMock{
			WriteMock: func(path string, content []byte) error {
				return nil
			},
		},
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	if cmd == nil {
		t.Errorf("NewInitCmd got %v", cmd)
		return
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewInitStdin(t *testing.T) {
	cmd := NewInitCmd(
		defaultRepoAdderMock,
		defaultGitRepositoryMock,
		TutorialFinderMock{},
		inputListCustomMock{
			list: func(name string, items []string) (string, error) {
				if name == AddCommonsQuestion {
					return "yes", nil
				}
				return AcceptMetrics, nil
			},
		},
		sMocks.FileWriteReadExisterCustomMock{
			WriteMock: func(path string, content []byte) error {
				return nil
			},
		},
	)
	cmd.PersistentFlags().Bool("stdin", true, "input by stdin")

	input := "{\"addCommons\": \"yes\"}\n"
	newReader := strings.NewReader(input)
	cmd.SetIn(newReader)

	if cmd == nil {
		t.Errorf("NewInitCmd got %v", cmd)
		return
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func Test_initCmd_runAnyEntry(t *testing.T) {
	someError := errors.New("some error")
	type fields struct {
		repo      formula.RepositoryAdder
		git       git.Repositories
		inputList prompt.InputList
		find      rtutorial.Finder
	}

	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		inputStdin string
	}{
		{
			name: "Run With Success when add commons",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				find: TutorialFinderMock{},
				inputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == AddCommonsQuestion {
							return "yes", nil
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr:    false,
			inputStdin: "{\"addCommons\": \"yes\"}\n",
		},
		{
			name: "Run With Success when not add commons",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				inputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == AddCommonsQuestion {
							return "no", nil
						}
						return AcceptMetrics, nil
					},
				},
				find: TutorialFinderMock{},
			},
			wantErr:    false,
			inputStdin: "{\"addCommons\": \"no\"}\n",
		},
		{
			name: "Warning when call git.LatestTag",
			fields: fields{
				repo: defaultRepoAdderMock,
				git: GitRepositoryMock{
					latestTag: func(info git.RepoInfo) (git.Tag, error) {
						return git.Tag{}, someError
					},
				},
				inputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == AddCommonsQuestion {
							return "yes", nil
						}
						return AcceptMetrics, nil
					},
				},
				find: TutorialFinderMock{},
			},
			wantErr:    false,
			inputStdin: "{\"addCommons\": \"yes\"}\n",
		},
		{
			name: "Warning when call repo.Add",
			fields: fields{
				repo: repoListerAdderCustomMock{
					add: func(d formula.Repo) error {
						return someError
					},
				},
				git: defaultGitRepositoryMock,
				inputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == AddCommonsQuestion {
							return "yes", nil
						}
						return AcceptMetrics, nil
					},
				},
				find: TutorialFinderMock{},
			},
			wantErr:    false,
			inputStdin: "{\"addCommons\": \"yes\"}\n",
		},
		{
			name: "Error when find returns err",
			fields: fields{
				repo: repoListerAdderCustomMock{
					add: func(d formula.Repo) error {
						return errors.New("some error")
					},
				},
				git:  defaultGitRepositoryMock,
				find: TutorialFinderMock{},
				inputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == AddCommonsQuestion {
							return "yes", nil
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr:    false,
			inputStdin: "{\"addCommons\": \"yes\"}\n",
		},
		{
			name: "Error in select response of metrics",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				inputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == AddCommonsQuestion {
							return "yes", nil
						}
						return "any", someError
					},
				},
				find: TutorialFinderMock{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamMock := sMocks.FileWriteReadExisterCustomMock{
				WriteMock: func(path string, content []byte) error {
					return nil
				},
			}

			initPrompt := NewInitCmd(tt.fields.repo, tt.fields.git, tt.fields.find, tt.fields.inputList, streamMock)
			initStdin := NewInitCmd(tt.fields.repo, tt.fields.git, tt.fields.find, tt.fields.inputList, streamMock)

			initPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")
			initStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

			newReader := strings.NewReader(tt.inputStdin)
			initStdin.SetIn(newReader)

			if err := initPrompt.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("init_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
