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
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	ErrInvalidWorkspace = prompt.NewError("the formula workspace does not exist, please enter a valid workspace")
)

type Manager struct {
	workspaceFile 		string
	defaultWorkspaceDir string
	file           		stream.FileWriteReadExister
}

func New(ritchieHome string, userHome string, fileManager stream.FileWriteReadExister) Manager {
	workspaceFile := filepath.Join(ritchieHome, formula.WorkspacesFile)
	workspaceHome := filepath.Join(userHome, formula.DefaultWorkspaceDir)
	return Manager{workspaceFile: workspaceFile, defaultWorkspaceDir: workspaceHome, file: fileManager}
}

func (m Manager) Add(workspace formula.Workspace) error {
	if workspace.Dir == m.defaultWorkspaceDir {
		return nil
	}
	if !m.file.Exists(workspace.Dir) {
		return ErrInvalidWorkspace
	}

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

func (m Manager) Delete(workspace formula.Workspace) error {
	workspaces, err := m.List()
	if err != nil {
		return err
	}

	if _, exists := workspaces[workspace.Name]; exists {
		delete(workspaces, workspace.Name)
		content, err := json.Marshal(workspaces)
		if err != nil {
			return err
		}

		if err := m.file.Write(m.workspaceFile, content); err != nil {
			return err
		}

		return nil
	}

	return ErrInvalidWorkspace
}

func (m Manager) List() (formula.Workspaces, error) {
	workspaces := formula.Workspaces{}
	workspaces[formula.DefaultWorkspaceName] = m.defaultWorkspaceDir
	if !m.file.Exists(m.workspaceFile) {
		return workspaces, nil
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
