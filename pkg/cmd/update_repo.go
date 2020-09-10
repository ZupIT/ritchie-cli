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

type updateRepoCmd struct {
	repo          formula.RepositoryListUpdater
	repoProviders formula.RepoProviders
	prompt.InputList
}

func NewUpdateRepoCmd(
	repo formula.RepositoryListUpdater,
	repoProviders formula.RepoProviders,
	inList prompt.InputList,
) *cobra.Command {
	updateRepo := updateRepoCmd{
		repo:          repo,
		repoProviders: repoProviders,
		InputList:     inList,
	}

	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Update a repository.",
		Example: "rit update repo",
		RunE:    RunFuncE(updateRepo.runStdin(), updateRepo.runPrompt()),
	}
	cmd.LocalFlags()

	return cmd
}

func (up updateRepoCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := up.repo.List()
		if err != nil {
			return err
		}

		var reposName []string
		for i := range repos {
			reposName = append(reposName, repos[i].Name.String())
		}

		name, err := up.List("Select a repository to update: ", reposName)
		if err != nil {
			return err
		}

		var repo formula.Repo
		for i := range repos {
			if repos[i].Name == formula.RepoName(name) {
				repo = repos[i]
				break
			}
		}

		git := up.repoProviders.Resolve(repo.Provider)

		repoInfo := git.NewRepoInfo(repo.Url, repo.Token)
		tags, err := git.Repos.Tags(repoInfo)
		if err != nil {
			return err
		}

		version, err := up.List("Select your new version:", tags.Names())
		if err != nil {
			return err
		}

		if err := up.repo.Update(formula.RepoName(name), formula.RepoVersion(version)); err != nil {
			return err
		}

		successMsg := fmt.Sprintf("The %q repository was updated with success to version %q", name, version)
		prompt.Success(successMsg)

		return nil
	}
}

func (up updateRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		r := formula.Repo{}

		err := stdin.ReadJson(os.Stdin, &r)
		if err != nil {
			return err
		}

		if err := up.repo.Update(r.Name, r.Version); err != nil {
			return err
		}

		successMsg := fmt.Sprintf("The %q repository was updated with success to version %q", r.Name, r.Version)
		prompt.Success(successMsg)

		return nil
	}
}
