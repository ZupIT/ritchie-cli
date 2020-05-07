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
	fLBin    = "fLBin"
	fMBin    = "fMBin"
	fWBin    = "fWBin"
	fBundle  = "fBundle"
	fConfig  = "fConfig"
	fRepoURL = "fRepoURL"
	subcmd   = " SUBCOMMAND"
	//Group formulas group
	Group = "group"
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
	annotations[fLBin] = frm.LBin
	annotations[fMBin] = frm.MBin
	annotations[fWBin] = frm.WBin
	annotations[fBundle] = frm.Bundle
	annotations[fConfig] = frm.Config
	annotations[fRepoURL] = frm.RepoURL

	c := &cobra.Command{
		Annotations: annotations,
		Use:         cmd.Usage,
		Short:       cmd.Help,
		Long:        cmd.Help,
		RunE:        execFormulaFunc(f.formulaRunner),
	}
	c.LocalFlags()
	return c
}

func execFormulaFunc(formulaRunner formula.Runner) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fPath := cmd.Annotations[fPath]
		fBin := cmd.Annotations[fBin]
		fLBin := cmd.Annotations[fLBin]
		fMBin := cmd.Annotations[fMBin]
		fWBin := cmd.Annotations[fWBin]
		fBundle := cmd.Annotations[fBundle]
		fConf := cmd.Annotations[fConfig]
		fRepoURL := cmd.Annotations[fRepoURL]
		frm := formula.Definition{
			Path:    fPath,
			Bin:     fBin,
			LBin:    fLBin,
			MBin:    fMBin,
			WBin:    fWBin,
			Bundle:  fBundle,
			Config:  fConf,
			RepoURL: fRepoURL,
		}
		stdin, err := cmd.Flags().GetBool(api.Stdin.ToLower())
		if err != nil {
			return err
		}
		inputType := api.Prompt
		if stdin {
			inputType = api.Stdin
		}
		return formulaRunner.Run(frm, inputType)
	}
}

func newSubCmd(cmd api.Command) *cobra.Command {
	group := ""
	if cmd.Parent == "root" {
		group = fmt.Sprintf("%s commands:", cmd.Repo)
	}

	return &cobra.Command{
		Use:         cmd.Usage + subcmd,
		Short:       cmd.Help,
		Long:        cmd.Help,
		Annotations: map[string]string{"group": group},
	}
}
