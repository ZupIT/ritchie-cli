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
		defValue:    "",
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
		_, err := up.resolveInput(cmd)
		if err != nil {
			return err
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

func (up updateRepoCmd) resolveInput(cmd *cobra.Command) (updateRepoCmd, error) {
	if IsFlagInput(cmd) {
		return up.resolveFlags(cmd)
	}
	return up.resolvePrompt()
}

func (up updateRepoCmd) resolvePrompt() (updateRepoCmd, error) {
	repos, err := up.repo.List()
	if err != nil {
		return updateRepoCmd{}, err
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
		return updateRepoCmd{}, err
	}

	flagAll := name == updateOptionAll

	var repoToUpdate []formula.Repo

	if flagAll {
		repoToUpdate = externalRepos
	} else {
		for i := range externalRepos {
			if externalRepos[i].Name == formula.RepoName(name) {
				repoToUpdate = append(repoToUpdate, externalRepos[i])
				break
			}
		}
	}

	for _, currRepo := range repoToUpdate {

		repoInfo, err := up.getRepoInfo(currRepo)

		if err != nil {
			return updateRepoCmd{}, err
		}

		currRepoVersion := fmt.Sprintf("Select your new version for %q:", currRepo.Name.String())

		version, err := up.List(currRepoVersion, repoInfo)
		if err != nil {
			return updateRepoCmd{}, err
		}

		currRepoName := string(currRepo.Name)

		if err := up.repo.Update(formula.RepoName(currRepoName), formula.RepoVersion(version)); err != nil {
			return updateRepoCmd{}, err
		}

		successMsg := fmt.Sprintf("The %q repository was updated with success to version %q\n", currRepo.Name, version)
		prompt.Success(successMsg)
	}


	// 	repositoryList, err := up.repo.List()
	// 	if err != nil {
	// 		return formula.Repo{}, err
	// 	}
	//
	// 	for index := range repositoryList {
	// 		if repositoryList[index].Name.String() == name {
	// 			tagsx, err := up.getRepoInfo(repositoryList[index])
	// 			for idx := range tagsx {
	// 				if tagsx[idx] == version {
	// 					return formula.Repo{Name: formula.RepoName(name), Version: formula.RepoVersion(version)}, err
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	return updateRepoCmd{}, nil

}

func (up *updateRepoCmd) resolveFlags(cmd *cobra.Command) (updateRepoCmd, error) {
	name, err := cmd.Flags().GetString(repoName)
	if err != nil {
		return updateRepoCmd{}, err
	}
	version, err := cmd.Flags().GetString(repoVersion)
	if err != nil {
		return updateRepoCmd{}, err
	}

	if name == "" {
		return updateRepoCmd{}, errors.New("please provide a value for 'name'")
	} else if version == "" {
		return updateRepoCmd{}, errors.New("please provide a value for 'version'")
	}
	repoTarget := formula.Repo{Name: formula.RepoName(name), Version: formula.RepoVersion(version)}

	updateStart := up.repo.Update(repoTarget.Name, repoTarget.Version)
	if updateStart != nil {
		return updateRepoCmd{}, updateStart
	}
	successMsg := fmt.Sprintf("The %q repository was updated with success to version %q",repoTarget.Name, repoTarget.Version)
	prompt.Success(successMsg)
	return updateRepoCmd{}, err

}

func (up *updateRepoCmd) getRepoInfo(repoToUpdate formula.Repo) ([]string, error) {
	gitResp := up.repoProviders.Resolve(repoToUpdate.Provider)
	repoInfo := gitResp.NewRepoInfo(repoToUpdate.Url, repoToUpdate.Token)
	tags, err := gitResp.Repos.Tags(repoInfo)
	fmt.Printf("%v", repoInfo)
	if err != nil {
		return nil, err
	}
	stringTags := tags.Names()
	return stringTags, err
}
