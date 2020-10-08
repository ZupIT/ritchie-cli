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
	"os"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
)

const (
	messageBuilding     = "Building formula..."
	messageBuildSuccess = "Build completed!\n"
	messageChangeError  = "Failed to detect formula changes, executing the last build"
	messageBuildError   = "Failed to build formula"
	messageChangePrompt = "This formula has changed since the last run, would you like to rebuild?"
	messageYes          = "yes"
	messageNo           = "no"

	sourceDir  = "src"
	hashesPath = "hashes"
	hashesExt  = ".txt"
)

type PreRunBuilderManager struct {
	ritchieHome string
	workspace   formula.WorkspaceLister
	builder     formula.LocalBuilder
	dir         stream.DirCreateHasher
	file        stream.FileReadWriter
	inBool      prompt.InputBool
}

func NewPreRunBuilder(
	ritchieHome string,
	workspace formula.WorkspaceLister,
	builder formula.LocalBuilder,
	dir stream.DirCreateHasher,
	file stream.FileReadWriter,
	inBool prompt.InputBool,
) PreRunBuilderManager {
	return PreRunBuilderManager{
		ritchieHome: ritchieHome,
		workspace:   workspace,
		builder:     builder,
		dir:         dir,
		file:        file,
		inBool:      inBool,
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

func (b PreRunBuilderManager) hasFormulaChanged(path string) (bool, error) {
	currentHash, err := b.currentHash(path)

	// Formula doesn't exist on this workspace
	if err != nil {
		return false, nil
	}

	previousHash, err := b.previousHash(path)
	if err != nil || previousHash != currentHash {
		updateErr := b.updateHash(path, currentHash)
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

func (b PreRunBuilderManager) previousHash(formulaPath string) (string, error) {
	filePath := b.hashPath(formulaPath)

	hash, err := b.file.Read(filePath)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (b PreRunBuilderManager) currentHash(formulaPath string) (string, error) {
	return b.dir.Hash(filepath.Join(formulaPath, sourceDir))
}

func (b PreRunBuilderManager) updateHash(formulaPath string, hash string) error {
	filePath := b.hashPath(formulaPath)

	hashDir := filepath.Join(b.ritchieHome, hashesPath)
	_ = b.dir.Create(hashDir)
	return b.file.Write(filePath, []byte(hash))
}

func (b PreRunBuilderManager) hashPath(formulaPath string) string {
	fileName := strings.ReplaceAll(formulaPath, string(os.PathSeparator), "-") + hashesExt
	return filepath.Join(b.ritchieHome, hashesPath, fileName)
}

func (b PreRunBuilderManager) buildOnWorkspace(workspace formula.Workspace, relativePath string) error {
	formulaAbsolutePath := filepath.Join(workspace.Dir, relativePath)
	s := spinner.StartNew(messageBuilding)
	if err := b.builder.Build(workspace.Dir, formulaAbsolutePath); err != nil {
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
