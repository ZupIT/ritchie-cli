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

package renamer

import (
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
	dir stream.DirCreateCheckerCopy
}

func NewRenamer(
	dir stream.DirCreateCheckerCopy,
) RenameManager {
	return RenameManager{dir: dir}
}

func (r *RenameManager) Rename(fr formula.Rename) error {
	if err := r.createNewFormua(fr); err != nil {
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

func (r *RenameManager) createNewFormua(fr formula.Rename) error {
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
