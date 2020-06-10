package workspace

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	msgErrMakefileNotFound = fmt.Sprintf(prompt.Error, "Makefile not found, a valid formula workspace must have a Makefile")
	msgErrTreeJsonNotFound = fmt.Sprintf(prompt.Error, "tree.json not found, a valid formula workspace must have a tree.json")
	msgErrInvalidWorkspace = fmt.Sprintf(prompt.Error, "the formula workspace does not exist, please enter a valid workspace")
	ErrInvalidWorkspace    = errors.New(msgErrInvalidWorkspace)
	ErrTreeJsonNotFound    = errors.New(msgErrTreeJsonNotFound)
	ErrMakefileNotFound    = errors.New(msgErrMakefileNotFound)
)

type Manager struct {
	workspaceFile string
	file          stream.FileWriteReadExister
}

func New(ritchieHome string, fileManager stream.FileWriteReadExister) Manager {
	workspaceFile := fmt.Sprintf(workspacesPattern, ritchieHome)
	return Manager{workspaceFile: workspaceFile, file: fileManager}
}

func (m Manager) Add(workspace Workspace) error {
	if err := m.Validate(workspace); err != nil {
		return err
	}

	workspaces := Workspaces{}
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

func (m Manager) List() (Workspaces, error) {
	workspaces := Workspaces{}
	if !m.file.Exists(m.workspaceFile) {
		return Workspaces{}, nil
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

func (m Manager) Validate(workspace Workspace) error {
	dir := workspace.Dir
	if !m.file.Exists(dir) {
		return ErrInvalidWorkspace
	}

	makefilePath := fmt.Sprintf(formula.MakefileCreatePathPattern, dir, formula.Makefile)
	if !m.file.Exists(makefilePath) {
		return ErrMakefileNotFound
	}

	treePath := fmt.Sprintf(formula.TreeCreatePathPattern, dir)
	if !m.file.Exists(treePath) {
		return ErrTreeJsonNotFound
	}

	return nil
}
