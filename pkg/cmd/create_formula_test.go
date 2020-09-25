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
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewCreateFormulaCmd(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tplM := template.NewManager("../../testdata", dirManager)
	cmd := NewCreateFormulaCmd(
		os.TempDir(),
		formCreator{},
		tplM,
		workspaceForm{},
		inputTextMock{},
		inputTextValidatorMock{},
		inputListMock{},
		TutorialFinderMock{},
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewCreateFormulaCmd got %v", cmd)
		return
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
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
			name: "error with a not allowed char on formula cmd",
			in: in{
				inTextValidator: inputTextValidatorCustomMock{
					text: func(name string, validate func(interface{}) error, helper ...string) (string, error) {
						return "rit@", nil
					}},
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
			//inputTextValidator := prompt.NewSurveyTextValidator()

			createFormulaCmd := NewCreateFormulaCmd(
				os.TempDir(),
				formCreator{},
				tt.in.tm,
				workspaceForm{},
				tt.in.inText,
				 tt.in.inTextValidator,
				//inputTextValidator,
				tt.in.inList,
				TutorialFinderMock{},
			)

			createFormulaCmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := createFormulaCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("%s = %v, want %v", createFormulaCmd.Use, err, nil)
			}

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
