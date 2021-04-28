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
	"os"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/deleter"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	ErrRepeatedCommand = prompt.NewError("this command already exists")
)

type RenameManager struct {
	dir        stream.DirCreateCheckerCopy
	file       stream.FileWriteRemover
	formula    formula.CreateBuilder
	workspace  formula.WorkspaceHasher
	ritHomeDir string
	treeGen    formula.TreeGenerator
	deleter    deleter.DeleteManager
}

func NewRenamer(
	dir stream.DirCreateCheckerCopy,
	file stream.FileWriteRemover,
	formula formula.CreateBuilder,
	workspace formula.WorkspaceHasher,
	ritHomeDir string,
	treeGen formula.TreeGenerator,
	deleter deleter.DeleteManager,
) RenameManager {
	return RenameManager{dir, file, formula, workspace, ritHomeDir, treeGen, deleter}
}

func (r *RenameManager) Rename(fr formula.Rename) error {
	fr.NewFormulaCmd = cleanSuffix(fr.NewFormulaCmd)
	fr.OldFormulaCmd = cleanSuffix(fr.OldFormulaCmd)

	if err := r.renameFormula(fr); err != nil {
		return err
	}

	info := formula.BuildInfo{FormulaPath: fr.FNewPath, Workspace: fr.Workspace}
	if err := r.formula.Build(info); err != nil {
		return err
	}

	hashNew, err := r.workspace.CurrentHash(fr.FNewPath)
	if err != nil {
		return err
	}

	if err := r.workspace.UpdateHash(fr.FNewPath, hashNew); err != nil {
		return err
	}

	return nil
}

func (r *RenameManager) renameFormula(fr formula.Rename) error {
	fOldPath := formulaPath(fr.Workspace.Dir, fr.OldFormulaCmd)
	fNewPath := formulaPath(fr.Workspace.Dir, fr.NewFormulaCmd)

	if err := r.isAvailableCmd(fNewPath); err != nil {
		return err
	}

	tmp := filepath.Join(os.TempDir(), "rit_oldFormula")
	if err := r.dir.Create(tmp); err != nil {
		return err
	}

	if err := r.dir.Copy(fOldPath, tmp); err != nil {
		return err
	}

	groupsOld := strings.Split(fr.OldFormulaCmd, " ")[1:]
	delOld := formula.Delete{
		GroupsFormula: groupsOld,
		Workspace:     fr.Workspace,
	}
	if err := r.deleter.Delete(delOld); err != nil {
		return err
	}

	if err := r.dir.Create(fNewPath); err != nil {
		return err
	}

	if err := r.dir.Copy(tmp, fNewPath); err != nil {
		return err
	}

	if err := os.RemoveAll(tmp); err != nil {
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

func cleanSuffix(cmd string) string {
	if strings.HasSuffix(cmd, "rit") {
		return cmd[4:]
	}
	return cmd
}
