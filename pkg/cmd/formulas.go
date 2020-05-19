package cmd

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/api"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
)

const (
	subCommand = " SUBCOMMAND"
	// Group formulas group
	Group = "group"
)

type FormulaCommand struct {
	coreCmds      api.Commands
	treeManager   formula.TreeManager
	defaultRunner formula.Runner
	dockerRunner  formula.Runner
}

func NewFormulaCommand(
	coreCmds api.Commands,
	treeManager formula.TreeManager,
	defaultRunner formula.Runner,
	dockerRunner formula.Runner) *FormulaCommand {
	return &FormulaCommand{
		coreCmds:      coreCmds,
		treeManager:   treeManager,
		defaultRunner: defaultRunner,
		dockerRunner:  dockerRunner,
	}
}

func (f FormulaCommand) Add(rootCmd *cobra.Command) error {
	treeRep := f.treeManager.MergedTree(false)
	commands := make(map[string]*cobra.Command)
	commands["root"] = rootCmd

	for _, cmd := range treeRep.Commands {
		cmdPath := api.Command{Parent: cmd.Parent, Usage: cmd.Usage}
		if !sliceutil.ContainsCmd(f.coreCmds, cmdPath) {
			var newCmd *cobra.Command
			if cmd.Formula.Path != "" {
				newCmd = f.newFormulaCmd(cmd)
			} else {
				newCmd = newSubCmd(cmd)
			}

			parentCmd := commands[cmd.Parent]
			parentCmd.AddCommand(newCmd)
			cmdKey := fmt.Sprintf("%s_%s", cmdPath.Parent, cmdPath.Usage)
			commands[cmdKey] = newCmd
		}
	}

	return nil
}

func newSubCmd(cmd api.Command) *cobra.Command {
	var group string
	if cmd.Parent == "root" {
		group = fmt.Sprintf("%s commands:", cmd.Repo)
	}

	return &cobra.Command{
		Use:         cmd.Usage + subCommand,
		Short:       cmd.Help,
		Long:        cmd.Help,
		Annotations: map[string]string{Group: group},
	}
}

func (f FormulaCommand) newFormulaCmd(cmd api.Command) *cobra.Command {
	formulaCmd := &cobra.Command{
		Use:   cmd.Usage,
		Short: cmd.Help,
		Long:  cmd.Help,
	}

	var docker bool
	formulaFlags := formulaCmd.Flags()
	formulaFlags.BoolVar(&docker, "docker", false, "Use to run formulas inside a docker container")
	formulaCmd.RunE = f.execFormulaFunc(cmd.Repo, cmd.Formula, &docker)

	return formulaCmd
}

func (f FormulaCommand) execFormulaFunc(repo string, form api.Formula, docker *bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		d := formula.Definition{
			Path:     form.Path,
			Bin:      form.Bin,
			LBin:     form.LBin,
			MBin:     form.MBin,
			WBin:     form.WBin,
			Bundle:   form.Bundle,
			Config:   form.Config,
			RepoUrl:  form.RepoURL,
			RepoName: repo,
		}

		if *docker {
			return f.dockerRunner.Run(d)
		}

		return f.defaultRunner.Run(d)
	}
}
