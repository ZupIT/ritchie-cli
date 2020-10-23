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
	formula.RepositoryListerLocal
	prompt.InputList
	formula.RepositoryDeleter
	formula.RepositoryLocalDeleter
}

// NewDeleteRepoCmd is the constructor for delete repo command
func NewDeleteRepoCmd(
	rl formula.RepositoryLister,
	rll formula.RepositoryListerLocal,
	il prompt.InputList,
	rd formula.RepositoryDeleter,
	rld formula.RepositoryLocalDeleter,
) *cobra.Command {
	dr := deleteRepoCmd{rl, rll, il, rd, rld}
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
		repoLocal, err := dr.RepositoryListerLocal.List()
		repoLocalExists := err == nil

		if len(repos) == 0 && !repoLocalExists {
			prompt.Warning("You don't have any repositories")
			return nil
		}

		var reposNames []string
		for _, r := range repos {
			reposNames = append(reposNames, r.Name.String())
		}

		if repoLocalExists {
			reposNames = append(reposNames, repoLocal.String())
		}

		repo, err := dr.InputList.List("Repository:", reposNames)
		if err != nil {
			return err
		}

		selectedRepoName := formula.RepoName(repo)

		if selectedRepoName == repoLocal {
			err = dr.RepositoryLocalDeleter.Delete()
		} else {
			err = dr.RepositoryDeleter.Delete(selectedRepoName)
		}

		if err != nil {
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

		repoLocal, _ := dr.RepositoryListerLocal.List()

		if repo.Name == repoLocal {
			err = dr.RepositoryLocalDeleter.Delete()
		} else {
			err = dr.RepositoryDeleter.Delete(repo.Name)
		}

		if err != nil {
			return err
		}

		prompt.Success(fmt.Sprintf(deleteSuccessMsg, repo.Name))
		return nil
	}
}
