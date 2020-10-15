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

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"k8s.io/kubectl/pkg/util/templates"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/flag"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/stdin"
	fprompt "github.com/ZupIT/ritchie-cli/pkg/formula/input/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner/docker"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner/local"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/git/gitlab"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"

	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/env/envcredential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/watcher"
	fworkspace "github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func main() {
	startTime := time.Now()
	rootCmd := buildCommands()
	err := rootCmd.Execute()
	if err != nil {
		sendMetric(executionTime(startTime), err.Error())
		errFmt := fmt.Sprintf("%+v", err)
		errFmt = prompt.Red(errFmt)
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", errFmt)
		os.Exit(1)
	}
	sendMetric(executionTime(startTime))
}

func executionTime(startTime time.Time) float64 {
	endTime := time.Now()
	return endTime.Sub(startTime).Seconds()
}

var Data metric.DataCollectorManager
var MetricSender = metric.NewHttpSender(metric.ServerRestURL, http.DefaultClient)

func buildCommands() *cobra.Command {
	userHomeDir := api.UserHomeDir()
	ritchieHomeDir := api.RitchieHomeDir()
	isRootCommand := len(os.Args[1:]) == 0

	// prompt
	inputText := prompt.NewSurveyText()
	inputTextValidator := prompt.NewSurveyTextValidator()
	inputInt := prompt.NewSurveyInt()
	inputBool := prompt.NewSurveyBool()
	inputPassword := prompt.NewSurveyPassword()
	inputList := prompt.NewSurveyList()
	inputURL := prompt.NewSurveyURL()

	// deps
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	githubRepo := github.NewRepoManager(http.DefaultClient)
	gitlabRepo := gitlab.NewRepoManager(http.DefaultClient)

	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: githubRepo, NewRepoInfo: github.NewRepoInfo})
	repoProviders.Add("Gitlab", formula.Git{Repos: gitlabRepo, NewRepoInfo: gitlab.NewRepoInfo})

	treeGen := tree.NewGenerator(dirManager, fileManager)

	userIdManager := metric.NewUserIdGenerator()
	Data = metric.NewDataCollector(userIdManager, ritchieHomeDir, fileManager)

	repoCreator := repo.NewCreator(ritchieHomeDir, repoProviders, dirManager, fileManager)
	repoLister := repo.NewLister(ritchieHomeDir, fileManager)
	repoAdder := repo.NewAdder(ritchieHomeDir, repoCreator, treeGen, dirManager, fileManager)
	repoListCreator := repo.NewListCreator(repoLister, repoCreator)
	repoUpdater := repo.NewUpdater(ritchieHomeDir, repoListCreator, treeGen, fileManager)
	repoAddLister := repo.NewListAdder(repoLister, repoAdder)
	repoListUpdater := repo.NewListUpdater(repoLister, repoUpdater)
	repoDeleter := repo.NewDeleter(ritchieHomeDir, fileManager, dirManager)
	repoPrioritySetter := repo.NewPrioritySetter(ritchieHomeDir, fileManager)

	tplManager := template.NewManager(api.RitchieHomeDir(), dirManager)
	ctxFinder := rcontext.NewFinder(ritchieHomeDir, fileManager)
	ctxSetter := rcontext.NewSetter(ritchieHomeDir, ctxFinder)
	ctxRemover := rcontext.NewRemover(ritchieHomeDir, ctxFinder)
	ctxFindSetter := rcontext.NewFindSetter(ritchieHomeDir, ctxFinder, ctxSetter)
	ctxFindRemover := rcontext.NewFindRemover(ritchieHomeDir, ctxFinder, ctxRemover)
	credSetter := credential.NewSetter(ritchieHomeDir, ctxFinder)
	credFinder := credential.NewFinder(ritchieHomeDir, ctxFinder, fileManager)
	treeManager := tree.NewTreeManager(ritchieHomeDir, repoLister, api.CoreCmds, fileManager, repoProviders, isRootCommand)
	credSettings := credential.NewSettings(fileManager, dirManager, userHomeDir)
	autocompleteGen := autocomplete.NewGenerator(treeManager)
	credResolver := envcredential.NewResolver(credFinder, credSetter, inputPassword)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver
	tutorialFinder := rtutorial.NewFinder(ritchieHomeDir, fileManager)
	tutorialSetter := rtutorial.NewSetter(ritchieHomeDir, fileManager)
	tutorialFindSetter := rtutorial.NewFindSetter(ritchieHomeDir, tutorialFinder, tutorialSetter)
	formBuildMake := builder.NewBuildMake()
	formBuildSh := builder.NewBuildShell()
	formBuildBat := builder.NewBuildBat(fileManager)
	formBuildDocker := builder.NewBuildDocker(fileManager)
	formulaLocalBuilder := builder.NewBuildLocal(ritchieHomeDir, dirManager, fileManager, treeGen)

	postRunner := runner.NewPostRunner(fileManager, dirManager)

	promptInManager := fprompt.NewInputManager(envResolvers, fileManager, inputList, inputText, inputTextValidator, inputBool, inputPassword)
	stdinInManager := stdin.NewInputManager(envResolvers)
	flagInManager := flag.NewInputManager(envResolvers, promptInManager)
	termInputTypes := formula.TermInputTypes{
		api.Prompt: promptInManager,
		api.Stdin:  stdinInManager,
		api.Flag:   flagInManager,
	}

	inputResolver := runner.NewInputResolver(termInputTypes)

	formulaLocalPreRun := local.NewPreRun(ritchieHomeDir, formBuildMake, formBuildBat, formBuildSh, dirManager, fileManager)
	formulaLocalRun := local.NewRunner(postRunner, inputResolver, formulaLocalPreRun, fileManager, ctxFinder, userHomeDir)

	formulaDockerPreRun := docker.NewPreRun(ritchieHomeDir, formBuildDocker, dirManager, fileManager)
	formulaDockerRun := docker.NewRunner(postRunner, inputResolver, formulaDockerPreRun, fileManager, ctxFinder, userHomeDir)

	runners := formula.Runners{
		formula.LocalRun:  formulaLocalRun,
		formula.DockerRun: formulaDockerRun,
	}

	configManager := runner.NewConfigManager(ritchieHomeDir, fileManager)
	formulaExec := runner.NewExecutor(runners, configManager)

	formulaCreator := creator.NewCreator(treeManager, dirManager, fileManager, tplManager)
	formulaWorkspace := fworkspace.New(ritchieHomeDir, fileManager)

	watchManager := watcher.New(formulaLocalBuilder, dirManager, sendMetric)
	createBuilder := formula.NewCreateBuilder(formulaCreator, formulaLocalBuilder)

	versionManager := version.NewManager(
		version.StableVersionUrl,
		fileManager,
	)
	upgradeDefaultUpdater := upgrade.NewDefaultUpdater()
	upgradeManager := upgrade.NewDefaultManager(upgradeDefaultUpdater)
	defaultUrlFinder := upgrade.NewDefaultUrlFinder(versionManager)
	rootCmd := cmd.NewRootCmd(ritchieHomeDir, dirManager, tutorialFinder, versionManager)

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	initCmd := cmd.NewInitCmd(repoAdder, githubRepo, tutorialFinder, configManager, fileManager, inputList, inputBool, MetricSender)
	listCmd := cmd.NewListCmd()
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()
	buildCmd := cmd.NewBuildCmd()
	upgradeCmd := cmd.NewUpgradeCmd(versionManager, upgradeManager, defaultUrlFinder, inputList, fileManager)
	metricsCmd := cmd.NewMetricsCmd(fileManager, inputList)
	tutorialCmd := cmd.NewTutorialCmd(ritchieHomeDir, inputList, tutorialFindSetter)

	// level 2
	setCredentialCmd := cmd.NewSetCredentialCmd(
		credSetter,
		credSettings,
		fileManager,
		inputText,
		inputBool,
		inputList,
		inputPassword)
	listCredentialCmd := cmd.NewListCredentialCmd(credSettings)

	deleteCtxCmd := cmd.NewDeleteContextCmd(ctxFindRemover, inputBool, inputList)
	setCtxCmd := cmd.NewSetContextCmd(ctxFindSetter, inputText, inputList)
	showCtxCmd := cmd.NewShowContextCmd(ctxFinder)
	addRepoCmd := cmd.NewAddRepoCmd(repoAddLister, repoProviders, inputTextValidator, inputPassword, inputURL, inputList, inputBool, inputInt, tutorialFinder)
	updateRepoCmd := cmd.NewUpdateRepoCmd(http.DefaultClient, repoListUpdater, repoProviders, inputText, inputPassword, inputURL, inputList, inputBool, inputInt)
	listRepoCmd := cmd.NewListRepoCmd(repoLister, repoProviders, tutorialFinder)
	deleteRepoCmd := cmd.NewDeleteRepoCmd(repoLister, inputList, repoDeleter)
	setPriorityCmd := cmd.NewSetPriorityCmd(inputList, inputInt, repoLister, repoPrioritySetter)
	autocompleteZsh := cmd.NewAutocompleteZsh(autocompleteGen)
	autocompleteBash := cmd.NewAutocompleteBash(autocompleteGen)
	autocompleteFish := cmd.NewAutocompleteFish(autocompleteGen)
	autocompletePowerShell := cmd.NewAutocompletePowerShell(autocompleteGen)
	deleteWorkspaceCmd := cmd.NewDeleteWorkspaceCmd(userHomeDir, formulaWorkspace, dirManager, inputList, inputBool)
	deleteFormulaCmd := cmd.NewDeleteFormulaCmd(userHomeDir, ritchieHomeDir, formulaWorkspace, dirManager, inputBool, inputText, inputList, treeGen, fileManager)

	createFormulaCmd := cmd.NewCreateFormulaCmd(userHomeDir, createBuilder, tplManager, formulaWorkspace, inputText, inputTextValidator, inputList, tutorialFinder)
	buildFormulaCmd := cmd.NewBuildFormulaCmd(userHomeDir, formulaLocalBuilder, formulaWorkspace, watchManager, dirManager, inputText, inputList, tutorialFinder)
	showFormulaRunnerCmd := cmd.NewShowFormulaRunnerCmd(configManager)
	setFormulaRunnerCmd := cmd.NewSetFormulaRunnerCmd(configManager, inputList)

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash, autocompleteFish, autocompletePowerShell)
	addCmd.AddCommand(addRepoCmd)
	updateCmd.AddCommand(updateRepoCmd)
	createCmd.AddCommand(createFormulaCmd)
	deleteCmd.AddCommand(deleteCtxCmd, deleteRepoCmd, deleteFormulaCmd, deleteWorkspaceCmd)
	listCmd.AddCommand(listRepoCmd)
	listCmd.AddCommand(listCredentialCmd)
	setCmd.AddCommand(setCredentialCmd, setCtxCmd, setPriorityCmd, setFormulaRunnerCmd)
	showCmd.AddCommand(showCtxCmd, showFormulaRunnerCmd)
	buildCmd.AddCommand(buildFormulaCmd)

	formulaCmd := cmd.NewFormulaCommand(api.CoreCmds, treeManager, formulaExec, fileManager)
	if err := formulaCmd.Add(rootCmd); err != nil {
		panic(err)
	}

	groups := templates.CommandGroups{
		{
			Message: api.CoreCmdsDesc,
			Commands: []*cobra.Command{
				addCmd,
				autocompleteCmd,
				createCmd,
				deleteCmd,
				initCmd,
				listCmd,
				setCmd,
				showCmd,
				updateCmd,
				buildCmd,
				upgradeCmd,
				tutorialCmd,
				metricsCmd,
			},
		},
	}

	cmds := rootCmd.Commands()
	for _, c := range cmds {
		exists := false
		g := c.Annotations[cmd.Group]
		for i, v := range groups {
			if v.Message == g {
				v.Commands = append(v.Commands, c)
				groups[i] = v
				exists = true
			}
		}
		if !exists {
			cg := templates.CommandGroup{
				Message:  c.Annotations[cmd.Group],
				Commands: []*cobra.Command{c},
			}
			groups = append(groups, cg)
		}
	}

	rootCmd.ResetCommands()
	groups.Add(rootCmd)
	templates.ActsAsRootCommand(rootCmd, nil, groups...)

	return rootCmd
}

func sendMetric(commandExecutionTime float64, err ...string) {
	metricEnable := metric.NewChecker(stream.NewFileManager())
	if metricEnable.Check() {
		var collectData metric.APIData
		collectData, _ = Data.Collect(commandExecutionTime, cmd.Version, err...)
		MetricSender.Send(collectData)
	}
}
