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

	"github.com/ZupIT/ritchie-cli/internal/pkg/config"
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
	addRepoMsg                = "Run \"rit add repo\" to add a new repository manually.\n"
	AddMetricsQuestion        = "To help us improve and deliver more value to the community,\nwe will collect anonymous data about product and feature\nuse statistics and crash reports."
	AcceptOpt                 = "✅ Yes"
	DeclineOpt                = "❌ No"
	LocalRunType              = "🏠 local"
	DockerRunType             = "🐳 docker"
	SelectFormulaTypeQuestion = "Select a default formula run type:"
	FormulaLocalRunWarning    = `
In order to run formulas locally, you must have the formula language installed on your machine,
if you don't want to install choose to run the formulas inside the docker.
`
	addRepoInfo = `
You can keep the configuration without adding the community repository,
but you will need to provide a git repo with the formulas templates and add them with
"rit add repo" command, naming this repository obligatorily as "commons".

See how to do this on the example: 
[https://github.com/ZupIT/ritchie-formulas/blob/master/templates/create_formula/README.md]

`
	CommonsRepoURL = "https://github.com/ZupIT/ritchie-formulas"
)

var (
	errMsg             = prompt.Yellow("It was not possible to add the commons repository at this time, please try again later.")
	ErrInitCommonsRepo = errors.New(errMsg)
	ErrInvalidRunType  = fmt.Errorf("invalid formula run type, these run types are enabled [%v]",
		strings.Join(formula.RunnerTypes, ", "))
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
	file     stream.FileWriter
	prompt.InputList
	prompt.InputBool
	prompt.InputMultiselect
	metricSender metric.SendManagerHttp
	ritConfig    config.Writer
}

func NewInitCmd(
	repo formula.RepositoryAdder,
	git git.Repositories,
	tutorial rtutorial.Finder,
	config formula.ConfigRunner,
	file stream.FileWriter,
	inList prompt.InputList,
	inBool prompt.InputBool,
	metricSender metric.SendManagerHttp,
	ritConfig config.Writer,
) *cobra.Command {
	o := initCmd{
		repo:         repo,
		git:          git,
		tutorial:     tutorial,
		config:       config,
		file:         file,
		InputList:    inList,
		InputBool:    inBool,
		metricSender: metricSender,
		ritConfig:    ritConfig,
	}

	cmd := &cobra.Command{
		Use:       "init",
		Short:     "Initialize rit configuration",
		Long:      "Initialize rit configuration",
		RunE:      RunFuncE(o.runStdin(), o.runPrompt()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	return cmd
}

func (in initCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		in.Welcome()

		metrics, err := in.metricsAuthorization()
		if err != nil {
			return err
		}

		if err := in.addCommonsRepo(); err != nil {
			return err
		}

		runType, err := in.setRunnerType()
		if err != nil {
			return err
		}

		prompt.Success("\n✅  Initialization successful!\n")

		configs := config.Configs{
			Language: "English",
			Metrics:  metrics,
			RunType:  runType,
			Tutorial: tutorialStatusEnabled,
		}

		if err := in.ritConfig.Write(configs); err != nil {
			return err
		}

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
			prompt.Warning(addRepoInfo)
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

			s.Success(prompt.Green("✅ Commons repository added successfully!\n"))
		}

		runType := formula.DefaultRun
		for i := range formula.RunnerTypes {
			if formula.RunnerTypes[i] == init.RunType {
				runType = formula.RunnerType(i)
				break
			}
		}

		if runType == formula.DefaultRun {
			return ErrInvalidRunType
		}

		if err := in.config.Create(runType); err != nil {
			return err
		}

		if runType == formula.LocalRun {
			prompt.Warning("\n\t\t\t⚠️  WARNING ⚠️")
			fmt.Print(FormulaLocalRunWarning)
		}

		prompt.Success("\n✅  Initialization successful!\n")

		if err := in.tutorialInit(); err != nil {
			return err
		}

		return nil
	}
}

