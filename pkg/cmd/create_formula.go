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
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const (
	newWorkspace     = "Type new formula workspace?"
	formulaCmdLabel  = "Enter the new formula command: "
	formulaCmdHelper = "You must create your command based in this example [rit group verb noun]"
)

var (
	ErrFormulaCmdNotBeEmpty        = errors.New("this input must not be empty")
	ErrFormulaCmdMustStartWithRit  = errors.New("rit formula's command needs to start with \"rit\" [ex.: rit group verb <noun>]")
	ErrInvalidFormulaCmdSize       = errors.New("rit formula's command needs at least 2 words following \"rit\" [ex.: rit group verb]")
	ErrInvalidCharactersFormulaCmd = errors.New(`these characters are not allowed in the formula command [\ /,> <@ -]`)
)

// createFormulaCmd type for add formula command.
type createFormulaCmd struct {
	homeDir         string
	formula         formula.CreateBuilder
	workspace       formula.WorkspaceAddListHasher
	inText          prompt.InputText
	inTextValidator prompt.InputTextValidator
	inList          prompt.InputList
	inPath          prompt.InputPath
	template        template.Manager
	tutorial        rtutorial.Finder
	tree            formula.TreeChecker
}

// CreateFormulaCmd creates a new cmd instance.
func NewCreateFormulaCmd(
	homeDir string,
	formula formula.CreateBuilder,
	tplM template.Manager,
	workspace formula.WorkspaceAddListHasher,
	inText prompt.InputText,
	inTextValidator prompt.InputTextValidator,
	inList prompt.InputList,
	inPath prompt.InputPath,
	rtf rtutorial.Finder,
	treeChecker formula.TreeChecker,
) *cobra.Command {
	c := createFormulaCmd{
		homeDir:         homeDir,
		formula:         formula,
		workspace:       workspace,
		inText:          inText,
		inTextValidator: inTextValidator,
		inList:          inList,
		inPath:          inPath,
		template:        tplM,
		tutorial:        rtf,
		tree:            treeChecker,
	}

	cmd := &cobra.Command{
		Use:       "formula",
		Short:     "Create a new formula",
		Example:   "rit create formula",
		RunE:      RunFuncE(c.runStdin(), c.runPrompt()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	cmd.LocalFlags()

	return cmd
}

func (c createFormulaCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		formulaCmd, err := c.inTextValidator.Text(formulaCmdLabel, c.surveyCmdValidator, formulaCmdHelper)
		if err != nil {
			return err
		}

		if err := c.template.Validate(); err != nil {
			return err
		}

		languages, err := c.template.Languages()
		if err != nil {
			return err
		}

		lang, err := c.inList.List("Choose the language: ", languages)
		if err != nil {
			return err
		}

		workspaces, err := c.workspace.List()
		if err != nil {
			return err
		}

		wspace, err := FormulaWorkspaceInput(workspaces, c.inList, c.inText, c.inPath)
		if err != nil {
			return err
		}

		if err := c.workspace.Add(wspace); err != nil {
			return err
		}

		formulaPath := formulaPath(wspace.Dir, formulaCmd)

		cf := formula.Create{
			FormulaCmd:  formulaCmd,
			Lang:        lang,
			Workspace:   wspace,
			FormulaPath: formulaPath,
		}

		check := c.tree.Check()

		printConflictingCommandsWarning(check)

		if err := c.create(cf); err != nil {
			return err
		}

		return nil
	}
}

func (c createFormulaCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		var cf formula.Create

		if err := stdin.ReadJson(os.Stdin, &cf); err != nil {
			return err
		}

		if err := formulaCommandValidator(cf.FormulaCmd); err != nil {
			return err
		}

		if err := c.create(cf); err != nil {
			return err
		}

		return nil
	}
}

