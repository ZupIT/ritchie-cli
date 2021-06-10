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
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDeleteFormulaCmd(t *testing.T) {
	home := os.TempDir()
	ritHome := filepath.Join(home, ".rit")

	defer os.RemoveAll(home)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	treeGen := tree.NewGenerator(dirManager, fileManager)

	reposPath := filepath.Join(ritHome, "repos")
	repoPathLocalDefault := filepath.Join(reposPath, "local-default")
	repoPathWS := filepath.Join(home, "ritchie-formulas-local")

	workspaces := formula.Workspaces{}
	workspaces["Default"] = repoPathWS

	fileInfo := func(path string) (string, error) {
		fileManager := stream.NewFileManager()
		b, err := fileManager.Read(path)
		return string(b), err
	}

	tests := []struct {
		name           string
		workspaces     formula.Workspaces
		selectedWspace string
		group          string
		groupErr       error
		formula        string
		delete         bool
		stdin          string
		args           []string
		wspaceErr      error
		iListWSErr     error
		want           error
	}{
		{
			name:           "run with success when the execution type is prompt",
			workspaces:     workspaces,
			selectedWspace: "Default (" + repoPathWS + ")",
			group:          "testing",
			formula:        "rit testing delete-formula",
			delete:         true,
			want:           nil,
		},
		{
			name:       "run with success when the execution type is stdin",
			workspaces: workspaces,
			stdin:      `{"workspace_path":"/tmp/ritchie-formulas-local","formula":"rit testing delete-formula"}`,
			delete:     true,
			want:       nil,
		},
		{
			name:       "run with success when the execution type is flag",
			workspaces: workspaces,
			args:       []string{"--workspace=default", "--formula=rit testing delete-formula"},
			delete:     true,
			want:       nil,
		},
		{
			name:       "run with an error when the formula flag contains an incorrect value",
			workspaces: workspaces,
			args:       []string{"--workspace=default", "--formula=\"rit testing delete-formula\""},
			delete:     true,
			want:       errors.New("formula name is incorrect"),
		},
		{
			name:       "run with an error when empty formula flag",
			workspaces: workspaces,
			args:       []string{"--workspace=default", "--formula="},
			delete:     true,
			want:       errors.New("please provide a value for 'formula'"),
		},
		{
			name:       "run with an error when the workspace flag contains the nonexistent value",
			workspaces: workspaces,
			args:       []string{"--workspace=personal", "--formula=rit testing delete-formula"},
			delete:     true,
			want:       errors.New("no workspace found with this name"),
		},
		{
			name:       "run with an error when the workspace flag is empty",
			workspaces: workspaces,
			args:       []string{"--workspace=", "--formula=rit testing delete-formula"},
			delete:     true,
			want:       errors.New("please provide a value for 'workspace'"),
		},
		{
			name:      "run with error when workspace list returns err",
			args:      []string{"--workspace=default", "--formula=rit testing delete-formula"},
			wspaceErr: errors.New("workspace list error"),
			want:      errors.New("workspace list error"),
		},
		{
			name:           "run with error when readFormulas returns err",
			workspaces:     workspaces,
			selectedWspace: "Test (" + "/tmp/rit-delete-formula/" + ")",
			iListWSErr:     errors.New("no such file or directory"),
			want:           errors.New("no such file or directory"),
		},
		{
			name:           "run with error when question about select formula or group returns err",
			workspaces:     workspaces,
			selectedWspace: "Default (" + repoPathWS + ")",
			groupErr:       errors.New("group error"),
			want:           errors.New("group error"),
		},
		{
			name:      "run with error when add new workspace",
			wspaceErr: errors.New("workspace add error"),
			want:      errors.New("workspace add error"),
		},
		{
			name:           "run with success when choose not to delete formula",
			workspaces:     workspaces,
			selectedWspace: "Default (" + repoPathWS + ")",
			group:          "testing",
			formula:        "rit testing delete-formula",
			delete:         false,
			want:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(home)

			_ = dirManager.Remove(ritHome)

			createSaved := func(path string) {
				_ = dirManager.Remove(path)
				_ = dirManager.Create(path)
			}

			createSaved(repoPathLocalDefault)
			createSaved(repoPathWS)

			zipFile := filepath.Join("..", "..", "testdata", "ritchie-formulas-test.zip")
			zipRepositories := filepath.Join("..", "..", "testdata", "repositories.zip")
			zipTree := filepath.Join("..", "..", "testdata", "tree.zip")
			_ = streams.Unzip(zipRepositories, reposPath)
			_ = streams.Unzip(zipFile, repoPathLocalDefault)
			_ = streams.Unzip(zipFile, repoPathWS)
			_ = streams.Unzip(zipTree, repoPathLocalDefault)

			createTree(repoPathWS, repoPathLocalDefault, treeGen, fileManager)
			setWorkspace(workspaces, ritHome)

			workspaceMock := &mocks.WorkspaceMock{}
			workspaceMock.On("Add", mock.Anything).Return(tt.wspaceErr)
			workspaceMock.On("List").Return(tt.workspaces, tt.wspaceErr)

			pathWS := filepath.Join(repoPathWS)
			dir, _ := dirManager.List(pathWS, false)

			pathWithFormulas := filepath.Join(pathWS, "testing")
			formulasDir, _ := dirManager.List(pathWithFormulas, false)

			inList := new(mocks.InputListMock)
			inList.On("List", "Select a formula workspace: ", mock.Anything, mock.Anything).Return(tt.selectedWspace, tt.iListWSErr)
			inList.On("List", mock.Anything, dir, mock.Anything).Return(tt.group, tt.groupErr)
			inList.On("List", mock.Anything, formulasDir, mock.Anything).Return("delete-formula", tt.groupErr)
			inList.On("List", foundFormulaQuestion, mock.Anything, mock.Anything).Return(tt.formula, tt.groupErr)
			inBool := new(mocks.InputBoolMock)
			inBool.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(tt.delete, nil)

			inPath := &mocks.InputPathMock{}
			inText := &mocks.InputTextMock{}

			cmd := NewDeleteFormulaCmd(
				home,
				ritHome,
				workspaceMock,
				dirManager,
				inBool,
				inText,
				inList,
				inPath,
				treeGen,
				fileManager,
			)

			if tt.stdin != "" {
				newReader := strings.NewReader(tt.stdin)
				cmd.SetIn(newReader)
				cmd.PersistentFlags().Bool("stdin", true, "input by stdin")
			} else {
				cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			}

			cmd.SetArgs(tt.args)
			got := cmd.Execute()

			pathWSDir, pathLocalDir, treePath := "", "", ""
			pathWSDir = filepath.Join(repoPathWS, "testing", "delete-formula", "src")
			pathLocalDir = filepath.Join(repoPathLocalDefault, "testing", "delete-formula", "src")
			treePath = filepath.Join(repoPathLocalDefault, "tree.json")

			if tt.want != nil || !tt.delete {
				assert.Equal(t, tt.want, got)

				bTree, err := fileInfo(treePath)
				assert.Nil(t, err)
				tree, err := getTree([]byte(bTree))
				assert.Nil(t, err)

				assert.DirExists(t, pathWSDir)
				assert.DirExists(t, pathLocalDir)
				assert.NotEmpty(t, tree.Commands[api.CommandID("root_testing_delete-formula")])
			} else {
				assert.Nil(t, got)

				bTree, err := fileInfo(treePath)
				assert.Nil(t, err)
				tree, err := getTree([]byte(bTree))
				assert.Nil(t, err)

				assert.NoDirExists(t, pathWSDir)
				assert.NoDirExists(t, pathLocalDir)
				assert.Empty(t, tree.Commands[api.CommandID("root_testing_delete-formula")])
			}
		})
	}
}
