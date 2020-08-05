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

package builder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	msgBuildOnWindows = prompt.Yellow("This formula cannot be built on Windows.")
	ErrBuildOnWindows = errors.New(msgBuildOnWindows)
)

type LocalManager struct {
	ritHome string
	dir     stream.DirCreateListCopyRemover
	file    stream.FileWriteReadExister
	tree    formula.TreeGenerator
}

func NewBuildLocal(
	ritHome string,
	dir stream.DirCreateListCopyRemover,
	file stream.FileWriteReadExister,
	tree formula.TreeGenerator,
) formula.LocalBuilder {
	return LocalManager{ritHome: ritHome, dir: dir, file: file, tree: tree}
}

func (m LocalManager) Build(workspacePath, formulaPath string) error {

	dest := filepath.Join(m.ritHome, "repos", "local")

	if err := m.dir.Create(dest); err != nil {
		return err
	}

	if err := m.copyWorkSpace(workspacePath, dest); err != nil {
		return err
	}

	if err := m.generateTree(dest); err != nil {
		return err
	}

	if err := m.buildFormulaBin(workspacePath, formulaPath, dest); err != nil {
		return err
	}

	return nil
}

func (m LocalManager) buildFormulaBin(workspacePath, formulaPath, dest string) error {
	formulaSrc := strings.ReplaceAll(formulaPath, workspacePath, dest)
	formulaBin := filepath.Join(formulaSrc, "bin")

	if err := m.dir.Remove(formulaBin); err != nil {
		return err
	}

	if err := os.Chdir(formulaSrc); err != nil {
		return err
	}

	so := runtime.GOOS
	var cmd *exec.Cmd
	switch so {
	case osutil.Windows:
		winBuild := filepath.Join(formulaPath, "build.bat")
		if !m.file.Exists(winBuild) {
			return ErrBuildOnWindows
		}
		cmd = exec.Command(winBuild)
	default:
		cmd = exec.Command("make", "build")
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Bytes() != nil {
			errMsg := fmt.Sprintf("Build error: \n%s \n%s", stderr.String(), err)
			return errors.New(errMsg)
		}

		return err
	}

	return nil
}

func (m LocalManager) generateTree(dest string) error {
	tree, err := m.tree.Generate(dest)
	if err != nil {
		return err
	}

	treeFilePath := filepath.Join(dest, "tree.json")
	treeIndented, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		return err
	}

	if err := m.file.Write(treeFilePath, treeIndented); err != nil {
		return err
	}
	return nil
}

func (m LocalManager) copyWorkSpace(workspacePath string, dest string) error {
	return m.dir.Copy(workspacePath, dest)
}
