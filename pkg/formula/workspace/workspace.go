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
	"reflect"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo/repoutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const invalidCharacters = `\/><,@#%!&*()=+§£¢¬ªº°"^~;.?`

var (
	ErrInvalidWorkspace         = prompt.NewError("the formula workspace does not exist, please enter a valid workspace")
	ErrInvalidWorkspaceName     = prompt.NewError(`the workspace name must not contain spaces or invalid characters (\/><,@#%!&*()=+§£¢¬ªº°"^~;.?)`)
	ErrInvalidWorkspaceNameType = prompt.NewError("the input type is invalid for the workspace name")

	hashesPath = "hashes"
	hashesExt  = ".txt"
)

type Manager struct {
	ritchieHome         string
	workspaceFile       string
	defaultWorkspaceDir string
	dir                 stream.DirCreateHasher
	local               builder.Initializer
	tree                formula.TreeGenerator
}

func New(
	ritchieHome string,
	userHome string,
	dirManager stream.DirCreateHasher,
	local builder.Initializer,
	tree formula.TreeGenerator,
) Manager {
	workspaceFile := filepath.Join(ritchieHome, formula.WorkspacesFile)
	workspaceHome := filepath.Join(userHome, formula.DefaultWorkspaceDir)
	return Manager{
		ritchieHome:         ritchieHome,
		workspaceFile:       workspaceFile,
		defaultWorkspaceDir: workspaceHome,
		dir:                 dirManager,
		local:               local,
		tree:                tree,
	}
}

func (m Manager) Add(workspace formula.Workspace) error {
	if workspace.Dir == m.defaultWorkspaceDir {
		return nil
	}

	err := WorkspaceNameValidator(workspace.Name)
	if err != nil {
		return err
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

func (m Manager) Update(workspace formula.Workspace) error {
	workspaces, err := m.List()
	if err != nil {
		return err
	}

	if _, exists := workspaces[workspace.Name]; !exists {
		return ErrInvalidWorkspace
	}

	workspaceLocalName := repoutil.LocalName(workspace.Name)
	workflowRitFolderPath := filepath.Join(m.ritchieHome, formula.ReposDir, workspaceLocalName.String())

	delete(workspaces, workspace.Name)

	if _, err := m.local.Init(workspace.Dir, workspace.Name); err != nil {
		return err
	}

	treeData, err := m.tree.Generate(workflowRitFolderPath)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(treeData, "", "\t")
	if err != nil {
		return err
	}

	treeFilePath := filepath.Join(workflowRitFolderPath, tree.FileName)
	if err := ioutil.WriteFile(treeFilePath, bytes, os.ModePerm); err != nil {
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

func WorkspaceNameValidator(cmd interface{}) error {
	if reflect.TypeOf(cmd).Kind() != reflect.String {
		return ErrInvalidWorkspaceNameType
	}

	workspaceName := cmd.(string)
	isWithSpaces := strings.Contains(workspaceName, " ")
	isWithInvalidCharacters := strings.ContainsAny(workspaceName, invalidCharacters)
	if isWithSpaces || isWithInvalidCharacters {
		return ErrInvalidWorkspaceName
	}
	return nil
}
