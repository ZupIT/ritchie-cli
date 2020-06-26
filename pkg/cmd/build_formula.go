package cmd

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
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
	workspace   formula.WorkspaceAddListValidator
	formula     formula.Builder
	watcher     formula.Watcher
	directory   stream.DirListChecker
	prompt.InputText
	prompt.InputList
}

func NewBuildFormulaCmd(
	userHomeDir string,
	formula formula.Builder,
	workManager formula.WorkspaceAddListValidator,
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

		defaultWorkspace := path.Join(b.userHomeDir, formula.DefaultWorkspaceDir)
		if b.directory.Exists(defaultWorkspace) {
			workspaces[formula.DefaultWorkspaceName] = defaultWorkspace
		}

		wspace, err := FormulaWorkspaceInput(workspaces, b.InputList, b.InputText)
		if err != nil {
			return err
		}

		if wspace.Dir != defaultWorkspace {
			if err := b.workspace.Validate(wspace); err != nil {
				return err
			}

			if err := b.workspace.Add(wspace); err != nil {
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
	buildInfo := fmt.Sprintf(prompt.Teal, "Building formula...")
	s := spinner.StartNew(buildInfo)
	time.Sleep(2 * time.Second)

	if err := b.formula.Build(workspacePath, formulaPath); err != nil {
		errorMsg := fmt.Sprintf(prompt.Red, err)
		s.Error(errors.New(errorMsg))
		return
	}

	success := fmt.Sprintf(prompt.Green, "✔ Build completed!")
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
