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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	ErrInvalidWorkspace = prompt.NewError("the formula workspace does not exist, please enter a valid workspace")

	hashesPath = "hashes"
	hashesExt  = ".txt"
)

type Manager struct {
	ritchieHome         string
	workspaceFile       string
	defaultWorkspaceDir string
	dir                 stream.DirCreateHasher
	local               builder.Initializer
}

func New(
	ritchieHome string,
	userHome string,
	dirManager stream.DirCreateHasher,
	local builder.Initializer,
) Manager {
	workspaceFile := filepath.Join(ritchieHome, formula.WorkspacesFile)
	workspaceHome := filepath.Join(userHome, formula.DefaultWorkspaceDir)
	return Manager{
		ritchieHome:         ritchieHome,
		workspaceFile:       workspaceFile,
		defaultWorkspaceDir: workspaceHome,
		dir:                 dirManager,
		local:               local,
	}
}

func (m Manager) Add(workspace formula.Workspace) error {
	if workspace.Dir == m.defaultWorkspaceDir {
		return nil
	}

	// Avoid finishing separators
	if last := len(workspace.Dir) - 1; last >= 0 && workspace.Dir[last] == filepath.Separator {
		workspace.Dir = workspace.Dir[:last]
	}
	if _, err := os.Stat(workspace.Dir); os.IsNotExist(err) {
		return ErrInvalidWorkspace
	}

	workspaces := formula.Workspaces{}
	file, err := ioutil.ReadFile(m.workspaceFile)
	if err == nil {
		if err := json.Unmarshal(file, &workspaces); err != nil {
			return err
		}
	}

	if _, ok := workspaces[workspace.Name]; ok {
		return nil
	}

	if _, err := m.local.Init(workspace.Dir, workspace.Name); err != nil {
		return err
	}

	workspaces[workspace.Name] = workspace.Dir
	content, err := json.MarshalIndent(workspaces, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(m.workspaceFile, content, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (m Manager) Delete(workspace formula.Workspace) error {
	workspaces, err := m.List()
	if err != nil {
		return err
	}

	if _, exists := workspaces[workspace.Name]; !exists {
		return ErrInvalidWorkspace
	}

	delete(workspaces, workspace.Name)
	content, err := json.MarshalIndent(workspaces, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(m.workspaceFile, content, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (m Manager) List() (formula.Workspaces, error) {
	workspaces := formula.Workspaces{}
	workspaces[formula.DefaultWorkspaceName] = m.defaultWorkspaceDir

	file, err := ioutil.ReadFile(m.workspaceFile)
	if err != nil {
		return workspaces, nil
	}

	if err := json.Unmarshal(file, &workspaces); err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (m Manager) PreviousHash(formulaPath string) (string, error) {
	filePath := m.hashPath(formulaPath)

	hash, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (m Manager) CurrentHash(formulaPath string) (string, error) {
	return m.dir.Hash(formulaPath)
}

func (m Manager) UpdateHash(formulaPath string, hash string) error {
	filePath := m.hashPath(formulaPath)

	hashDir := filepath.Join(m.ritchieHome, hashesPath)
	_ = m.dir.Create(hashDir)
	return ioutil.WriteFile(filePath, []byte(hash), os.ModePerm)
}

func (m Manager) hashPath(formulaPath string) string {
	fileName := strings.ReplaceAll(formulaPath, string(os.PathSeparator), "-") + hashesExt
	return filepath.Join(m.ritchieHome, hashesPath, fileName)
}
