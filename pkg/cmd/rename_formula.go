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
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	workspaceFlagName        = "workspace"
	workspaceFlagDescription = "Workspace to rename"
)

var renameWorkspaceFlags = flags{
	{
		name:        workspaceFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: workspaceFlagDescription,
	},
}

// renameFormulaCmd type for add formula command.
type renameFormulaCmd struct {
	workspace formula.WorkspaceAddListHasher
	inText    prompt.InputText
	inList    prompt.InputList
	inPath    prompt.InputPath
	directory stream.DirListChecker
}

// New renameFormulaCmd rename a cmd instance.
func NewRenameFormulaCmd(
	workspace formula.WorkspaceAddListHasher,
	inText prompt.InputText,
	inList prompt.InputList,
	inPath prompt.InputPath,
	directory stream.DirListChecker,
) *cobra.Command {
	r := renameFormulaCmd{
		workspace: workspace,
		inText:    inText,
		inList:    inList,
		inPath:    inPath,
		directory: directory,
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
		wspace, formula, err := r.resolveInput(cmd)
		if err != nil {
			return err
		}

		fmt.Println(wspace, formula)

		return nil
	}
}

func (r *renameFormulaCmd) resolveInput(cmd *cobra.Command) (formula.Workspace, string, error) {
	if IsFlagInput(cmd) {
		return r.resolveFlags(cmd)
	}
	return r.resolvePrompt()
}

func (r *renameFormulaCmd) resolveFlags(cmd *cobra.Command) (formula.Workspace, string, error) {
	return formula.Workspace{}, "", nil
}

func (r *renameFormulaCmd) resolvePrompt() (formula.Workspace, string, error) {
	workspaces, err := r.workspace.List()
	if err != nil {
		return formula.Workspace{}, "", err
	}
	wspace, err := FormulaWorkspaceInput(workspaces, r.inList, r.inText, r.inPath)
	if err != nil {
		return formula.Workspace{}, "", err
	}

	groups, err := r.readFormulas(wspace.Dir, "rit")
	if err != nil {
		return formula.Workspace{}, "", err
	}

	if groups == nil {
		return formula.Workspace{}, "", ErrCouldNotFindFormula
	}

	return wspace, strings.Join(groups, " "), nil
}

func (r *renameFormulaCmd) readFormulas(dir string, currentFormula string) ([]string, error) {
	dirs, err := r.directory.List(dir, false)
	if err != nil {
		return nil, err
	}

	dirs = removeFromArray(dirs, docsDir)

	var groups []string
	var formulaOptions []string
	var response string

	if isFormula(dirs) {
		if !hasFormulaInDir(dirs) {
			return groups, nil
		}

		formulaOptions = append(formulaOptions, currentFormula, optionOtherFormula)

		response, err = r.inList.List(foundFormulaQuestion, formulaOptions)
		if err != nil {
			return nil, err
		}
		if response == currentFormula {
			return groups, nil
		}
		dirs = removeFromArray(dirs, srcDir)
	}

	selected, err := r.inList.List(questionSelectFormulaGroup, dirs)
	if err != nil {
		return nil, err
	}

	newFormulaSelected := fmt.Sprintf("%s %s", currentFormula, selected)

	var aux []string
	aux, err = r.readFormulas(filepath.Join(dir, selected), newFormulaSelected)
	if err != nil {
		return nil, err
	}

	aux = append([]string{selected}, aux...)
	groups = append(groups, aux...)

	return groups, nil
}
