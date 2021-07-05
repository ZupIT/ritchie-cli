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
	"io/ioutil"
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
	repoListWriteCreator := repo.NewCreateWriteListDetailDeleter(
		repoLister, repoCreator, repoWriter, repoDetail, repoDeleter,
	)

	treeGen := tree.NewGenerator(dirManager, fileManager)
	repoAdder := repo.NewAdder(ritHome, repoListWriteCreator, treeGen)
	formBuildLocal := builder.NewBuildLocal(ritHome, dirManager, repoAdder)

	formulaWorkspace := workspace.New(ritHome, home, dirManager, formBuildLocal, treeGen)

	reposPath := filepath.Join(ritHome, "repos")
	repoPathLocalDefault := filepath.Join(reposPath, "local-default")
	repoPathWS := filepath.Join(home, "ritchie-formulas-local")
	repoPathWSCustom := filepath.Join(home, "custom")
	repoPathLocalCustom := filepath.Join(reposPath, "local-custom")

	repoListDetailWriter := repo.NewListDetailWrite(repoLister, repoDetail, repoWriter)
	treeManager := tree.NewTreeManager(ritHome, repoListDetailWriter, api.CoreCmds)
	tplManager := template.NewManager(api.RitchieHomeDir(), dirManager)
	formulaCreator := creator.NewCreator(treeManager, dirManager, fileManager, tplManager)
	createBuilder := formula.NewCreateBuilder(formulaCreator, formBuildLocal)

	validator := validator.New()
	deleter := deleter.NewDeleter(dirManager, fileManager, treeGen, ritHome)

	fileInfo := func(path string) (string, error) {
		fileManager := stream.NewFileManager()
		b, err := fileManager.Read(path)
		return string(b), err
	}

	type in struct {
		inputOldFormula         string
		inputNewFormula         string
		workspaceSelected       string
		customWorkspaceSelected bool
		args                    []string
		approveOnConfirmation   bool
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
				inputOldFormula:       "rit testing formula",
				inputNewFormula:       "rit testing new-formula",
				workspaceSelected:     "Default (" + repoPathWS + ")",
				approveOnConfirmation: true,
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
				inputOldFormula:       "",
				inputNewFormula:       "rit testing new-formula",
				workspaceSelected:     "Default (" + repoPathWS + ")",
				approveOnConfirmation: true,
			},
			out: out{
				want: errors.New("this input must not be empty"),
			},
		},
		{
			name: "error on prompt input when oldName flag dont exists in workspace",
			in: in{
				inputOldFormula:       "rit testing other",
				inputNewFormula:       "rit testing new-formula",
				workspaceSelected:     "Default (" + repoPathWS + ")",
				approveOnConfirmation: true,
			},
			out: out{
				want: errors.New("formula \"rit testing other\" wasn't found in the workspaces"),
			},
		},
		{
			name: "error on prompt input when newName flag exists in workspace",
			in: in{
				inputOldFormula:       "rit testing formula",
				inputNewFormula:       "rit testing formula",
				workspaceSelected:     "Default (" + repoPathWS + ")",
				approveOnConfirmation: true,
			},
			out: out{
				want: errors.New("formula \"rit testing formula\" already exists on this workspace = \"Default\""),
			},
		},
		{
			name: "no confirmation on prompt",
			in: in{
				inputOldFormula:       "rit testing formula",
				inputNewFormula:       "rit testing formula new",
				approveOnConfirmation: false,
			},
		},
		{
			name: "success on flag input",
			in: in{
				args: []string{
					"--oldName=rit testing formula",
					"--newName=rit testing new-formula",
				},
			},
			out: out{
				formulaPathExpected: filepath.Join("testing", "new-formula"),
				formulaToBeCreated:  "root_testing_new-formula",
				formulaToBeEmpty:    "root_testing_formula",
			},
		},
		{
			name: "success on flag input when workspace flag is nil",
			in: in{
				args: []string{
					"--oldName=rit testing formula",
					"--newName=rit testing other",
				},
			},
			out: out{
				formulaPathExpected: filepath.Join("testing", "other"),
				formulaToBeCreated:  "root_testing_other",
				formulaToBeEmpty:    "root_testing_formula",
			},
		},
		{
			name: "error on flag input when oldName flag is nil",
			in: in{
				args: []string{
					"--oldName=",
					"--newName=rit testing formula new",
				},
			},
			out: out{
				want: errors.New("please provide a value for 'oldName'"),
			},
		},
		{
			name: "error on flag input when newNameFormula flag is nil",
			in: in{
				args: []string{
					"--oldName=rit testing formula",
					"--newName=",
				},
			},
			out: out{
				want: errors.New("please provide a value for 'newName'"),
			},
		},
		{
			name: "error on flag input when oldName flag dont exists in workspace",
			in: in{
				args: []string{
					"--oldName=rit testing other",
					"--newName=rit testing formula new",
				},
			},
			out: out{
				want: errors.New("formula \"rit testing other\" wasn't found in the workspaces"),
			},
		},
		{
			name: "error on flag input when newName flag exists in workspace",
			in: in{
				args: []string{
					"--oldName=rit testing formula",
					"--newName=rit testing formula",
				},
			},
			out: out{
				want: errors.New("formula \"rit testing formula\" already exists on this workspace = \"Default\""),
			},
		},
		{
			name: "success when new formula is added a higher level of the tree",
			in: in{
				inputOldFormula:       "rit testing formula",
				inputNewFormula:       "rit testing formula new",
				approveOnConfirmation: true,
			},
			out: out{
				formulaPathExpected: filepath.Join("testing", "formula", "new"),
				formulaToBeCreated:  "root_testing_formula_new",
			},
		},
		{
			name: "success when new formula is added a lower level of the tree",
			in: in{
				inputOldFormula:       "rit testing withOneMoreLevel level",
				inputNewFormula:       "rit testing level",
				approveOnConfirmation: true,
			},
			out: out{
				formulaPathExpected: filepath.Join("testing", "level"),
				formulaToBeCreated:  "root_testing_level",
				formulaToBeEmpty:    "root_testing_withOneMoreLevel_level",
			},
		},
		{
			name: "success when new formula exists in two workspaces",
			in: in{
				inputOldFormula:         "rit testing formula",
				inputNewFormula:         "rit testing formulaCustom",
				workspaceSelected:       "Custom (" + repoPathWSCustom + ")",
				customWorkspaceSelected: true,
				approveOnConfirmation:   true,
			},
			out: out{
				formulaPathExpected: filepath.Join("testing", "formulaCustom"),
				formulaToBeCreated:  "root_testing_formulaCustom",
				formulaToBeEmpty:    "root_testing_formula",
			},
		},
		{
			name: "err when invalid workspace flag",
			in: in{
				args: []string{
					"--oldName=rit testing formula",
					"--newName=rit testing other",
					"--workspace=test",
				},
			},
			out: out{
				want: errors.New("workspace \"test\" was not found"),
			},
		},
		{
			name: "err when new formula exists in two workspaces and workspace flag is empty",
			in: in{
				customWorkspaceSelected: true,
				args: []string{
					"--oldName=rit testing formula",
					"--newName=rit testing other",
				},
			},
			out: out{
				want: errors.New("formula \"rit testing formula\" was found in 2 workspaces. Please enter a value for the 'workspace' flag"),
			},
		},
		{
			name: "err when new formula exists in two workspaces and workspace flag is invalid",
			in: in{
				customWorkspaceSelected: true,
				args: []string{
					"--oldName=rit testing formula",
					"--newName=rit testing other",
					"--workspace=test",
				},
			},
			out: out{
				want: errors.New("workspace \"test\" was not found"),
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
			_ = streams.Unzip(zipFile, repoPathWS)
			_ = streams.Unzip(zipTree, repoPathLocalDefault)

			createTree(repoPathWS, repoPathLocalDefault, treeGen, fileManager)

			if tt.in.customWorkspaceSelected {
				createSaved(repoPathWSCustom)
				createSaved(repoPathLocalCustom)
				_ = streams.Unzip(zipFile, repoPathWSCustom)
				_ = streams.Unzip(zipFile, repoPathLocalCustom)
				_ = streams.Unzip(zipTree, repoPathLocalCustom)

				createTree(repoPathWSCustom, repoPathLocalCustom, treeGen, fileManager)

				workspaces := formula.Workspaces{}
				workspaces["Default"] = repoPathWS
				workspaces["Custom"] = repoPathWSCustom

				setWorkspace(workspaces, ritHome)
			}

			inputTextValidatorMock := new(mocks.InputTextValidatorMock)
			inputTextValidatorMock.On("Text", formulaOldCmdLabel, mock.Anything, mock.Anything).Return(
				tt.in.inputOldFormula, nil,
			)
			inputTextValidatorMock.On("Text", formulaNewCmdLabel, mock.Anything, mock.Anything).Return(
				tt.in.inputNewFormula, nil,
			)

			inputListMock := new(mocks.InputListMock)
			inputListMock.On(
				"List", "We found the old formula \"rit testing formula\" in 2 workspaces. Select the workspace:",
				mock.Anything,
				mock.Anything,
			).Return(tt.in.workspaceSelected, nil)

			inputBoolMock := new(mocks.InputBoolMock)
			inputBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.approveOnConfirmation, nil)

			cmd := NewRenameFormulaCmd(formulaWorkspace, inputListMock, inputTextValidatorMock, inputBoolMock,
				dirManager, validator, createBuilder, treeGen, deleter, home, ritHome)

			cmd.SetArgs(tt.in.args)

			got := cmd.Execute()

			if tt.out.want != nil {
				assert.Equal(t, tt.out.want.Error(), got.Error())
			} else if tt.in.approveOnConfirmation == false {
				assert.Nil(t, got)
			} else {
				assert.Nil(t, got)

				pathWSDir, pathLocalDir, treePath := "", "", ""
				if tt.in.customWorkspaceSelected {
					pathWSDir = filepath.Join(repoPathWSCustom, tt.out.formulaPathExpected, "src")
					pathLocalDir = filepath.Join(repoPathLocalCustom, tt.out.formulaPathExpected, "src")
					treePath = filepath.Join(repoPathLocalCustom, "tree.json")
				} else {
					pathWSDir = filepath.Join(repoPathWS, tt.out.formulaPathExpected, "src")
					pathLocalDir = filepath.Join(repoPathLocalDefault, tt.out.formulaPathExpected, "src")
					treePath = filepath.Join(repoPathLocalDefault, "tree.json")
				}

				bTree, err := fileInfo(treePath)
				assert.Nil(t, err)
				tree, err := getTree([]byte(bTree))
				assert.Nil(t, err)

				assert.DirExists(t, pathWSDir)
				assert.DirExists(t, pathLocalDir)

				assert.True(t, tree.Commands[api.CommandID(tt.out.formulaToBeCreated)].Formula)
				assert.Empty(t, tree.Commands[api.CommandID(tt.out.formulaToBeEmpty)])

			}
		})
	}

}

func createTree(ws, pathLocal string, tg formula.TreeGenerator, fm stream.FileWriteRemover) {
	localTree, _ := tg.Generate(ws)

	jsonString, _ := json.MarshalIndent(localTree, "", "\t")
	pathLocalTreeJSON := filepath.Join(pathLocal, "tree.json")
	_ = ioutil.WriteFile(pathLocalTreeJSON, jsonString, os.ModePerm)
}

func getTree(f []byte) (formula.Tree, error) {
	tree := formula.Tree{}
	if err := json.Unmarshal(f, &tree); err != nil {
		return formula.Tree{}, err
	}
	return tree, nil
}

func setWorkspace(workspaces formula.Workspaces, ritHome string) {
	wsFile := filepath.Join(ritHome, formula.WorkspacesFile)

	content, _ := json.MarshalIndent(workspaces, "", "\t")
	_ = ioutil.WriteFile(wsFile, content, os.ModePerm)
}
