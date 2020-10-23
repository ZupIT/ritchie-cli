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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type fieldsTestDeleteRepoCmd struct {
	repositoryLister      formula.RepositoryLister
	repositoryListerLocal formula.RepositoryListerLocal
	inputList             prompt.InputList
	setup                 testBeforeSetupMock
}

func TestNewDeleteRepoCmd(t *testing.T) {
	someError := errors.New("some error")

	tests := []struct {
		name    string
		fields  fieldsTestDeleteRepoCmd
		wantErr bool
	}{
		{
			name: "Run with empty repository",
			fields: fieldsTestDeleteRepoCmd{
				repositoryLister: repoListerMock{},
				repositoryListerLocal: repoListerLocalCustomMock{
					list: func() (formula.RepoName, error) {
						return "", someError
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Run delete local repository success",
			fields: fieldsTestDeleteRepoCmd{
				repositoryLister:      repoListerMock{},
				repositoryListerLocal: repoListerLocalMock{},
				setup: testBeforeSetupMock{
					setup: func() {
						setupMockLocalFolders()
					},
				},
				inputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "local", nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Run repository success",
			fields: fieldsTestDeleteRepoCmd{
				repositoryLister:      repoListerMock{},
				repositoryListerLocal: repoListerLocalMock{},
				setup: testBeforeSetupMock{
					setup: func() {
						setupMockFolders("someRepo")
					},
				},
				inputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "someRepo", nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Try delete unexistent repo",
			fields: fieldsTestDeleteRepoCmd{
				repositoryLister: repoListerCustomMock{
					list: func() (formula.Repos, error) {
						return nil, someError
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDeleteRepoCmd(
				tt.fields.repositoryLister,
				tt.fields.repositoryListerLocal,
				inputListMock{},
				repositoryDeleteMock{},
				repositoryDeleteLocalMock{},
			)
			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

			if cmd == nil {
				t.Errorf("TestNewDeleteRepoCmd got %v", cmd)
				return
			}

			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", cmd.Use, err, tt.wantErr)
			}
		})
	}
}

func setupMockFolders(repoName string) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	ritHomePath := filepath.Join(os.TempDir(), "TestDeleteRepo")
	repoPath := filepath.Join(ritHomePath, "repos", repoName)

	repoJSON := formula.Repos{
		{
			Name: formula.RepoName(repoName),
		},
	}

	_ = dirManager.Remove(ritHomePath)
	_ = dirManager.Create(ritHomePath)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)

	repoData, _ := json.Marshal(repoJSON)
	reposFileName := "repositories.json"
	repoFilePath := filepath.Join(ritHomePath, "repos", reposFileName)
	_ = fileManager.Remove(repoFilePath)
	_ = fileManager.Write(repoFilePath, repoData)
}

func setupMockLocalFolders() {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	ritHomePath := filepath.Join(os.TempDir(), "TestDeleteRepo")
	repoPath := filepath.Join(ritHomePath, "repos", "local")

	_ = dirManager.Remove(ritHomePath)
	_ = dirManager.Create(ritHomePath)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
}
