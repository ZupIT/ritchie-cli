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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type FileManagerMock struct{}

func (f FileManagerMock) Write(path string, content []byte) error {
	return errors.New("some error")
}

func (f FileManagerMock) Remove(path string) error {
	return nil
}

type fieldsTestDeleteFormulaCmd struct {
	workspaceManager formula.WorkspaceAddLister
	directory        stream.DirListChecker
	inList           prompt.InputList
	fileManager      stream.FileWriteRemover
	inBool           prompt.InputBool
}

func TestNewDeleteFormulaCmd(t *testing.T) {
	userHomeDir := os.TempDir()
	ritchieHomeDir := filepath.Join(os.TempDir(), ".rit")
	defaultWorkspace := filepath.Join(userHomeDir, formula.DefaultWorkspaceDir)
	someError := errors.New("some error")

	workspacePath := filepath.Join(os.TempDir(), "ritchie-formulas-local")
	if err := os.MkdirAll(filepath.Join(workspacePath, "group", "verb", "src"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaCmd got error %v", err)
	}
	if err := os.MkdirAll(filepath.Join(workspacePath, "group", "verb2", "src"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaCmd got error %v", err)
	}
	if err := os.MkdirAll(filepath.Join(os.TempDir(), ".rit", "repos", "local", "group", "verb", "src"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaCmd got error %v", err)
	}

	var fieldsDefault fieldsTestDeleteFormulaCmd = fieldsTestDeleteFormulaCmd{
		workspaceManager: WorkspaceAddListerCustomMock{
			list: func() (formula.Workspaces, error) {
				return formula.Workspaces{
					"Default": defaultWorkspace,
				}, nil
			},
			add: func(workspace formula.Workspace) error {
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
				return fmt.Sprintf("Default (%s)", defaultWorkspace), nil
			},
		},
		fileManager: stream.FileManager{},
		inBool:      inputTrueMock{},
	}

	tests := []struct {
		name    string
		fields  fieldsTestDeleteFormulaCmd
		wantErr bool
	}{
		{
			name:    "Run with success",
			fields:  fieldsTestDeleteFormulaCmd{},
			wantErr: false,
		},
		{
			name: "Run with error when workspace list returns err",
			fields: fieldsTestDeleteFormulaCmd{
				workspaceManager: WorkspaceAddListerCustomMock{
					list: func() (formula.Workspaces, error) {
						return formula.Workspaces{}, someError
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Run with error when readFormulas returns err",
			fields: fieldsTestDeleteFormulaCmd{
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
			fields: fieldsTestDeleteFormulaCmd{
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
			fields: fieldsTestDeleteFormulaCmd{
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
			fields: fieldsTestDeleteFormulaCmd{
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
			fields: fieldsTestDeleteFormulaCmd{
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
		{
			name: "Run with error when add new workspace",
			fields: fieldsTestDeleteFormulaCmd{
				workspaceManager: WorkspaceAddListerCustomMock{
					list: func() (formula.Workspaces, error) {
						return formula.Workspaces{}, nil
					},
					add: func(workspace formula.Workspace) error {
						return someError
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == questionSelectFormulaGroup {
							return "any", someError
						}
						return "Ritchie-Formulas (/tmp/ritchie-formulas)", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Run with error when recreate tree.json",
			fields: fieldsTestDeleteFormulaCmd{
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
				fileManager: FileManagerMock{},
			},
			wantErr: true,
		},
		{
			name: "Run with success when choose not to delete formula",
			fields: fieldsTestDeleteFormulaCmd{
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
				inBool: inputFalseMock{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := getFieldsDeleteFormula(fieldsDefault, tt.fields)

			if err := os.MkdirAll(filepath.Join(workspacePath, "group", "verb", "src"), os.ModePerm); err != nil {
				t.Errorf("TestNewDeleteFormulaCmd got error %v", err)
			}
			if err := os.MkdirAll(filepath.Join(os.TempDir(), ".rit", "repos", "local", "group", "verb", "src"), os.ModePerm); err != nil {
				t.Errorf("TestNewDeleteFormulaCmd got error %v", err)
			}

			cmd := NewDeleteFormulaCmd(
				userHomeDir,
				ritchieHomeDir,
				fields.workspaceManager,
				fields.directory,
				fields.inBool,
				inputTextMock{},
				fields.inList,
				treeGeneratorMock{},
				fields.fileManager,
			)
			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

			if cmd == nil {
				t.Errorf("TestNewDeleteFormulaCmd got %v", cmd)
				return
			}

			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", cmd.Use, err, tt.wantErr)
			}
		})
	}
}

func getFieldsDeleteFormula(fieldsDefault fieldsTestDeleteFormulaCmd, fieldsTest fieldsTestDeleteFormulaCmd) fieldsTestDeleteFormulaCmd {
	var fields fieldsTestDeleteFormulaCmd = fieldsDefault

	if fieldsTest.directory != nil {
		fields.directory = fieldsTest.directory
	}
	if fieldsTest.inList != nil {
		fields.inList = fieldsTest.inList
	}
	if fieldsTest.workspaceManager != nil {
		fields.workspaceManager = fieldsTest.workspaceManager
	}
	if fieldsTest.fileManager != nil {
		fields.fileManager = fieldsTest.fileManager
	}
	if fieldsTest.inBool != nil {
		fields.inBool = fieldsTest.inBool
	}

	return fields
}

func TestNewDeleteFormulaStdin(t *testing.T) {
	workspacePath := filepath.Join(os.TempDir(), "ritchie-formulas-local")
	if err := os.MkdirAll(filepath.Join(workspacePath, "mock", "test", "src"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaStdin got error %v", err)
	}

	if err := ioutil.WriteFile(filepath.Join(workspacePath, "mock", "test", "help.txt"), []byte{'a'}, os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaStdin got error %v", err)
	}

	if err := os.MkdirAll(filepath.Join(workspacePath, "mock", "test", "nested", "src"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaStdin got error %v", err)
	}

	if err := os.MkdirAll(filepath.Join(os.TempDir(), ".rit", "repos", "local", "mock", "test", "src"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaStdin got error %v", err)
	}

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	treeGen := tree.NewGenerator(dirManager, fileManager)

	cmd := NewDeleteFormulaCmd(
		os.TempDir(),
		filepath.Join(os.TempDir(), ".rit"),
		workspaceForm{},
		dirManager,
		inputTrueMock{},
		inputTextMock{},
		inputListMock{},
		treeGen,
		stream.FileManager{},
	)
	cmd.PersistentFlags().Bool("stdin", true, "input by stdin")

	json := fmt.Sprintf("{\"workspace_path\": \"%s\", \"formula\": \"rit mock test\"}\n", workspacePath)
	newReader := strings.NewReader(json)
	cmd.SetIn(newReader)

	if cmd == nil {
		t.Errorf("NewDeleteFormulaCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}

	json = fmt.Sprintf("{\"workspace_path\": \"%s\", \"formula\": \"mock test\"}\n", workspacePath)
	newReader = strings.NewReader(json)
	cmd.SetIn(newReader)

	if err := cmd.Execute(); err == nil {
		t.Errorf("%s = %v, want %v", cmd.Use, nil, ErrIncorrectFormulaName)
	}

	if err := os.MkdirAll(filepath.Join(workspacePath, "mock", "test", "src"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaStdin got error %v", err)
	}

	if err := os.MkdirAll(filepath.Join(os.TempDir(), ".rit", "repos", "local", "mock", "test", "src"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaStdin got error %v", err)
	}

	cmd = NewDeleteFormulaCmd(
		os.TempDir(),
		filepath.Join(os.TempDir(), ".rit"),
		workspaceForm{},
		dirManager,
		inputTrueMock{},
		inputTextMock{},
		inputListMock{},
		treeGen,
		FileManagerMock{},
	)
	cmd.PersistentFlags().Bool("stdin", true, "input by stdin")

	json = fmt.Sprintf("{\"workspace_path\": \"%s\", \"formula\": \"rit mock test\"}\n", workspacePath)
	newReader = strings.NewReader(json)
	cmd.SetIn(newReader)

	if cmd == nil {
		t.Errorf("NewDeleteFormulaCmd got %v", cmd)
	}

	if err := cmd.Execute(); err == nil {
		t.Errorf("%s = %v, want %v", cmd.Use, nil, errors.New("some error"))
	}
}
