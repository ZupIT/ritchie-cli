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
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewDeleteFormulaCmd(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
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
	}

	cmd := NewDeleteFormulaCmd(
		os.TempDir(),
		os.TempDir()+"/.rit",
		workspaceForm{},
		dirManager,
		inputTrueMock{},
		inputTextCustomWithoutValidateMock{text: os.TempDir() + "/ritchie-formulas-local"},
		inputListCustomMock{name: os.TempDir() + "/ritchie-formulas-local"},
		treeMock,
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewDeleteFormulaCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
