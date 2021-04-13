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

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	workspaceFlagName        = "env"
	workspaceFlagDescription = "Env name to delete"
)

// renameFormulaCmd type for add formula command.
type renameFormulaCmd struct {
	workspace formula.WorkspaceAddListHasher
	inText    prompt.InputText
	inList    prompt.InputList
	inPath    prompt.InputPath
}

// New renameFormulaCmd rename a cmd instance.
func NewRenameFormulaCmd(
	workspace formula.WorkspaceAddListHasher,
	inText prompt.InputText,
	inList prompt.InputList,
	inPath prompt.InputPath,
) *cobra.Command {
	r := renameFormulaCmd{
		workspace: workspace,
		inText:    inText,
		inList:    inList,
		inPath:    inPath,
	}

	cmd := &cobra.Command{
		Use:       "formula",
		Short:     "Rename a formula",
		Example:   "rit rename formula",
		RunE:      r.runFormula(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	cmd.LocalFlags()

	return cmd
}

func (r *renameFormulaCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		wspace, err := r.resolveInput(cmd)
		if err != nil {
			return err
		}

		fmt.Println(wspace)

		return nil

	}
}

func (r *renameFormulaCmd) resolveInput(cmd *cobra.Command) (formula.Workspace, error) {
	if IsFlagInput(cmd) {
		return r.resolveFlags(cmd)
	}
	return r.resolvePrompt()
}

func (r *renameFormulaCmd) resolveFlags(cmd *cobra.Command) (formula.Workspace, error) {
	return formula.Workspace{}, nil
}

func (r *renameFormulaCmd) resolvePrompt() (formula.Workspace, error) {
	workspaces, err := r.workspace.List()
	if err != nil {
		return formula.Workspace{}, err
	}
	wspace, err := FormulaWorkspaceInput(workspaces, r.inList, r.inText, r.inPath)
	if err != nil {
		return formula.Workspace{}, err
	}

	return wspace, nil
}
