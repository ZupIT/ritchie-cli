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
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type listMock struct{}

func (listMock) List(name string, items []string) (string, error) {
	workspace := filepath.Join(os.TempDir(), "ritchie-formulas-local")
	return fmt.Sprintf("Default (%s)", workspace), nil
}

func TestDeleteWorkspaceCmd(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	cmd := NewDeleteWorkspaceCmd(
		os.TempDir(),
		workspaceForm{},
		dirManager,
		listMock{},
		inputTrueMock{},
	)

	workspace := filepath.Join(os.TempDir(), "ritchie-formulas-local")
	if err := os.MkdirAll(filepath.Join(workspace, "mock", "test", "scr"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteWorkspaceCmd got error %v", err)
	}

	if cmd == nil {
		t.Errorf("NewDeleteWorkspaceCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
