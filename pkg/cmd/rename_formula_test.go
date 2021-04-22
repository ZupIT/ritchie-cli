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
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/renamer"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/formula/validator"
	"github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
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
	repoPathLocal := filepath.Join(home, "ritchie-formulas-local")

	validator := validator.NewValidator()
	renamer := renamer.NewRenamer(dirManager)

	type in struct {
		inputText         string
		inputTextVal      string
		workspaceSelected string
		formulaSelected   string
		args              []string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success on prompt input",
			in: in{
				inputText:         "rit testing formula",
				inputTextVal:      "rit testing formula new",
				workspaceSelected: "Default (" + repoPathLocal + ")",
				formulaSelected:   "rit testing formula",
			},
		},
		{
			name: "success on flag input",
			in: in{
				args: []string{
					"--workspace=Default",
					"--oldNameFormula=rit testing formula",
					"--newNameFormula=rit testing formula new",
				},
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
			want: errors.New("please provide a value for 'workspace'"),
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
			want: errors.New("please provide a value for 'oldNameFormula'"),
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
			want: errors.New("please provide a value for 'newNameFormula'"),
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
			want: errors.New("the formula workspace does not exist, please enter a valid workspace"),
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
			want: errors.New("This formula 'rit testing other' dont's exists on this workspace = 'Default'"),
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
			want: errors.New("This formula 'rit testing formula' already exists on this workspace = 'Default'"),
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
			want: errors.New("This formula 'rit other formula' dont's exists on this workspace = 'Default'"),
		},
		{
			name: "error on flag input when new formula exists in workspace",
			in: in{
				args: []string{
					"--workspace=Default",
					"--oldNameFormula=rit testing formula",
					"--newNameFormula=rit testing formula",
				},
			},
			want: errors.New("This formula 'rit testing formula' already exists on this workspace = 'Default'"),
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
			createSaved(repoPath)
			createSaved(repoPathLocal)

			zipFile := filepath.Join("..", "..", "testdata", "ritchie-formulas-test.zip")
			zipRepositories := filepath.Join("..", "..", "testdata", "repositories.zip")
			_ = streams.Unzip(zipFile, repoPath)
			_ = streams.Unzip(zipRepositories, reposPath)
			_ = streams.Unzip(zipFile, repoPathLocal)

			inputTextMock := new(mocks.InputTextMock)
			inputTextMock.On("Text", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputText, nil)

			inputTextValidatorMock := new(mocks.InputTextValidatorMock)
			inputTextValidatorMock.On("Text", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputTextVal, nil)

			inPath := &mocks.InputPathMock{}
			inPath.On("Read", "Workspace path (e.g.: /home/user/github): ").Return("", nil)
			inputListMock := insertListMock(tt.in.workspaceSelected, tt.in.formulaSelected)

			inputFormula := prompt.NewInputFormula(inputListMock, dirManager)

			cmd := NewRenameFormulaCmd(
				formulaWorkspace,
				inputTextMock,
				inputListMock,
				inPath,
				inputTextValidatorMock,
				dirManager,
				home,
				validator,
				inputFormula,
				renamer,
			)

			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			cmd.SetArgs(tt.in.args)
			got := cmd.Execute()

			assert.Equal(t, tt.want, got)
		})
	}

}

func insertListMock(workspace, formula string) *mocks.InputListMock {
	if workspace == "" || formula == "" {
		return new(mocks.InputListMock)
	}

	firstGroupFormulas := []string{"testing"}
	secondGroupFormulas := []string{"formula", "invalid-volumes-config", "withLatestVersionRequired", "without-build-files", "without-build-sh", "without-dockerfile", "without-dockerimg"}

	formulaSplited := strings.Split(formula, " ")

	inputListMock := new(mocks.InputListMock)
	inputListMock.On("List", "Select a formula workspace: ", mock.Anything, mock.Anything).Return(workspace, nil)
	inputListMock.On("List", "Select a formula or group: ", firstGroupFormulas, mock.Anything).Return(formulaSplited[1], nil)
	inputListMock.On("List", "Select a formula or group: ", secondGroupFormulas, mock.Anything).Return(formulaSplited[2], nil)
	inputListMock.On("List", prompt.InputFormulaQuestionFoundedFormula, mock.Anything, mock.Anything).Return(formula, nil)
	return inputListMock
}
