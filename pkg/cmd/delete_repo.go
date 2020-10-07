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
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const (
	deleteSuccessMsg = "Repository %q was deleted with success"
)

type deleteRepoCmd struct {
	formula.RepositoryLister
	prompt.InputList
	formula.RepositoryDeleter
	formula.RepositoryLocalDeleter
}

func NewDeleteRepoCmd(rl formula.RepositoryLister, il prompt.InputList, rd formula.RepositoryDeleter, rld formula.RepositoryLocalDeleter) *cobra.Command {
	dr := deleteRepoCmd{rl, il, rd, rld}
	cmd := &cobra.Command{
		Use:       "repo",
		Short:     "Delete a repository",
		Example:   "rit delete repo",
		RunE:      RunFuncE(dr.runStdin(), dr.runFunc()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
	return cmd
}

func (dr deleteRepoCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := dr.RepositoryLister.List()
		if err != nil {
			return err
		}

		if len(repos) == 0 {
			prompt.Warning("You don't have any repositories")
			return nil
		}

		var reposNames []string
		for _, r := range repos {
			reposNames = append(reposNames, r.Name.String())
		}

		repoLocal, err := dr.RepositoryLister.ListLocal()
		if err == nil {
			reposNames = append(reposNames, repoLocal.String())
		}

		repo, err := dr.InputList.List("Repository:", reposNames)
		if err != nil {
			return err
		}

		selectedRepoName := formula.RepoName(repo)

		var errReturn error
		switch selectedRepoName {
		case repoLocal:
			errReturn = dr.DeleteLocal()
		default:
			errReturn = dr.Delete(selectedRepoName)
		}

		if errReturn != nil {
			return err
		}

		prompt.Success(fmt.Sprintf(deleteSuccessMsg, repo))
		return nil
	}
}

func (dr deleteRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		repo := formula.Repo{}

		err := stdin.ReadJson(os.Stdin, &repo)
		if err != nil {
			return err
		}

		if err := dr.Delete(repo.Name); err != nil {
			return err
		}

		prompt.Success(fmt.Sprintf(deleteSuccessMsg, repo.Name))
		return nil
	}
}
