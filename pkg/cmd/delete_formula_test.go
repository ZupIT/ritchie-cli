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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewDeleteFormulaCmdStdin(t *testing.T) {
	workspace := filepath.Join(os.TempDir(), "ritchie-formulas-local")
	if err := os.MkdirAll(filepath.Join(workspace, "mock", "test", "scr"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaCmdStdin got error %v", err)
	}

	if err := os.MkdirAll(filepath.Join(os.TempDir(), ".rit", "repos", "local", "mock", "test", "scr"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteFormulaCmdStdin got error %v", err)
	}

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	treeGen := tree.NewGenerator(dirManager, fileManager)

	cmd := NewDeleteFormulaCmd(
		os.TempDir(),
		filepath.Join(os.TempDir(), ".rit"),
		workspaceForm{},
		dirMock{},
		inputTrueMock{},
		inputTextMock{},
		inputListMock{},
		treeGen,
	)
	cmd.PersistentFlags().Bool("stdin", true, "input by stdin")

	json := fmt.Sprintf("{\"workspace\": \"%s\", \"groups\": [\"mock\", \"test\"]}\n", workspace)
	fmt.Println("JSON: " + json)
	newReader := strings.NewReader(json)
	cmd.SetIn(newReader)

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
