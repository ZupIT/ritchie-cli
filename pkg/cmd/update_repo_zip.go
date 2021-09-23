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

/*
	Regras:
		- Todos os campos são de texto
		- Quando informo ALL no update repo, qual deveria ser o comportamento?
		- Validar:
			- Se o repo existe no repositories.json (nome)
			- Se a versão é diferente da informada no arquivo repositories.json
*/

var updateRepoZipFlags = flags{
	{
		name:        repoName,
		kind:        reflect.String,
		defValue:    "",
		description: "repository name",
	},
	{
		name:        repoUrlFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: "repository url",
	},
	{
		name:        repoVersion,
		kind:        reflect.String,
		defValue:    "",
		description: "repository version",
	},
}

var ErrSameVersion = errors.New("the informed version cannot be the same as the current repository")

type updateRepoZipCmd struct {
	repo formula.RepositoryListUpdater
	prompt.InputTextValidator
	prompt.InputURL
	prompt.InputList
	prompt.InputBool
}

func NewUpdateRepoZipCmd(
	repo formula.RepositoryListUpdater,
	inText prompt.InputTextValidator,
	inURL prompt.InputURL,
	inList prompt.InputList,
) *cobra.Command {
	updateRepoZip := updateRepoZipCmd{
		repo:               repo,
		InputTextValidator: inText,
		InputURL:           inURL,
		InputList:          inList,
	}

	cmd := &cobra.Command{
		Use:       "repo-zip",
		Short:     "Update a repository by URL.",
		Example:   "rit update repo-zip",
		RunE:      updateRepoZip.runCmd(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), updateRepoZipFlags)

	return cmd
}

func (up *updateRepoZipCmd) runCmd() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repoToUpdate, err := up.resolveInput(cmd)
		if err != nil {
			return err
		}

		if err := up.repo.Update(repoToUpdate.Name, repoToUpdate.Version, repoToUpdate.Url); err != nil {
			return err
		}

		s := fmt.Sprintf(successUpdate, repoToUpdate.Name, repoToUpdate.Version)
		prompt.Success(s)

		return nil
	}
}

func (up *updateRepoZipCmd) resolveInput(cmd *cobra.Command) (formula.Repo, error) {
	if IsFlagInput(cmd) {
		return up.resolveFlags(cmd)
	}
	return up.resolvePrompt()
}

func (up *updateRepoZipCmd) resolvePrompt() (formula.Repo, error) {
	repos, err := up.repo.List()
	if err != nil {
		return formula.Repo{}, err
	}

	var reposName []string
	for i := range repos {
		if repos[i].Provider == "ZipRemote" {
			reposName = append(reposName, repos[i].Name.String())
		}
	}

	helper := "Select a repository to update your version"
	name, err := up.List(questionSelectARepo, reposName, helper)
	if err != nil {
		return formula.Repo{}, err
	}

	var repoToUpdate formula.Repo
	for i := range repos {
		if repos[i].Name == formula.RepoName(name) {
			repoToUpdate = repos[i]
		}
	}

	questionTypeNewVersion := fmt.Sprintf("Type your new version for %q:", name)
	version, err := up.Text(questionTypeNewVersion, func(i interface{}) error {
		in := i.(string)
		if in == "" {
			return ErrVersionNotEmpty
		}

		if repoToUpdate.Version == formula.RepoVersion(in) {
			return ErrSameVersion
		}

		return nil
	})
	if err != nil {
		return formula.Repo{}, err
	}

	questionTypeNewURL := fmt.Sprintf("Type your new URL for %q:", name)
	repoURL, err := up.URL(questionTypeNewURL, defaultRepoURL)
	if err != nil {
		return formula.Repo{}, err
	}

	repoToUpdate.Version = formula.RepoVersion(version)
	repoToUpdate.Url = repoURL

	return repoToUpdate, nil
}

func (up *updateRepoZipCmd) resolveFlags(cmd *cobra.Command) (formula.Repo, error) {
	name, err := cmd.Flags().GetString(nameFlagName)
	if err != nil || name == "" {
		return formula.Repo{}, errors.New(missingFlagText(nameFlagName))
	}

	repoURL, err := cmd.Flags().GetString(repoUrlFlagName)
	if err != nil || repoURL == "" {
		return formula.Repo{}, errors.New(missingFlagText(repoUrlFlagName))
	}

	version, err := cmd.Flags().GetString(versionFlagName)
	if err != nil || version == "" {
		return formula.Repo{}, errors.New(missingFlagText(versionFlagName))
	}

	repos, err := up.repo.List()
	if err != nil {
		return formula.Repo{}, err
	}

	var repoToUpdate formula.Repo
	for i := range repos {
		if repos[i].Name == formula.RepoName(name) {
			repoToUpdate = repos[i]
		}
	}

	if repoToUpdate == (formula.Repo{}) {
		repoToUpdate.Name = formula.RepoName(name)
		return repoToUpdate, nil
	}

	if repoToUpdate.Version == formula.RepoVersion(version) {
		return formula.Repo{}, ErrSameVersion
	}

	repoToUpdate.Version = formula.RepoVersion(version)
	repoToUpdate.Url = repoURL

	return repoToUpdate, nil
}
