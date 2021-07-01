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
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/formula/validator"
	"github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/git/gitlab"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestCreateFormulaCmd(t *testing.T) {
	validator := validator.New()

	type in struct {
		inputBool        bool
		inputBoolErr     error
		inputText        string
		inputTextErr     error
		inputTextVal     string
		inputTextValErr  error
		tempValErr       error
		tempLanguages    []string
		tempLanguagesErr error
		inputList        string
		inputListErr     error
		wspaceList       formula.Workspaces
		wspaceListErr    error
		wspaceAddErr     error
		dirCreate        error
		createErr        error
		inputPath        string
		inputPathErr     error
		inputWspace      string
	}

	tests := []struct {
		name       string
		in         in
		want       error
		inputFlags []string
	}{
		{
			name: "success",
			in: in{
				inputTextVal:  "rit test test",
				tempLanguages: []string{"go", "rust", "java", "kotlin"},
				inputList:     "go",
			},
		},
		{
			name: "success with new workspace dir",
			in: in{
				inputTextVal:  "rit test test",
				wspaceAddErr:  workspace.ErrInvalidWorkspace,
				inputBool:     true,
				tempLanguages: []string{"go", "rust", "java", "kotlin"},
				dirCreate:     nil,
			},
		},
		{
			name: "add dir error",
			in: in{
				inputTextVal:  "rit test test",
				tempLanguages: []string{"go", "rust", "java", "kotlin"},
				inputList:     "go",
				inputBool:     true,
				wspaceAddErr:  workspace.ErrInvalidWorkspace,
				dirCreate:     errors.New("failed to create dir"),
			},
			want: errors.New("failed to create dir"),
		},
		{
			name: "error on input text validator",
			in: in{
				inputTextValErr: errors.New("error on input text"),
			},
			want: ErrFormulaCmdNotBeEmpty,
		},
		{
			name: "error on template manager Validate func",
			in: in{
				inputTextVal: "rit test test",
				tempValErr:   errors.New("error on validate func"),
			},
			want: errors.New("error on validate func"),
		},
		{
			name: "error on template manager Languages func",
			in: in{
				inputTextVal:     "rit test test",
				tempLanguagesErr: errors.New("error on language func"),
			},
			want: errors.New("error on language func"),
		},
		{
			name: "error on input list",
			in: in{
				inputTextVal:  "rit test test",
				tempLanguages: []string{"go", "java", "c", "rust"},
				inputListErr:  errors.New("error to list languages"),
			},
			want: errors.New("error to list languages"),
		},
		{
			name: "list workspace error",
			in: in{
				inputTextVal:  "rit test test",
				tempLanguages: []string{"go", "rust", "java", "kotlin"},
				inputList:     "go",
				wspaceListErr: errors.New("error to list workspaces"),
			},
			want: errors.New("error to list workspaces"),
		},
		{
			name: "error on workspace list with flags",
			in: in{
				tempLanguages: []string{"go", "java", "c", "rust"},
				wspaceListErr: errors.New("error to list workspaces"),
			},
			want:       errors.New("error to list workspaces"),
			inputFlags: []string{"--name=rit test test", "--language=go", "--workspace=default"},
		},
		{
			name: "success with new workspace in function FormulaWorkspaceInput",
			in: in{
				inputTextVal: "rit test test",
				inputList:    newWorkspace,
				inputWspace:  "Test",
			},
		},
		{
			name: "error function FormulaWorkspaceInput",
			in: in{
				inputTextVal:  "rit test test",
				tempLanguages: []string{"go", "rust", "java", "kotlin"},
				inputList:     newWorkspace,
				inputWspace:   "Invalid workspace",
			},
			want: workspace.ErrInvalidWorkspaceName,
		},
		{
			name: "add workspace error",
			in: in{
				inputTextVal:  "rit test test",
				tempLanguages: []string{"go", "rust", "java", "kotlin"},
				inputList:     "go",
				wspaceAddErr:  errors.New("error to add workspace"),
			},
			want: errors.New("error to add workspace"),
		},
		{
			name: "formula create error",
			in: in{
				inputTextVal:  "rit test test",
				tempLanguages: []string{"go", "rust", "java", "kotlin"},
				inputList:     "go",
				createErr:     errors.New("error to create formula"),
			},
			want: errors.New("error to create formula"),
		},
		{
			name: "success with flags",
			in: in{
				wspaceList: formula.Workspaces{
					"Default": "C:\\Users\\mauri\\ritchie-formulas-local",
				},
				tempLanguages: []string{"go", "rust", "java", "kotlin"},
			},
			inputFlags: []string{"--name=rit test test", "--language=go", "--workspace=Default"},
		},
		{
			name: "err invalidWorkspace",
			in: in{
				wspaceList: formula.Workspaces{
					"Default": "C:\\Users\\mauri\\ritchie-formulas-local",
				},
			},
			want:       errors.New(InvalidWorkspace),
			inputFlags: []string{"--name=rit test test", "--language=go", "--workspace=invalidWorkspace"},
		},
		{
			name: "err invalidLanguage",
			in: in{
				tempLanguages: []string{"go", "rust", "java", "kotlin"},
			},
			want:       errors.New("language not found"),
			inputFlags: []string{"--name=rit test test", "--language=invalidLanguage", "--workspace=Default"},
		},
		{
			name: "err empty formula command",
			in: in{
				inputTextVal: "",
			},
			want: errors.New("this input must not be empty"),
		},
		{
			name: "error with new workspace with wrong path",
			in: in{
				inputTextVal: "rit test test",
				inputList:    newWorkspace,
				inputWspace:  "Default",
				inputPathErr: errors.New("invalid input path"),
				dirCreate:    nil,
			},
			want: errors.New("invalid input path"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceMock := new(mocks.WorkspaceForm)
			if tt.name == "success with new workspace dir" {
				workspaceMock.On("Add", mock.Anything).Return(tt.in.wspaceAddErr).Once()
				workspaceMock.On("Add", mock.Anything).Return(nil).Once()

			} else {
				workspaceMock.On("Add", mock.Anything).Return(tt.in.wspaceAddErr)
			}

			workspaceMock.On("List").Return(tt.in.wspaceList, tt.in.wspaceListErr)
			workspaceMock.On("CurrentHash", mock.Anything).Return("48d47029-2abf-4a2e-b5f2-f5b60471423e", nil)
			workspaceMock.On("UpdateHash", mock.Anything, mock.Anything).Return(nil)

			templateManagerMock := new(mocks.TemplateManagerMock)
			templateManagerMock.On("Validate").Return(tt.in.tempValErr)
			templateManagerMock.On("Languages").Return(tt.in.tempLanguages, tt.in.tempLanguagesErr)

			formulaCreatorMock := new(mocks.FormCreator)
			formulaCreatorMock.On("Create", mock.Anything).Return(tt.in.createErr)
			formulaCreatorMock.On("Build", mock.Anything).Return(nil)

			inputTextMock := new(mocks.InputTextMock)
			inputTextMock.On("Text", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputText, tt.in.inputTextErr)

			inputTextValidatorMock := new(mocks.InputTextValidatorMock)
			inputTextValidatorMock.On("Text", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputTextVal, tt.in.inputTextValErr).Once()
			inputTextValidatorMock.On("Text", "Workspace name: ", mock.Anything).Return(tt.in.inputWspace, tt.in.inputTextValErr)

			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputList, tt.in.inputListErr)

			inPath := &mocks.InputPathMock{}
			inPath.On("Read", "Workspace path (e.g.: /home/user/github): ").Return(tt.in.inputPath, tt.in.inputPathErr)

			tutorialMock := new(mocks.TutorialFindSetterMock)
			tutorialMock.On("Find").Return(rtutorial.TutorialHolder{Current: "enabled"}, nil)

			treeMock := new(mocks.TreeManager)
			treeMock.On("Check").Return([]api.CommandID{})

			inputBoolM := new(mocks.InputBoolMock)
			inputBoolM.On("Bool", InvalidWorkspace, []string{"no", "yes"}, mock.Anything).Return(tt.in.inputBool, tt.in.inputBoolErr)

			directoryMock := new(mocks.DirManagerMock)
			directoryMock.On("Create", mock.Anything).Return(tt.in.dirCreate)

			createFormulaCmd := NewCreateFormulaCmd(
				os.TempDir(),
				formulaCreatorMock,
				templateManagerMock,
				workspaceMock,
				inputTextMock,
				inputTextValidatorMock,
				inputListMock,
				inPath,
				tutorialMock,
				treeMock,
				validator,
				inputBoolM,
				directoryMock,
			)
			createFormulaCmd.SetArgs([]string{})
			// TODO: remove it after being deprecated
			createFormulaCmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			if len(tt.inputFlags) > 1 {
				createFormulaCmd.SetArgs(tt.inputFlags)
			}
			got := createFormulaCmd.Execute()
			assert.Equal(t, tt.want, got)
		})
	}

}

