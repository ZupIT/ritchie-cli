package workspace

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
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
