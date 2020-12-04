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
	"net/http"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const (
	questionSelectARepo = "Select a repository to update: "
	updateOptionAll     = "ALL"
)

type updateRepoCmd struct {
	client        *http.Client
	repo          formula.RepositoryListUpdater
	repoProviders formula.RepoProviders
	prompt.InputText
	prompt.InputPassword
	prompt.InputURL
	prompt.InputList
	prompt.InputBool
	prompt.InputInt
}

func NewUpdateRepoCmd(
	client *http.Client,
	repo formula.RepositoryListUpdater,
	repoProviders formula.RepoProviders,
	inText prompt.InputText,
	inPass prompt.InputPassword,
	inURL prompt.InputURL,
	inList prompt.InputList,
	inBool prompt.InputBool,
	inInt prompt.InputInt,
) *cobra.Command {
	updateRepo := updateRepoCmd{
		client:        client,
		repo:          repo,
		repoProviders: repoProviders,
		InputText:     inText,
		InputURL:      inURL,
		InputList:     inList,
		InputBool:     inBool,
		InputInt:      inInt,
		InputPassword: inPass,
	}

	cmd := &cobra.Command{
		Use:       "repo",
		Short:     "Update a repository.",
		Example:   "rit update repo",
		RunE:      RunFuncE(updateRepo.runStdin(), updateRepo.runPrompt()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
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
		var externalRepos formula.Repos
		reposName = append(reposName, updateOptionAll)
		for i := range repos {
			if !repos[i].IsLocal {
				externalRepos = append(externalRepos, repos[i])
				reposName = append(reposName, repos[i].Name.String())
			}
		}

		helper := "Select a repository to update your version. P.S. Local repositories cannot be updated."
		name, err := up.List(questionSelectARepo, reposName, helper)
		if err != nil {
			return err
		}

		flagAll := (name == updateOptionAll)

		var repoToUpdate []formula.Repo

		if flagAll {
			repoToUpdate = externalRepos
		} else {
			for i := range externalRepos {
				if repos[i].Name == formula.RepoName(name) {
					repoToUpdate = append(repoToUpdate, repos[i])
					break
				}
			}
		}

		for _, currRepo := range repoToUpdate {

			git := up.repoProviders.Resolve(currRepo.Provider)

			repoInfo := git.NewRepoInfo(currRepo.Url, currRepo.Token)
			tags, err := git.Repos.Tags(repoInfo)
			if err != nil {
				return err
			}

			currRepoVersion := fmt.Sprintf("Select your new version for %q:", currRepo.Name.String())

			version, err := up.List(currRepoVersion, tags.Names())
			if err != nil {
				return err
			}

			currRepoName := string(currRepo.Name)

			if err := up.repo.Update(formula.RepoName(currRepoName), formula.RepoVersion(version)); err != nil {
				return err
			}

			successMsg := fmt.Sprintf("The %q repository was updated with success to version %q\n", currRepo.Name, version)
			prompt.Success(successMsg)
		}

		return nil
	}
}

func (up updateRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		r := formula.Repo{}

		if err := stdin.ReadJson(cmd.InOrStdin(), &r); err != nil {
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
