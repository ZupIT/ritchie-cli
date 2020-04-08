package cmd

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/api"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/spf13/cobra"
)

const (
	fPath    = "fPath"
	fBin     = "fBin"
	fConfig  = "fConfig"
	fRepoURL = "fRepoURL"
	subcmd   = " SUBCOMMAND"
)

type FormulaCommand struct {
	coreCmds      []api.Command
	treeManager   formula.TreeManager
	formulaRunner formula.Runner
}

func NewFormulaCommand(
	coreCmds []api.Command,
	treeManager formula.TreeManager,
	formulaRunner formula.Runner) *FormulaCommand {
	return &FormulaCommand{
		coreCmds:      coreCmds,
		treeManager:   treeManager,
		formulaRunner: formulaRunner,
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

func (f FormulaCommand) newFormulaCmd(cmd api.Command) *cobra.Command {
	frm := cmd.Formula
	annotations := make(map[string]string)
	annotations[fPath] = frm.Path
	annotations[fBin] = frm.Bin
	annotations[fConfig] = frm.Config
	annotations[fRepoURL] = frm.RepoURL

	return &cobra.Command{
		Annotations: annotations,
		Use:         cmd.Usage,
		Short:       cmd.Help,
		Long:        cmd.Help,
		RunE:        execFormulaFunc(f.formulaRunner),
	}
}

func execFormulaFunc(formulaRunner formula.Runner) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fPath := cmd.Annotations[fPath]
		fBin := cmd.Annotations[fBin]
		fConf := cmd.Annotations[fConfig]
		fRepoURL := cmd.Annotations[fRepoURL]
		frm := formula.Definition{
			Path:    fPath,
			Bin:     fBin,
			Config:  fConf,
			RepoUrl: fRepoURL,
		}
		return formulaRunner.Run(frm)
	}
}

func newSubCmd(cmd api.Command) *cobra.Command {
	return &cobra.Command{
		Use:   cmd.Usage + subcmd,
		Short: cmd.Help,
		Long:  cmd.Help,
	}
}
