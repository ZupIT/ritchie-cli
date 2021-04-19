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
	"fmt"
	"sort"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

const (
	listOptionAll      = "ALL"
	rootString         = "root"
	rootCommand        = "rit"
	commandSeparator   = "_"
	totalFormulasMsg   = "There are %v formulas"
	totalOneFormulaMsg = "There is 1 formula"
)

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
		RunE:    lf.runFunc(),
	}
	return cmd
}

func (lr listFormulaCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := lr.RepositoryLister.List()
		if err != nil {
			return err
		}

		if len(repos) == 0 {
			prompt.Warning("You don't have any repositories")
			return nil
		}

		var reposNames []string
		reposNames = append(reposNames, listOptionAll)
		for _, r := range repos {
			reposNames = append(reposNames, r.Name.String())
		}

		repoName, err := lr.InputList.List("Repository:", reposNames)
		if err != nil {
			return err
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

		formulaCount, err := lr.printFormulas(repoToListFormulas)
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

func (lr listFormulaCmd) printFormulas(repos formula.Repos) (formulaCount int, err error) {
	table := uitable.New()
	table.AddRow("COMMAND", "DESCRIPTION")
	var allFormulas []formulaDefinition
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
