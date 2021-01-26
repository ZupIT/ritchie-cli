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
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo/repoutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var ErrEmptyWorkspaces = errors.New("there are no workspaces to delete")

type deleteWorkspaceCmd struct {
	userHomeDir string
	workspace   formula.WorkspaceListDeleter
	repo        formula.RepositoryDeleter
	inList      prompt.InputList
	inBool      prompt.InputBool
}

var deleteWorkspaceFlags = flags{
	{
		name:        nameFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: "workspace name",
	},
}

func NewDeleteWorkspaceCmd(
	userHomeDir string,
	workspace formula.WorkspaceListDeleter,
	repo formula.RepositoryDeleter,
	inList prompt.InputList,
	inBool prompt.InputBool,
) *cobra.Command {
	d := deleteWorkspaceCmd{
		userHomeDir: userHomeDir,
		workspace:   workspace,
		repo:        repo,
		inList:      inList,
		inBool:      inBool,
	}

	cmd := &cobra.Command{
		Use:     "workspace",
		Short:   "Delete a workspace",
		Example: "rit delete workspace",
		RunE:    d.runFormula(),
	}

	addReservedFlags(cmd.Flags(), deleteWorkspaceFlags)

	return cmd
}

func (d *deleteWorkspaceCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		workspace, err := d.resolveInput(cmd)
		if err != nil {
			return err
		}

		repoLocalName := repoutil.LocalName(workspace.Name)
		if err := d.repo.Delete(repoLocalName); err != nil {
			return err
		}

		if workspace.Name == formula.DefaultWorkspaceName {
			return errors.New("cannot delete default workspace")
		}

		if err := d.workspace.Delete(workspace); err != nil {
			return err
		}

		prompt.Success("Workspace successfully deleted!")

		return nil
	}
}

func (d *deleteWorkspaceCmd) resolveInput(cmd *cobra.Command) (formula.Workspace, error) {
	if IsFlagInput(cmd) {
		return d.resolveFlags(cmd)
	}
	return d.resolvePrompt()
}

func (d *deleteWorkspaceCmd) resolvePrompt() (formula.Workspace, error) {
	workspaces, err := d.workspace.List()
	if err != nil {
		return formula.Workspace{}, err
	}

	if len(workspaces) == 0 {
		return formula.Workspace{}, ErrEmptyWorkspaces
	}

	wspace, err := WorkspaceListInput(workspaces, d.inList)
	if err != nil {
		return formula.Workspace{}, err
	}

	question := fmt.Sprintf("Are you sure you want to delete the workspace: rit %s", wspace.Dir)
	ans, err := d.inBool.Bool(question, []string{"no", "yes"})
	if err != nil {
		return formula.Workspace{}, err
	}
	if !ans {
		return formula.Workspace{}, nil
	}
	return wspace, nil
}

func (d *deleteWorkspaceCmd) resolveFlags(cmd *cobra.Command) (formula.Workspace, error) {
	name, err := cmd.Flags().GetString(nameFlagName)
	if err != nil {
		return formula.Workspace{}, err
	} else if name == "" {
		return formula.Workspace{}, errors.New("please provide a value for 'name'")
	}

	workspaces, err := d.workspace.List()
	workspaces[formula.DefaultWorkspaceName] = filepath.Join(d.userHomeDir, formula.DefaultWorkspaceDir)
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

func WorkspaceListInput(
	workspaces formula.Workspaces,
	inList prompt.InputList,
) (formula.Workspace, error) {
	items := make([]string, 0, len(workspaces))
	for k, v := range workspaces {
		kv := fmt.Sprintf("%s (%s)", k, v)
		items = append(items, kv)
	}

	selected, err := inList.List("Select the workspace: ", items)
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
