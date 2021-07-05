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

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/deleter"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo/repoutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula/validator"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	oldFormulaFlagName = "oldName"
	oldFormulaFlagDesc = "Old name of formula to rename"
	newFormulaFlagName = "newName"
	newFormulaFlagDesc = "New name of formula to rename"

	formulaOldCmdLabel  = "Enter the formula command to rename:"
	formulaOldCmdHelper = "Enter the existing formula in the informed workspace to rename it"
	formulaNewCmdLabel  = "Enter the new formula command:"
	formulaNewCmdHelper = "You must create your command based in this example [rit group verb noun]"

	questionConfirmation = "Are you sure you want to rename the formula from %q to %q?"

	errNonExistFormula = "formula %q wasn't found in the workspaces"
	errFormulaExists   = "formula %q already exists on this workspace = %q"
	errFormulaInManyWS = "formula %q was found in %d workspaces. Please enter a value for the 'workspace' flag"
	errInvalidWS       = "workspace %q was not found"
	renameSuccessMsg   = "The formula was renamed with success"
)

var renameWorkspaceFlags = flags{
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
	{
		name:        workspaceFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: workspaceFlagDesc,
	},
}

// renameFormulaCmd type for add formula command.
type renameFormulaCmd struct {
	workspace       formula.WorkspaceAddListHasher
	inList          prompt.InputList
	inTextValidator prompt.InputTextValidator
	inBool          prompt.InputBool
	directory       stream.DirManager
	validator       validator.Manager
	formula         formula.CreateBuilder
	treeGen         formula.TreeGenerator
	deleter         deleter.DeleteManager
	userHomeDir     string
	ritHomeDir      string
}

