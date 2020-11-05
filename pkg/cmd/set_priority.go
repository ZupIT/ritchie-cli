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

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	newRepositoryPriority = "Now %q repository has priority %v"
)

type setPriorityCmd struct {
	prompt.InputList
	prompt.InputInt
	formula.RepositoryLister
	formula.RepositoryPrioritySetter
}

func NewSetPriorityCmd(
	inList prompt.InputList,
	inInt prompt.InputInt,
	repoLister formula.RepositoryLister,
	repoPriority formula.RepositoryPrioritySetter,
) *cobra.Command {
	s := setPriorityCmd{
		InputList:                inList,
		InputInt:                 inInt,
		RepositoryLister:         repoLister,
		RepositoryPrioritySetter: repoPriority,
	}

	cmd := &cobra.Command{
		Use:       "repo-priority",
		Short:     "Set a repository priority",
		Example:   "rit set repo-priority",
		RunE:      s.runFunc(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	return cmd
}

func (s setPriorityCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repositories, err := s.RepositoryLister.List()
		if err != nil {
			return err
		}

		if len(repositories) == 0 {
			prompt.Warning("You should add a repository first")
			prompt.Info("Command: rit add repo")
			return nil
		}

		var reposNames []string
		for _, r := range repositories {
			reposNames = append(reposNames, r.Name.String())
		}

		repoName, err := s.InputList.List("Repository:", reposNames)
		if err != nil {
			return err
		}

		priority, err := s.InputInt.Int("New priority:")
		if err != nil {
			return err
		}

		var repo formula.Repo
		for _, r := range repositories {
			if r.Name == formula.RepoName(repoName) {
				repo = r
				break
			}
		}

		if err := s.SetPriority(repo.Name, int(priority)); err != nil {
			return err
		}


		successMsg := fmt.Sprintf(newRepositoryPriority, repoName, priority)
		if int(priority) > repositories.Len() {
			successMsg = fmt.Sprintf("Now %q repository has the least priority", repoName)
		}

		prompt.Success(successMsg)
		return nil
	}
}
