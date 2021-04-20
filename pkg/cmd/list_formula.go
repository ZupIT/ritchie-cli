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
	"reflect"
	"sort"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

const (
	listOptionAll       = "ALL"
	rootString          = "root"
	rootCommand         = "rit"
	commandSeparator    = "_"
	totalFormulasMsg    = "There are %v formulas"
	totalOneFormulaMsg  = "There is 1 formula"
	repoFlagDescription = "Repository name to list formulas"
)

var listFormulaFlags = flags{
	{
		name:        nameFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: repoFlagDescription,
	},
}

type formulaDefinition struct {
	Cmd  string
	Desc string
}

type ByCmd []formulaDefinition

func (a ByCmd) Len() int           { return len(a) }
func (a ByCmd) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCmd) Less(i, j int) bool { return a[i].Cmd < a[j].Cmd }

type listFormulaCmd struct {
	formula.RepositoryLister
	prompt.InputList
	treeManager formula.TreeManager
	rt          rtutorial.Finder
}

func NewListFormulaCmd(
	rl formula.RepositoryLister,
	il prompt.InputList,
	tm formula.TreeManager,
	rtf rtutorial.Finder,
) *cobra.Command {
	lf := listFormulaCmd{rl, il, tm, rtf}
	cmd := &cobra.Command{
		Use:     "formula",
		Short:   "Show a list with available formulas from a specific repository",
		Example: "rit list formula",
		RunE:    lf.runFormula(),
	}

	addReservedFlags(cmd.Flags(), listFormulaFlags)

	return cmd
}

func (lr listFormulaCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := lr.resolveInput(cmd)
		if err != nil {
			return err
		}

		formulaCount, err := lr.printFormulas(repos)
		if err != nil {
			return err
		} else if formulaCount != 1 {
			prompt.Info(fmt.Sprintf(totalFormulasMsg, formulaCount))
		} else {
			prompt.Info(totalOneFormulaMsg)
		}

		tutorialHolder, err := lr.rt.Find()
		if err != nil {
			return err
		}
		tutorialListFormulas(tutorialHolder.Current)
		return nil
	}
}

func (lr *listFormulaCmd) resolveInput(cmd *cobra.Command) (formula.Repos, error) {
	if IsFlagInput(cmd) {
		return lr.resolveFlags(cmd)
	}
	return lr.resolvePrompt()
}

func (lr *listFormulaCmd) resolvePrompt() (formula.Repos, error) {
	repos, err := lr.RepositoryLister.List()
	if err != nil {
		return formula.Repos{}, err
	}

	if len(repos) == 0 {
		prompt.Warning("You don't have any repositories")
		return formula.Repos{}, nil
	}

	reposNames := make([]string, 0, len(repos)+1)
	reposNames = append(reposNames, listOptionAll)
	for _, r := range repos {
		reposNames = append(reposNames, r.Name.String())
	}

	repoName, err := lr.InputList.List("Repository:", reposNames)
	if err != nil {
		return formula.Repos{}, err
	}

	var repoToListFormulas []formula.Repo
	if repoName == listOptionAll {
		repoToListFormulas = repos
	} else {
		for i := range repos {
			if repos[i].Name == formula.RepoName(repoName) {
				repoToListFormulas = append(repoToListFormulas, repos[i])
				break
			}
		}
	}
	return repoToListFormulas, nil
}

func (lr listFormulaCmd) resolveFlags(cmd *cobra.Command) (formula.Repos, error) {
	name, err := cmd.Flags().GetString(nameFlagName)
	if err != nil {
		return formula.Repos{}, err
	} else if name == "" {
		return formula.Repos{}, errors.New("please provide a value for 'name'")
	}

	if name == listOptionAll {
		repos, err := lr.RepositoryLister.List()
		if err != nil {
			return formula.Repos{}, err
		}
		return repos, nil
	} else {
		return formula.Repos{formula.Repo{Name: formula.RepoName(name)}}, nil
	}
}

func (lr listFormulaCmd) printFormulas(repos formula.Repos) (formulaCount int, err error) {
	table := uitable.New()
	table.AddRow("COMMAND", "DESCRIPTION")
	allFormulas := make([]formulaDefinition, 0)
	for _, r := range repos {
		repoFormulas, err := lr.formulasByRepo(r.Name)
		if err != nil {
			return 0, err
		}
		allFormulas = append(allFormulas, repoFormulas...)
	}

	sort.Sort(ByCmd(allFormulas))
	for _, fm := range allFormulas {
		table.AddRow(fm.Cmd, fm.Desc)
	}
	raw := table.Bytes()
	raw = append(raw, []byte("\n")...)
	fmt.Println(string(raw))

	return len(table.Rows) - 1, nil
}

func (lr listFormulaCmd) formulasByRepo(repoName formula.RepoName) ([]formulaDefinition, error) {
	tree, err := lr.treeManager.TreeByRepo(repoName)
	if err != nil {
		return []formulaDefinition{}, err
	} else if tree.Commands == nil {
		return []formulaDefinition{}, errors.New("no repository with this name")
	}

	var repoFormulas []formulaDefinition
	for _, cmd := range tree.Commands {
		if cmd.Formula {
			replacer := strings.NewReplacer(rootString, rootCommand, commandSeparator, " ")
			parentFormula := replacer.Replace(cmd.Parent)

			var fd formulaDefinition
			fd.Cmd = strings.Join([]string{parentFormula, cmd.Usage}, " ")
			fd.Desc = cmd.Help
			repoFormulas = append(repoFormulas, fd)
		}
	}

	return repoFormulas, nil
}

func tutorialListFormulas(tutorialStatus string) {
	const tagTutorial = "\n[TUTORIAL]"
	const MessageTitle = "To delete a formula repository:"
	const MessageBody = ` ∙ Run "rit delete repo"`

	if tutorialStatus == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
	}
}
