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
	"strings"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

const (
	addRepoMsg         = "Run \"rit add repo\" to add a new repository manually."
	AddCommonsQuestion = "Would you like to add the community repository? [https://github.com/ZupIT/ritchie-formulas]"
	AddMetricsQuestion = `To help us improve and deliver more value to the community, 
do you agree to let us collect anonymous data about product 
and feature use statistics and crash reports?`
	AcceptMetrics             = "Yes, I agree to contribute with data anonymously"
	DoNotAcceptMetrics        = "No, not for now."
	SelectFormulaTypeQuestion = "Select a default formula run type:"
)

var (
	addRepoInfo = `You can keep the configuration without adding the community repository,
 but you will need to provide a git repo with the formulas templates and add them with 
 rit add repo command, naming this repository obligatorily as "commons".
 
 See how to do this on the example: [https://github.com/ZupIT/ritchie-formulas/blob/master/templates/create_formula/README.md]`
	errMsg             = prompt.Yellow("It was not possible to add the commons repository at this time, please try again later.")
	ErrInitCommonsRepo = errors.New(errMsg)
	CommonsRepoURL     = "https://github.com/ZupIT/ritchie-formulas"
	ErrInvalidRunType  = fmt.Errorf("invalid formula run type, these run types are enabled [%v]", strings.Join(formula.RunnerTypes, ", "))
)

type initStdin struct {
	AddCommons  bool   `json:"addCommons"`
	SendMetrics bool   `json:"sendMetrics"`
	RunType     string `json:"runType"`
}

type initCmd struct {
	repo     formula.RepositoryAdder
	git      git.Repositories
	tutorial rtutorial.Finder
	config   formula.ConfigRunner
	file     stream.FileWriteReadExister
	prompt.InputList
	prompt.InputBool
}

func NewInitCmd(
	repo formula.RepositoryAdder,
	git git.Repositories,
	tutorial rtutorial.Finder,
	config formula.ConfigRunner,
	file stream.FileWriteReadExister,
	inList prompt.InputList,
	inBool prompt.InputBool,
) *cobra.Command {
	o := initCmd{repo: repo, git: git, tutorial: tutorial, config: config, file: file, InputList: inList, InputBool: inBool}

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
		if err := in.metricsAuthorization(); err != nil {
			return err
		}

		if err := in.addCommonsRepo(); err != nil {
			return err
		}

		if err := in.setRunnerType(); err != nil {
			return err
		}

		prompt.Success("\nInitialization successful!\n")

		if err := in.tutorialInit(); err != nil {
			return err
		}

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

		sendMetrics := "no"
		if init.SendMetrics {
			sendMetrics = "yes"
		}

		if err = in.file.Write(metric.FilePath, []byte(sendMetrics)); err != nil {
			return err
		}

		if !init.AddCommons {
			fmt.Println()
			prompt.Warning(addRepoInfo)
			fmt.Println()
			fmt.Println(addRepoMsg)
		} else {
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
			}

			repo.Version = formula.RepoVersion(tag.Name)

			if err := in.repo.Add(repo); err != nil {
				s.Error(ErrInitCommonsRepo)
				fmt.Println(addRepoMsg)
				return nil
			}

			s.Success(prompt.Green("Commons repository added successfully!\n"))
		}

		runType := formula.RunnerType(-1)
		for i := range formula.RunnerTypes {
			if formula.RunnerTypes[i] == init.RunType {
				runType = formula.RunnerType(i)
				break
			}
		}

		if runType == -1 {
			return ErrInvalidRunType
		}

		if err := in.config.Create(runType); err != nil {
			return err
		}

		if runType == formula.Local {
			prompt.Warning(`
In order to run formulas locally, you must have the formula language installed on your machine,
if you don't want to install choose to run the formulas inside the docker.
`)
		}

		prompt.Success("Initialization successful!")

		if err := in.tutorialInit(); err != nil {
			return err
		}

		return nil
	}
}

func (in initCmd) metricsAuthorization() error {
	const welcome = "Welcome to Ritchie!\n"
	const header = `Ritchie is a platform that helps you and your team to save time by 
giving you the power to create powerful templates to execute important 
tasks across your team and organization with minimum time and with standards, 
delivering autonomy to developers with security.

You can view our Privacy Policy (http://insights.zup.com.br/politica-privacidade) to better understand our commitment.
`
	const footer = "\nYou can always modify your choice using the \"rit metrics\" command.\n"
	options := []string{AcceptMetrics, DoNotAcceptMetrics}

	prompt.Info(welcome)
	fmt.Println(header)

	choose, err := in.InputList.List(AddMetricsQuestion, options)
	if err != nil {
		return err
	}
	fmt.Println(footer)

	responseToWrite := "yes"
	if choose == DoNotAcceptMetrics {
		responseToWrite = "no"
	}

	if err = in.file.Write(metric.FilePath, []byte(responseToWrite)); err != nil {
		return err
	}

	return nil
}

func (in initCmd) setRunnerType() error {
	selected, err := in.List(SelectFormulaTypeQuestion, formula.RunnerTypes)
	if err != nil {
		return err
	}

	runType := formula.RunnerType(-1)
	for i := range formula.RunnerTypes {
		if formula.RunnerTypes[i] == selected {
			runType = formula.RunnerType(i)
			break
		}
	}

	if runType == -1 {
		return ErrInvalidRunType
	}

	if err := in.config.Create(runType); err != nil {
		return err
	}

	if runType == formula.Local {
		prompt.Warning(`
In order to run formulas locally, you must have the formula language installed on your machine,
if you don't want to install choose to run the formulas inside the docker.
`)
	}

	return nil
}

func (in initCmd) addCommonsRepo() error {
	choose, err := in.Bool(AddCommonsQuestion, []string{"yes", "no"})
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

	s.Success(prompt.Green("Commons repository added successfully!\n"))

	return nil
}

func (in initCmd) tutorialInit() error {
	tutorialHolder, err := in.tutorial.Find()
	if err != nil {
		return err
	}

	const tagTutorial = "\n[TUTORIAL]"
	const MessageTitle = "How to create new formulas:"
	const MessageBody = ` ∙ Run "rit create formula"
  ∙ Open the project with your favorite text editor.` + "\n"

	if tutorialHolder.Current == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
	}

	return nil
}