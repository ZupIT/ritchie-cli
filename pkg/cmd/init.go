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
	"time"

	"github.com/ZupIT/ritchie-cli/internal/pkg/config"
	"github.com/ZupIT/ritchie-cli/internal/pkg/i18n"
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

const TemplatesRepoURL = "https://github.com/ZupIT/ritchie-templates"

var (
	addRepoInfo               = i18n.T("init.add.commons.repo.info")
	FormulaLocalRunWarning    = i18n.T("init.run.type.local.warning")
	AddTheCommunityRepo       = i18n.T("init.add.commons.repo.question")
	SelectFormulaTypeQuestion = i18n.T("init.run.type.question")
	AgreeSendMetrics          = i18n.T("init.add.metric.question")
	addRepoMsg                = i18n.T("init.add.commons.repo.help")
	AcceptOpt                 = i18n.T("input.accept.opt")
	DeclineOpt                = i18n.T("input.decline.opt")
	LocalRunType              = i18n.T("input.run.type.local")
	DockerRunType             = i18n.T("input.run.type.docker")
	AcceptDeclineOpts         = []string{AcceptOpt, DeclineOpt}
	RunTypes                  = []string{LocalRunType, DockerRunType}
	errMsg                    = prompt.Yellow(i18n.T("init.add.commons.repo.error"))
	ErrInitCommonsRepo        = errors.New(errMsg)
	ErrInvalidRunType         = errors.New(i18n.T("init.invalid.run.type.error", strings.Join(formula.RunnerTypes, ", ")))
	metricsFlag               = "sendMetrics"
	commonsFlag               = "addCommons"
	runnerFlag                = "runType"
	provideValidValue         = "provide a valid value to the flag %q"
)

var initFlags = flags{
	{
		name:        metricsFlag,
		kind:        reflect.String,
		defValue:    "",
		description: "Do you accept to submit anonymous metrics? (ie: yes, no)",
	},
	{
		name:        commonsFlag,
		kind:        reflect.String,
		defValue:    "",
		description: "Do you want to download the commons repository? (ie: yes, no)",
	},
	{
		name:        runnerFlag,
		kind:        reflect.String,
		defValue:    "local",
		description: "Which default runner do you want to use? (ie: local, docker)",
	},
}

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
		Short:     i18n.T("init.cmd.description"),
		Long:      i18n.T("init.cmd.description"),
		RunE:      RunFuncE(o.runStdin(), o.runCmd()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
	cmd.LocalFlags()
	addReservedFlags(cmd.Flags(), initFlags)
	return cmd
}

func (in initCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		init := initStdin{}

		if err := stdin.ReadJson(cmd.InOrStdin(), &init); err != nil {
			return err
		}

		in.welcome()

		sendMetrics := "no"
		if init.SendMetrics {
			sendMetrics = "yes"
		}

		if err := in.file.Write(metric.FilePath, []byte(sendMetrics)); err != nil {
			return err
		}

		if !init.AddCommons {
			in.commonsWarning()
		} else {
			repo := formula.Repo{
				Provider: "Github",
				Name:     "commons",
				Url:      TemplatesRepoURL,
				Priority: 0,
			}

			s := spinner.StartNew(i18n.T("init.adding.commons.repo"))

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

			in.commonsSuccess(s)
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
			in.warning()
			fmt.Print(FormulaLocalRunWarning)
		}

		configs := config.Configs{
			Language: "English",
			Metrics:  sendMetrics,
			RunType:  runType,
			Tutorial: tutorialStatusEnabled,
		}

		if err := in.ritConfig.Write(configs); err != nil {
			return err
		}

		in.initSuccess()

		if err := in.tutorialInit(); err != nil {
			return err
		}

		return nil
	}
}

func (in initCmd) runCmd() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		in.welcome()

		configs, err := in.resolveInput(cmd)
		if err != nil {
			return err
		}

		if err := in.ritConfig.Write(configs); err != nil {
			return err
		}

		in.initSuccess()

		if err := in.tutorialInit(); err != nil {
			return err
		}
		return nil
	}
}

func (in initCmd) resolveInput(cmd *cobra.Command) (config.Configs, error) {
	if IsFlagInput(cmd) {
		return in.runFlags(cmd)
	}
	return in.runPrompt()
}

func (in *initCmd) runFlags(cmd *cobra.Command) (config.Configs, error) {
	metrics, err := cmd.Flags().GetString(metricsFlag)
	if err != nil {
		return config.Configs{}, err
	} else if metrics == "" {
		return config.Configs{}, errors.New(missingFlagText(metricsFlag))
	}
	commons, err := cmd.Flags().GetString(commonsFlag)
	if err != nil {
		return config.Configs{}, err
	} else if commons == "" {
		return config.Configs{}, errors.New(missingFlagText(commonsFlag))
	}
	runner, err := cmd.Flags().GetString(runnerFlag)
	if err != nil {
		return config.Configs{}, err
	} else if runner == "" {
		return config.Configs{}, errors.New(missingFlagText(runnerFlag))
	}

	metricBool, err := in.flagToBool(metrics, metricsFlag)
	if err != nil {
		return config.Configs{}, err
	}
	commonsBool, err := in.flagToBool(commons, commonsFlag)
	if err != nil {
		return config.Configs{}, err
	}

	switch metricBool {
	case false:
		{
			in.metricSender.Send(metric.APIData{
				Id:        "rit_init",
				UserId:    "",
				Timestamp: time.Now(),
				Data: metric.Data{
					MetricsAcceptance: metrics,
				},
			})
		}
	case true:
		{
			if err = in.file.Write(metric.FilePath, []byte(metrics)); err != nil {
				return config.Configs{}, err
			}
		}
	}

	switch commonsBool {
	case false:
		{
			in.commonsWarning()
			metric.CommonsRepoAdded = "no"
		}
	case true:
		{
			repo := formula.Repo{
				Provider: "Github",
				Name:     "commons",
				Url:      TemplatesRepoURL,
				Priority: 0,
			}
			s := spinner.StartNew(i18n.T("init.adding.commons.repo"))
			repoInfo := github.NewRepoInfo(repo.Url, repo.Token)
			tag, err := in.git.LatestTag(repoInfo)
			if err != nil {
				s.Error(ErrInitCommonsRepo)
				fmt.Println(addRepoMsg)
				return config.Configs{}, err
			}
			repo.Version = formula.RepoVersion(tag.Name)
			if err := in.repo.Add(repo); err != nil {
				s.Error(ErrInitCommonsRepo)
				fmt.Println(addRepoMsg)
				return config.Configs{}, err
			}
			in.commonsSuccess(s)
		}
	}

	var runType formula.RunnerType
	switch runner {
	case "local":
		runType = formula.LocalRun
	case "docker":
		runType = formula.DockerRun
	default:
		return config.Configs{}, fmt.Errorf(provideValidValue, runnerFlag)
	}

	configs := config.Configs{
		Language: config.DefaultLang,
		Metrics:  metrics,
		RunType:  runType,
		Tutorial: tutorialStatusEnabled,
	}

	return configs, nil
}

