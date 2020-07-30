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

package workspace

import (
	"encoding/json"
	"path"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	ErrInvalidWorkspace = prompt.NewError("the formula workspace does not exist, please enter a valid workspace")
	ErrTreeJsonNotFound = prompt.NewError("tree.json not found, a valid formula workspace must have a tree.json")
	ErrMakefileNotFound = prompt.NewError("MakefilePath not found, a valid formula workspace must have a MakefilePath")
)

type Manager struct {
	workspaceFile string
	file          stream.FileWriteReadExister
}

func New(ritchieHome string, fileManager stream.FileWriteReadExister) Manager {
	workspaceFile := path.Join(ritchieHome, formula.WorkspacesFile)
	return Manager{workspaceFile: workspaceFile, file: fileManager}
}

func (m Manager) Add(workspace formula.Workspace) error {
	workspaces := formula.Workspaces{}
	if m.file.Exists(m.workspaceFile) {
		file, err := m.file.Read(m.workspaceFile)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(file, &workspaces); err != nil {
			return err
		}
	}

	workspaces[workspace.Name] = workspace.Dir
	content, err := json.Marshal(workspaces)
	if err != nil {
		return err
	}

	if err := m.file.Write(m.workspaceFile, content); err != nil {
		return err
	}

	return nil
}

func (m Manager) List() (formula.Workspaces, error) {
	workspaces := formula.Workspaces{}
	if !m.file.Exists(m.workspaceFile) {
		return formula.Workspaces{}, nil
	}

	file, err := m.file.Read(m.workspaceFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &workspaces); err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (m Manager) Validate(workspace formula.Workspace) error {
	dir := workspace.Dir
	if !m.file.Exists(dir) {
		return ErrInvalidWorkspace
	}

	return nil
}
