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
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type listErrorMock struct{}

func (listErrorMock) List(name string, items []string) (string, error) {
	workspace := filepath.Join(os.TempDir(), "formulas-ritchie")
	return fmt.Sprintf("Formulas-Ritchie (%s)", workspace), nil
}

type listMock struct{}

func (listMock) List(name string, items []string) (string, error) {
	workspace := filepath.Join(os.TempDir(), "ritchie-formulas-local")
	return fmt.Sprintf("Default (%s)", workspace), nil
}

type workspaceMock struct{}

func (workspaceMock) List() (formula.Workspaces, error) {
	m := formula.Workspaces{"Formulas-Ritchie": filepath.Join(os.TempDir(), "formulas-ritchie")}
	return m, nil
}

func (workspaceMock) Delete(workspace formula.Workspace) error {
	return errors.New("Some error")
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

	workspace := filepath.Join(os.TempDir(), "formulas-ritchie")
	if err := os.MkdirAll(filepath.Join(workspace, "mock", "test", "scr"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteWorkspaceCmd got error %v", err)
	}

	if cmd == nil {
		t.Errorf("NewDeleteWorkspaceCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}

	cmd = NewDeleteWorkspaceCmd(
		os.TempDir(),
		workspaceForm{},
		dirManager,
		inputListCustomMock{
			list: func(name string, items []string) (string, error) {
				return "", errors.New("Some error")
			},
		},
		inputTrueMock{},
	)

	workspace = filepath.Join(os.TempDir(), "ritchie-formulas-local")
	if err := os.MkdirAll(filepath.Join(workspace, "mock", "test", "scr"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteWorkspaceCmd got error %v", err)
	}

	if cmd == nil {
		t.Errorf("NewDeleteWorkspaceCmd got %v", cmd)
	}

	if err := cmd.Execute(); err == nil {
		t.Errorf("%s = %v, want %v", cmd.Use, nil, err)
	}

	cmd = NewDeleteWorkspaceCmd(
		os.TempDir(),
		workspaceForm{},
		DirManagerCustomMock{
			exists: func(dir string) bool {
				return false
			},
		},
		listMock{},
		inputTrueMock{},
	)

	workspace = filepath.Join(os.TempDir(), "ritchie-formulas-local")
	if err := os.MkdirAll(filepath.Join(workspace, "mock", "test", "scr"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteWorkspaceCmd got error %v", err)
	}

	if cmd == nil {
		t.Errorf("NewDeleteWorkspaceCmd got %v", cmd)
	}

	if err := cmd.Execute(); err == nil {
		t.Errorf("%s = %v, want %v", cmd.Use, nil, err)
	}

	cmd = NewDeleteWorkspaceCmd(
		os.TempDir(),
		workspaceMock{},
		dirManager,
		listErrorMock{},
		inputTrueMock{},
	)

	workspace = filepath.Join(os.TempDir(), "formulas-ritchie")
	if err := os.MkdirAll(filepath.Join(workspace, "mock", "test", "scr"), os.ModePerm); err != nil {
		t.Errorf("TestNewDeleteWorkspaceCmd got error %v", err)
	}

	if cmd == nil {
		t.Errorf("NewDeleteWorkspaceCmd got %v", cmd)
	}

	if err := cmd.Execute(); err == nil {
		t.Errorf("%s = %v, want %v", cmd.Use, nil, err)
	}
}
