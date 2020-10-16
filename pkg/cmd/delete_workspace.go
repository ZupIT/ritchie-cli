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
	"os"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/spf13/cobra"
)

const (
	msgWorkspaceIsNotValid = "The workspace informed is not valid"
	msgEmptyWorkspaces     = "There are no workspaces to delete"
)

var (
	ErrWorkspaceIsNotValid = errors.New(msgWorkspaceIsNotValid)
	ErrEmptyWorkspaces     = errors.New(msgEmptyWorkspaces)
)

type deleteWorkspaceCmd struct {
	userHomeDir string
	workspace   formula.WorkspaceListDeleter
	directory   stream.DirListChecker
	inList      prompt.InputList
	inBool      prompt.InputBool
}

func NewDeleteWorkspaceCmd(
	userHomeDir string,
	workspace formula.WorkspaceListDeleter,
	directory stream.DirListChecker,
	inList prompt.InputList,
	inBool prompt.InputBool,
) *cobra.Command {
	d := deleteWorkspaceCmd{
		userHomeDir,
		workspace,
		directory,
		inList,
		inBool,
	}

	cmd := &cobra.Command{
		Use:     "workspace",
		Short:   "Delete a workspace",
		Example: "rit delete workspace",
		RunE:    d.runPrompt(),
	}

	return cmd
}

func (d deleteWorkspaceCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		workspaces, err := d.workspace.List()
		if err != nil {
			return err
		}

		defaultWorkspace := filepath.Join(d.userHomeDir, formula.DefaultWorkspaceDir)
		if d.directory.Exists(defaultWorkspace) {
			workspaces[formula.DefaultWorkspaceName] = defaultWorkspace
		}

		if len(workspaces) == 0 {
			return ErrEmptyWorkspaces
		}

		wspace, err := WorkspaceListInput(workspaces, d.inList)
		if err != nil {
			return err
		}

		question := fmt.Sprintf("Are you sure you want to delete the workspace: rit %s", wspace.Dir)
		ans, err := d.inBool.Bool(question, []string{"no", "yes"})
		if err != nil {
			return err
		}
		if !ans {
			return nil
		}

		if err := d.deleteWorkspace(wspace.Dir); err != nil {
			return err
		}

		if wspace.Dir != defaultWorkspace {
			if err := d.workspace.Delete(wspace); err != nil {
				return err
			}
		}

		prompt.Success("Workspace successfully deleted!")

		return nil
	}
}

func (d deleteWorkspaceCmd) deleteWorkspace(workspace string) error {
	if d.directory.Exists(workspace) {
		return os.RemoveAll(workspace)
	}

	return ErrWorkspaceIsNotValid
}

func WorkspaceListInput(
	workspaces formula.Workspaces,
	inList prompt.InputList,
) (formula.Workspace, error) {
	var items []string
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
