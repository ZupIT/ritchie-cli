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
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	subCommand  = " SUBCOMMAND"
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
	}
)

type flag struct {
	name        string
	shortName   string
	kind        reflect.Kind
	defValue    interface{}
	description string
}

type flags []flag

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
	addReservedFlags(flags)
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

func addReservedFlags(flags *pflag.FlagSet) {
	for _, flag := range reservedFlags {
		switch flag.kind {
		case reflect.String:
			flags.StringP(flag.name, flag.shortName, flag.defValue.(string), flag.description)
		case reflect.Bool:
			flags.BoolP(flag.name, flag.shortName, flag.defValue.(bool), flag.description)
		case reflect.Int:
			flags.IntP(flag.name, flag.shortName, flag.defValue.(int), flag.description)
		default:
			prompt.Warning("this type of flag is not supported")
		}
	}
}

func (f FormulaCommand) addInputFlags(def formula.Definition, flags *pflag.FlagSet) {
	s := def.FormulaPath(api.RitchieHomeDir())
	configPath := def.ConfigPath(s)
	file, _ := f.file.Read(configPath)
	var config formula.Config
	_ = json.Unmarshal(file, &config)

	for _, in := range config.Inputs {
		switch in.Type {
		case input.TextType, input.PassType:
			flags.String(in.Name, in.Default, in.Tutorial)
		case input.BoolType:
			flags.Bool(in.Name, false, in.Tutorial)
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
		if changed := flags.Changed(flag.name); changed {
			c++
		}
	}

	return flags.NFlag() > c
}

func path(cmd api.Command) string {
	path := strings.ReplaceAll(strings.Replace(cmd.Parent, "root", "", 1), "_", string(os.PathSeparator))
	return filepath.Join(path, cmd.Usage)
}