func (in initCmd) metricsAuthorization() (string, error) {
	prompt.Info("📊 Metrics 📊\n")
	options := []string{AcceptOpt, DeclineOpt}

	fmt.Println(AddMetricsQuestion)

	choose, err := in.Bool("Do you agree?", options)
	if err != nil {
		return "", err
	}

	const footer = "\nYou can always modify your choice using the \"rit metrics\" command.\n"
	fmt.Println(footer)

	responseToWrite := "yes"
	if !choose {
		responseToWrite = "no"
		in.metricSender.Send(metric.APIData{
			Id:        "rit_init",
			Timestamp: time.Now(),
			Data: metric.Data{
				MetricsAcceptance: responseToWrite,
			},
		})
	}

	if err = in.file.Write(metric.FilePath, []byte(responseToWrite)); err != nil {
		return "", err
	}

	return responseToWrite, nil
}

func (in initCmd) setRunnerType() (formula.RunnerType, error) {
	prompt.Info("🏃 FORMULA RUN TYPE 🏃\n")
	runTypes := []string{LocalRunType, DockerRunType}
	selected, err := in.List(SelectFormulaTypeQuestion, runTypes)
	if err != nil {
		return formula.DefaultRun, err
	}

	runType := formula.DefaultRun
	for i := range runTypes {
		if runTypes[i] == selected {
			runType = formula.RunnerType(i)
			break
		}
	}

	if runType == formula.DefaultRun {
		return runType, ErrInvalidRunType
	}

	if err := in.config.Create(runType); err != nil {
		return runType, err
	}

	if runType == formula.LocalRun {
		prompt.Warning("\n\t\t\t⚠️  WARNING ⚠️")
		fmt.Print(FormulaLocalRunWarning)
	}

	return runType, nil
}

func (in initCmd) addCommonsRepo() error {
	prompt.Info("⭐ Commons repository ⭐\n")
	choose, err := in.Bool("Would you like to add the community repository?", []string{AcceptOpt, DeclineOpt})
	if err != nil {
		return err
	}
	metric.CommonsRepoAdded = "yes"
	if !choose {
		in.CommonsWarning()
		metric.CommonsRepoAdded = "no"
		return nil
	}

	repo := formula.Repo{
		Provider: "Github",
		Name:     "commons",
		Url:      CommonsRepoURL,
		Priority: 0,
	}

	s := spinner.StartNew("Adding the commons repository...")

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

	s.Success(prompt.Green("✅ Commons repository added successfully!\n"))

	return nil
}

func (in initCmd) tutorialInit() error {
	tutorialHolder, err := in.tutorial.Find()
	if err != nil {
		return err
	}

	const tagTutorial = "\n📖 TUTORIAL 📖"
	const MessageTitle = "How to create new formulas:"
	const MessageBody = `
 ∙ Run "rit create formula"
 ∙ Open the project with your favorite text editor.
`
	const MessageCommons = "Take a look at the formulas you can run and" +
		" test to see what you can with Ritchie using \"rit\"\n"

	if tutorialHolder.Current == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
		prompt.Info(MessageCommons)
	}

	return nil
}

func (in initCmd) Welcome() {
	const welcome = " _______       _     _             __         _                           ______    _____      _____  \n|_   __ \\     (_)   / |_          [  |       (_)                        .' ___  |  |_   _|    |_   _| \n  | |__) |    __   `| |-'  .---.   | |--.    __    .---.     ______    / .'   \\_|    | |        | |   \n  |  __ /    [  |   | |   / /'`\\]  | .-. |  [  |  / /__\\\\   |______|   | |           | |   _    | |   \n _| |  \\ \\_   | |   | |,  | \\__.   | | | |   | |  | \\__.,              \\ `.___.'\\   _| |__/ |  _| |_  \n|____| |___| [___]  \\__/  '.___.' [___]|__] [___]  '.__.'               `.____ .'  |________| |_____| \n\n"
	const header = `Ritchie is a platform that helps you and your team to save time by
giving you the power to create powerful templates to execute important
tasks across your team and organization with minimum time and with standards,
delivering autonomy to developers with security.

You can view our Privacy Policy (http://insights.zup.com.br/politica-privacidade) to better understand our commitment.

`
	prompt.Info(welcome)
	fmt.Print(header)
}

func (in initCmd) CommonsWarning() {
	prompt.Warning("\n\t\t\t⚠️  WARNING ⚠️")
	fmt.Print(addRepoInfo)
	fmt.Println(addRepoMsg)
}
