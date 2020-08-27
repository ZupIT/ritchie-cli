/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	newWorkspace = "Type new formula workspace?"
	docsDir      = "docs"
	srcDir       = "src"
)

type buildFormulaCmd struct {
	userHomeDir string
	workspace   formula.WorkspaceAddListValidator
	formula     formula.LocalBuilder
	watcher     formula.Watcher
	directory   stream.DirListChecker
	prompt.InputText
	prompt.InputList
	rt rtutorial.Finder
}

func NewBuildFormulaCmd(
	userHomeDir string,
	formula formula.LocalBuilder,
	workManager formula.WorkspaceAddListValidator,
	watcher formula.Watcher,
	directory stream.DirListChecker,
	inText prompt.InputText,
	inList prompt.InputList,
	rtf rtutorial.Finder,
) *cobra.Command {
	s := buildFormulaCmd{
		userHomeDir: userHomeDir,
		workspace:   workManager,
		formula:     formula,
		watcher:     watcher,
		directory:   directory,
		InputText:   inText,
		InputList:   inList,
		rt:          rtf,
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

		defaultWorkspace := filepath.Join(b.userHomeDir, formula.DefaultWorkspaceDir)
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

		formulaPath, err := b.readFormulas(wspace.Dir, "rit")
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

		tutorialHolder, err := b.rt.Find()
		if err != nil {
			return err
		}
		tutorialBuildFormula(tutorialHolder.Current)

		return nil
	}
}

func (b buildFormulaCmd) build(workspacePath, formulaPath string) {
	buildInfo := prompt.Red("Building formula...")
	s := spinner.StartNew(buildInfo)
	time.Sleep(2 * time.Second)

	if err := b.formula.Build(workspacePath, formulaPath); err != nil {
		errorMsg := prompt.Red(err.Error())
		s.Error(errors.New(errorMsg))
		return
	}

	success := prompt.Green("✔ Build completed!")
	s.Success(success)
}

func (b buildFormulaCmd) readFormulas(dir string, currentFormula string) (string, error) {
	dirs, err := b.directory.List(dir, false)
	if err != nil {
		return "", err
	}

	dirs = sliceutil.Remove(dirs, docsDir)

	var formulaOptions []string
	var response string
	otherFormula := "Another formula"

	if isFormula(dirs) {
		if !hasFormulaInDir(dirs) {
			return dir, nil
		}

		formulaOptions = append(formulaOptions, currentFormula, otherFormula)

		response, err = b.List("We found a formula, which one do you want to run the build: ", formulaOptions)
		if err != nil {
			return "", err
		}
		if response == currentFormula {
			return dir, nil
		}
		dirs = sliceutil.Remove(dirs, srcDir)
	}

	selected, err := b.List("Select a formula or group: ", dirs)
	if err != nil {
		return "", err
	}

	newFormulaSelected := fmt.Sprintf("%s %s", currentFormula, selected)
	dir, err = b.readFormulas(filepath.Join(dir, selected), newFormulaSelected)
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

func hasFormulaInDir(dirs []string) bool {
	dirs = sliceutil.Remove(dirs, docsDir)
	dirs = sliceutil.Remove(dirs, srcDir)

	return len(dirs) > 0
}

func tutorialBuildFormula(tutorialStatus string) {
	const tagTutorial = "\n[TUTORIAL]"
	const titleNewRepositories = "To add a new repository of formulas:"
	const bodyNewRepositories = ` ∙ Run "rit add repo"`

	const titlePublishFormula = "To publish your formula:"
	const bodyPublishFormula = ` ∙ Create a git repo
 ∙ Commit and push your formula in repo created
 ∙ Create a GitHub or Gitlab release
 ∙ Run "rit add repo"`

	if tutorialStatus == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(titleNewRepositories)
		fmt.Println(bodyNewRepositories)
		prompt.Info(titlePublishFormula)
		fmt.Println(bodyPublishFormula)
	}
}
