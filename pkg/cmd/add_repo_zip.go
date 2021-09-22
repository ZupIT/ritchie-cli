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
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

const versionFlagName = "version"

var addRepoZipFlags = flags{
	{
		name:        nameFlagName,
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
		name:        versionFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: "repository version",
	},
}

var ErrVersionNotEmpty = errors.New("the version field must not be empty")

type addRepoZipCmd struct {
	repo formula.RepositoryAddLister
	prompt.InputTextValidator
	prompt.InputURL
	prompt.InputList
	prompt.InputBool
	tutorial rtutorial.Finder
	tree     tree.CheckerManager
}

func NewAddRepoZipCmd(
	repo formula.RepositoryAddLister,
	inText prompt.InputTextValidator,
	inURL prompt.InputURL,
	inBool prompt.InputBool,
	rtf rtutorial.Finder,
	treeChecker tree.CheckerManager,
) *cobra.Command {
	addRepoZip := addRepoZipCmd{
		repo:               repo,
		InputTextValidator: inText,
		InputURL:           inURL,
		InputBool:          inBool,
		tutorial:           rtf,
		tree:               treeChecker,
	}
	cmd := &cobra.Command{
		Use:       "repo-zip",
		Short:     "Add a repository by URL",
		Example:   "rit add repo-zip",
		RunE:      addRepoZip.runFormula(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
		Hidden:    true,
	}

	addReservedFlags(cmd.Flags(), addRepoZipFlags)

	return cmd
}

func (ar *addRepoZipCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repo, err := ar.resolveInput(cmd)
		if err != nil {
			return err
		}

		if err := ar.repo.Add(repo); err != nil {
			return err
		}

		successMsg := fmt.Sprintf(
			"The %q repository was added with success, now you can use your formulas with the Ritchie!",
			repo.Name,
		)
		prompt.Success(successMsg)

		tutorialHolder, err := ar.tutorial.Find()
		if err != nil {
			return err
		}
		tutorialAddRepo(tutorialHolder.Current)
		conflictCmds := ar.tree.Check()

		printConflictingCommandsWarning(conflictCmds)

		return nil
	}
}

func (ar *addRepoZipCmd) resolveInput(cmd *cobra.Command) (formula.Repo, error) {
	if IsFlagInput(cmd) {
		return ar.resolveFlags(cmd)
	}
	return ar.resolvePrompt()
}

func (ar *addRepoZipCmd) resolvePrompt() (formula.Repo, error) {
	name, err := ar.Text("Repository name:", ar.repoNameValidator)
	if err != nil {
		return formula.Repo{}, err
	}

	repos, err := ar.repo.List()
	if err != nil {
		return formula.Repo{}, err
	}

	for i := range repos {
		repo := repos[i]
		if repo.Name == formula.RepoCommonsName && formula.RepoName(name) == formula.RepoCommonsName {
			prompt.Warning("You are trying to replace the \"commons\" repository!")
			choice, _ := ar.Bool("Do you want to proceed?", []string{"yes", "no"})
			if !choice {
				prompt.Info("Operation cancelled")
				return formula.Repo{}, nil
			}
		}

		if repo.Name == formula.RepoName(name) {
			prompt.Warning(fmt.Sprintf("Your repository %q is gonna be overwritten.", repo.Name))
			choice, _ := ar.Bool("Do you want to proceed?", []string{"yes", "no"})
			if !choice {
				prompt.Info("Operation cancelled")
				return formula.Repo{}, nil
			}
		}
	}

	url, err := ar.URL("Repository URL zip:", defaultRepoURL)
	if err != nil {
		return formula.Repo{}, err
	}

	version, err := ar.Text("Version:", ar.repoVersionValidator)
	if err != nil {
		return formula.Repo{}, err
	}

	return formula.Repo{
		Provider:      formula.RepoProvider("ZipRemote"),
		Name:          formula.RepoName(name),
		Version:       formula.RepoVersion(version),
		Url:           url,
		IsLocal:       true,
		LatestVersion: formula.RepoVersion(version),
	}, nil
}

func (ar *addRepoZipCmd) resolveFlags(cmd *cobra.Command) (formula.Repo, error) {
	name, err := cmd.Flags().GetString(nameFlagName)
	if err != nil || name == "" {
		return formula.Repo{}, errors.New(missingFlagText(nameFlagName))
	}

	repoUrl, err := cmd.Flags().GetString(repoUrlFlagName)
	if err != nil || repoUrl == "" {
		return formula.Repo{}, errors.New(missingFlagText(repoUrlFlagName))
	}

	version, err := cmd.Flags().GetString(versionFlagName)
	if err != nil || version == "" {
		return formula.Repo{}, errors.New(missingFlagText(versionFlagName))
	}

	return formula.Repo{
		Provider:      formula.RepoProvider("ZipRemote"),
		Name:          formula.RepoName(name),
		Version:       formula.RepoVersion(version),
		Url:           repoUrl,
		IsLocal:       true,
		LatestVersion: formula.RepoVersion(version),
	}, nil
}

func (ar addRepoZipCmd) repoNameValidator(text interface{}) error {
	in := text.(string)
	if in == "" {
		return ErrRepoNameNotEmpty
	}

	return nil
}

func (ar addRepoZipCmd) repoVersionValidator(text interface{}) error {
	in := text.(string)
	if in == "" {
		return ErrVersionNotEmpty
	}

	return nil
}
