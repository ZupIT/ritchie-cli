package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	workspaceNameFlag            = "name"
	workspaceNameFlagDescription = "Workspace name"
	workspacePathFlag            = "path"
	workspacePathFlagDescription = "Workspace path"
)

type addWorkspaceCmd struct {
	workspace formula.WorkspaceAddLister
	input     prompt.InputText
	inPath    prompt.InputPath
}

func NewAddWorkspaceCmd(
	workspace formula.WorkspaceAddLister,
	input prompt.InputText,
	inPath prompt.InputPath,
) *cobra.Command {
	a := addWorkspaceCmd{
		workspace: workspace,
		input:     input,
		inPath:    inPath,
	}

	cmd := &cobra.Command{
		Use:       "workspace",
		Short:     "Add new workspace",
		Example:   "rit add workspace",
		RunE:      a.runFormula(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	flags := cmd.Flags()
	flags.String(workspaceNameFlag, "", workspaceNameFlagDescription)
	flags.String(workspacePathFlag, "", workspacePathFlagDescription)

	return cmd
}

func (a addWorkspaceCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().NFlag() == 0 {
			wspace, err := addWorkspaceFromPrompt(a)
			if err != nil {
				return err
			}

			return a.workspace.Add(wspace)
		}

		wspace, err := addWorkspaceFromFlags(cmd)
		if err != nil {
			return err
		}

		return a.workspace.Add(wspace)
	}
}

func addWorkspaceFromFlags(cmd *cobra.Command) (formula.Workspace, error) {
	workspaceName, _ := cmd.Flags().GetString(workspaceNameFlag)
	workspacePath, _ := cmd.Flags().GetString(workspacePathFlag)

	if len(workspaceName) == 0 || len(workspacePath) == 0 {
		return formula.Workspace{}, errors.New("all flags need to be filled")
	}

	wspace := formula.Workspace{
		Name: strings.Title(workspaceName),
		Dir:  workspacePath,
	}

	return wspace, nil
}

func addWorkspaceFromPrompt(a addWorkspaceCmd) (formula.Workspace, error) {
	workspaceName, err := a.input.Text("Enter the name of workspace", true)
	if err != nil {
		return formula.Workspace{}, err
	}

	workspacePath, err := a.inPath.Read("Enter the path of workspace (e.g.: /home/user/github) ")
	if err != nil {
		return formula.Workspace{}, err
	}

	wspace := formula.Workspace{
		Name: strings.Title(workspaceName),
		Dir:  workspacePath,
	}

	return wspace, nil
}
