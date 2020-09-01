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

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type fieldsTestBuildFormulaCmd struct {
	localBuilder     formula.LocalBuilder
	workspaceManager formula.WorkspaceAddListValidator
	directory        stream.DirListChecker
	inList           prompt.InputList
}

func TestBuildFormulaCmd(t *testing.T) {
	userHomeDir := os.TempDir()
	defaultWorkspace := filepath.Join(userHomeDir, formula.DefaultWorkspaceDir)
	someError := errors.New("some error")

	var fieldsDefault fieldsTestBuildFormulaCmd = fieldsTestBuildFormulaCmd{
		localBuilder: LocalBuilderMock{
			build: func(workspacePath, formulaPath string) error {
				return nil
			},
		},
		workspaceManager: WorkspaceAddListValidatorCustomMock{
			list: func() (formula.Workspaces, error) {
				return formula.Workspaces{}, nil
			},
			validate: func(workspace formula.Workspace) error {
				return nil
			},
		},
		directory: DirManagerCustomMock{
			exists: func(dir string) bool {
				return true
			},
			list: func(dir string, hiddenDir bool) ([]string, error) {
				switch dir {
				case defaultWorkspace:
					return []string{"group"}, nil
				case defaultWorkspace + "/group":
					return []string{"verb"}, nil
				case defaultWorkspace + "/group/verb":
					return []string{"src"}, nil
				default:
					return []string{"any"}, nil
				}
			},
		},
		inList: inputListCustomMock{
			list: func(name string, items []string) (string, error) {
				if name == questionSelectFormulaGroup {
					return items[0], nil
				}
				return "Default (/tmp/ritchie-formulas-local)", nil
			},
		},
	}

	tests := []struct {
		name    string
		fields  fieldsTestBuildFormulaCmd
		wantErr bool
	}{
		{
			name:    "Run with success",
			fields:  fieldsTestBuildFormulaCmd{},
			wantErr: false,
		},
		{
			name: "Run with error when workspace list returns err",
			fields: fieldsTestBuildFormulaCmd{
				workspaceManager: WorkspaceAddListValidatorCustomMock{
					list: func() (formula.Workspaces, error) {
						return formula.Workspaces{}, someError
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Run with error when readFormulas returns err",
			fields: fieldsTestBuildFormulaCmd{
				directory: DirManagerCustomMock{
					exists: func(dir string) bool {
						return true
					},
					list: func(dir string, hiddenDir bool) ([]string, error) {
						switch dir {
						case defaultWorkspace:
							return []string{"group"}, someError
						default:
							return []string{"any"}, nil
						}
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Run with error when question about select formula or group returns err",
			fields: fieldsTestBuildFormulaCmd{
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectFormulaGroup {
							return "any", someError
						}
						return "Default (/tmp/ritchie-formulas-local)", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Run with sucess when the selected formula is deeper in the tree",
			fields: fieldsTestBuildFormulaCmd{
				directory: DirManagerCustomMock{
					exists: func(dir string) bool {
						return true
					},
					list: func(dir string, hiddenDir bool) ([]string, error) {
						switch dir {
						case defaultWorkspace:
							return []string{"group"}, nil
						case defaultWorkspace + "/group":
							return []string{"verb", "src"}, nil
						case defaultWorkspace + "/group/verb":
							return []string{"src"}, nil
						default:
							return []string{"any"}, nil
						}
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectFormulaGroup {
							return items[0], nil
						}
						if name == questionAboutFoundedFormula {
							return optionOtherFormula, nil
						}
						return "Default (/tmp/ritchie-formulas-local)", nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Run with sucess when selected formula is less deep in the tree",
			fields: fieldsTestBuildFormulaCmd{
				directory: DirManagerCustomMock{
					exists: func(dir string) bool {
						return true
					},
					list: func(dir string, hiddenDir bool) ([]string, error) {
						switch dir {
						case defaultWorkspace:
							return []string{"group"}, nil
						case defaultWorkspace + "/group":
							return []string{"verb", "src"}, nil
						case defaultWorkspace + "/group/verb":
							return []string{"src"}, nil
						default:
							return []string{"any"}, nil
						}
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectFormulaGroup {
							return items[0], nil
						}
						if name == questionAboutFoundedFormula {
							return "rit group", nil
						}
						return "Default (/tmp/ritchie-formulas-local)", nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Run with error when readFormula returns error on second call",
			fields: fieldsTestBuildFormulaCmd{
				directory: DirManagerCustomMock{
					exists: func(dir string) bool {
						return true
					},
					list: func(dir string, hiddenDir bool) ([]string, error) {
						switch dir {
						case defaultWorkspace:
							return []string{"group"}, nil
						case defaultWorkspace + "/group":
							return []string{"verb"}, someError
						default:
							return []string{"any"}, nil
						}
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := getFields(fieldsDefault, tt.fields)

			commandBuildFormula := NewBuildFormulaCmd(
				userHomeDir,
				fields.localBuilder,
				fields.workspaceManager,
				WatcherMock{},
				fields.directory,
				inputTextMock{},
				fields.inList,
				tutorialFindWithReturnDisabled(),
			)
			commandBuildFormula.PersistentFlags().Bool("stdin", false, "input by stdin")

			if commandBuildFormula == nil {
				t.Errorf("commandBuildFormula got %v", commandBuildFormula)
				return
			}

			if err := commandBuildFormula.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", commandBuildFormula.Use, err, tt.wantErr)
			}
		})
	}
}

func tutorialFindWithReturnDisabled() rtutorial.Finder {
	return TutorialFindSetterCustomMock{
		find: func() (rtutorial.TutorialHolder, error) {
			return rtutorial.TutorialHolder{Current: "disabled"}, nil
		},
	}
}

func getFields(fieldsDefault fieldsTestBuildFormulaCmd, fieldsTest fieldsTestBuildFormulaCmd) fieldsTestBuildFormulaCmd {
	var fields fieldsTestBuildFormulaCmd = fieldsDefault

	if fieldsTest.directory != nil {
		fields.directory = fieldsTest.directory
	}
	if fieldsTest.inList != nil {
		fields.inList = fieldsTest.inList
	}
	if fieldsTest.localBuilder != nil {
		fields.localBuilder = fieldsTest.localBuilder
	}
	if fieldsTest.workspaceManager != nil {
		fields.workspaceManager = fieldsTest.workspaceManager
	}

	return fields
}
