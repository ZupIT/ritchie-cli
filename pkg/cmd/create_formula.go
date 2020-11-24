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
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const newWorkspace = "Type new formula workspace?"

// createFormulaCmd type for add formula command.
type createFormulaCmd struct {
	homeDir         string
	formula         formula.CreateBuilder
	workspace       formula.WorkspaceAddLister
	inText          prompt.InputText
	inTextValidator prompt.InputTextValidator
	inList          prompt.InputList
	tplM            template.Manager
	tutorial        rtutorial.Finder
	tree            tree.CheckerManager
}

// CreateFormulaCmd creates a new cmd instance.
func NewCreateFormulaCmd(
	homeDir string,
	formula formula.CreateBuilder,
	tplM template.Manager,
	workspace formula.WorkspaceAddLister,
	inText prompt.InputText,
	inTextValidator prompt.InputTextValidator,
	inList prompt.InputList,
	rtf rtutorial.Finder,
	treeChecker tree.CheckerManager,
) *cobra.Command {
	c := createFormulaCmd{
		homeDir:         homeDir,
		formula:         formula,
		workspace:       workspace,
		inText:          inText,
		inTextValidator: inTextValidator,
		inList:          inList,
		tplM:            tplM,
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
		formulaCmd, err := c.inTextValidator.Text(
			"Enter the new formula command: ",
			c.surveyCmdValidator,
			"You must create your command based in this example [rit group verb noun]",
		)
		if err != nil {
			return err
		}

		if err := c.tplM.Validate(); err != nil {
			return err
		}

		languages, err := c.tplM.Languages()
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

		wspace, err := FormulaWorkspaceInput(workspaces, c.inList, c.inText)
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

		c.tree.Check()
		c.create(cf)

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

		c.create(cf)
		return nil
	}
}

func (c createFormulaCmd) create(cf formula.Create) {
	buildInfo := prompt.Bold("Creating and building formula...")
	s := spinner.StartNew(buildInfo)
	time.Sleep(2 * time.Second)

	if err := c.formula.Create(cf); err != nil {
		err := prompt.NewError(err.Error())
		s.Error(err)
		return
	}

	info := formula.BuildInfo{FormulaPath: cf.FormulaPath, Workspace: cf.Workspace}
	if err := c.formula.Build(info); err != nil {
		err := prompt.NewError(err.Error())
		s.Error(err)
		return
	}

	tutorialHolder, err := c.tutorial.Find()
	if err != nil {
		s.Error(err)
		return
	}
	createSuccess(s, cf.Lang)
	buildSuccess(cf.FormulaPath, cf.FormulaCmd, tutorialHolder.Current)
}

func createSuccess(s *spinner.Spinner, lang string) {
	msg := fmt.Sprintf("%s formula successfully created!", lang)
	success := prompt.Green(msg)
	s.Success(success)
}

func buildSuccess(formulaPath, formulaCmd, tutorialStatus string) {
	prompt.Info(fmt.Sprintf("Formula path is %s", formulaPath))

	if tutorialStatus == tutorialStatusEnabled {
		tutorialCreateFormula(formulaCmd)
	} else {
		prompt.Info(fmt.Sprintf("Now you can run your formula with the following command %q", formulaCmd))
	}
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
		return prompt.
			NewError("this input must not be empty")
	}

	s := strings.Split(formulaCmd, " ")
	if s[0] != "rit" {
		return prompt.
			NewError("Rit formula's command needs to start with \"rit\" [ex.: rit group verb <noun>]")
	}

	if len(s) <= 2 {
		return prompt.
			NewError("Rit formula's command needs at least 2 words following \"rit\" [ex.: rit group verb]")
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
		return prompt.NewError(`not allowed character on formula name \/,><@-`)
	}
	return nil
}

func FormulaWorkspaceInput(
	workspaces formula.Workspaces,
	inList prompt.InputList,
	inText prompt.InputText,
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

		workspacePath, err = inText.Text("Workspace path (e.g.: /home/user/github):", true)
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
	const messageBody = ` ∙ Run %q
 ∙ Run "rit build formula" to update your changes
 ∙ Run "rit build formula --watch" to have automatic updates` + "\n"

	prompt.Info(tagTutorial)
	prompt.Info(messageTitle)
	fmt.Println(fmt.Sprintf(messageBody, formulaCmd))
}
