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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestRenameFormulaCmd(t *testing.T) {
	tmp := os.TempDir()
	home := filepath.Join(tmp, "rit_test-renameFormula")
	ritHome := filepath.Join(home, ".rit")

	defer os.RemoveAll(home)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	githubRepo := github.NewRepoManager(http.DefaultClient)
	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: githubRepo, NewRepoInfo: github.NewRepoInfo})

	repoCreator := repo.NewCreator(ritHome, repoProviders, dirManager, fileManager)
	repoLister := repo.NewLister(ritHome, fileManager)
	repoWriter := repo.NewWriter(ritHome, fileManager)
	repoDetail := repo.NewDetail(repoProviders)
	repoListWriter := repo.NewListWriter(repoLister, repoWriter)
	repoDeleter := repo.NewDeleter(ritHome, repoListWriter, dirManager)
	repoListWriteCreator := repo.NewCreateWriteListDetailDeleter(repoLister, repoCreator, repoWriter, repoDetail, repoDeleter)

	treeGen := tree.NewGenerator(dirManager, fileManager)
	repoAdder := repo.NewAdder(ritHome, repoListWriteCreator, treeGen)
	formBuildLocal := builder.NewBuildLocal(ritHome, dirManager, repoAdder)

	formulaWorkspace := workspace.New(ritHome, home, dirManager, formBuildLocal)

	reposPath := filepath.Join(ritHome, "repos")
	repoPath := filepath.Join(reposPath, "commons")
	_ = dirManager.Remove(ritHome)

	createSaved := func(path string) {
		_ = dirManager.Remove(path)
		_ = dirManager.Create(path)
	}
	createSaved(repoPath)

	zipFile := filepath.Join("..", "..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	zipRepositories := filepath.Join("..", "..", "..", "..", "testdata", "repositories.zip")
	_ = streams.Unzip(zipFile, repoPath)
	_ = streams.Unzip(zipRepositories, reposPath)

	type in struct {
		inputText    string
		inputTextErr error
		inputTextVal string
		inputList    string
		inputListErr error
		args         []string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				inputTextVal: "rit test test",
				inputList:    "default",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputTextMock := new(mocks.InputTextMock)
			inputTextMock.On("Text", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputText, tt.in.inputTextErr)

			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputList, tt.in.inputListErr)

			inPath := &mocks.InputPathMock{}
			inPath.On("Read", "Workspace path (e.g.: /home/user/github): ").Return("", nil)

			cmd := NewRenameFormulaCmd(
				formulaWorkspace,
				inputTextMock,
				inputListMock,
				inPath,
				dirManager,
			)
			// TODO: remove stdin flag after  deprecation
			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			cmd.SetArgs(tt.in.args)

			got := cmd.Execute()

			assert.Equal(t, tt.want, got)
		})
	}

}
