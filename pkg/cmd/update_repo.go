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
	"net/http"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const (
	questionSelectARepo = "Select a repository to update: "
	updateOptionAll     = "ALL"
	repoName            = "name"
	repoVersion         = "version"
	successUpdate       = "The %q repository was updated with success to version %q\n"
)

var updateRepoFlags = flags{
	{
		name:        repoName,
		kind:        reflect.String,
		defValue:    "",
		description: "repository name",
	},
	{
		name:        repoVersion,
		kind:        reflect.String,
		defValue:    "latest",
		description: "repository version",
	},
}

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
		RunE:      RunFuncE(updateRepo.runStdin(), updateRepo.runCmd()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
	cmd.LocalFlags()
	addReservedFlags(cmd.Flags(), updateRepoFlags)

	return cmd
}

func (up updateRepoCmd) runCmd() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		reposToUpdate, err := up.resolveInput(cmd)
		if err != nil {
			return err
		}

		for _, value := range reposToUpdate {
			err := up.repo.Update(value.Name, value.Version)
			if err != nil {
				return err
			}
			s := fmt.Sprintf(successUpdate, value.Name, value.Version)
			prompt.Success(s)
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

		s := fmt.Sprintf(successUpdate, r.Name, r.Version)
		prompt.Success(s)

		return nil
	}
}

func (up updateRepoCmd) resolveInput(cmd *cobra.Command) (formula.Repos, error) {
	if IsFlagInput(cmd) {
		return up.resolveFlags(cmd)
	}
	return up.resolvePrompt()
}

func (up updateRepoCmd) resolvePrompt() (formula.Repos, error) {
	repos, err := up.repo.List()
	if err != nil {
		return formula.Repos{}, err
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
		return formula.Repos{}, err
	}

	var repoToUpdate []formula.Repo

	if name == updateOptionAll {
		repoToUpdate = externalRepos
	} else {
		for i := range externalRepos {
			if externalRepos[i].Name == formula.RepoName(name) {
				repoToUpdate = append(repoToUpdate, externalRepos[i])
				break
			}
		}
	}

	for i, currRepo := range repoToUpdate {
		repoInfo, err := up.getRepoInfo(currRepo)
		if err != nil {
			return formula.Repos{}, err
		}

		currRepoVersion := fmt.Sprintf("Select your new version for %q:", currRepo.Name.String())

		version, err := up.List(currRepoVersion, repoInfo)
		if err != nil {
			return formula.Repos{}, err
		}
		repoToUpdate[i].Version = formula.RepoVersion(version)
	}
	return repoToUpdate, nil
}

func (up *updateRepoCmd) resolveFlags(cmd *cobra.Command) (formula.Repos, error) {
	name, err := cmd.Flags().GetString(repoName)
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, errors.New(missingFlagText(repoName))
	}

	version, err := cmd.Flags().GetString(repoVersion)
	if err != nil {
		return nil, err
	}
	if version == "" {
		return formula.Repos{}, errors.New(missingFlagText(repoVersion))
	}

	repoTarget := formula.Repo{Name: formula.RepoName(name), Version: formula.RepoVersion(version)}
	var repoToUpdate []formula.Repo

	repos, err := up.repo.List()
	if err != nil {
		return nil, err
	}

	var externalRepos formula.Repos
	for i := range repos {
		if !repos[i].IsLocal {
			externalRepos = append(externalRepos, repos[i])
		}
	}

	for _, currRepo := range externalRepos {
		if repoTarget.Name == currRepo.Name {
			info, _ := up.getRepoInfo(currRepo)
			if version == "latest" {
				repoTarget.Version = currRepo.LatestVersion
				repoToUpdate = append(repoToUpdate, repoTarget)
				return repoToUpdate, nil
			} else if findVersion(info, repoTarget.Version) {
				repoToUpdate = append(repoToUpdate, repoTarget)
				return repoToUpdate, nil
			} else {
				errorMsg := fmt.Sprintf("The version %q of repository %q was not found.\n", repoTarget.Version, repoTarget.Name)
				return repoToUpdate, errors.New(errorMsg)
			}
		}
	}

	if len(repoToUpdate) == 0 {
		errorMsg := fmt.Sprintf("The repository %q was not found.\n", repoTarget.Name)
		return formula.Repos{}, errors.New(errorMsg)
	}

	return repoToUpdate, nil
}

func (up *updateRepoCmd) getRepoInfo(repoToUpdate formula.Repo) ([]string, error) {
	gitResp := up.repoProviders.Resolve(repoToUpdate.Provider)
	repoInfo := gitResp.NewRepoInfo(repoToUpdate.Url, repoToUpdate.Token)
	tags, err := gitResp.Repos.Tags(repoInfo)
	if err != nil {
		return nil, err
	}
	stringTags := tags.Names()
	return stringTags, nil
}

func findVersion(source []string, value formula.RepoVersion) bool {
	for _, item := range source {
		if item == string(value) {
			return true
		}
	}
	return false
}
