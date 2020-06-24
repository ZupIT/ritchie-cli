package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	newWorkspace = "Type new formula workspace?"
	dirPattern   = "%s/%s"
	treeDir      = "tree"
	srcDir       = "src"
)

type buildFormulaCmd struct {
	userHomeDir string
	workspace   workspace.AddListValidator
	formula     formula.Builder
	watcher     formula.Watcher
	directory   stream.DirListChecker
	prompt.InputText
	prompt.InputList
}

func NewBuildFormulaCmd(
	userHomeDir string,
	workManager workspace.AddListValidator,
	formula formula.Builder,
	watcher formula.Watcher,
	directory stream.DirListChecker,
	inText prompt.InputText,
	inList prompt.InputList,
) *cobra.Command {
	s := buildFormulaCmd{
		userHomeDir: userHomeDir,
		workspace:   workManager,
		formula:     formula,
		watcher:     watcher,
		directory:   directory,
		InputText:   inText,
		InputList:   inList,
	}

	cmd := &cobra.Command{
		Use:   "formula",
		Short: "Build your formulas locally. Use --watch flag and get real-time updates.",
		Long: `Use this command to build your formulas locally. To make formulas development easier, you can run 
the command with the --watch flag and get real-time updates.`,
		RunE: s.runFunc(),
	}
	cmd.Flags().BoolP("watch", "w", false, "Use this flag to watch your developing formulas")

	return cmd
}

func (b buildFormulaCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		workspaces, err := b.workspace.List()
		if err != nil {
			return err
		}

		defaultWorkspace := fmt.Sprintf(workspace.DefaultWorkspaceDirPattern, b.userHomeDir)
		if b.directory.Exists(defaultWorkspace) {
			workspaces[workspace.DefaultWorkspaceName] = defaultWorkspace
		}

		var items []string
		for k, v := range workspaces {
			kv := fmt.Sprintf("%s (%s)", k, v)
			items = append(items, kv)
		}

		items = append(items, newWorkspace)
		selected, err := b.List("Select a formula workspace: ", items)
		if err != nil {
			return err
		}

		var workspaceName string
		var workspacePath string
		var wspace workspace.Workspace
		if selected == newWorkspace {
			workspaceName, err = b.Text("Workspace name: ", true)
			if err != nil {
				return err
			}

			workspacePath, err = b.Text("Workspace path (e.g.: /home/user/github):", true)
			if err != nil {
				return err
			}

			wspace = workspace.Workspace{
				Name: strings.Title(workspaceName),
				Dir:  workspacePath,
			}

			if err := b.workspace.Add(wspace); err != nil {
				return err
			}
		} else {
			split := strings.Split(selected, " ")
			workspaceName = split[0]
			workspacePath = workspaces[workspaceName]
			wspace = workspace.Workspace{
				Name: strings.Title(workspaceName),
				Dir:  workspacePath,
			}

			if err := b.workspace.Validate(wspace); err != nil {
				return err
			}
		}

		formulaPath, err := b.readFormulas(wspace.Dir)
		if err != nil {
			return err
		}

		watch, err := cmd.Flags().GetBool("watch")
		if err != nil {
			return err
		}

		if watch {
			b.watcher.Watch(wspace.Dir, formulaPath)
			return nil
		}

		b.build(wspace.Dir, formulaPath)

		return nil
	}
}

func (b buildFormulaCmd) build(workspacePath, formulaPath string) {
	buildInfo := color.FgRed.Render("Building formula...")
	s := spinner.StartNew(buildInfo)
	time.Sleep(2 * time.Second)

	if err := b.formula.Build(workspacePath, formulaPath); err != nil {
		errorMsg := color.FgRed.Render(err)
		s.Error(errors.New(errorMsg))
		return
	}

	success := color.FgGreen.Render("✔ Build completed!")
	s.Success(success)
	prompt.Info("Now you can run your formula with Ritchie!")
}

func (b buildFormulaCmd) readFormulas(dir string) (string, error) {
	dirs, err := b.directory.List(dir, false)
	if err != nil {
		return "", err
	}

	dirs = sliceutil.Remove(dirs, treeDir)

	if isFormula(dirs) {
		return dir, nil
	}

	selected, err := b.List("Select a formula or group: ", dirs)
	if err != nil {
		return "", err
	}

	dir, err = b.readFormulas(fmt.Sprintf(dirPattern, dir, selected))
	if err != nil {
		return "", err
	}

	return dir, nil
}

func isFormula(dirs []string) bool {
	for _, dir := range dirs {
		if dir == srcDir {
			return true
		}
	}

	return false
}
