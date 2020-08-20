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
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

func TestNewDeleteFormulaCmdStdin(t *testing.T) {
	treeMock := treeMock{
		tree: formula.Tree{
			Commands: api.Commands{
				{
					Id:     "root_mock",
					Parent: "root",
					Usage:  "mock",
					Help:   "mock for add",
				},
				{
					Id:      "root_mock_test",
					Parent:  "root_mock",
					Usage:   "test",
					Help:    "test for add",
					Formula: true,
				},
			},
		},
		value: "LOCAL",
	}

	workspace := os.TempDir() + "/ritchie-formulas-local"

	if err := os.MkdirAll(workspace+"/mock/test/scr", os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaCmdStdin got error %v", err)
	}

	if err := os.MkdirAll(os.TempDir()+"/.rit/repos/local/mock/test/scr", os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaCmdStdin got error %v", err)
	}

	json := `{"workspace": "` + workspace + `", "groups": ["mock", "test"]}`
	tmpfile, oldStdin, err := stdin.WriteToStdin(json)
	defer os.Remove(tmpfile.Name())
	defer func() { os.Stdin = oldStdin }()
	if err != nil {
		t.Errorf("TestNewDeleteFormulaCmdStdin got error %v", err)
	}

	cmd := NewDeleteFormulaCmd(
		os.TempDir(),
		os.TempDir()+"/.rit",
		workspaceForm{},
		dirMock{},
		inputTrueMock{},
		inputTextCustomMock{text: func(name string, required bool) (string, error) {
			return workspace, nil
		}},
		inputListCustomMock{name: "Default (" + workspace + ")"},
		treeMock,
	)
	cmd.PersistentFlags().Bool("stdin", true, "input by stdin")
	if cmd == nil {
		t.Errorf("NewDeleteFormulaCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

type dirMock struct{}

func (d dirMock) List(_ string, _ bool) ([]string, error) {
	return []string{""}, nil
}
