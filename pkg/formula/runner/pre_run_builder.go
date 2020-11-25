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

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	messageBuilding        = "Building formula..."
	messageBuildSuccess    = "Build completed!\n"
	messageChangeError     = "Failed to detect formula changes, executing the last build"
	messageListError       = "Failed to load workspaces, ignoring the --force-build"
	messageForceBuildError = "Failed to build formula, ignoring the --force-build"
	messageBuildError      = "Failed to build formula"
	messageChangePrompt    = "This formula has changed since the last run, would you like to rebuild?"
	messageYes             = "yes"
	messageNo              = "no"
)

type PreRunBuilderManager struct {
	workspace formula.WorkspaceListHasher
	builder   formula.Builder
	dir       stream.DirRemoveChecker
	inBool    prompt.InputBool
	ritHome   string
}

func NewPreRunBuilder(
	workspace formula.WorkspaceListHasher,
	builder formula.Builder,
	dir stream.DirRemoveChecker,
	inBool prompt.InputBool,
	ritHome string,
) PreRunBuilderManager {
	return PreRunBuilderManager{
		workspace: workspace,
		builder:   builder,
		dir:       dir,
		inBool:    inBool,
		ritHome:   ritHome,
	}
}

func (b PreRunBuilderManager) ForceBuild(def formula.Definition) {
	workspace, err := b.anyWorkspace(def.Path)
	if err != nil {
		fmt.Println(prompt.Yellow(messageListError))
		return
	}

	// None of the workspaces have this formula, deleting
	// the bin to rebuild what's already on .rit
	if workspace == nil {
		formulaPath := def.FormulaPath(b.ritHome)
		binPath := def.BinPath(formulaPath)
		err = b.dir.Remove(binPath)
		if err != nil {
			fmt.Println(prompt.Yellow(messageForceBuildError))
		}
		return
	}

	// One or more workspaces have this formula, building from there
	if err = b.buildOnWorkspace(*workspace, def.Path); err != nil {
		fmt.Println(prompt.Red(messageBuildError))
		return
	}
}

func (b PreRunBuilderManager) Build(relativePath string) {
	workspace, err := b.modifiedWorkspace(relativePath)
	if err != nil {
		fmt.Println(prompt.Yellow(messageChangeError))
		return
	}

	// No modifications on any workspace, skip
	if workspace == nil {
		return
	}

	// User chose not to rebuild
	if !b.mustBuild() {
		return
	}

	if err = b.buildOnWorkspace(*workspace, relativePath); err != nil {
		fmt.Println(prompt.Red(messageBuildError))
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

func (b PreRunBuilderManager) anyWorkspace(relativePath string) (*formula.Workspace, error) {
	workspaces, err := b.workspace.List()
	if err != nil {
		return nil, err
	}

	for workspaceName, workspacePath := range workspaces {
		formulaAbsolutePath := filepath.Join(workspacePath, relativePath)

		if b.dir.Exists(formulaAbsolutePath) {
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
	s := spinner.StartNew(messageBuilding)
	info := formula.BuildInfo{FormulaPath: formulaAbsolutePath, Workspace: workspace}
	if err := b.builder.Build(info); err != nil {
		s.Error(err)
		return err
	}

	s.Success(prompt.Green(messageBuildSuccess))
	return nil
}

func (b PreRunBuilderManager) mustBuild() bool {
	ans, err := b.inBool.Bool(messageChangePrompt, []string{messageYes, messageNo})
	if err != nil {
		return false // Don't rebuild when Ctrl+C on question
	}

	return ans
}
