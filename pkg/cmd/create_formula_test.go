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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewCreateFormulaCmd(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tplM := template.NewManager("../../testdata", dirManager)
	tChecker := tree.NewChecker(treeMock{})
	workspaceMock := new(mocks.WorkspaceForm)
	workspaceMock.On("List").Return(formula.Workspaces{}, nil)
	workspaceMock.On("Add", mock.Anything).Return(nil)
	workspaceMock.On("CurrentHash", mock.Anything).Return("dsadasdas", nil)
	workspaceMock.On("UpdateHash", mock.Anything, mock.Anything).Return(nil)

	cmd := NewCreateFormulaCmd(
		os.TempDir(),
		formCreator{},
		tplM,
		workspaceMock,
		inputTextMock{},
		inputTextValidatorMock{},
		inputListMock{},
		TutorialFinderMock{},
		tChecker,
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewCreateFormulaCmd got %v", cmd)
		return
	}

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestCreateFormulaCmd(t *testing.T) {
	type in struct {
		tm              template.Manager
		inText          prompt.InputText
		inTextValidator prompt.InputTextValidator
		inList          prompt.InputList
	}

	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "error on input text validator",
			in: in{
				inTextValidator: inputTextValidatorErrorMock{},
			},
			wantErr: true,
		},
		{
			name: "error on template manager Validate func",
			in: in{
				inTextValidator: inputTextValidatorMock{},
				tm: TemplateManagerCustomMock{ValidateMock: func() error {
					return errors.New("error on validate func")
				}},
			},
			wantErr: true,
		},
		{
			name: "error on template manager Languages func",
			in: in{
				inTextValidator: inputTextValidatorMock{},
				tm: TemplateManagerCustomMock{
					ValidateMock: func() error {
						return nil
					},
					LanguagesMock: func() ([]string, error) {
						return []string{}, errors.New("error on language func")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error on input list",
			in: in{
				inTextValidator: inputTextValidatorMock{},
				tm: TemplateManagerCustomMock{
					ValidateMock: func() error {
						return nil
					},
					LanguagesMock: func() ([]string, error) {
						return []string{}, nil
					},
				},
				inList: inputListErrorMock{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createFormulaCmd := NewCreateFormulaCmd(
				os.TempDir(),
				formCreator{},
				tt.in.tm,
				new(mocks.WorkspaceForm),
				tt.in.inText,
				tt.in.inTextValidator,
				tt.in.inList,
				TutorialFinderMock{},
				tree.CheckerManager{},
			)
			createFormulaCmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := createFormulaCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("%s = %v, want %v", createFormulaCmd.Use, err, nil)
			}
		})
	}

}

func TestCreateFormula(t *testing.T) {
	workDir := filepath.Join(os.TempDir(), ".ritchie-formulas-local")
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
		name string
		in   in
		want error
	}{
		{
			name: "success",
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
			creatorMock := new(mocks.FormCreator)
			creatorMock.On("Create", mock.Anything).Return(tt.in.createFormErr)
			creatorMock.On("Build", mock.Anything).Return(tt.in.buildFormErr)
			workspaceMock := new(mocks.WorkspaceForm)
			workspaceMock.On("CurrentHash", mock.Anything).Return(tt.in.currentHash, tt.in.currentHashErr)
			workspaceMock.On("UpdateHash", mock.Anything, mock.Anything).Return(tt.in.updateHashErr)
			tutorialMock := new(mocks.TutorialFindSetterMock)
			tutorialMock.On("Find").Return(tt.in.tutorialHolder, tt.in.tutorialErr)

			createForm := createFormulaCmd{
				formula:   creatorMock,
				workspace: workspaceMock,
				tutorial:  tutorialMock,
			}

			got := createForm.create(cf)

			assert.Equal(t, tt.want, got)
		})
	}
}

type TemplateManagerCustomMock struct {
	LanguagesMock         func() ([]string, error)
	LangTemplateFilesMock func(lang string) ([]template.File, error)
	ResolverNewPathMock   func(oldPath, newDir, lang, workspacePath string) (string, error)
	ValidateMock          func() error
}

func (tm TemplateManagerCustomMock) Languages() ([]string, error) {
	return tm.LanguagesMock()
}

func (tm TemplateManagerCustomMock) LangTemplateFiles(lang string) ([]template.File, error) {
	return tm.LangTemplateFilesMock(lang)
}
func (tm TemplateManagerCustomMock) ResolverNewPath(oldPath, newDir, lang, workspacePath string) (string, error) {
	return tm.ResolverNewPathMock(oldPath, newDir, lang, workspacePath)
}
func (tm TemplateManagerCustomMock) Validate() error {
	return tm.ValidateMock()
}
