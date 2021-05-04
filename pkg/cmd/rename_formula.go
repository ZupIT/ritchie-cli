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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/deleter"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo/repoutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula/validator"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	wsFlagName         = "workspace"
	wsFlagDesc         = "Name of workspace to rename"
	oldFormulaFlagName = "oldName"
	oldFormulaFlagDesc = "Old name of formula to rename"
	newFormulaFlagName = "newName"
	newFormulaFlagDesc = "New name of formula to rename"

	formulaOldCmdLabel  = "Enter the formula command to rename:"
	formulaOldCmdHelper = "Enter the existing formula in the informed workspace to rename it"
	formulaNewCmdLabel  = "Enter the new formula command:"
	formulaNewCmdHelper = "You must create your command based in this example [rit group verb noun]"

	questionConfirmation = "Are you sure you want to rename the formula from '%s' to '%s'?"

	errNonExistFormula   = "This formula '%s' does not exist on this workspace = '%s'"
	errFormulaExists     = "This formula '%s' already exists on this workspace = '%s'"
	errNonExistWorkspace = "The formula workspace '%s' does not exist, please enter a valid workspace"

	renameSuccessMsg = "The formula was renamed with success"
)

var renameWorkspaceFlags = flags{
	{
		name:        wsFlagName,
		kind:        reflect.String,
		defValue:    formula.DefaultWorkspaceName,
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
	inBool          prompt.InputBool
	directory       stream.DirManager
	validator       validator.ValidatorManager
	formula         formula.CreateBuilder
	treeGen         formula.TreeGenerator
	deleter         deleter.DeleteManager
	userHomeDir     string
	ritHomeDir      string
}

// New renameFormulaCmd rename a cmd instance.
func NewRenameFormulaCmd(
	workspace formula.WorkspaceAddListHasher,
	inText prompt.InputText,
	inList prompt.InputList,
	inPath prompt.InputPath,
	inTextValidator prompt.InputTextValidator,
	inBool prompt.InputBool,
	directory stream.DirManager,
	validator validator.ValidatorManager,
	formula formula.CreateBuilder,
	treeGen formula.TreeGenerator,
	deleter deleter.DeleteManager,
	userHomeDir string,
	ritHomeDir string,

) *cobra.Command {
	r := renameFormulaCmd{
		workspace:       workspace,
		inText:          inText,
		inList:          inList,
		inPath:          inPath,
		inTextValidator: inTextValidator,
		inBool:          inBool,
		directory:       directory,
		validator:       validator,
		formula:         formula,
		treeGen:         treeGen,
		deleter:         deleter,
		userHomeDir:     userHomeDir,
		ritHomeDir:      ritHomeDir,
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

		question := fmt.Sprintf(questionConfirmation, result.OldFormulaCmd, result.NewFormulaCmd)
		ans, err := r.inBool.Bool(question, []string{"no", "yes"})
		if err != nil {
			return err
		}
		if !ans {
			return nil
		}

		result.FOldPath = fPath(result.Workspace.Dir, result.OldFormulaCmd)
		result.FNewPath = fPath(result.Workspace.Dir, result.NewFormulaCmd)

		if err := r.Rename(result); err != nil {
			return err
		}

		prompt.Success(renameSuccessMsg)

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
		wsName = formula.DefaultWorkspaceName
	}
	wspaces[formula.DefaultWorkspaceName] = filepath.Join(r.userHomeDir, formula.DefaultWorkspaceDir)
	dir, exists := wspaces[wsName]
	if !exists {
		return result, fmt.Errorf(errNonExistWorkspace, wsName)
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
		return result, fmt.Errorf(errNonExistFormula, oldFormula, result.Workspace.Name)
	}
	result.OldFormulaCmd = oldFormula

	newFormula, err := cmd.Flags().GetString(newFormulaFlagName)
	if err != nil {
		return result, err
	} else if newFormula == "" {
		return result, errors.New(missingFlagText(newFormulaFlagName))
	}
	if r.formulaExistsInWorkspace(result.Workspace.Dir, newFormula) {
		return result, fmt.Errorf(errFormulaExists, newFormula, result.Workspace.Name)
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
		return result, fmt.Errorf(errNonExistFormula, oldFormula, result.Workspace.Name)
	}
	result.OldFormulaCmd = oldFormula

	newFormula, err := r.inTextValidator.Text(formulaNewCmdLabel, r.surveyCmdValidator, formulaNewCmdHelper)
	if err != nil {
		return result, err
	}
	if r.formulaExistsInWorkspace(result.Workspace.Dir, newFormula) {
		return result, fmt.Errorf(errFormulaExists, newFormula, result.Workspace.Name)
	}
	result.NewFormulaCmd = newFormula

	return result, nil
}

func (r *renameFormulaCmd) formulaExistsInWorkspace(path string, formula string) bool {
	formulaSplited := strings.Split(formula, " ")
	if formulaSplited[0] == "rit" {
		formulaSplited = formulaSplited[1:]
	}

	for _, group := range formulaSplited {
		path = filepath.Join(path, group)
	}

	path = filepath.Join(path, "src")

	return r.directory.Exists(path)
}

func (r *renameFormulaCmd) surveyCmdValidator(cmd interface{}) error {
	return r.validator.FormulaCommmandValidator(cmd.(string))
}

func (r *renameFormulaCmd) Rename(fr formula.Rename) error {
	fr.NewFormulaCmd = cleanSuffix(fr.NewFormulaCmd)
	fr.OldFormulaCmd = cleanSuffix(fr.OldFormulaCmd)

	if err := r.changeFormulaToNewDir(fr); err != nil {
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

	repoNameStandard := repoutil.LocalName(fr.Workspace.Name)
	repoNameStandardPath := filepath.Join(r.ritHomeDir, "repos", repoNameStandard.String())
	if err := r.recreateTreeJSON(repoNameStandardPath); err != nil {
		return err
	}
	return nil
}

func (r *renameFormulaCmd) changeFormulaToNewDir(fr formula.Rename) error {
	fOldPath := fPath(fr.Workspace.Dir, fr.OldFormulaCmd)
	fNewPath := fPath(fr.Workspace.Dir, fr.NewFormulaCmd)

	tmp := filepath.Join(os.TempDir(), "rit_oldFormula")
	if err := r.directory.Create(tmp); err != nil {
		return err
	}

	if err := r.directory.Copy(fOldPath, tmp); err != nil {
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

	if err := r.directory.Create(fNewPath); err != nil {
		return err
	}

	if err := r.directory.Copy(tmp, fNewPath); err != nil {
		return err
	}

	if err := os.RemoveAll(tmp); err != nil {
		return err
	}

	return nil
}

func (r *renameFormulaCmd) recreateTreeJSON(pathLocalWS string) error {
	localTree, err := r.treeGen.Generate(pathLocalWS)
	if err != nil {
		return err
	}

	jsonString, _ := json.MarshalIndent(localTree, "", "\t")
	pathLocalTreeJSON := filepath.Join(pathLocalWS, "tree.json")
	if err := ioutil.WriteFile(pathLocalTreeJSON, jsonString, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func fPath(workspacePath, cmd string) string {
	cc := strings.Split(cmd, " ")
	path := strings.Join(cc[1:], string(os.PathSeparator))
	return filepath.Join(workspacePath, path)
}

func cleanSuffix(cmd string) string {
	if strings.HasSuffix(cmd, "rit") {
		return cmd[4:]
	}
	return cmd
}