func TestCreateFormula(t *testing.T) {
	tmpDir := os.TempDir()
	ritchieHomeDir := filepath.Join(tmpDir, ".rit_create_formula")
	workDir := filepath.Join(tmpDir, ".ritchie-formulas-local")
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	_ = dirManager.Remove(ritchieHomeDir)
	_ = dirManager.Remove(workDir)

	cf := formula.Create{
		FormulaCmd: "rit test test",
		Lang:       "go",
		Workspace: formula.Workspace{
			Name: "default",
			Dir:  workDir,
		},
		FormulaPath: filepath.Join(workDir, "test", "test"),
	}

	type in struct {
		createFormErr  error
		buildFormErr   error
		currentHash    string
		currentHashErr error
		updateHashErr  error
		tutorialHolder rtutorial.TutorialHolder
		tutorialErr    error
	}

	tests := []struct {
		name        string
		in          in
		withoutMock bool
		want        error
	}{
		{
			name:        "success without mock",
			withoutMock: true,
		},
		{
			name: "success mocked",
			in: in{
				currentHash: "a6edc906-2f9f-5fb2-a373-efac406f0ef2",
			},
		},
		{
			name: "success with tutorial enabled",
			in: in{
				currentHash: "a6edc906-2f9f-5fb2-a373-efac406f0ef2",
				tutorialHolder: rtutorial.TutorialHolder{
					Current: "enabled",
				},
			},
		},
		{
			name: "create formula error",
			in: in{
				createFormErr: errors.New("error to create formula"),
			},
			want: errors.New("error to create formula"),
		},
		{
			name: "build formula error",
			in: in{
				buildFormErr: errors.New("error to build formula"),
			},
			want: errors.New("error to build formula"),
		},
		{
			name: "current hash error",
			in: in{
				currentHashErr: errors.New("error to create current formula hash"),
			},
			want: errors.New("error to create current formula hash"),
		},
		{
			name: "update hash error",
			in: in{
				updateHashErr: errors.New("error to update formula hash"),
			},
			want: errors.New("error to update formula hash"),
		},
		{
			name: "tutorial disable",
			in: in{
				tutorialHolder: rtutorial.TutorialHolder{
					Current: "disable",
				},
			},
		},
		{
			name: "tutorial error",
			in: in{
				tutorialErr: errors.New("error to find tutorial"),
			},
			want: errors.New("error to find tutorial"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var createForm createFormulaCmd
			if !tt.withoutMock {
				creatorMock := new(mocks.FormCreator)
				creatorMock.On("Create", mock.Anything).Return(tt.in.createFormErr)
				creatorMock.On("Build", mock.Anything).Return(tt.in.buildFormErr)
				workspaceMock := new(mocks.WorkspaceForm)
				workspaceMock.On("CurrentHash", mock.Anything).Return(tt.in.currentHash, tt.in.currentHashErr)
				workspaceMock.On("UpdateHash", mock.Anything, mock.Anything).Return(tt.in.updateHashErr)
				tutorialMock := new(mocks.TutorialFindSetterMock)
				tutorialMock.On("Find").Return(tt.in.tutorialHolder, tt.in.tutorialErr)

				createForm = createFormulaCmd{
					formula:   creatorMock,
					workspace: workspaceMock,
					tutorial:  tutorialMock,
				}
			} else {
				createForm = createFormulaCmdDeps(ritchieHomeDir, dirManager, fileManager)
			}

			got := createForm.create(cf)
			assert.Equal(t, tt.want, got)

			if tt.withoutMock {
				assert.DirExists(t, ritchieHomeDir)

				hashesDir := filepath.Join(ritchieHomeDir, "hashes")
				assert.DirExists(t, hashesDir)

				reposDir := filepath.Join(ritchieHomeDir, "repos")
				assert.DirExists(t, reposDir)
				assert.DirExists(t, filepath.Join(reposDir, "local-default"))
				assert.FileExists(t, filepath.Join(reposDir, "local-default", "tree.json"))

				fileName := strings.ReplaceAll(cf.FormulaPath, string(os.PathSeparator), "-") + ".txt"
				assert.FileExists(t, filepath.Join(hashesDir, fileName))

				assert.FileExists(t, filepath.Join(reposDir, "repositories.json"))
			}
		})
	}
}

