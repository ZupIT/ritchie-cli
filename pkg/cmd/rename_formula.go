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
	"errors"
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

	formulaOldCmdLabel  = "Enter the formula command to rename:"
	formulaOldCmdHelper = "Enter the existing formula in the informed workspace to rename it"
	formulaNewCmdLabel  = "Enter the new formula command:"
	formulaNewCmdHelper = "You must create your command based in this example [rit group verb noun]"

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
	inputFormula    prompt.InputFormula
	directory       stream.DirListChecker
	userHomeDir     string
	validator       validator.ValidatorManager
	renamer         renamer.RenameManager
}

// New renameFormulaCmd rename a cmd instance.
func NewRenameFormulaCmd(
	workspace formula.WorkspaceAddListHasher,
	inText prompt.InputText,
	inList prompt.InputList,
	inPath prompt.InputPath,
	inTextValidator prompt.InputTextValidator,
	inputFormula prompt.InputFormula,
	directory stream.DirListChecker,
	userHomeDir string,
	validator validator.ValidatorManager,
	renamer renamer.RenameManager,
) *cobra.Command {
	r := renameFormulaCmd{workspace, inText, inList, inPath, inTextValidator, inputFormula, directory, userHomeDir, validator,
		renamer}

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
		fmt.Println(result)

		if err := r.renamer.Rename(result); err != nil {
			return err
		}

		fmt.Println(result)

		return nil
	}
}

func (r *renameFormulaCmd) resolveInput(cmd *cobra.Command) (formula.Rename, error) {
	if IsFlagInput(cmd) {
		return r.resolveFlags(cmd)
	}
	return r.resolvePrompt()
}

func (r *renameFormulaCmd) resolveFlags(cmd *cobra.Command) (formula.Rename, error) {
	var result formula.Rename
	wspaces, err := r.workspace.List()
	if err != nil {
		return result, err
	}

	wsName, err := cmd.Flags().GetString(wsFlagName)
	if err != nil {
		return result, err
	} else if wsName == "" {
		return result, errors.New(missingFlagText(wsFlagName))
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
		return result, errors.New(missingFlagText(oldFormulaFlagName))
	}
	if !r.formulaExistsInWorkspace(result.Workspace.Dir, oldFormula) {
		return result, fmt.Errorf(ErrFormulaDontExists, oldFormula, result.Workspace.Name)
	}
	result.OldFormulaCmd = oldFormula

	newFormula, err := cmd.Flags().GetString(newFormulaFlagName)
	if err != nil {
		return result, err
	} else if newFormula == "" {
		return result, errors.New(missingFlagText(newFormulaFlagName))
	}
	if r.formulaExistsInWorkspace(result.Workspace.Dir, newFormula) {
		return result, fmt.Errorf(ErrFormulaExists, newFormula, result.Workspace.Name)
	}
	result.NewFormulaCmd = newFormula

	return result, nil
}

func (r *renameFormulaCmd) resolvePrompt() (formula.Rename, error) {
	var result formula.Rename
	wspaces, err := r.workspace.List()
	if err != nil {
		return result, err
	}

	ws, err := FormulaWorkspaceInput(wspaces, r.inList, r.inText, r.inPath)
	if err != nil {
		return result, err
	}
	result.Workspace = ws

	oldFormula, err := r.inTextValidator.Text(formulaOldCmdLabel, r.surveyCmdValidator, formulaOldCmdHelper)
	if err != nil {
		return result, err
	}
	if !r.formulaExistsInWorkspace(result.Workspace.Dir, oldFormula) {
		return result, fmt.Errorf(ErrFormulaDontExists, oldFormula, result.Workspace.Name)
	}
	result.OldFormulaCmd = oldFormula

	newFormula, err := r.inTextValidator.Text(formulaNewCmdLabel, r.surveyCmdValidator, formulaNewCmdHelper)
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

	path = filepath.Join(path, "src")

	return r.directory.Exists(path)
}

func cleanFormula(formula string) []string {
	formulaSplited := strings.Split(formula, " ")
	return formulaSplited[1:]
}

func (r *renameFormulaCmd) surveyCmdValidator(cmd interface{}) error {
	return r.validator.FormulaCommmandValidator(cmd.(string))
}
