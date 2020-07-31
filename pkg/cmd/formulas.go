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
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
)

const (
	subCommand  = " SUBCOMMAND"
	Group       = "group"
	dockerFlag  = "docker"
	rootCmdName = "root"
)

type FormulaCommand struct {
	coreCmds    api.Commands
	treeManager formula.TreeManager
	formula     formula.Runner
}

func NewFormulaCommand(
	coreCmds api.Commands,
	treeManager formula.TreeManager,
	formula formula.Runner) *FormulaCommand {
	return &FormulaCommand{
		coreCmds:    coreCmds,
		treeManager: treeManager,
		formula:     formula,
	}
}

func (f FormulaCommand) Add(root *cobra.Command) error {
	treeRep := f.treeManager.MergedTree(false)
	commands := make(map[string]*cobra.Command)
	commands[rootCmdName] = root

	for _, cmd := range treeRep.Commands {
		cmdPath := api.Command{Id: cmd.Id, Parent: cmd.Parent, Usage: cmd.Usage}
		if !sliceutil.ContainsCmd(f.coreCmds, cmdPath) {
			var newCmd *cobra.Command
			if cmd.Formula {
				newCmd = f.newFormulaCmd(cmd)
			} else {
				newCmd = newSubCmd(cmd)
			}

			parentCmd := commands[cmd.Parent]
			parentCmd.AddCommand(newCmd)
			commands[cmdPath.Id] = newCmd
		}
	}

	return nil
}

func newSubCmd(cmd api.Command) *cobra.Command {
	var group string
	if cmd.Parent == rootCmdName {
		group = fmt.Sprintf("%s repo commands:", cmd.Repo)
	}

	c := &cobra.Command{
		Use:         cmd.Usage + subCommand,
		Short:       cmd.Help,
		Long:        cmd.LongHelp,
		Annotations: map[string]string{Group: group},
	}
	c.LocalFlags()
	return c
}

func (f FormulaCommand) newFormulaCmd(cmd api.Command) *cobra.Command {
	formulaCmd := &cobra.Command{
		Use:   cmd.Usage,
		Short: cmd.Help,
		Long:  cmd.LongHelp,
	}

	addFlags(formulaCmd)
	path := strings.ReplaceAll(strings.Replace(cmd.Parent, "root", "", 1), "_", string(os.PathSeparator))
	path = fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), cmd.Usage)
	formulaCmd.RunE = f.execFormulaFunc(cmd.Repo, path)

	return formulaCmd
}

func (f FormulaCommand) execFormulaFunc(repo, path string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		d := formula.Definition{
			Path:     path,
			RepoName: repo,
		}

		stdin, err := cmd.Flags().GetBool(api.Stdin.ToLower())
		if err != nil {
			return err
		}
		inputType := api.Prompt
		if stdin {
			inputType = api.Stdin
		}

		docker, err := cmd.Flags().GetBool(dockerFlag)
		if err != nil {
			return err
		}

		if err := f.formula.Run(d, inputType, docker); err != nil {
			return err
		}

		return nil
	}
}

func addFlags(cmd *cobra.Command) {
	formulaFlags := cmd.Flags()
	formulaFlags.BoolP(dockerFlag, "d", false, "Use to run formulas inside docker")
}
