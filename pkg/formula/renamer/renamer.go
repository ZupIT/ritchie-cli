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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implier.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package renamer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	ErrRepeatedCommand = prompt.NewError("this command already exists")
)

type RenameManager struct {
	dir  stream.DirCreateCheckerCopy
	file stream.FileWriteRemover
}

func NewRenamer(
	dir stream.DirCreateCheckerCopy,
	file stream.FileWriteRemover,
) RenameManager {
	return RenameManager{dir, file}
}

func (r *RenameManager) Rename(fr formula.Rename) error {
	if err := r.createNewFormula(fr); err != nil {
		return err
	}

	groups := strings.Split(fr.OldFormulaCmd, " ")[1:]
	fmt.Println(fr.Workspace.Dir, groups, 0)
	if err := r.deleteFormula(fr.Workspace.Dir, groups, 0); err != nil {
		return err
	}

	return nil
}

func (r *RenameManager) createNewFormula(fr formula.Rename) error {
	fOldPath := formulaPath(fr.Workspace.Dir, fr.OldFormulaCmd)
	fNewPath := formulaPath(fr.Workspace.Dir, fr.NewFormulaCmd)

	if err := r.isAvailableCmd(fNewPath); err != nil {
		return err
	}

	if err := r.dir.Create(fNewPath); err != nil {
		return err
	}

	if err := r.dir.Copy(fOldPath, fNewPath); err != nil {
		return err
	}

	return nil
}

func (r *RenameManager) isAvailableCmd(fPath string) error {
	if r.dir.Exists(fPath) {
		return ErrRepeatedCommand
	}

	return nil
}

func formulaPath(workspacePath, cmd string) string {
	cc := strings.Split(cmd, " ")
	formulaPath := strings.Join(cc[1:], string(os.PathSeparator))
	return filepath.Join(workspacePath, formulaPath)
}

func (r *RenameManager) deleteFormula(path string, groups []string, index int) error {
	if index == len(groups) {
		nested, err := nestedFormula(path)
		if err != nil {
			return err
		}

		if nested {
			return r.safeRemoveFormula(path)
		}

		return os.RemoveAll(path)
	}

	newPath := filepath.Join(path, groups[index])
	if err := r.deleteFormula(newPath, groups, index+1); err != nil {
		return err
	}

	if index == 0 {
		return nil
	}

	ok, err := canDelete(path)
	if err != nil {
		return err
	}

	if ok {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}

func (r *RenameManager) safeRemoveFormula(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() && (file.Name() == "src" || file.Name() == "bin") {
			pathToDelete := filepath.Join(path, file.Name())
			if err := os.RemoveAll(pathToDelete); err != nil {
				return err
			}
		} else if !file.IsDir() {
			pathToDelete := filepath.Join(path, file.Name())
			if err := r.file.Remove(pathToDelete); err != nil {
				return err
			}
		}
	}

	return nil
}

func nestedFormula(path string) (bool, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != "src" && file.Name() != "bin" {
			return true, nil
		}
	}

	return false, nil
}

func canDelete(path string) (bool, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if file.IsDir() {
			return false, nil
		}
	}

	return true, nil
}
