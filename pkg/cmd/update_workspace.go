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
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type updateWorkspaceCmd struct {
	workspace formula.WorkspaceListUpdater
	inList    prompt.InputList
}

var updateWorkspaceFlags = flags{
	{
		name:        nameFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: "workspace name",
	},
}

func NewUpdateWorkspaceCmd(
	workspace formula.WorkspaceListUpdater,
	inList prompt.InputList,
) *cobra.Command {
	u := updateWorkspaceCmd{
		workspace: workspace,
		inList:    inList,
	}

	cmd := &cobra.Command{
		Use:       "workspace",
		Short:     "Update a workspace",
		Example:   "rit update workspace",
		RunE:      u.runFormula(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), updateWorkspaceFlags)

	return cmd
}

func (u *updateWorkspaceCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		workspace, err := u.resolveInput(cmd)
		if err != nil {
			return err
		}

		if err := u.workspace.Update(workspace); err != nil {
			return err
		}

		prompt.Success("Workspace successfully updated!")

		return nil
	}
}

func (u *updateWorkspaceCmd) resolveInput(cmd *cobra.Command) (formula.Workspace, error) {
	if IsFlagInput(cmd) {
		return u.resolveFlags(cmd)
	}
	return u.resolvePrompt()
}

func (u *updateWorkspaceCmd) resolvePrompt() (formula.Workspace, error) {
	workspaces, err := u.workspace.List()
	if err != nil {
		return formula.Workspace{}, err
	}

	if len(workspaces) == 0 {
		return formula.Workspace{}, ErrEmptyWorkspaces
	}

	items := make([]string, 0, len(workspaces))
	for k, v := range workspaces {
		kv := fmt.Sprintf("%s (%s)", k, v)
		items = append(items, kv)
	}

	selected, err := u.inList.List("Select the workspace to update: ", items)
	if err != nil {
		return formula.Workspace{}, err
	}

	split := strings.Split(selected, " (")
	workspaceName := split[0]
	workspacePath := workspaces[workspaceName]
	wspace := formula.Workspace{
		Name: strings.Title(workspaceName),
		Dir:  workspacePath,
	}

	return wspace, nil
}

func (u *updateWorkspaceCmd) resolveFlags(cmd *cobra.Command) (formula.Workspace, error) {
	name, err := cmd.Flags().GetString(nameFlagName)
	if err != nil {
		return formula.Workspace{}, err
	} else if name == "" {
		return formula.Workspace{}, errors.New(missingFlagText(nameFlagName))
	}

	workspaces, err := u.workspace.List()
	if err != nil {
		return formula.Workspace{}, err
	}
	for workspaceName, path := range workspaces {
		if workspaceName == name {
			return formula.Workspace{Name: workspaceName, Dir: path}, nil
		}
	}

	return formula.Workspace{}, errors.New("no workspace found with this name")
}
