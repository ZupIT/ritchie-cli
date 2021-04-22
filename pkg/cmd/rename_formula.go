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

package cmd

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/renamer"
	"github.com/ZupIT/ritchie-cli/pkg/formula/validator"
	work "github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	wsFlagName         = "workspace"
	wsFlagDesc         = "name of workspace to rename"
	oldFormulaFlagName = "oldNameFormula"
	oldFormulaFlagDesc = "old name of formula to rename"
	newFormulaFlagName = "newNameFormula"
	newFormulaFlagDesc = "new name of formula to rename"

	ErrFormulaDontExists = "This formula '%s' dont's exists on this workspace = '%s'"
	ErrFormulaExists     = "This formula '%s' already exists on this workspace = '%s'"
)

var renameWorkspaceFlags = flags{
	{
		name:        wsFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: wsFlagDesc,
	},
	{
		name:        oldFormulaFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: oldFormulaFlagDesc,
	},
	{
		name:        newFormulaFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: newFormulaFlagDesc,
	},
}

// renameFormulaCmd type for add formula command.
type renameFormulaCmd struct {
	workspace       formula.WorkspaceAddListHasher
	inText          prompt.InputText
	inList          prompt.InputList
	inPath          prompt.InputPath
	inTextValidator prompt.InputTextValidator
	directory       stream.DirListChecker
	userHomeDir     string
	validator       validator.ValidatorManager
	inputFormula    prompt.InputFormula
	renamer         renamer.RenameManager
}

// New renameFormulaCmd rename a cmd instance.
func NewRenameFormulaCmd(
	workspace formula.WorkspaceAddListHasher,
	inText prompt.InputText,
	inList prompt.InputList,
	inPath prompt.InputPath,
	inTextValidator prompt.InputTextValidator,
	directory stream.DirListChecker,
	userHomeDir string,
	validator validator.ValidatorManager,
	inputFormula prompt.InputFormula,
	renamer renamer.RenameManager,
) *cobra.Command {
	r := renameFormulaCmd{
		workspace:       workspace,
		inText:          inText,
		inList:          inList,
		inPath:          inPath,
		inTextValidator: inTextValidator,
		directory:       directory,
		userHomeDir:     userHomeDir,
		validator:       validator,
		inputFormula:    inputFormula,
		renamer:         renamer,
	}

	cmd := &cobra.Command{
		Use:       "formula",
		Short:     "Rename a formula",
		Example:   "rit rename formula",
		RunE:      r.runFormula(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), renameWorkspaceFlags)

	return cmd
}

func (r *renameFormulaCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		result, err := r.resolveInput(cmd)
		if err != nil {
			return err
		}
		result.FOldPath = formulaPath(result.Workspace.Dir, result.OldFormulaCmd)
		result.FNewPath = formulaPath(result.Workspace.Dir, result.NewFormulaCmd)

		if err := r.renamer.Rename(result); err != nil {
			return err
		}

		// info := formula.BuildInfo{FormulaPath: cf.FormulaPath, Workspace: cf.Workspace}
		// if err := c.formula.Build(info); err != nil {
		// 	return err
		// }

		// hash, err := c.workspace.CurrentHash(cf.FormulaPath)
		// if err != nil {
		// 	return err
		// }

		// if err := c.workspace.UpdateHash(cf.FormulaPath, hash); err != nil {
		// 	return err
		// }

		fmt.Println(result.Workspace, result.NewFormulaCmd, result.OldFormulaCmd)

		return nil
	}
}

func (r *renameFormulaCmd) resolveInput(cmd *cobra.Command) (formula.Rename, error) {
	var result formula.Rename
	ws, err := r.workspace.List()
	if err != nil {
		return result, err
	}
	if IsFlagInput(cmd) {
		return r.resolveFlags(cmd, ws)
	}
	return r.resolvePrompt(ws)
}

func (r *renameFormulaCmd) resolveFlags(cmd *cobra.Command, wspaces formula.Workspaces) (formula.Rename, error) {
	// Default (/home/bruna/ritchie-formulas-local)
	// rit test sandokan
	var result formula.Rename
	flagError := "please provide a value for '%s'"

	wsName, err := cmd.Flags().GetString(wsFlagName)
	if err != nil {
		return result, err
	} else if wsName == "" {
		return result, fmt.Errorf(flagError, wsFlagName)
	}
	wspaces[formula.DefaultWorkspaceName] = filepath.Join(r.userHomeDir, formula.DefaultWorkspaceDir)
	dir, exists := wspaces[wsName]
	if !exists {
		return result, work.ErrInvalidWorkspace
	}
	result.Workspace.Dir = dir
	result.Workspace.Name = wsName

	oldFormula, err := cmd.Flags().GetString(oldFormulaFlagName)
	if err != nil {
		return result, err
	} else if oldFormula == "" {
		return result, fmt.Errorf(flagError, oldFormulaFlagName)
	}
	if !r.formulaExistsInWorkspace(result.Workspace.Dir, oldFormula) {
		return result, fmt.Errorf(ErrFormulaDontExists, oldFormula, result.Workspace.Name)
	}
	result.OldFormulaCmd = oldFormula

	newFormula, err := cmd.Flags().GetString(newFormulaFlagName)
	if err != nil {
		return result, err
	} else if newFormula == "" {
		return result, fmt.Errorf(flagError, newFormulaFlagName)
	}
	if r.formulaExistsInWorkspace(result.Workspace.Dir, newFormula) {
		return result, fmt.Errorf(ErrFormulaExists, newFormula, result.Workspace.Name)
	}
	result.NewFormulaCmd = newFormula

	return result, nil
}

func (r *renameFormulaCmd) resolvePrompt(workspaces formula.Workspaces) (formula.Rename, error) {
	var result formula.Rename

	ws, err := FormulaWorkspaceInput(workspaces, r.inList, r.inText, r.inPath)
	if err != nil {
		return result, err
	}
	result.Workspace = ws

	oldFormula, err := r.inputFormula.Select(ws.Dir, "rit")
	if err != nil {
		return result, err
	} else if oldFormula == "" {
		return result, ErrCouldNotFindFormula
	}
	result.OldFormulaCmd = oldFormula

	newFormula, err := r.inTextValidator.Text(formulaCmdLabel, r.surveyCmdValidator, formulaCmdHelper)
	if err != nil {
		return result, err
	}
	if r.formulaExistsInWorkspace(result.Workspace.Dir, newFormula) {
		return result, fmt.Errorf(ErrFormulaExists, newFormula, result.Workspace.Name)
	}
	result.NewFormulaCmd = newFormula

	return result, nil
}

func (r *renameFormulaCmd) formulaExistsInWorkspace(path string, formula string) bool {
	fc := cleanFormula(formula)
	for _, group := range fc {
		path = filepath.Join(path, group)
	}

	return r.directory.Exists(path)
}

func cleanFormula(formula string) []string {
	formulaSplited := strings.Split(formula, " ")
	return formulaSplited[1:]
}

func (r *renameFormulaCmd) surveyCmdValidator(cmd interface{}) error {
	if err := r.validator.FormulaCommmandValidator(cmd.(string)); err != nil {
		return err
	}

	return nil
}