// New renameFormulaCmd rename a cmd instance.
func NewRenameFormulaCmd(
	workspace formula.WorkspaceAddListHasher,
	inList prompt.InputList,
	inTextValidator prompt.InputTextValidator,
	inBool prompt.InputBool,
	directory stream.DirManager,
	validator validator.Manager,
	formula formula.CreateBuilder,
	treeGen formula.TreeGenerator,
	deleter deleter.DeleteManager,
	userHomeDir string,
	ritHomeDir string,

) *cobra.Command {
	r := renameFormulaCmd{
		workspace:       workspace,
		inList:          inList,
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

		if isEmpty(result) {
			return nil
		}

		wsOldFormula := r.getWorkspace(result.OldFormulaCmd)
		wsNewFormula := r.getWorkspace(result.NewFormulaCmd)

		wspaceName := result.Workspace.Name
		ws, err := r.cleanWorkspace(wsOldFormula, wsNewFormula, result, cmd, wspaceName)
		if err != nil {
			return err
		}
		result.Workspace.Dir = ws.Dir
		result.Workspace.Name = ws.Name

		result.NewFormulaCmd = cleanPreffix(result.NewFormulaCmd)
		result.OldFormulaCmd = cleanPreffix(result.OldFormulaCmd)

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

	oldFormula, err := cmd.Flags().GetString(oldFormulaFlagName)
	if err != nil {
		return result, err
	} else if oldFormula == "" {
		return result, errors.New(missingFlagText(oldFormulaFlagName))
	}
	result.OldFormulaCmd = oldFormula

	newFormula, err := cmd.Flags().GetString(newFormulaFlagName)
	if err != nil {
		return result, err
	} else if newFormula == "" {
		return result, errors.New(missingFlagText(newFormulaFlagName))
	}
	result.NewFormulaCmd = newFormula

	workspaceName, err := cmd.Flags().GetString(workspaceFlagName)
	if err != nil {
		return result, err
	}
	result.Workspace.Name = workspaceName

	return result, nil
}

func (r *renameFormulaCmd) resolvePrompt() (formula.Rename, error) {
	var result formula.Rename

	oldFormula, err := r.inTextValidator.Text(formulaOldCmdLabel, r.surveyCmdValidator, formulaOldCmdHelper)
	if err != nil {
		return result, err
	}
	result.OldFormulaCmd = oldFormula

	newFormula, err := r.inTextValidator.Text(formulaNewCmdLabel, r.surveyCmdValidator, formulaNewCmdHelper)
	if err != nil {
		return result, err
	}
	result.NewFormulaCmd = newFormula

	question := fmt.Sprintf(questionConfirmation, result.OldFormulaCmd, result.NewFormulaCmd)
	ans, err := r.inBool.Bool(question, []string{"no", "yes"})
	if err != nil {
		return result, err
	}
	if !ans {
		return formula.Rename{}, nil
	}

	return result, nil
}

func (r *renameFormulaCmd) formulaExistsInWorkspace(path string, formula string) bool {
	formulaSplited := strings.Split(formula, " ")
	if formulaSplited[0] == api.RootName {
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

func (r *renameFormulaCmd) getWorkspace(cmd string) formula.Workspaces {
	wsWithFormula := formula.Workspaces{}
	wspaces, err := r.workspace.List()
	if err != nil {
		return wsWithFormula
	}
	wspaces[formula.DefaultWorkspaceName] = filepath.Join(r.userHomeDir, formula.DefaultWorkspaceDir)

	for name := range wspaces {
		dir := wspaces[name]
		if r.formulaExistsInWorkspace(dir, cmd) {
			wsWithFormula[name] = dir
		}
	}
	return wsWithFormula
}

func (r *renameFormulaCmd) cleanWorkspace(
	workspacesOld, workspacesNew formula.Workspaces,
	result formula.Rename,
	cmd *cobra.Command,
	wspaceName string,
) (formula.Workspace, error) {
	wsCleaned := formula.Workspace{}

	if len(workspacesOld) == 0 {
		return formula.Workspace{}, fmt.Errorf(errNonExistFormula, result.OldFormulaCmd)
	}
	if len(workspacesOld) > 1 {
		var items []string
		for k, v := range workspacesOld {
			kv := fmt.Sprintf("%s (%s)", k, v)
			items = append(items, kv)
		}

		var name string
		if !isInputFlag(cmd) {
			question := fmt.Sprintf("We found the old formula %q in %d workspaces. Select the workspace:",
				result.OldFormulaCmd, len(workspacesOld),
			)
			selected, err := r.inList.List(question, items)
			if err != nil {
				return formula.Workspace{}, err
			}
			name = strings.Split(selected, " (")[0]
		} else {
			name = wspaceName
		}

		if name != "" {
			name = strings.Title(name)
			wsCleaned.Name = name
			wsCleaned.Dir = workspacesOld[name]
		} else {
			return formula.Workspace{}, fmt.Errorf(errFormulaInManyWS, result.OldFormulaCmd, len(workspacesOld))
		}
	} else {
		if wspaceName != "" {
			wsCleaned.Name = strings.Title(wspaceName)
			wsCleaned.Dir = workspacesOld[wspaceName]
		} else {
			for n, d := range workspacesOld {
				wsCleaned.Name = n
				wsCleaned.Dir = d
			}
		}
	}

	for n := range workspacesNew {
		if n == wsCleaned.Name {
			return formula.Workspace{}, fmt.Errorf(errFormulaExists, result.NewFormulaCmd, wsCleaned.Name)
		}
	}

	if wsCleaned.Dir == "" {
		return formula.Workspace{}, fmt.Errorf(errInvalidWS, wspaceName)
	}

	return wsCleaned, nil
}

func (r *renameFormulaCmd) Rename(fr formula.Rename) error {
	if err := r.moveFormulaToNewDir(fr); err != nil {
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

func (r *renameFormulaCmd) moveFormulaToNewDir(fr formula.Rename) error {
	tmpName := api.RootName + "_oldFormula"
	tmp := filepath.Join(os.TempDir(), tmpName)
	if err := r.directory.Create(tmp); err != nil {
		return err
	}

	//nolint:errcheck
	defer os.RemoveAll(tmp)

	if err := r.directory.Copy(fr.FOldPath, tmp); err != nil {
		return err
	}

	groupsOld := strings.Split(fr.OldFormulaCmd, " ")
	delOld := formula.Delete{
		GroupsFormula: groupsOld,
		Workspace:     fr.Workspace,
	}

	if err := r.deleter.Delete(delOld); err != nil {
		return err
	}

	if err := r.directory.Create(fr.FNewPath); err != nil {
		return err
	}

	if err := r.directory.Copy(tmp, fr.FNewPath); err != nil {
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

func isEmpty(fr formula.Rename) bool {
	return fr.OldFormulaCmd == "" || fr.NewFormulaCmd == ""
}

func fPath(workspacePath, cmd string) string {
	cc := strings.Split(cmd, " ")
	path := strings.Join(cc, string(os.PathSeparator))

	return filepath.Join(workspacePath, path)
}

func cleanPreffix(cmd string) string {
	if strings.HasPrefix(cmd, api.RootName) {
		return cmd[4:]
	}
	return cmd
}
