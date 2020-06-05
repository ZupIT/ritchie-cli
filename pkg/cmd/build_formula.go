package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	newWorkspace = "Type new formula workspace?"
	dirPattern   = "%s/%s"
	treeDir      = "tree"
	srcDir       = "src"
)

type buildFormulaCmd struct {
	userHomeDir string
	workspace   workspace.AddLister
	formula     formula.Builder
	watcher     formula.Watcher
	prompt.InputText
	prompt.InputList
}

func NewBuildFormulaCmd(
	userHomeDir string,
	workManager workspace.AddLister,
	formula formula.Builder,
	watcher formula.Watcher,
	inText prompt.InputText,
	inList prompt.InputList,
) *cobra.Command {
	s := buildFormulaCmd{
		userHomeDir: userHomeDir,
		workspace:   workManager,
		watcher:     watcher,
		formula:     formula,
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

		workspaces[workspace.DefaultWorkspaceName] = fmt.Sprintf(workspace.DefaultWorkspaceDirPattern, b.userHomeDir)

		var test []string
		for k, v := range workspaces {
			kv := fmt.Sprintf("%s (%s)", k, v)
			test = append(test, kv)
		}

		test = append(test, newWorkspace)
		selected, err := b.List("Select a formula workspace: ", test)
		if err != nil {
			return err
		}

		var workspaceName string
		var workspacePath string
		var wspace workspace.Workspace
		if selected == newWorkspace {
			workspaceName, err = b.Text("Type a new formula workspace name: ", true)
			if err != nil {
				return err
			}

			workspacePath, err = b.Text("Type a new formula workspace path: ", true)
			if err != nil {
				return err
			}

			wspace = workspace.Workspace{
				Name: strings.Title(workspaceName),
				Dir:  fmt.Sprintf(dirPattern, b.userHomeDir, workspacePath),
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
			b.watcher.Watch(formulaPath)
			return nil
		}

		fmt.Printf(prompt.Info, "Building formula... \n")
		if err := b.formula.Build(wspace.Dir, formulaPath); err != nil {
			return err
		}

		fmt.Printf(prompt.Success, "Formula built with success \\o/ \n")
		return nil
	}
}

func (b buildFormulaCmd) readFormulas(dir string) (string, error) {
	open, err := os.Open(dir)
	if err != nil {
		return "", err
	}

	fileInfos, err := open.Readdir(0)
	if err != nil {
		return "", err
	}

	formulas, isFormula := filterDir(fileInfos)

	if isFormula {
		return dir, nil
	}

	selected, err := b.List("Select a formula you want to build: ", formulas)
	if err != nil {
		return "", err
	}

	dir, err = b.readFormulas(fmt.Sprintf(dirPattern, dir, selected))
	if err != nil {
		return "", err
	}

	return dir, nil
}

func filterDir(fileInfos []os.FileInfo) ([]string, bool) {
	var dirs []string
	var isFormula bool
	for _, fileInfo := range fileInfos {
		n := fileInfo.Name()
		if n == srcDir {
			isFormula = true
			break
		}

		if fileInfo.IsDir() && n != treeDir && !strings.ContainsAny(n, ".") {
			dirs = append(dirs, n)
		}
	}

	return dirs, isFormula
}
