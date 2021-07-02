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
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const (
	defaultRepoURL   = "https://github.com/ZupIT/ritchie-formulas"
	messageExisting  = "This formula repository already exists, check using \"rit list repo\""
	repoUrlFlagName  = "repoUrl"
	priorityFlagName = "priority"
	tokenFlagName    = "token"
	tagFlagName      = "tag"
	permissionError  = `
	permission error:
	You must overwrite the current token (%s-add-repo) with command:
		rit set credential
	Or switch to a new environment with the command:
		rit set env`
)

var ErrRepoNameNotEmpty = errors.New("the field repository name must not be empty")

var addRepoFlags = flags{
	{
		name:        providerFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: "provider name (Github|Gitlab|Bitbucket)",
	},
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
		name:        tokenFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: "access token",
	},
	{
		name:        tagFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: "repository tag version",
	},
	{
		name:        priorityFlagName,
		kind:        reflect.Int,
		defValue:    1000,
		description: "repository priority (0 is highest)",
	},
}

type addRepoCmd struct {
	repo          formula.RepositoryAddLister
	repoProviders formula.RepoProviders
	cred          credential.Resolver
	prompt.InputTextValidator
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
	resolver credential.Resolver,
	inText prompt.InputTextValidator,
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
		tutorial:           rtf,
		tree:               treeChecker,
		detail:             rd,
		cred:               resolver,
	}
	cmd := &cobra.Command{
		Use:       "repo",
		Short:     "Add a repository",
		Example:   "rit add repo",
		RunE:      RunFuncE(addRepo.runStdin(), addRepo.runFormula()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), addRepoFlags)

	return cmd
}

func (ar *addRepoCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repo, err := ar.resolveInput(cmd)
		if err != nil || repo.Provider == "" {
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

func (ar *addRepoCmd) resolveInput(cmd *cobra.Command) (formula.Repo, error) {
	if IsFlagInput(cmd) {
		return ar.resolveFlags(cmd)
	}
	return ar.resolvePrompt()
}

func (ar *addRepoCmd) resolvePrompt() (formula.Repo, error) {
	provider, err := ar.List("Select your provider:", ar.repoProviders.List())
	if err != nil {
		return formula.Repo{}, err
	}

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
			break
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

	url, err := ar.URL("Repository URL:", defaultRepoURL)
	if err != nil {
		return formula.Repo{}, err
	}

	isPrivate, err := ar.Bool("Is a private repository?", []string{"no", "yes"})
	if err != nil {
		return formula.Repo{}, err
	}

	var token string
	if isPrivate {
		token, err = ar.cred.Resolve("CREDENTIAL_" + provider + "-ADD-REPO_TOKEN")
		if err != nil {
			return formula.Repo{}, err
		}
	}

	git := ar.repoProviders.Resolve(formula.RepoProvider(provider))

	gitRepoInfo := git.NewRepoInfo(url, token)
	tags, err := git.Repos.Tags(gitRepoInfo)
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			errorString := fmt.Sprintf(permissionError, provider)
			return formula.Repo{}, errors.New(errorString)
		}
		return formula.Repo{}, err
	}

	if len(tags) <= 0 {
		return formula.Repo{}, fmt.Errorf("please, generate a release to add your repository")
	}

	tagNames := make([]string, 0, len(tags))
	for i := range tags {
		tagNames = append(tagNames, tags[i].Name)
	}

	version, err := ar.List("Select a tag version:", tagNames)
	if err != nil {
		return formula.Repo{}, err
	}

	if existsRepo(url, version, repos) {
		prompt.Info(messageExisting)
		return formula.Repo{}, nil
	}

	priority, err := ar.Int("Set the priority:", "0 is highest")
	if err != nil {
		return formula.Repo{}, err
	}

	return formula.Repo{
		Provider: formula.RepoProvider(provider),
		Name:     formula.RepoName(name),
		Version:  formula.RepoVersion(version),
		Token:    token,
		Url:      url,
		Priority: int(priority),
	}, nil
}

func (ar *addRepoCmd) resolveFlags(cmd *cobra.Command) (formula.Repo, error) {
	provider, err := cmd.Flags().GetString(providerFlagName)
	if err != nil || provider == "" {
		return formula.Repo{}, errors.New(missingFlagText(providerFlagName))
	}

	providers := ar.repoProviders.List()
	providerValid := false
	for _, repoProvider := range providers {
		if repoProvider == provider {
			providerValid = true
			break
		}
	}
	if !providerValid {
		return formula.Repo{}, errors.New("please select a provider from " + strings.Join(providers, ", "))
	}

	name, err := cmd.Flags().GetString(nameFlagName)
	if err != nil || name == "" {
		return formula.Repo{}, errors.New(missingFlagText(nameFlagName))
	}

	repoUrl, err := cmd.Flags().GetString(repoUrlFlagName)
	if err != nil || repoUrl == "" {
		return formula.Repo{}, errors.New(missingFlagText(repoUrlFlagName))
	}

	tag, err := cmd.Flags().GetString(tagFlagName)
	if err != nil {
		return formula.Repo{}, errors.New(missingFlagText(tagFlagName))
	}
	token, err := cmd.Flags().GetString(tokenFlagName)
	if err != nil {
		return formula.Repo{}, err
	}

	priority, err := cmd.Flags().GetInt(priorityFlagName)
	if err != nil {
		return formula.Repo{}, err
	}

	repo := formula.Repo{
		Provider: formula.RepoProvider(provider),
		Name:     formula.RepoName(name),
		Version:  formula.RepoVersion(tag),
		Token:    token,
		Url:      repoUrl,
		Priority: priority,
	}

	if repo.EmptyVersion() {
		latestTag := ar.detail.LatestTag(repo)
		repo.Version = formula.RepoVersion(latestTag)
	}

	return repo, nil
}

func printConflictingCommandsWarning(conflictingCommands []api.CommandID) {
	if len(conflictingCommands) <= 0 {
		return
	}

	lastCommandIndex := len(conflictingCommands) - 1
	lastCommand := conflictingCommands[lastCommandIndex].String()
	lastCommand = strings.Replace(lastCommand, "root", "rit", 1)
	lastCommand = strings.ReplaceAll(lastCommand, "_", " ")
	msg := fmt.Sprintf("There's a total of %d formula conflicting commands, like:\n %s", len(conflictingCommands), lastCommand)
	msg = prompt.Yellow(msg)
	fmt.Println(msg)
}

func (ar addRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		r := formula.Repo{}

		err := stdin.ReadJson(cmd.InOrStdin(), &r)
		if err != nil {
			return err
		}

		if r.EmptyVersion() {
			latestTag := ar.detail.LatestTag(r)
			r.Version = formula.RepoVersion(latestTag)
		}

		repos, _ := ar.repo.List()
		if existsRepo(r.Url, r.Version.String(), repos) {
			prompt.Info(messageExisting)
			return nil
		}

		if err := ar.repo.Add(r); err != nil {
			return err
		}

		successMsg := fmt.Sprintf(
			"The %q repository was added with success, now you can use your formulas with the Ritchie!",
			r.Name,
		)
		prompt.Success(successMsg)

		tutorialHolder, err := ar.tutorial.Find()
		if err != nil {
			return err
		}
		tutorialAddRepo(tutorialHolder.Current)
		return nil
	}
}

func (ar addRepoCmd) repoNameValidator(text interface{}) error {
	in := text.(string)
	if in == "" {
		return ErrRepoNameNotEmpty
	}

	return nil
}

func existsRepo(urlToAdd, versionToAdd string, repos formula.Repos) bool {
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
