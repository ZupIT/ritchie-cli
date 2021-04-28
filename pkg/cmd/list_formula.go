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
	repoFlagDescription = "Repository name to list formulas, use 'ALL' to list formulas from all repositories."
	noRepoFoundMsg      = "You don't have any repositories"
	failedRepoMsg       = "Formulas from %q could not be retrieved."
	emptyRepoMsg        = "Repo %q has no formulas."
)

var (
	errEmptyTree    = errors.New("no formula found in selected repo")
	errRepoNotFound = errors.New("no repository with this name")
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
		RunE:    lf.runCmd(),
	}

	addReservedFlags(cmd.Flags(), listFormulaFlags)

	return cmd
}

func (lr *listFormulaCmd) runCmd() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := lr.resolveInput(cmd)
		if err != nil {
			return err
		} else if len(repos) == 0 {
			return nil
		}

		noFormula := false
		formulaCount, err := lr.printFormulas(repos)
		if err != nil {
			return err
		} else if formulaCount > 1 {
			prompt.Info(fmt.Sprintf(totalFormulasMsg, formulaCount))
		} else if formulaCount == 0 {
			noFormula = true
		} else {
			prompt.Info(totalOneFormulaMsg)
		}

		tutorialHolder, err := lr.rt.Find()
		if err != nil {
			return err
		}
		tutorialListFormulas(tutorialHolder.Current, noFormula)
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
		prompt.Warning(noRepoFoundMsg)
		return formula.Repos{}, nil
	}

	reposNames := []string{listOptionAll}
	for _, r := range repos {
		reposNames = append(reposNames, r.Name.String())
	}

	repoName, err := lr.InputList.List("Repository:", reposNames)
	if err != nil {
		return formula.Repos{}, err
	}

	if repoName != listOptionAll {
		for i := range repos {
			if repos[i].Name == formula.RepoName(repoName) {
				return formula.Repos{repos[i]}, nil
			}
		}
	}

	return repos, nil
}

func (lr *listFormulaCmd) resolveFlags(cmd *cobra.Command) (formula.Repos, error) {
	name, err := cmd.Flags().GetString(nameFlagName)
	if err != nil {
		return formula.Repos{}, err
	} else if name == "" {
		return formula.Repos{}, errors.New(missingFlagText(nameFlagName))
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
	failedRepos := make([]string, 0)
	emptyRepos := make([]string, 0)
	for _, r := range repos {
		repoFormulas, err := lr.formulasByRepo(r.Name)
		if err != nil {
			if len(repos) == 1 {
				if err != errEmptyTree {
					return 0, err
				}
				prompt.Warning(fmt.Sprintf(emptyRepoMsg, r.Name.String()))
			}

			if err != errEmptyTree {
				failedRepos = append(failedRepos, r.Name.String())
				continue
			}
			emptyRepos = append(emptyRepos, r.Name.String())
		}
		allFormulas = append(allFormulas, repoFormulas...)
	}

	if len(allFormulas) == 0 {
		return 0, nil
	}

	sort.Sort(ByCmd(allFormulas))
	for _, fm := range allFormulas {
		table.AddRow(fm.Cmd, fm.Desc)
	}
	raw := table.Bytes()
	raw = append(raw, []byte("\n")...)
	fmt.Println(string(raw))

	for _, r := range failedRepos {
		prompt.Warning(fmt.Sprintf(failedRepoMsg, r))
	}

	for _, r := range emptyRepos {
		prompt.Warning(fmt.Sprintf(emptyRepoMsg, r))
	}

	return len(table.Rows) - 1, nil
}

func (lr listFormulaCmd) formulasByRepo(repoName formula.RepoName) ([]formulaDefinition, error) {
	tree, err := lr.treeManager.TreeByRepo(repoName)
	if err != nil {
		return []formulaDefinition{}, err
	} else if tree.Commands == nil {
		return []formulaDefinition{}, errRepoNotFound
	} else if len(tree.Commands) == 0 {
		return []formulaDefinition{}, errEmptyTree
	}

	var repoFormulas []formulaDefinition
	replacer := strings.NewReplacer(rootString, rootCommand, commandSeparator, " ")
	for key, cmd := range tree.Commands {
		if cmd.Formula {
			fd := formulaDefinition{
				Cmd:  replacer.Replace(key.String()),
				Desc: cmd.Help,
			}
			repoFormulas = append(repoFormulas, fd)
		}
	}

	return repoFormulas, nil
}

func tutorialListFormulas(tutorialStatus string, emptyRepo bool) {
	const tagTutorial = "\n[TUTORIAL]"
	var MessageTitle, MessageBody string
	if emptyRepo {
		MessageTitle = "To create a formula:"
		MessageBody = ` ∙ Run "rit create formula"`
	} else {
		MessageTitle = "To delete a formula repository:"
		MessageBody = ` ∙ Run "rit delete repo"`
	}

	if tutorialStatus == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
	}
}
