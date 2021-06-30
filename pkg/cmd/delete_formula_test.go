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
	home := filepath.Join(os.TempDir(), "rit-delete-formula")
	ritHome := filepath.Join(home, ".rit")

	defer os.RemoveAll(home)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	treeGen := tree.NewGenerator(dirManager, fileManager)

	reposPath := filepath.Join(ritHome, "repos")
	repoPathLocalDefault := filepath.Join(reposPath, "local-default")
	repoPathLocalEmpty := filepath.Join(reposPath, "local-empty")
	repoPathWS := filepath.Join(home, "ritchie-formulas-local")
	repoPathWSEmpty := filepath.Join(home, "ritchie-formulas-empty")

	workspaces := formula.Workspaces{}
	workspaces["Default"] = repoPathWS
	workspaces["Empty"] = repoPathWSEmpty

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
		subgroup       string
		formula        string
		treeCmd        string
		wspacePath     string
		repoPath       string
		stdin          string
		args           []string
		delete         bool
		nested         bool
		groupErr       error
		wspaceErr      error
		iListWSErr     error
		iAddWSErr      error
		iBoolErr       error
		want           error
	}{
		{
			name:           "run with success when the execution type is prompt",
			workspaces:     workspaces,
			selectedWspace: "Default (" + repoPathWS + ")",
			group:          "testing",
			subgroup:       "delete-formula",
			formula:        "rit testing delete-formula",
			treeCmd:        "root_testing_delete-formula",
			wspacePath:     repoPathWS,
			repoPath:       repoPathLocalDefault,
			delete:         true,
			want:           nil,
		},
		{
			name:           "run with success when the execution type is prompt and nested formula",
			workspaces:     workspaces,
			selectedWspace: "Default (" + repoPathWS + ")",
			group:          "testing",
			subgroup:       "nested-formula",
			formula:        "rit testing nested-formula",
			treeCmd:        "root_testing_nested-formula",
			wspacePath:     repoPathWS,
			repoPath:       repoPathLocalDefault,
			delete:         true,
			nested:         true,
			want:           nil,
		},
		{
			name:       "run with success when the execution type is stdin",
			workspaces: workspaces,
			stdin:      `{"workspace_path":"/tmp/rit-delete-formula/ritchie-formulas-local","formula":"rit testing delete-formula"}`,
			treeCmd:    "root_testing_delete-formula",
			wspacePath: repoPathWS,
			repoPath:   repoPathLocalDefault,
			delete:     true,
			want:       nil,
		},
		{
			name:       "run with success when the execution type is flag",
			workspaces: workspaces,
			args:       []string{"--workspace=default", "--formula=rit testing delete-formula"},
			treeCmd:    "root_testing_delete-formula",
			wspacePath: repoPathWS,
			repoPath:   repoPathLocalDefault,
			delete:     true,
			want:       nil,
		},
		{
			name:       "run with an error when the formula flag contains an incorrect value",
			workspaces: workspaces,
			args:       []string{"--workspace=default", "--formula=\"rit testing delete-formula\""},
			treeCmd:    "root_testing_delete-formula",
			wspacePath: repoPathWS,
			repoPath:   repoPathLocalDefault,
			delete:     true,
			want:       errors.New("formula name is incorrect"),
		},
		{
			name:       "run with an error when empty formula flag",
			workspaces: workspaces,
			args:       []string{"--workspace=default", "--formula="},
			treeCmd:    "root_testing_delete-formula",
			wspacePath: repoPathWS,
			repoPath:   repoPathLocalDefault,
			delete:     true,
			want:       errors.New("please provide a value for 'formula'"),
		},
		{
			name:       "run with an error when the workspace flag contains the nonexistent value",
			workspaces: workspaces,
			args:       []string{"--workspace=personal", "--formula=rit testing delete-formula"},
			treeCmd:    "root_testing_delete-formula",
			wspacePath: repoPathWS,
			repoPath:   repoPathLocalDefault,
			delete:     true,
			want:       errors.New("no workspace found with this name"),
		},
		{
			name:       "run with an error when the workspace flag is empty",
			workspaces: workspaces,
			args:       []string{"--workspace=", "--formula=rit testing delete-formula"},
			treeCmd:    "root_testing_delete-formula",
			wspacePath: repoPathWS,
			repoPath:   repoPathLocalDefault,
			delete:     true,
			want:       errors.New("please provide a value for 'workspace'"),
		},
		{
			name:       "run with error when workspace list returns err",
			args:       []string{"--workspace=default", "--formula=rit testing delete-formula"},
			treeCmd:    "root_testing_delete-formula",
			wspacePath: repoPathWS,
			repoPath:   repoPathLocalDefault,
			wspaceErr:  errors.New("workspace list error"),
			want:       errors.New("workspace list error"),
		},
		{
			name:           "run with error when readFormulas returns err",
			workspaces:     workspaces,
			selectedWspace: "Test (" + "/tmp/rit-delete-formula/" + ")",
			treeCmd:        "root_testing_delete-formula",
			wspacePath:     repoPathWS,
			repoPath:       repoPathLocalDefault,
			iListWSErr:     errors.New("no such file or directory"),
			want:           errors.New("no such file or directory"),
		},
		{
			name:           "run with error when question about select formula or group returns err",
			workspaces:     workspaces,
			selectedWspace: "Default (" + repoPathWS + ")",
			treeCmd:        "root_testing_delete-formula",
			wspacePath:     repoPathWS,
			repoPath:       repoPathLocalDefault,
			groupErr:       errors.New("group error"),
			want:           errors.New("group error"),
		},
		{
			name:       "run with error when add new workspace",
			treeCmd:    "root_testing_delete-formula",
			wspacePath: repoPathWS,
			repoPath:   repoPathLocalDefault,
			wspaceErr:  errors.New("workspace add error"),
			want:       errors.New("workspace add error"),
		},
		{
			name:           "run with success when choose not to delete formula",
			workspaces:     workspaces,
			selectedWspace: "Default (" + repoPathWS + ")",
			group:          "testing",
			subgroup:       "delete-formula",
			treeCmd:        "root_testing_delete-formula",
			formula:        "rit testing delete-formula",
			wspacePath:     repoPathWS,
			repoPath:       repoPathLocalDefault,
			delete:         false,
			want:           nil,
		},
		{
			name:           "run with error when workspace is empty",
			workspaces:     workspaces,
			selectedWspace: "Empty (" + repoPathWSEmpty + ")",
			treeCmd:        "",
			wspacePath:     repoPathWSEmpty,
			repoPath:       repoPathLocalEmpty,
			want:           errors.New("could not find formula"),
		},
		{
			name:       "run add error",
			treeCmd:    "root_testing_delete-formula",
			wspacePath: repoPathWS,
			repoPath:   repoPathLocalDefault,
			iAddWSErr:  errors.New("workspace add error"),
			want:       errors.New("workspace add error"),
		},
		{
			name:           "run with error when invalid input bool",
			workspaces:     workspaces,
			selectedWspace: "Default (" + repoPathWS + ")",
			group:          "testing",
			subgroup:       "delete-formula",
			treeCmd:        "root_testing_delete-formula",
			formula:        "rit testing delete-formula",
			wspacePath:     repoPathWS,
			repoPath:       repoPathLocalDefault,
			delete:         false,
			iBoolErr:       errors.New("invalid input bool"),
			want:           errors.New("invalid input bool"),
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

			createSaved(tt.repoPath)
			createSaved(tt.wspacePath)

			zipFile := filepath.Join("..", "..", "testdata", "ritchie-formulas-test.zip")
			emptyZipFile := filepath.Join("..", "..", "testdata", "ritchie-delete-test.zip")
			zipRepositories := filepath.Join("..", "..", "testdata", "repositories.zip")
			zipTree := filepath.Join("..", "..", "testdata", "tree.zip")
			_ = streams.Unzip(zipRepositories, reposPath)
			_ = streams.Unzip(zipFile, repoPathLocalDefault)
			_ = streams.Unzip(zipFile, repoPathWS)
			_ = streams.Unzip(emptyZipFile, repoPathWSEmpty)
			_ = streams.Unzip(zipTree, repoPathLocalEmpty)

			createTree(tt.wspacePath, tt.repoPath, treeGen, fileManager)
			setWorkspace(workspaces, ritHome)

			workspaceMock := &mocks.WorkspaceMock{}
			workspaceMock.On("Add", mock.Anything).Return(tt.iAddWSErr)
			workspaceMock.On("List").Return(tt.workspaces, tt.wspaceErr)

			pathWS := filepath.Join(tt.wspacePath)
			dir, _ := dirManager.List(pathWS, false)

			pathWithFormulas := filepath.Join(pathWS, "testing")
			formulasDir, _ := dirManager.List(pathWithFormulas, false)

			question := fmt.Sprintf("Are you sure you want to delete the formula: %s", tt.formula)

			inList := new(mocks.InputListMock)
			inList.On("List", "Select a formula workspace: ", mock.Anything, mock.Anything).Return(tt.selectedWspace, tt.iListWSErr)
			inList.On("List", mock.Anything, dir, mock.Anything).Return(tt.group, tt.groupErr)
			inList.On("List", mock.Anything, formulasDir, mock.Anything).Return(tt.subgroup, tt.groupErr)
			inList.On("List", foundFormulaQuestion, mock.Anything, mock.Anything).Return(tt.formula, tt.groupErr)
			inList.On("List", mock.Anything, []string{"test", "test"}, mock.Anything).Return("test", tt.groupErr)
			inBool := new(mocks.InputBoolMock)
			inBool.On("Bool", question, mock.Anything, mock.Anything).Return(tt.delete, tt.iBoolErr)

			inPath := &mocks.InputPathMock{}
			inTextValidator := &mocks.InputTextValidatorMock{}

			cmd := NewDeleteFormulaCmd(
				home,
				ritHome,
				workspaceMock,
				dirManager,
				inBool,
				inTextValidator,
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
			if tt.formula == "rit testing nested-formula" {
				pathWSDir = filepath.Join(repoPathWS, "testing", "nested-formula", "test")
				pathLocalDir = filepath.Join(repoPathLocalDefault, "testing", "nested-formula", "test")
				treePath = filepath.Join(repoPathLocalDefault, "tree.json")
			} else if tt.wspacePath == repoPathWSEmpty {
				pathWSDir = filepath.Join(tt.wspacePath, "src")
				pathLocalDir = filepath.Join(tt.repoPath)
				treePath = filepath.Join(tt.repoPath, "tree.json")
			} else {
				pathWSDir = filepath.Join(repoPathWS, "testing", "delete-formula", "src")
				pathLocalDir = filepath.Join(repoPathLocalDefault, "testing", "delete-formula", "src")
				treePath = filepath.Join(repoPathLocalDefault, "tree.json")
			}

			if tt.want != nil || !tt.delete || tt.nested {
				assert.Equal(t, tt.want, got)

				bTree, err := fileInfo(treePath)
				assert.Nil(t, err)
				tree, err := getTree([]byte(bTree))
				assert.Nil(t, err)

				assert.DirExists(t, pathWSDir)
				assert.DirExists(t, pathLocalDir)

				if tt.treeCmd == "" {
					assert.Empty(t, tree.Commands[api.CommandID(tt.treeCmd)])
				} else {
					assert.NotEmpty(t, tree.Commands[api.CommandID(tt.treeCmd)])
				}
			} else {
				assert.Nil(t, got)

				bTree, err := fileInfo(treePath)
				assert.Nil(t, err)
				tree, err := getTree([]byte(bTree))
				assert.Nil(t, err)

				assert.NoDirExists(t, pathWSDir)
				assert.NoDirExists(t, pathLocalDir)
				assert.Empty(t, tree.Commands[api.CommandID(tt.treeCmd)])
			}
		})
	}
}
