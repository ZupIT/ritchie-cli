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

package commands

import (
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/ZupIT/ritchie-cli/internal/pkg/config"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/deleter"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/flag"
	fprompt "github.com/ZupIT/ritchie-cli/pkg/formula/input/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/stdin"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner/docker"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner/local"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/formula/validator"
	fworkspace "github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/git/bitbucket"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/git/gitlab"
	"github.com/ZupIT/ritchie-cli/pkg/git/zipremote"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

func ExecutionTime(startTime time.Time) float64 {
	endTime := time.Now()
	return endTime.Sub(startTime).Seconds()
}

var Data metric.DataCollectorManager
var MetricSender = metric.NewHttpSender(metric.ServerRestURL, http.DefaultClient)

func SendMetric(commandExecutionTime float64, err ...string) {
	metricEnable := metric.NewChecker(stream.NewFileManager())
	if metricEnable.Check() {
		var collectData metric.APIData
		collectData, _ = Data.Collect(commandExecutionTime, cmd.Version, err...)
		MetricSender.Send(collectData)
	}
}

func Build() *cobra.Command {
	userHomeDir := api.UserHomeDir()
	ritchieHomeDir := api.RitchieHomeDir()
	ritConfigManager := config.NewManager(ritchieHomeDir)

	// prompt
	inputText := prompt.NewSurveyText()
	inputTextValidator := prompt.NewSurveyTextValidator()
	inputTextDefault := fprompt.NewSurveyDefault()
	inputInt := prompt.NewSurveyInt()
	inputBool := prompt.NewSurveyBool()
	inputPassword := prompt.NewSurveyPassword()
	inputList := prompt.NewSurveyList()
	inputURL := prompt.NewSurveyURL()
	inputMultiselect := prompt.NewSurveyMultiselect()
	inputAutocomplete := prompt.NewInputAutocomplete()

	// deps
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	githubRepo := github.NewRepoManager(http.DefaultClient)
	gitlabRepo := gitlab.NewRepoManager(http.DefaultClient)
	bitbucketRepo := bitbucket.NewRepoManager(http.DefaultClient)
	zipremoteRepo := zipremote.NewRepoManager(http.DefaultClient)

	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: githubRepo, NewRepoInfo: github.NewRepoInfo})
	repoProviders.Add("Gitlab", formula.Git{Repos: gitlabRepo, NewRepoInfo: gitlab.NewRepoInfo})
	repoProviders.Add("Bitbucket", formula.Git{Repos: bitbucketRepo, NewRepoInfo: bitbucket.NewRepoInfo})
	repoProviders.Add("ZipRemote", formula.Git{Repos: zipremoteRepo, NewRepoInfo: zipremote.NewRepoInfo})

	treeGen := tree.NewGenerator(dirManager, fileManager)

	userIdManager := metric.NewUserIdGenerator()
	Data = metric.NewDataCollector(userIdManager, ritchieHomeDir, fileManager)

	repoCreator := repo.NewCreator(ritchieHomeDir, repoProviders, dirManager, fileManager)
	repoLister := repo.NewLister(ritchieHomeDir, fileManager)
	repoWriter := repo.NewWriter(ritchieHomeDir, fileManager)
	repoDetail := repo.NewDetail(repoProviders)
	repoListWriter := repo.NewListWriter(repoLister, repoWriter)
	repoDeleter := repo.NewDeleter(ritchieHomeDir, repoListWriter, dirManager)
	repoListWriteCreator := repo.NewCreateWriteListDetailDeleter(repoLister, repoCreator, repoWriter, repoDetail,
		repoDeleter)
	repoUpdater := repo.NewUpdater(ritchieHomeDir, repoListWriteCreator, treeGen)
	repoListUpdater := repo.NewListUpdater(repoLister, repoUpdater)

	repoAdder := repo.NewAdder(ritchieHomeDir, repoListWriteCreator, treeGen)
	repoAddLister := repo.NewListAdder(repoLister, repoAdder)
	repoPrioritySetter := repo.NewPrioritySetter(repoListWriter)
	repoListDetailWriter := repo.NewListDetailWrite(repoLister, repoDetail, repoWriter)

	tplManager := template.NewManager(api.RitchieHomeDir(), dirManager)
	envFinder := env.NewFinder(ritchieHomeDir, fileManager)
	envSetter := env.NewSetter(ritchieHomeDir, envFinder, fileManager)
	envRemover := env.NewRemover(ritchieHomeDir, envFinder, fileManager)
	envFindSetter := env.NewFindSetter(envFinder, envSetter)
	envFindRemover := env.NewFindRemover(envFinder, envRemover)
	credSetter := credential.NewSetter(ritchieHomeDir, envFinder, dirManager)
	credFinder := credential.NewFinder(ritchieHomeDir, envFinder)
	credDeleter := credential.NewCredDelete(ritchieHomeDir, envFinder)
	credSettings := credential.NewSettings(fileManager, dirManager, userHomeDir)

	treeManager := tree.NewTreeManager(ritchieHomeDir, repoListDetailWriter, api.CoreCmds)
	treeChecker := tree.NewChecker(treeManager)
	autocompleteGen := autocomplete.NewGenerator(treeManager)
	credResolver := credential.NewResolver(credFinder, credSetter, inputPassword)
	tutorialFinder := rtutorial.NewFinder(ritchieHomeDir)
	tutorialSetter := rtutorial.NewSetter(ritchieHomeDir)
	tutorialFindSetter := rtutorial.NewFindSetter(tutorialFinder, tutorialSetter)

	formBuildMake := builder.NewBuildMake()
	formBuildSh := builder.NewBuildShell()
	formBuildBat := builder.NewBuildBat(fileManager)
	formBuildDocker := builder.NewBuildDocker(fileManager)
	formBuildLocal := builder.NewBuildLocal(ritchieHomeDir, dirManager, repoAdder)

	postRunner := runner.NewPostRunner(fileManager, dirManager)

	promptInManager := fprompt.NewInputManager(credResolver, inputList, inputText, inputTextValidator, inputTextDefault, inputBool, inputPassword, inputMultiselect, inputAutocomplete)
	stdinInManager := stdin.NewInputManager(credResolver)
	flagInManager := flag.NewInputManager(credResolver)
	termInputTypes := formula.TermInputTypes{
		api.Prompt: promptInManager,
		api.Stdin:  stdinInManager,
		api.Flag:   flagInManager,
	}

	inputResolver := runner.NewInputResolver(termInputTypes)

	preRunChecker := runner.NewPreRunBuilderChecker(repoLister)
	formulaLocalPreRun := local.NewPreRun(ritchieHomeDir, formBuildMake, formBuildBat, formBuildSh, dirManager, fileManager, preRunChecker)
	formulaLocalRun := local.NewRunner(postRunner, inputResolver, formulaLocalPreRun, fileManager, envFinder, userHomeDir)

	formulaDockerPreRun := docker.NewPreRun(ritchieHomeDir, formBuildDocker, dirManager, fileManager, preRunChecker)
	formulaDockerRun := docker.NewRunner(postRunner, inputResolver, formulaDockerPreRun, fileManager, envFinder, userHomeDir)

	runners := formula.Runners{
		formula.LocalRun:  formulaLocalRun,
		formula.DockerRun: formulaDockerRun,
	}

	formulaCreator := creator.NewCreator(treeManager, dirManager, fileManager, tplManager)
	formulaWorkspace := fworkspace.New(ritchieHomeDir, userHomeDir, dirManager, formBuildLocal, treeGen)

	preRunBuilder := runner.NewPreRunBuilder(formulaWorkspace, formBuildLocal)
	configManager := runner.NewConfigManager(ritchieHomeDir)
	formulaExec := runner.NewExecutor(runners, preRunBuilder, configManager)

	createBuilder := formula.NewCreateBuilder(formulaCreator, formBuildLocal)

	versionManager := version.NewManager(
		version.StableVersionURL,
		fileManager,
	)
	upgradeDefaultUpdater := upgrade.NewDefaultUpdater()
	upgradeManager := upgrade.NewDefaultManager(upgradeDefaultUpdater)
	defaultUrlFinder := upgrade.NewDefaultUrlFinder(versionManager)
	rootCmd := cmd.NewRootCmd(ritchieHomeDir, dirManager, fileManager, tutorialFinder, versionManager, treeGen, repoListWriter)

	validator := validator.New()
	deleter := deleter.NewDeleter(dirManager, fileManager, treeGen, ritchieHomeDir)

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	initCmd := cmd.NewInitCmd(repoAdder, githubRepo, tutorialFinder, configManager, fileManager, inputList, inputBool, MetricSender, ritConfigManager)
	listCmd := cmd.NewListCmd()
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()
	buildCmd := cmd.NewBuildCmd()
	upgradeCmd := cmd.NewUpgradeCmd(
		versionManager,
		upgradeManager,
		defaultUrlFinder,
		inputBool,
		fileManager,
		githubRepo,
	)
	metricsCmd := cmd.NewMetricsCmd(fileManager, inputList)
	tutorialCmd := cmd.NewTutorialCmd(inputList, tutorialFindSetter)
	renameCmd := cmd.NewRenameCmd()

	// level 2
	setCredentialCmd := cmd.NewSetCredentialCmd(
		credSetter,
		credSettings,
		inputText,
		inputBool,
		inputList,
		inputPassword,
	)
	listCredentialCmd := cmd.NewListCredentialCmd(credSettings)
	deleteCredentialCmd := cmd.NewDeleteCredentialCmd(
		credDeleter,
		credSettings,
		envFinder,
		inputBool,
		inputList,
	)

	deleteEnvCmd := cmd.NewDeleteEnvCmd(envFindRemover, inputBool, inputList)
	setEnvCmd := cmd.NewSetEnvCmd(envFindSetter, inputText, inputList)
	showEnvCmd := cmd.NewShowEnvCmd(envFinder)
	addRepoCmd := cmd.NewAddRepoCmd(repoAddLister, repoProviders, credResolver, inputTextValidator, inputURL, inputList, inputBool, inputInt, tutorialFinder, treeChecker, repoDetail)
	addRepoZipCmd := cmd.NewAddRepoZipCmd(repoAddLister, inputTextValidator, inputURL, inputBool, tutorialFinder, treeChecker)
	updateRepoCmd := cmd.NewUpdateRepoCmd(http.DefaultClient, repoListUpdater, repoProviders, inputText, inputPassword, inputURL, inputList, inputBool, inputInt)
	listRepoCmd := cmd.NewListRepoCmd(repoLister, tutorialFinder)
	deleteRepoCmd := cmd.NewDeleteRepoCmd(repoLister, inputList, inputBool, repoDeleter)
	listWorkspaceCmd := cmd.NewListWorkspaceCmd(formulaWorkspace, tutorialFinder)
	setPriorityCmd := cmd.NewSetPriorityCmd(inputList, inputInt, repoLister, repoPrioritySetter)
	autocompleteZsh := cmd.NewAutocompleteZsh(autocompleteGen)
	autocompleteBash := cmd.NewAutocompleteBash(autocompleteGen)
	autocompleteFish := cmd.NewAutocompleteFish(autocompleteGen)
	autocompletePowerShell := cmd.NewAutocompletePowerShell(autocompleteGen)
	deleteWorkspaceCmd := cmd.NewDeleteWorkspaceCmd(userHomeDir, formulaWorkspace, repoDeleter, inputList, inputBool)
	deleteFormulaCmd := cmd.NewDeleteFormulaCmd(userHomeDir, ritchieHomeDir, formulaWorkspace, dirManager, inputBool, inputTextValidator, inputList, inputAutocomplete, treeGen, fileManager)
	addWorkspaceCmd := cmd.NewAddWorkspaceCmd(formulaWorkspace, inputTextValidator, inputAutocomplete)
	updateWorkspaceCmd := cmd.NewUpdateWorkspaceCmd(formulaWorkspace, inputList)

	createFormulaCmd := cmd.NewCreateFormulaCmd(userHomeDir, createBuilder, tplManager, formulaWorkspace, inputText, inputTextValidator, inputList, inputAutocomplete, tutorialFinder, treeChecker, validator, inputBool, dirManager)
	buildFormulaCmd := cmd.NewBuildFormulaCmd()
	showFormulaRunnerCmd := cmd.NewShowFormulaRunnerCmd(configManager)
	setFormulaRunnerCmd := cmd.NewSetFormulaRunnerCmd(configManager, inputList)
	renameFormulaCmd := cmd.NewRenameFormulaCmd(formulaWorkspace, inputList, inputTextValidator, inputBool, dirManager,
		validator, createBuilder, treeGen, deleter, userHomeDir, ritchieHomeDir)
	listFormulaCmd := cmd.NewListFormulaCmd(repoLister, inputList, treeManager, tutorialFinder)

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash, autocompleteFish, autocompletePowerShell)
	addCmd.AddCommand(addRepoCmd, addRepoZipCmd, addWorkspaceCmd)
	updateCmd.AddCommand(updateRepoCmd, updateWorkspaceCmd)
	createCmd.AddCommand(createFormulaCmd)
	deleteCmd.AddCommand(deleteEnvCmd, deleteRepoCmd, deleteFormulaCmd, deleteWorkspaceCmd, deleteCredentialCmd)
	listCmd.AddCommand(listRepoCmd)
	listCmd.AddCommand(listCredentialCmd)
	listCmd.AddCommand(listWorkspaceCmd)
	listCmd.AddCommand(listFormulaCmd)
	setCmd.AddCommand(setCredentialCmd, setEnvCmd, setPriorityCmd, setFormulaRunnerCmd)
	showCmd.AddCommand(showEnvCmd, showFormulaRunnerCmd)
	buildCmd.AddCommand(buildFormulaCmd)
	renameCmd.AddCommand(renameFormulaCmd)

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
				renameCmd,
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
