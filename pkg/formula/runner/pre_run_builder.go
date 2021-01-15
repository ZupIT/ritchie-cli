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

package runner

import (
	"fmt"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	messageChangeError = "Failed to detect formula changes, executing the last build: %s"
	messageBuildError  = "Failed to build formula: %s"
)

type PreRunBuilderManager struct {
	workspace formula.WorkspaceListHasher
	builder   formula.Builder
}

func NewPreRunBuilder(
	workspace formula.WorkspaceListHasher,
	builder formula.Builder,
) PreRunBuilderManager {
	return PreRunBuilderManager{
		workspace: workspace,
		builder:   builder,
	}
}

func (b PreRunBuilderManager) Build(relativePath string) {
	workspace, err := b.modifiedWorkspace(relativePath)
	if err != nil {
		msg := fmt.Sprintf(messageChangeError, err.Error())
		fmt.Println(prompt.Yellow(msg))
		return
	}

	// No modifications on any workspace, skip
	if workspace == nil {
		return
	}

	if err = b.buildOnWorkspace(*workspace, relativePath); err != nil {
		msg := fmt.Sprintf(messageBuildError, err.Error())
		fmt.Println(prompt.Red(msg))
		return
	}
}

func (b PreRunBuilderManager) modifiedWorkspace(relativePath string) (*formula.Workspace, error) {
	workspaces, err := b.workspace.List()
	if err != nil {
		return nil, err
	}

	for workspaceName, workspacePath := range workspaces {
		formulaAbsolutePath := filepath.Join(workspacePath, relativePath)
		hasChanged, err := b.hasFormulaChanged(formulaAbsolutePath)
		if err != nil {
			return nil, err
		}
		if hasChanged {
			return &formula.Workspace{
				Name: workspaceName,
				Dir:  workspacePath,
			}, nil
		}
	}

	return nil, nil
}

func (b PreRunBuilderManager) hasFormulaChanged(path string) (bool, error) {
	currentHash, err := b.workspace.CurrentHash(path)

	// Formula doesn't exist on this workspace
	if err != nil {
		return false, nil
	}

	previousHash, err := b.workspace.PreviousHash(path)
	if err != nil || previousHash != currentHash {
		updateErr := b.workspace.UpdateHash(path, currentHash)
		if updateErr != nil {
			return false, updateErr
		}
	}

	// First time hashing this formula
	if err != nil {
		return false, nil
	}

	return previousHash != currentHash, nil
}

func (b PreRunBuilderManager) buildOnWorkspace(workspace formula.Workspace, relativePath string) error {
	formulaAbsolutePath := filepath.Join(workspace.Dir, relativePath)
	info := formula.BuildInfo{FormulaPath: formulaAbsolutePath, Workspace: workspace}
	if err := b.builder.Build(info); err != nil {
		return err
	}

	return nil
}