func (c createFormulaCmd) create(cf formula.Create) error {
	if err := c.formula.Create(cf); err != nil {
		return err
	}

	info := formula.BuildInfo{FormulaPath: cf.FormulaPath, Workspace: cf.Workspace}
	if err := c.formula.Build(info); err != nil {
		return err
	}

	hash, err := c.workspace.CurrentHash(cf.FormulaPath)
	if err != nil {
		return err
	}

	if err := c.workspace.UpdateHash(cf.FormulaPath, hash); err != nil {
		return err
	}

	successMsg := fmt.Sprintf("%s formula successfully created!", cf.Lang)
	prompt.Success(successMsg)

	tutorialHolder, err := c.tutorial.Find()
	if err != nil {
		return err
	}

	buildSuccess(cf.FormulaPath, cf.FormulaCmd, tutorialHolder.Current)

	return nil
}

func buildSuccess(formulaPath, formulaCmd, tutorialStatus string) {
	prompt.Info(fmt.Sprintf("Formula path is %s", formulaPath))

	if tutorialStatus == tutorialStatusEnabled {
		tutorialCreateFormula(formulaCmd)
		return
	}

	prompt.Info(fmt.Sprintf("Now you can run your formula with the following command %q", formulaCmd))
}

func formulaPath(workspacePath, cmd string) string {
	cc := strings.Split(cmd, " ")
	formulaPath := strings.Join(cc[1:], string(os.PathSeparator))
	return filepath.Join(workspacePath, formulaPath)
}

func (c createFormulaCmd) surveyCmdValidator(cmd interface{}) error {
	if err := formulaCommandValidator(cmd.(string)); err != nil {
		return err
	}

	return nil
}

func formulaCommandValidator(formulaCmd string) error {
	if len(strings.TrimSpace(formulaCmd)) < 1 {
		return ErrFormulaCmdNotBeEmpty
	}

	s := strings.Split(formulaCmd, " ")
	if s[0] != "rit" {
		return ErrFormulaCmdMustStartWithRit
	}

	if len(s) <= 2 {
		return ErrInvalidFormulaCmdSize
	}

	if err := characterValidator(formulaCmd); err != nil {
		return err
	}

	if err := coreCmdValidator(formulaCmd); err != nil {
		return err
	}

	return nil
}

func coreCmdValidator(formulaCmd string) error {
	wordAfterCore := strings.Split(formulaCmd, " ")[1]
	for i := range api.CoreCmds {
		if wordAfterCore == api.CoreCmds[i].Usage {
			errorString := fmt.Sprintf("core command verb %q after rit\n"+
				"Use your formula group before the verb\n"+
				"Example: rit aws list bucket\n",
				api.CoreCmds[i].Usage)

			return errors.New(errorString)
		}
	}
	return nil
}

func characterValidator(formula string) error {
	if strings.ContainsAny(formula, `\/><,@`) {
		return ErrInvalidCharactersFormulaCmd
	}
	return nil
}

func FormulaWorkspaceInput(
	workspaces formula.Workspaces,
	inList prompt.InputList,
	inText prompt.InputText,
	inPath prompt.InputPath,
) (formula.Workspace, error) {
	items := make([]string, 0, len(workspaces))
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

		workspacePath, err = inPath.Read("Workspace path (e.g.: /home/user/github): ")
		if err != nil {
			return formula.Workspace{}, err
		}

		wspace = formula.Workspace{
			Name: strings.Title(workspaceName),
			Dir:  workspacePath,
		}
	} else {
		split := strings.Split(selected, " (")
		workspaceName = split[0]
		workspacePath = workspaces[workspaceName]
		wspace = formula.Workspace{
			Name: strings.Title(workspaceName),
			Dir:  workspacePath,
		}
	}
	return wspace, nil
}

func tutorialCreateFormula(formulaCmd string) {
	const tagTutorial = "\n[TUTORIAL]"
	const messageTitle = "In order to test your new formula:"
	const messageBody = ` ∙ Simply edit the formula files and run %q again` + "\n"

	prompt.Info(tagTutorial)
	prompt.Info(messageTitle)
	fmt.Println(fmt.Sprintf(messageBody, formulaCmd))
}
