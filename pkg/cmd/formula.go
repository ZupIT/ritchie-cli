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
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	Group       = "group"
	verboseFlag = "verbose"
	rootCmdName = "root"
)

var (
	reservedFlags = flags{
		{
			name:        formula.DockerRun.String(),
			shortName:   "d",
			kind:        reflect.Bool,
			defValue:    false,
			description: "Use to run formulas inside docker",
		},
		{
			name:        formula.LocalRun.String(),
			shortName:   "l",
			kind:        reflect.Bool,
			defValue:    false,
			description: "Use to run formulas locally",
		},
		{
			name:        "verbose",
			shortName:   "a",
			kind:        reflect.Bool,
			defValue:    false,
			description: "Verbose mode (All). Indicate to a formula that it should show log messages in more detail",
		},
		{
			name:        "default",
			kind:        reflect.Bool,
			defValue:    false,
			description: "Use to automatically fill inputs with default value provided on config.json",
		},
	}
)

var ErrRunFormulaWithTwoFlag = errors.New("you cannot run formula with --docker and --local flags together")

type FormulaCommand struct {
	coreCmds    api.Commands
	treeManager formula.TreeManager
	formula     formula.Executor
	file        stream.FileReader
}

func NewFormulaCommand(
	coreCmds api.Commands,
	treeManager formula.TreeManager,
	formula formula.Executor,
	file stream.FileReader,
) *FormulaCommand {
	return &FormulaCommand{
		coreCmds:    coreCmds,
		treeManager: treeManager,
		formula:     formula,
		file:        file,
	}
}

func (f FormulaCommand) Add(root *cobra.Command) error {
	tree := f.treeManager.MergedTree(false)
	commands := make(map[string]*cobra.Command)
	commands[rootCmdName] = root

	for _, id := range tree.CommandsID {
		cmd := tree.Commands[id]

		if containsCmd(f.coreCmds, cmd) {
			continue
		}

		var newCmd *cobra.Command
		if cmd.Formula {
			newCmd = f.newFormulaCmd(cmd)
		} else {
			newCmd = newSubCmd(cmd)
		}

		parentCmd := commands[cmd.Parent]
		parentCmd.AddCommand(newCmd)
		commands[id.String()] = newCmd
	}

	return nil
}

func newSubCmd(cmd api.Command) *cobra.Command {
	var group string
	if cmd.Parent == rootCmdName {
		group = fmt.Sprintf("%s repo commands:", cmd.Repo)
		if cmd.RepoNewVersion != "" {
			group = fmt.Sprintf("%s repo commands: %s", cmd.Repo, prompt.Cyan("(New version available "+cmd.RepoNewVersion+")"))
		}
	}

	c := &cobra.Command{
		Use:   cmd.Usage,
		Short: cmd.Help,
		Long:  cmd.LongHelp,
		Annotations: map[string]string{
			Group: group,
		},
	}
	c.LocalFlags()
	return c
}

func (f FormulaCommand) newFormulaCmd(cmd api.Command) *cobra.Command {
	formulaPath := path(cmd)
	formulaCmd := &cobra.Command{
		Use:   cmd.Usage,
		Short: cmd.Help,
		Long:  cmd.LongHelp,
		RunE:  f.execFormulaFunc(cmd.Repo, formulaPath),
	}

	def := formula.Definition{
		Path:     formulaPath,
		RepoName: cmd.Repo,
	}

	flags := formulaCmd.Flags()
	addReservedFlags(flags, reservedFlags)
	f.addInputFlags(def, flags)

	return formulaCmd
}

func (f FormulaCommand) execFormulaFunc(repo, path string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		docker, err := flags.GetBool(formula.DockerRun.String())
		if err != nil {
			return err
		}

		local, err := flags.GetBool(formula.LocalRun.String())
		if err != nil {
			return err
		}

		verbose, err := flags.GetBool(verboseFlag)
		if err != nil {
			return err
		}

		if docker && local {
			return ErrRunFormulaWithTwoFlag
		}

		runType := formula.DefaultRun
		if docker {
			runType = formula.DockerRun
		}

		if local {
			runType = formula.LocalRun
		}

		inputType := inputResolver(cmd)

		exe := formula.ExecuteData{
			Def: formula.Definition{
				Path:     path,
				RepoName: repo,
			},
			InType:  inputType,
			RunType: runType,
			Verbose: verbose,
			Flags:   flags,
		}

		if err := f.formula.Execute(exe); err != nil {
			return err
		}

		return nil
	}
}

func (f FormulaCommand) addInputFlags(def formula.Definition, flags *pflag.FlagSet) {
	s := def.FormulaPath(api.RitchieHomeDir())
	configPath := def.ConfigPath(s)
	file, _ := f.file.Read(configPath)
	var config formula.Config
	_ = json.Unmarshal(file, &config)

	for _, in := range config.Inputs {
		if flags.Lookup(in.Name) != nil {
			continue
		}

		if in.Type == input.BoolType {
			flags.Bool(in.Name, false, in.Tutorial)
		} else {
			flags.String(in.Name, in.Default, in.Tutorial)
		}
	}
}

func inputResolver(cmd *cobra.Command) api.TermInputType {
	switch {
	case isInputStdin():
		return api.Stdin
	case isInputFlag(cmd):
		return api.Flag
	default:
		return api.Prompt
	}
}

func isInputStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return stat.Size() > 0
}

func isInputFlag(cmd *cobra.Command) bool {
	flags := cmd.Flags()
	c := 0
	for _, flag := range reservedFlags {
		if flags.Changed(flag.name) {
			c++
		}
	}

	return flags.NFlag() > c
}

func path(cmd api.Command) string {
	path := strings.ReplaceAll(strings.Replace(cmd.Parent, "root", "", 1), "_", string(os.PathSeparator))
	return filepath.Join(path, cmd.Usage)
}

func containsCmd(aa api.Commands, c api.Command) bool {
	for _, v := range aa {
		if c.Parent == v.Parent && c.Usage == v.Usage {
			return true
		}

		coreCmd := fmt.Sprintf("%s_%s", v.Parent, v.Usage)
		if c.Parent == coreCmd { // Ensures that no core commands will be added
			return true
		}
	}
	return false
}
