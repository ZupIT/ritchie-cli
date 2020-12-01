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

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const (
	defaultRepoURL  = "https://github.com/ZupIT/ritchie-formulas"
	messageExisting = "This formula repository already exists, check using \"rit list repo\""
)

var ErrRepoNameNotEmpty = errors.New("the field repository name must not be empty")

type addRepoCmd struct {
	repo          formula.RepositoryAddLister
	repoProviders formula.RepoProviders
	prompt.InputTextValidator
	prompt.InputPassword
	prompt.InputURL
	prompt.InputList
	prompt.InputBool
	prompt.InputInt
	tutorial rtutorial.Finder
	tree     tree.CheckerManager
	detail   formula.RepositoryDetail
}

func NewAddRepoCmd(
	repo formula.RepositoryAddLister,
	repoProviders formula.RepoProviders,
	inText prompt.InputTextValidator,
	inPass prompt.InputPassword,
	inURL prompt.InputURL,
	inList prompt.InputList,
	inBool prompt.InputBool,
	inInt prompt.InputInt,
	rtf rtutorial.Finder,
	treeChecker tree.CheckerManager,
	rd formula.RepositoryDetail,
) *cobra.Command {
	addRepo := addRepoCmd{
		repo:               repo,
		repoProviders:      repoProviders,
		InputTextValidator: inText,
		InputURL:           inURL,
		InputList:          inList,
		InputBool:          inBool,
		InputInt:           inInt,
		InputPassword:      inPass,
		tutorial:           rtf,
		tree:               treeChecker,
		detail:             rd,
	}
	cmd := &cobra.Command{
		Use:       "repo",
		Short:     "Add a repository",
		Example:   "rit add repo",
		RunE:      RunFuncE(addRepo.runStdin(), addRepo.runPrompt()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
	cmd.LocalFlags()

	return cmd
}

func (ad addRepoCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		provider, err := ad.List("Select your provider:", ad.repoProviders.List())
		if err != nil {
			return err
		}

		name, err := ad.Text("Repository name:", ad.repoNameValidator)
		if err != nil {
			return err
		}

		repos, err := ad.repo.List()
		if err != nil {
			return err
		}

		for i := range repos {
			repo := repos[i]
			if repo.Name == formula.RepoCommonsName && formula.RepoName(name) == formula.RepoCommonsName {
				prompt.Warning("You are trying to replace the \"commons\" repository!")
				choice, _ := ad.Bool("Do you want to proceed?", []string{"yes", "no"})
				if !choice {
					prompt.Info("Operation cancelled")
					return nil
				}
				break
			}

			if repo.Name == formula.RepoName(name) {
				prompt.Warning(fmt.Sprintf("Your repository %q is gonna be overwritten.", repo.Name))
				choice, _ := ad.Bool("Do you want to proceed?", []string{"yes", "no"})
				if !choice {
					prompt.Info("Operation cancelled")
					return nil
				}
			}
		}

		url, err := ad.URL("Repository URL:", defaultRepoURL)
		if err != nil {
			return err
		}

		isPrivate, err := ad.Bool("Is a private repository?", []string{"no", "yes"})
		if err != nil {
			return err
		}

		var token string
		if isPrivate {
			token, err = ad.Password("Personal access tokens:")
			if err != nil {
				return err
			}
		}

		git := ad.repoProviders.Resolve(formula.RepoProvider(provider))

		gitRepoInfo := git.NewRepoInfo(url, token)
		tags, err := git.Repos.Tags(gitRepoInfo)
		if err != nil {
			return err
		}

		if len(tags) <= 0 {
			return fmt.Errorf("please, generate a release to add your repository")
		}

		var tagNames []string
		for i := range tags {
			tagNames = append(tagNames, tags[i].Name)
		}

		version, err := ad.List("Select a tag version:", tagNames)
		if err != nil {
			return err
		}

		if exitsRepo(url, version, repos) {
			prompt.Info(messageExisting)
			return nil
		}

		priority, err := ad.Int("Set the priority:", "0 is higher priority, the lower higher the priority")
		if err != nil {
			return err
		}

		repository := formula.Repo{
			Provider: formula.RepoProvider(provider),
			Name:     formula.RepoName(name),
			Version:  formula.RepoVersion(version),
			Token:    token,
			Url:      url,
			Priority: int(priority),
		}

		if err := ad.repo.Add(repository); err != nil {
			return err
		}

		successMsg := fmt.Sprintf(
			"The %q repository was added with success, now you can use your formulas with the Ritchie!",
			repository.Name,
		)
		prompt.Success(successMsg)

		tutorialHolder, err := ad.tutorial.Find()
		if err != nil {
			return err
		}
		tutorialAddRepo(tutorialHolder.Current)
		ad.tree.Check()
		return nil
	}
}

func (ad addRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		r := formula.Repo{}

		err := stdin.ReadJson(cmd.InOrStdin(), &r)
		if err != nil {
			return err
		}

		if r.Version.String() == "" {
			latestTag := ad.detail.LatestTag(r)
			r.Version = formula.RepoVersion(latestTag)
		}

		repos, _ := ad.repo.List()
		if exitsRepo(r.Url, r.Version.String(), repos) {
			prompt.Info(messageExisting)
			return nil
		}

		if err := ad.repo.Add(r); err != nil {
			return err
		}

		successMsg := fmt.Sprintf(
			"The %q repository was added with success, now you can use your formulas with the Ritchie!",
			r.Name,
		)
		prompt.Success(successMsg)

		tutorialHolder, err := ad.tutorial.Find()
		if err != nil {
			return err
		}
		tutorialAddRepo(tutorialHolder.Current)
		return nil
	}
}

func (ad addRepoCmd) repoNameValidator(text interface{}) error {
	in := text.(string)
	if in == "" {
		return ErrRepoNameNotEmpty
	}

	return nil
}

func exitsRepo(urlToAdd, versionToAdd string, repos formula.Repos) bool {
	for i := range repos {
		if repos[i].Url == urlToAdd && repos[i].Version.String() == versionToAdd {
			return true
		}
	}
	return false
}

func tutorialAddRepo(tutorialStatus string) {
	const tagTutorial = "\n[TUTORIAL]"
	const MessageTitle = "To view your formula repositories:"
	const MessageBody = ` ∙ Run "rit list repo"` + "\n"

	if tutorialStatus == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
	}
}
