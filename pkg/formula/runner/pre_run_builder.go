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
	"runtime"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
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

func (b PreRunBuilderManager) Build(formulaRelativePath string) {
	workspace, err := b.getModifiedWorkspace(formulaRelativePath)
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

	if err = b.buildOnWorkspace(*workspace, formulaRelativePath); err != nil {
		fmt.Println(prompt.Red(messageBuildError))
		return
	}
}

func (b PreRunBuilderManager) getModifiedWorkspace(formulaRelativePath string) (*formula.Workspace, error) {
	workspaces, err := b.workspace.List()
	if err != nil {
		return nil, err
	}

	for workspaceName, workspacePath := range workspaces {
		formulaAbsolutePath := filepath.Join(workspacePath, formulaRelativePath)
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
	currentHash, err := b.getCurrentHash(path)

	// Formula doesn't exist on this workspace
	if err != nil {
		return false, nil
	}

	previousHash, err := b.getPreviousHash(path)
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

func (b PreRunBuilderManager) getPreviousHash(formulaPath string) (string, error) {
	filePath := b.getHashPath(formulaPath)

	hash, err := b.file.Read(filePath)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (b PreRunBuilderManager) getCurrentHash(formulaPath string) (string, error) {
	return b.dir.Hash(filepath.Join(formulaPath, sourceDir))
}

func (b PreRunBuilderManager) updateHash(formulaPath string, hash string) error {
	filePath := b.getHashPath(formulaPath)

	_ = b.dir.Create(b.getHashDir())
	return b.file.Write(filePath, []byte(hash))
}

func (b PreRunBuilderManager) getHashDir() string {
	return filepath.Join(b.ritchieHome, hashesPath)
}

func (b PreRunBuilderManager) getHashPath(formulaPath string) string {
	divider := "/"
	if runtime.GOOS == osutil.Windows {
		divider = "\\"
	}

	fileName := strings.ReplaceAll(formulaPath, divider, "-") + hashesExt
	return filepath.Join(b.ritchieHome, hashesPath, fileName)
}

func (b PreRunBuilderManager) buildOnWorkspace(workspace formula.Workspace, formulaRelativePath string) error {
	formulaAbsolutePath := filepath.Join(workspace.Dir, formulaRelativePath)
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
