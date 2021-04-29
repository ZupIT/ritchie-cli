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
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/deleter"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/formula/validator"
	"github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestRenameFormulaCmd(t *testing.T) {
	home := filepath.Join(os.TempDir(), "rit_test-renameFormula")
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
	repoPathLocalDefault := filepath.Join(reposPath, "local-default")
	repoPathWS := filepath.Join(home, "ritchie-formulas-local")

	repoListDetailWriter := repo.NewListDetailWrite(repoLister, repoDetail, repoWriter)
	treeManager := tree.NewTreeManager(ritHome, repoListDetailWriter, api.CoreCmds)
	tplManager := template.NewManager(api.RitchieHomeDir(), dirManager)
	formulaCreator := creator.NewCreator(treeManager, dirManager, fileManager, tplManager)
	createBuilder := formula.NewCreateBuilder(formulaCreator, formBuildLocal)

	validator := validator.NewValidator()
	deleter := deleter.NewDeleter(dirManager, fileManager, treeGen, ritHome)

	fileInfo := func(path string) (string, error) {
		fileManager := stream.NewFileManager()
		b, err := fileManager.Read(path)
		return string(b), err
	}

	type in struct {
		inputOldFormula   string
		inputNewFormula   string
		workspaceSelected string
		args              []string
	}

	type out struct {
		formulaToBeEmpty    string
		formulaToBeCreated  string
		formulaPathExpected string
		want                error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "success on prompt input",
			in: in{
				inputOldFormula:   "rit testing formula",
				inputNewFormula:   "rit testing new-formula",
				workspaceSelected: "Default (" + repoPathWS + ")",
			},
			out: out{
				formulaPathExpected: filepath.Join("testing", "new-formula"),
				formulaToBeCreated:  "root_testing_new-formula",
				formulaToBeEmpty:    "root_testing_formula",
			},
		},
		{
			name: "error on prompt input when oldName flag is nil",
			in: in{
				inputOldFormula:   "",
				inputNewFormula:   "rit testing new-formula",
				workspaceSelected: "Default (" + repoPathWS + ")",
			},
			out: out{
				want: errors.New("this input must not be empty"),
			},
		},
		{
			name: "error on prompt input when oldName flag dont exists in workspace",
			in: in{
				inputOldFormula:   "rit testing other",
				inputNewFormula:   "rit testing new-formula",
				workspaceSelected: "Default (" + repoPathWS + ")",
			},
			out: out{
				want: errors.New("This formula 'rit testing other' dont's exists on this workspace = 'Default'"),
			},
		},
		{
			name: "error on prompt input when newName flag exists in workspace",
			in: in{
				inputOldFormula:   "rit testing formula",
				inputNewFormula:   "rit testing formula",
				workspaceSelected: "Default (" + repoPathWS + ")",
			},
			out: out{
				want: errors.New("This formula 'rit testing formula' already exists on this workspace = 'Default'"),
			},
		},
		{
			name: "success on flag input",
			in: in{
				args: []string{
					"--workspace=Default",
					"--oldNameFormula=rit testing formula",
					"--newNameFormula=rit testing new-formula",
				},
			},
			out: out{
				formulaPathExpected: filepath.Join("testing", "new-formula"),
				formulaToBeCreated:  "root_testing_new-formula",
				formulaToBeEmpty:    "root_testing_formula",
			},
		},
		{
			name: "error on flag input when workspace flag is nil",
			in: in{
				args: []string{
					"--workspace=",
					"--oldNameFormula=rit testing formula",
					"--newNameFormula=rit testing formula new",
				},
			},
			out: out{
				want: errors.New("please provide a value for 'workspace'"),
			},
		},
		{
			name: "error on flag input when oldName flag is nil",
			in: in{
				args: []string{
					"--workspace=Default",
					"--oldNameFormula=",
					"--newNameFormula=rit testing formula new",
				},
			},
			out: out{
				want: errors.New("please provide a value for 'oldNameFormula'"),
			},
		},
		{
			name: "error on flag input when newNameFormula flag is nil",
			in: in{
				args: []string{
					"--workspace=Default",
					"--oldNameFormula=rit testing formula",
					"--newNameFormula=",
				},
			},
			out: out{
				want: errors.New("please provide a value for 'newNameFormula'"),
			},
		},
		{
			name: "error on flag input when workspace flag dont exists",
			in: in{
				args: []string{
					"--workspace=other",
					"--oldNameFormula=rit testing formula",
					"--newNameFormula=rit testing formula new",
				},
			},
			out: out{
				want: errors.New("The formula workspace 'other' does not exist, please enter a valid workspace"),
			},
		},
		{
			name: "error on flag input when oldName flag dont exists in workspace",
			in: in{
				args: []string{
					"--workspace=Default",
					"--oldNameFormula=rit testing other",
					"--newNameFormula=rit testing formula new",
				},
			},
			out: out{
				want: errors.New("This formula 'rit testing other' dont's exists on this workspace = 'Default'"),
			},
		},
		{
			name: "error on flag input when newName flag exists in workspace",
			in: in{
				args: []string{
					"--workspace=Default",
					"--oldNameFormula=rit testing formula",
					"--newNameFormula=rit testing formula",
				},
			},
			out: out{
				want: errors.New("This formula 'rit testing formula' already exists on this workspace = 'Default'"),
			},
		},
		{
			name: "error on flag input when old formula dont exists in workspace",
			in: in{
				args: []string{
					"--workspace=Default",
					"--oldNameFormula=rit other formula",
					"--newNameFormula=rit testing formula new",
				},
			},
			out: out{
				want: errors.New("This formula 'rit other formula' dont's exists on this workspace = 'Default'"),
			},
		},
		{
			name: "success when new formula is added a higher level of the tree",
			in: in{
				inputOldFormula:   "rit testing formula",
				inputNewFormula:   "rit testing formula new",
				workspaceSelected: "Default (" + repoPathWS + ")",
			},
			out: out{
				formulaPathExpected: filepath.Join("testing", "formula", "new"),
				formulaToBeCreated:  "root_testing_formula_new",
				formulaToBeEmpty:    "root_testing_formula",
			},
		},
		{
			name: "success when new formula is added a lower level of the tree",
			in: in{
				inputOldFormula:   "rit testing withOneMoreLevel level",
				inputNewFormula:   "rit testing leaft",
				workspaceSelected: "Default (" + repoPathWS + ")",
			},
			out: out{
				formulaPathExpected: filepath.Join("testing", "leaf"),
				formulaToBeCreated:  "root_testing_leaf",
				formulaToBeEmpty:    "root_testing_withOneMoreLevel_level",
			},
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
			_ = streams.Unzip(zipTree, repoPathLocalDefault)
			_ = streams.Unzip(zipFile, repoPathWS)

			createTree(ritHome, repoPathWS, treeGen, fileManager)

			inputTextMock := new(mocks.InputTextMock)

			inputTextValidatorMock := new(mocks.InputTextValidatorMock)
			inputTextValidatorMock.On("Text", formulaOldCmdLabel, mock.Anything, mock.Anything).Return(
				tt.in.inputOldFormula, nil,
			)
			inputTextValidatorMock.On("Text", formulaNewCmdLabel, mock.Anything, mock.Anything).Return(
				tt.in.inputNewFormula, nil,
			)

			inPath := &mocks.InputPathMock{}
			inPath.On("Read", "Workspace path (e.g.: /home/user/github): ").Return("", nil)

			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", "Select a formula workspace: ", mock.Anything, mock.Anything).Return(
				tt.in.workspaceSelected, nil,
			)

			inputBoolMock := new(mocks.InputBoolMock)
			inputBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

			cmd := NewRenameFormulaCmd(formulaWorkspace, inputTextMock, inputListMock, inPath, inputTextValidatorMock,
				inputBoolMock, dirManager, validator, createBuilder, treeGen, deleter, home, ritHome, fileManager)

			cmd.SetArgs(tt.in.args)

			got := cmd.Execute()

			if tt.out.want == nil {
				pathWSDir := filepath.Join(repoPathWS, tt.out.formulaPathExpected, "src")
				pathLocalDir := filepath.Join(repoPathLocalDefault, tt.out.formulaPathExpected, "src")
				treePath := filepath.Join(repoPathLocalDefault, "tree.json")

				bTree, err := fileInfo(treePath)
				assert.Nil(t, err)
				tree, err := getTree([]byte(bTree))
				assert.Nil(t, err)

				assert.Nil(t, got)

				assert.DirExists(t, pathWSDir)
				assert.DirExists(t, pathLocalDir)

				assert.True(t, tree.Commands[api.CommandID(tt.out.formulaToBeCreated)].Formula)
				assert.Empty(t, tree.Commands[api.CommandID(tt.out.formulaToBeEmpty)])
			} else {
				assert.Equal(t, tt.out.want.Error(), got.Error())
			}
		})
	}

}

func createTree(ritHome, ws string, tg formula.TreeGenerator, fm stream.FileWriteRemover) {
	localTree, _ := tg.Generate(ws)

	jsonString, _ := json.MarshalIndent(localTree, "", "\t")
	pathLocalTreeJSON := filepath.Join(ritHome, "repos", "local-default", "tree.json")
	_ = fm.Write(pathLocalTreeJSON, jsonString)
}

func getTree(f []byte) (formula.Tree, error) {
	tree := formula.Tree{}
	if err := json.Unmarshal(f, &tree); err != nil {
		return formula.Tree{}, err
	}
	return tree, nil
}