func createFormulaCmdDeps(ritchieHomeDir string, dirManager stream.DirManager, fileManager stream.FileManager) createFormulaCmd {
	treeGen := tree.NewGenerator(dirManager, fileManager)
	githubRepo := github.NewRepoManager(http.DefaultClient)
	gitlabRepo := gitlab.NewRepoManager(http.DefaultClient)

	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: githubRepo, NewRepoInfo: github.NewRepoInfo})
	repoProviders.Add("Gitlab", formula.Git{Repos: gitlabRepo, NewRepoInfo: gitlab.NewRepoInfo})
	repoCreator := repo.NewCreator(ritchieHomeDir, repoProviders, dirManager, fileManager)
	repoLister := repo.NewLister(ritchieHomeDir, fileManager)
	repoWriter := repo.NewWriter(ritchieHomeDir, fileManager)
	repoListWriter := repo.NewListWriter(repoLister, repoWriter)
	repoDeleter := repo.NewDeleter(ritchieHomeDir, repoListWriter, dirManager)
	repoDetail := repo.NewDetail(repoProviders)
	repoListWriteCreator := repo.NewCreateWriteListDetailDeleter(repoLister, repoCreator, repoWriter, repoDetail, repoDeleter)
	repoAdder := repo.NewAdder(ritchieHomeDir, repoListWriteCreator, treeGen)
	repoListDetailWriter := repo.NewListDetailWrite(repoLister, repoDetail, repoWriter)

	treeManager := tree.NewTreeManager(ritchieHomeDir, repoListDetailWriter, api.CoreCmds)
	tmpManager := template.NewManager("../../testdata", dirManager)
	createManager := creator.NewCreator(treeManager, dirManager, fileManager, tmpManager)
	formBuildLocal := builder.NewBuildLocal(ritchieHomeDir, dirManager, repoAdder)
	createBuilder := formula.NewCreateBuilder(createManager, formBuildLocal)
	buildLocal := builder.NewBuildLocal(ritchieHomeDir, dirManager, repoAdder)
	wspaceManager := workspace.New(ritchieHomeDir, os.TempDir(), dirManager, buildLocal, treeGen)
	tutorialFinder := rtutorial.NewFinder(ritchieHomeDir)

	return createFormulaCmd{
		formula:   createBuilder,
		workspace: wspaceManager,
		tutorial:  tutorialFinder,
	}
}
