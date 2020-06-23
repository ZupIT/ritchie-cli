package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
)

const (
	subCommand = " SUBCOMMAND"
	Group      = "group"
	dockerFlag = "docker"
	RootCmd    = "root"
)

type FormulaCommand struct {
	coreCmds      api.Commands
	treeManager   formula.Manager
	defaultRunner formula.Runner
	dockerRunner  formula.Runner
}

func NewFormulaCommand(
	coreCmds api.Commands,
	treeManager formula.Manager,
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
	commands[RootCmd] = rootCmd

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
	if cmd.Parent == RootCmd {
		group = fmt.Sprintf("%s commands:", cmd.Repo)
	}

	c := &cobra.Command{
		Use:         cmd.Usage + subCommand,
		Short:       cmd.Help,
		Long:        cmd.Help,
		Annotations: map[string]string{Group: group},
	}
	c.LocalFlags()
	return c
}

func (f FormulaCommand) newFormulaCmd(cmd api.Command) *cobra.Command {
	formulaCmd := &cobra.Command{
		Use:   cmd.Usage,
		Short: cmd.Help,
		Long:  cmd.Help,
	}

	addFlags(formulaCmd)
	formulaCmd.RunE = f.execFormulaFunc(cmd.Repo, cmd.Formula)

	return formulaCmd
}

func (f FormulaCommand) execFormulaFunc(repo string, form api.Formula) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		d := formula.Definition{
			Path:     form.Path,
			Bin:      form.Bin,
			LBin:     form.LBin,
			MBin:     form.MBin,
			WBin:     form.WBin,
			Bundle:   form.Bundle,
			Config:   form.Config,
			RepoURL:  form.RepoURL,
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

		if docker {
			return f.dockerRunner.Run(d, inputType)
		}

		return f.defaultRunner.Run(d, inputType)
	}
}

func addFlags(cmd *cobra.Command) {
	formulaFlags := cmd.Flags()
	formulaFlags.BoolP(dockerFlag, "d", false, "Use to run formulas inside a docker container")
}
