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
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

const (
	addRepoMsg = "Run \"rit add repo\" to add a new repository manually."
)

var (
	addRepoInfo = `You can keep the configuration without adding the community repository,
but you will need to provide a git repo with the formulas templates and add them with 
rit add repo command, naming this repository obligatorily as "commons".

See how to do this on the example: [https://github.com/ZupIT/ritchie-formulas/blob/master/templates/create_formula/README.md]`
	errMsg             = prompt.Yellow("It was not possible to add the commons repository at this time, please try again later.")
	ErrInitCommonsRepo = errors.New(errMsg)
	CommonsRepoURL     = "https://github.com/ZupIT/ritchie-formulas"
)

type initStdin struct {
	AddRepo string `json:"addRepo"`
}

type initCmd struct {
	repo formula.RepositoryAdder
	git  git.Repositories
	rt   rtutorial.Finder
	prompt.InputBool
}

func NewInitCmd(repo formula.RepositoryAdder, git git.Repositories, rtf rtutorial.Finder, inBool prompt.InputBool) *cobra.Command {
	o := initCmd{repo: repo, git: git, rt: rtf, InputBool: inBool}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize rit configuration",
		Long:  "Initialize rit configuration",
		RunE:  RunFuncE(o.runStdin(), o.runPrompt()),
	}

	return cmd
}

func (in initCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		label := "Would you like to add the community repository? [https://github.com/ZupIT/ritchie-formulas]"
		choose, err := in.Bool(label, []string{"yes", "no"})
		if err != nil {
			return err
		}

		if !choose {
			fmt.Println()
			prompt.Warning(addRepoInfo)
			fmt.Println()
			fmt.Println(addRepoMsg)
			return nil
		}

		repo := formula.Repo{
			Provider: "Github",
			Name:     "commons",
			Url:      CommonsRepoURL,
			Priority: 0,
		}

		s := spinner.StartNew("Adding the commons repository...")
		time.Sleep(time.Second * 2)

		repoInfo := github.NewRepoInfo(repo.Url, repo.Token)

		tag, err := in.git.LatestTag(repoInfo)
		if err != nil {
			s.Error(ErrInitCommonsRepo)
			fmt.Println(addRepoMsg)
			return nil
		}

		repo.Version = formula.RepoVersion(tag.Name)

		if err := in.repo.Add(repo); err != nil {
			s.Error(ErrInitCommonsRepo)
			fmt.Println(addRepoMsg)
			return nil
		}

		s.Success(prompt.Green("Initialization successful!"))

		tutorialHolder, err := in.rt.Find()
		if err != nil {
			return err
		}
		tutorialInit(tutorialHolder.Current)
		return nil
	}
}

func (in initCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		init := initStdin{}

		err := stdin.ReadJson(cmd.InOrStdin(), &init)
		if err != nil {
			return err
		}

		if init.AddRepo == "no" {
			fmt.Printf("\n%s\n", prompt.Yellow(addRepoInfo))

			fmt.Println(addRepoMsg)
		} else {
			s := spinner.StartNew("Adding the commons repository...")
			time.Sleep(time.Second * 2)

			repo := formula.Repo{
				Provider: "Github",
				Name:     "commons",
				Url:      CommonsRepoURL,
				Priority: 0,
			}

			repoInfo := github.NewRepoInfo(repo.Url, repo.Token)

			tag, err := in.git.LatestTag(repoInfo)
			if err != nil {
				s.Error(ErrInitCommonsRepo)
				fmt.Println(addRepoMsg)
				return nil
			}

			repo.Version = formula.RepoVersion(tag.Name)

			if err := in.repo.Add(repo); err != nil {
				s.Error(ErrInitCommonsRepo)
				fmt.Println(addRepoMsg)
				return nil
			}

			s.Success(prompt.Green("Initialization successful!"))
		}

		tutorialHolder, err := in.rt.Find()
		if err != nil {
			return err
		}
		tutorialInit(tutorialHolder.Current)
		return nil
	}
}

func tutorialInit(tutorialStatus string) {
	const tagTutorial = "\n[TUTORIAL]"
	const MessageTitle = "How to create new formulas:"
	const MessageBody = ` ∙ Run "rit create formula"
 ∙ Open the project with your favorite text editor.` + "\n"

	if tutorialStatus == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
	}
}