func (in initCmd) runPrompt() (config.Configs, error) {

	metrics, err := in.metricsAuthorization()
	if err != nil {
		return config.Configs{}, err
	}

	if err := in.addCommonsRepo(); err != nil {
		return config.Configs{}, err
	}

	runType, err := in.setRunnerType()
	if err != nil {
		return config.Configs{}, err
	}

	configs := config.Configs{
		Language: config.DefaultLang,
		Metrics:  metrics,
		RunType:  runType,
		Tutorial: tutorialStatusEnabled,
	}

	return configs, nil
}

func (in initCmd) metricsAuthorization() (string, error) {
	header := i18n.T("init.metric.header")

	prompt.Info(header)
	fmt.Println(i18n.T("init.add.metric.info"))

	choose, err := in.Bool(AgreeSendMetrics, AcceptDeclineOpts)
	if err != nil {
		return "", err
	}

	footer := i18n.T("init.metric.footer")
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
	header := i18n.T("init.run.type.header")
	prompt.Info(header)

	selected, err := in.List(SelectFormulaTypeQuestion, RunTypes)
	if err != nil {
		return formula.DefaultRun, err
	}

	runType := formula.DefaultRun
	for i := range RunTypes {
		if RunTypes[i] == selected {
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
		in.warning()
		fmt.Print(FormulaLocalRunWarning)
	}

	return runType, nil
}

func (in initCmd) addCommonsRepo() error {
	header := i18n.T("init.commons.header")
	prompt.Info(header)

	choose, err := in.Bool(AddTheCommunityRepo, AcceptDeclineOpts)
	if err != nil {
		return err
	}

	metric.CommonsRepoAdded = "yes"
	if !choose {
		in.commonsWarning()
		metric.CommonsRepoAdded = "no"
		return nil
	}

	repo := formula.Repo{
		Provider: "Github",
		Name:     "commons",
		Url:      TemplatesRepoURL,
		Priority: 0,
	}

	s := spinner.StartNew(i18n.T("init.adding.commons.repo"))

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

	in.commonsSuccess(s)

	return nil
}

func (in initCmd) tutorialInit() error {
	tutorialHolder, err := in.tutorial.Find()
	if err != nil {
		return err
	}

	if tutorialHolder.Current == tutorialStatusEnabled {
		header := i18n.T("init.tutorial.header")
		title := i18n.T("init.tutorial.title")
		body := i18n.T("init.tutorial.body")
		commons := i18n.T("init.tutorial.commons")
		prompt.Info(header)
		prompt.Info(title)
		fmt.Println(body)
		prompt.Info(commons)
	}

	return nil
}

func (in initCmd) welcome() {
	const welcome = " _______       _     _             __         _                           ______    _____      _____  \n|_   __ \\     (_)   / |_          [  |       (_)                        .' ___  |  |_   _|    |_   _| \n  | |__) |    __   `| |-'  .---.   | |--.    __    .---.     ______    / .'   \\_|    | |        | |   \n  |  __ /    [  |   | |   / /'`\\]  | .-. |  [  |  / /__\\\\   |______|   | |           | |   _    | |   \n _| |  \\ \\_   | |   | |,  | \\__.   | | | |   | |  | \\__.,              \\ `.___.'\\   _| |__/ |  _| |_  \n|____| |___| [___]  \\__/  '.___.' [___]|__] [___]  '.__.'               `.____ .'  |________| |_____| \n\n"

	prompt.Info(welcome)
	fmt.Print(i18n.T("init.welcome"))
}

func (in initCmd) commonsSuccess(s *spinner.Spinner) {
	success := i18n.T("init.commons.repo.success")
	s.Success(prompt.Green(success))
}

func (in initCmd) commonsWarning() {
	in.warning()
	fmt.Print(addRepoInfo)
	fmt.Println(addRepoMsg)
}

func (in initCmd) warning() {
	warningMsg := i18n.T("init.warning")
	prompt.Warning(warningMsg)
}

func (in initCmd) initSuccess() {
	success := i18n.T("init.successful")
	prompt.Success(success)
}

func (in initCmd) flagToBool(f string, fn string) (bool, error) {
	if result, found := prompt.BoolOpts[f]; found {
		return result, nil
	} else {
		return false, fmt.Errorf(provideValidValue, fn)
	}
}
