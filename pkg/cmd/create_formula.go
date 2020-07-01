package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

var (
	ErrNotAllowedCharacter = prompt.NewError(`not allowed character on formula name \/,><@-`)
	ErrDontStartWithRit    = prompt.NewError("Rit formula's command needs to start with \"rit\" [ex.: rit group verb <noun>]")
	ErrTooShortCommand     = prompt.NewError("Rit formula's command needs at least 2 words following \"rit\" [ex.: rit group verb]")
)

const notAllowedChars = `\/><,@-`

// createFormulaCmd type for add formula command
type createFormulaCmd struct {
	homeDir   string
	formula   formula.CreateBuilder
	workspace formula.WorkspaceAddListValidator
	inText    prompt.InputText
	inList    prompt.InputList
}

// CreateFormulaCmd creates a new cmd instance
func NewCreateFormulaCmd(
	homeDir string,
	formula formula.CreateBuilder,
	workspace formula.WorkspaceAddListValidator,
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
		formulaCmd, err := c.inText.TextWithValidate("Enter the new formula command [ex.: rit group verb noun]", c.surveyCmdValidator)
		if err != nil {
			return err
		}

		if strings.ContainsAny(formulaCmd, notAllowedChars) {
			return ErrNotAllowedCharacter
		}

		lang, err := c.inList.List("Choose the language: ", creator.Languages)
		if err != nil {
			return err
		}

		workspaces, err := c.workspace.List()
		if err != nil {
			return err
		}

		defaultWorkspace := path.Join(c.homeDir, formula.DefaultWorkspaceDir)
		workspaces[formula.DefaultWorkspaceName] = defaultWorkspace

		wspace, err := FormulaWorkspaceInput(workspaces, c.inList, c.inText)
		if err != nil {
			return err
		}

		if wspace.Dir != defaultWorkspace {
			if err := c.workspace.Add(wspace); err != nil {
				return err
			}
		}

		formulaPath := formulaPath(wspace.Dir, formulaCmd)

		cf := formula.Create{
			FormulaCmd:    formulaCmd,
			Lang:          lang,
			WorkspacePath: wspace.Dir,
			FormulaPath:   formulaPath,
		}

		c.create(cf, wspace.Dir, formulaPath)

		return nil
	}
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

func (c createFormulaCmd) create(cf formula.Create, workspacePath, formulaPath string) {
	buildInfo := prompt.Bold("Creating and building formula...")
	s := spinner.StartNew(buildInfo)
	time.Sleep(2 * time.Second)

	if err := c.formula.Create(cf); err != nil {
		err := prompt.NewError(err.Error())
		s.Error(err)
		return
	}

	// TODO: Remove this function after creating the build scripts for other languages ​​on windows
	os := runtime.GOOS
	if os == osutil.Windows && cf.Lang != creator.GoLang {
		createSuccess(s, cf.Lang)
		return
	}

	if err := c.formula.Build(workspacePath, formulaPath); err != nil {
		err := prompt.NewError(err.Error())
		s.Error(err)
		return
	}

	createSuccess(s, cf.Lang)
	buildSuccess(formulaPath, cf.FormulaCmd)
}

func createSuccess(s *spinner.Spinner, lang string) {
	msg := fmt.Sprintf("✔ %s formula successfully created!", lang)
	success := prompt.Green(msg)
	s.Success(success)
}

func buildSuccess(formulaPath, formulaCmd string) {
	prompt.Info(fmt.Sprintf("Formula path is %s", formulaPath))
	prompt.Info(fmt.Sprintf("Now you can run your formula with the following command %q", formulaCmd))
}

func formulaPath(workspacePath, cmd string) string {
	cc := strings.Split(cmd, " ")
	formulaPath := strings.Join(cc[1:], "/")
	return path.Join(workspacePath, formulaPath)
}

func (c createFormulaCmd) cmdValidator(cmd string) error {
	if len(strings.TrimSpace(cmd)) < 1 {
		return errors.New("this input must not be empty")
	}

	s := strings.Split(cmd, " ")
	if s[0] != "rit" {
		return ErrDontStartWithRit
	}

	if len(s) <= 2 {
		return ErrTooShortCommand
	}
	return nil
}

func (c createFormulaCmd) surveyCmdValidator(cmd interface{}) error {
	if len(strings.TrimSpace(cmd.(string))) < 1 {
		return errors.New("this input must not be empty")
	}

	s := strings.Split(cmd.(string), " ")
	if s[0] != "rit" {
		return ErrDontStartWithRit
	}

	if len(s) <= 2 {
		return ErrTooShortCommand
	}
	return nil
}


func FormulaWorkspaceInput(
	workspaces formula.Workspaces,
	inList prompt.InputList,
	inText prompt.InputText,
) (formula.Workspace, error) {
	var items []string
	for k, v := range workspaces {
		kv := fmt.Sprintf("%s (%s)", k, v)
		items = append(items, kv)
	}

	items = append(items, newWorkspace)
	selected, err := inList.List("Select a formula workspace: ", items)
	if err != nil {
		return formula.Workspace{}, err
	}

	var workspaceName string
	var workspacePath string
	var wspace formula.Workspace
	if selected == newWorkspace {
		workspaceName, err = inText.Text("Workspace name: ", true)
		if err != nil {
			return formula.Workspace{}, err
		}

		workspacePath, err = inText.Text("Workspace path (e.g.: /home/user/github):", true)
		if err != nil {
			return formula.Workspace{}, err
		}

		wspace = formula.Workspace{
			Name: strings.Title(workspaceName),
			Dir:  workspacePath,
		}
	} else {
		split := strings.Split(selected, " ")
		workspaceName = split[0]
		workspacePath = workspaces[workspaceName]
		wspace = formula.Workspace{
			Name: strings.Title(workspaceName),
			Dir:  workspacePath,
		}
	}
	return wspace, nil
}
