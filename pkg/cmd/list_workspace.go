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

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

const (
	totalWorkspacesMsg   = "There are %v workspaces"
	totalOneWorkspaceMsg = "There is 1 workspace"
)

type listWorkspaceCmd struct {
	formula.WorkspaceLister
	rt rtutorial.Finder
}

func NewListWorkspaceCmd(wm workspace.Manager, rtf rtutorial.Finder) *cobra.Command {
	lw := listWorkspaceCmd{wm, rtf}
	cmd := &cobra.Command{
		Use:     "workspace",
		Short:   "Show a list with all your available workspaces",
		Example: "rit list workspace",
		RunE:    lw.runFunc(),
	}
	return cmd
}

func (lr listWorkspaceCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		workspaces, err := lr.List()
		if err != nil {
			return err
		}

		printWorkspaces(workspaces)

		if len(workspaces) != 1 {
			prompt.Info(fmt.Sprintf(totalWorkspacesMsg, len(workspaces)))
		} else {
			prompt.Info(totalOneWorkspaceMsg)
		}

		tutorialHolder, err := lr.rt.Find()
		if err != nil {
			return err
		}
		tutorialListWorkspaces(tutorialHolder.Current)
		return nil
	}
}

func printWorkspaces(workspaces formula.Workspaces) {
	table := uitable.New()
	table.AddRow("NAME", "PATH")
	for k, v := range workspaces {
		table.AddRow(k, v)
	}
	raw := table.Bytes()
	raw = append(raw, []byte("\n")...)
	fmt.Println(string(raw))

}

func tutorialListWorkspaces(tutorialStatus string) {
	const tagTutorial = "\n[TUTORIAL]"
	const MessageTitle = "To delete a workspace:"
	const MessageBody = ` ∙ Run "rit delete workspace"`

	if tutorialStatus == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
	}
}

