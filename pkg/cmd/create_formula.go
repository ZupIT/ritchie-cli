package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

var (
	ErrNotAllowedCharacter = fmt.Errorf(prompt.Red, `not allowed character on formula name \/,><@-`)
)

const notAllowedChars = `\/><,@-`

// createFormulaCmd type for add formula command
type createFormulaCmd struct {
	homeDir   string
	formula   formula.CreateBuilder
	workspace workspace.AddListValidator
	inText    prompt.InputText
	inList    prompt.InputList
}

// CreateFormulaCmd creates a new cmd instance
func NewCreateFormulaCmd(
	homeDir string,
	formula formula.CreateBuilder,
	workspace workspace.AddListValidator,
	inText prompt.InputText,
	inList prompt.InputList,
) *cobra.Command {
	c := createFormulaCmd{
		homeDir,
		formula,
		workspace,
		inText,
		inList,
	}

	cmd := &cobra.Command{
		Use:     "formula",
		Short:   "Create a new formula",
		Example: "rit create formula",
		RunE:    RunFuncE(c.runStdin(), c.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (c createFormulaCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		formulaCmd, err := c.inText.Text("Enter the new formula command [ex.: rit group verb noun]", true)
		if err != nil {
			return err
		}

		if strings.ContainsAny(formulaCmd, notAllowedChars) {
			return ErrNotAllowedCharacter
		}

		lang, err := c.inList.List("Choose the language: ", []string{"Go", "Java", "Node", "Python", "Shell"})
		if err != nil {
			return err
		}

		workspaces, err := c.workspace.List()
		if err != nil {
			return err
		}

		wspace, err := FormulaWorkspaceInput(c.homeDir, workspaces, c.inList, c.inText)
		if err != nil {
			return err
		}

		if err := c.workspace.Add(wspace); err != nil {
			return err
		}

		formulaPath := formulaPath(wspace.Dir, formulaCmd)

		cf := formula.Create{
			FormulaCmd:    formulaCmd,
			Lang:          lang,
			WorkspacePath: wspace.Dir,
			FormulaPath:   formulaPath,
		}

		if err := c.formula.Create(cf); err != nil {
			return err
		}

		if err := c.formula.Build(wspace.Dir, formulaPath); err != nil {
			return err
		}

		prompt.Success(fmt.Sprintf("%s formula successfully created!", lang))
		prompt.Info(fmt.Sprintf("Formula path is %s", wspace.Dir))

		return nil
	}
}

func formulaPath(workspacePath, cmd string) string {
	cc := strings.Split(cmd, " ")
	formulaPath := strings.Join(cc[1:], "/")
	return path.Join(workspacePath, formulaPath)
}

func (c createFormulaCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		var cf formula.Create

		if err := stdin.ReadJson(os.Stdin, &cf); err != nil {
			prompt.Error(stdin.MsgInvalidInput)
			return err
		}

		if strings.ContainsAny(cf.FormulaCmd, notAllowedChars) {
			return ErrNotAllowedCharacter
		}

		if err := c.formula.Create(cf); err != nil {
			return err
		}

		prompt.Success(fmt.Sprintf("%s formula successfully created!\n", cf.Lang))
		prompt.Info(fmt.Sprintf("Formula path is %s \n", cf.WorkspacePath))
		return nil
	}
}

func FormulaWorkspaceInput(
	homeDir string,
	workspaces workspace.Workspaces,
	inList prompt.InputList,
	inText prompt.InputText,
) (workspace.Workspace, error) {
	defaultWorkspace := path.Join(homeDir, workspace.DefaultWorkspaceDir)
	workspaces[workspace.DefaultWorkspaceName] = defaultWorkspace

	var items []string
	for k, v := range workspaces {
		kv := fmt.Sprintf("%s (%s)", k, v)
		items = append(items, kv)
	}

	items = append(items, newWorkspace)
	selected, err := inList.List("Select a formula workspace: ", items)
	if err != nil {
		return workspace.Workspace{}, err
	}

	var workspaceName string
	var workspacePath string
	var wspace workspace.Workspace
	if selected == newWorkspace {
		workspaceName, err = inText.Text("Workspace name: ", true)
		if err != nil {
			return workspace.Workspace{}, err
		}

		workspacePath, err = inText.Text("Workspace path (e.g.: /home/user/github):", true)
		if err != nil {
			return workspace.Workspace{}, err
		}

		wspace = workspace.Workspace{
			Name: strings.Title(workspaceName),
			Dir:  workspacePath,
		}
	} else {
		split := strings.Split(selected, " ")
		workspaceName = split[0]
		workspacePath = workspaces[workspaceName]
		wspace = workspace.Workspace{
			Name: strings.Title(workspaceName),
			Dir:  workspacePath,
		}
	}
	return wspace, nil
}
