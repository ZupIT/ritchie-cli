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

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	deleteSuccessMsg    = "Repository %q was deleted with success"
	repoFlagDescription = "Repository name to delete"
)

var deleteRepoFlags = flags{
	{
		name:        nameFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: repoFlagDescription,
	},
}

type deleteRepoCmd struct {
	formula.RepositoryLister
	prompt.InputList
	prompt.InputBool
	formula.RepositoryDeleter
}

func NewDeleteRepoCmd(
	rl formula.RepositoryLister,
	il prompt.InputList,
	ib prompt.InputBool,
	rd formula.RepositoryDeleter,
) *cobra.Command {
	dr := deleteRepoCmd{rl, il, ib, rd}

	cmd := &cobra.Command{
		Use:       "repo",
		Short:     "Delete a repository",
		Example:   "rit delete repo",
		RunE:      dr.runFormula(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), deleteRepoFlags)

	return cmd
}

func (dr deleteRepoCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		name, err := dr.resolveInput(cmd)
		if err != nil {
			return err
		} else if name == "" {
			return nil
		}

		if err := dr.Delete(formula.RepoName(name)); err != nil {
			return err
		}

		prompt.Success(fmt.Sprintf(deleteSuccessMsg, name))
		return nil
	}
}

func (dr *deleteRepoCmd) resolveInput(cmd *cobra.Command) (string, error) {
	if IsFlagInput(cmd) {
		return dr.resolveFlags(cmd)
	}
	return dr.resolvePrompt()
}

func (dr *deleteRepoCmd) resolvePrompt() (string, error) {
	repos, err := dr.RepositoryLister.List()
	if err != nil {
		return "", err
	}

	if len(repos) == 0 {
		prompt.Warning("You don't have any repositories")
		return "", nil
	}

	reposNames := make([]string, 0, len(repos))
	for _, r := range repos {
		reposNames = append(reposNames, r.Name.String())
	}

	repo, err := dr.InputList.List("Repository:", reposNames)
	if err != nil {
		return "", err
	}

	question := "Are you sure you want to delete this repo?"
	ans, err := dr.InputBool.Bool(question, []string{"no", "yes"})
	if err != nil {
		return "", err
	}
	if !ans {
		return "", nil
	}
	return repo, nil
}

func (dr *deleteRepoCmd) resolveFlags(cmd *cobra.Command) (string, error) {
	name, err := cmd.Flags().GetString(nameFlagName)
	if err != nil {
		return "", err
	} else if name == "" {
		return "", errors.New("please provide a value for 'name'")
	}

	return name, nil
}
