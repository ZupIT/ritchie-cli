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
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDeleteRepo(t *testing.T) {
	tmpDir := os.TempDir()
	ritHomeName := ".rit"
	ritHome := filepath.Join(tmpDir, ritHomeName)
	reposPath := filepath.Join(ritHome, "repos")
	repoName := "repoName"
	repoPath := filepath.Join(reposPath, repoName)

	type in struct {
		args                  []string
		repoList              formula.Repos
		repoListErr           error
		inputListString       string
		inputListErr          error
		inputBool             bool
		existingRepoIsDeleted bool
		inputBoolErr          error
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success prompt",
			in: in{
				args: []string{},
				repoList: formula.Repos{
					{
						Name:     "repoName",
						Priority: 0,
					},
				},
				inputListString:       "repoName",
				inputBool:             true,
				existingRepoIsDeleted: true,
			},
		},
		{
			name: "success flag",
			in: in{
				args: []string{"--name=repoName"},
				repoList: formula.Repos{
					{
						Name:     "repoName",
						Priority: 0,
					},
				},
				existingRepoIsDeleted: true,
			},
		},
		{
			name: "error to list repos",
			in: in{
				args:                  []string{},
				repoListErr:           errors.New("error to list repos"),
				existingRepoIsDeleted: false,
			},
			want: errors.New("error to list repos"),
		},
		{
			name: "error to input list",
			in: in{
				args: []string{},
				repoList: formula.Repos{
					{
						Name:     "repoName",
						Priority: 0,
					},
				},
				inputListErr:          errors.New("error to input list"),
				existingRepoIsDeleted: false,
			},
			want: errors.New("error to input list"),
		},
		{
			name: "error to input bool",
			in: in{
				args: []string{},
				repoList: formula.Repos{
					{
						Name:     "repoName",
						Priority: 0,
					},
				},
				inputListString:       "repoName",
				inputBoolErr:          errors.New("error to input bool"),
				existingRepoIsDeleted: false,
			},
			want: errors.New("error to input bool"),
		},
		{
			name: "do not accept delete selected repo",
			in: in{
				args: []string{},
				repoList: formula.Repos{
					{
						Name:     "repoName",
						Priority: 0,
					},
				},
				inputListString:       "repoName",
				inputBool:             false,
				existingRepoIsDeleted: false,
			},
		},
		{
			name: "error on empty flag",
			in: in{
				args:                  []string{"--name="},
				existingRepoIsDeleted: false,
			},
			want: errors.New("please provide a value for 'name'"),
		},
		{
			name: "error to delete repo with wrong name",
			in: in{
				args: []string{"--name=wrongName"},
				repoList: formula.Repos{
					{
						Name:     "repoName",
						Priority: 0,
					},
				},
				existingRepoIsDeleted: false,
			},
			want: errors.New("no repository with this name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.MkdirAll(repoPath, os.ModePerm)
			defer os.RemoveAll(ritHome)
			fileManager := stream.NewFileManager()
			dirManager := stream.NewDirManager(fileManager)

			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputListString, tt.in.inputListErr)
			inputBoolMock := new(mocks.InputBoolMock)
			inputBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputBool, tt.in.inputBoolErr)
			repoManagerMock := new(mocks.RepoManager)
			repoManagerMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.repoList, tt.in.repoListErr)
			repoManagerMock.On("Write", mock.Anything).Return(nil)

			repoDeleter := repo.NewDeleter(ritHome, repoManagerMock, dirManager)

			cmd := NewDeleteRepoCmd(
				repoManagerMock,
				inputListMock,
				inputBoolMock,
				repoDeleter,
			)
			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			cmd.SetArgs(tt.in.args)

			got := cmd.Execute()
			if tt.in.existingRepoIsDeleted {
				assert.NoDirExists(t, repoPath)
			} else {
				assert.DirExists(t, repoPath)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
