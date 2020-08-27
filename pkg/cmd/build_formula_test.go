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
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestBuildFormulaCmd(t *testing.T) {
	userHomeDir := os.TempDir()
	defaultWorkspace := filepath.Join(userHomeDir, formula.DefaultWorkspaceDir)
	workspaceManager := WorkspaceAddListValidatorCustomMock{
		list: func() (formula.Workspaces, error) {
			return formula.Workspaces{}, nil
		},
		validate: func(workspace formula.Workspace) error {
			return nil
		},
	}
	localBuild := LocalBuilderMock{
		build: func(workspacePath, formulaPath string) error {
			return nil
		},
	}
	dirManager := DirManagerCustomMock{
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
	}

	commandBuildFormula := NewBuildFormulaCmd(
		userHomeDir,
		localBuild,
		workspaceManager,
		WatcherMock{},
		dirManager,
		inputTextMock{},
		inputListCustomMock{
			list: func(name string, items []string) (string, error) {
				if name == questionSelectFormulaGroup {
					return items[0], nil
				}
				return "Default (/tmp/ritchie-formulas-local)", nil
			},
		},
		TutorialFinderMock{},
	)
	commandBuildFormula.PersistentFlags().Bool("stdin", false, "input by stdin")
	if commandBuildFormula == nil {
		t.Errorf("buildFormulaCmd got %v", commandBuildFormula)
		return
	}

	if err := commandBuildFormula.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", commandBuildFormula.Use, err, nil)
	}
}
