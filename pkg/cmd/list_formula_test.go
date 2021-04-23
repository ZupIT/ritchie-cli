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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewListFormula(t *testing.T) {
	finderTutorial := TutorialFinderMock{}
	tmpDir := os.TempDir()
	ritHomeName := ".rit"
	ritHome := filepath.Join(tmpDir, ritHomeName)
	reposPath := filepath.Join(ritHome, "repos")

	repos := formula.Repos{
		{
			Name: "repoName",
		},
		{
			Name: "repoOtherName",
		},
	}

	emptyTree := `{
		"version": "v2",
		"commands": {}
	}
	`
	expectedOut :=
		`COMMAND                      	DESCRIPTION               
rit http generate http-config	Creates http-load template`

	type in struct {
		args            []string
		repoList        formula.Repos
		repoListErr     error
		inputListString string
		inputListErr    error
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success prompt",
			in: in{
				args:            []string{},
				repoList:        repos,
				inputListString: "repoName",
			},
		},
		{
			name: "success prompt option all",
			in: in{
				args:            []string{},
				repoList:        repos,
				inputListString: "ALL",
			},
		},
		{
			name: "success flag",
			in: in{
				args:     []string{"--name=repoName"},
				repoList: repos,
			},
		},
		{
			name: "success flag option all",
			in: in{
				args:     []string{"--name=ALL"},
				repoList: repos,
			},
		},
		{
			name: "error to list repos",
			in: in{
				args:        []string{},
				repoListErr: errors.New("error to list repos"),
			},
			want: errors.New("error to list repos"),
		},
		{
			name: "error to input list",
			in: in{
				args:         []string{},
				repoList:     repos,
				inputListErr: errors.New("error to input list"),
			},
			want: errors.New("error to input list"),
		},
		{
			name: "error on empty flag",
			in: in{
				args: []string{"--name="},
			},
			want: errors.New("please provide a value for 'name'"),
		},
		{
			name: "error to list formulas from repo with wrong name",
			in: in{
				args:     []string{"--name=wrongName"},
				repoList: repos,
			},
			want: errors.New("no repository with this name"),
		},
		{
			name: "error tree with no commands",
			in: in{
				args:     []string{"--name=repoOtherName"},
				repoList: repos,
			},
			want: errors.New("no formula found in selected repo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputListString, tt.in.inputListErr)
			repoManagerMock := new(mocks.RepoManager)
			repoManagerMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.repoList, tt.in.repoListErr)
			repoManagerMock.On("Write", mock.Anything).Return(nil)

			for _, r := range tt.in.repoList {
				repoName := r.Name.String()
				repoPath := filepath.Join(reposPath, repoName)
				_ = os.MkdirAll(repoPath, os.ModePerm)
				_ = streams.Unzip("../../testdata/tree.zip", repoPath)
				if tt.want != nil {
					emptyTreeData := []byte(emptyTree)
					_ = ioutil.WriteFile(filepath.Join(repoPath, "tree.json"), emptyTreeData, 0666)
				}
			}
			defer os.RemoveAll(ritHome)
			treeManager := tree.NewTreeManager(ritHome, repoManagerMock, api.CoreCmds)

			cmd := NewListFormulaCmd(
				repoManagerMock,
				inputListMock,
				treeManager,
				finderTutorial,
			)
			cmd.SetArgs(tt.in.args)

			rescueStdout := os.Stdout
			r, w, err := os.Pipe()
			assert.NoError(t, err)
			os.Stdout = w

			got := cmd.Execute()
			assert.Equal(t, tt.want, got)

			_ = w.Close()
			out, err := ioutil.ReadAll(r)
			assert.NoError(t, err)
			capturedOut := string(out)
			os.Stdout = rescueStdout
			if tt.want == nil {
				assert.Contains(t, capturedOut, expectedOut)
			} else {
				assert.NotContains(t, capturedOut, expectedOut)
			}
		})
	}
}
