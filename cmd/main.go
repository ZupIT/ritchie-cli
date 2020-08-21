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
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
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
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/watcher"
	fworkspace "github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func main() {
	rootCmd := buildCommands()
	err := rootCmd.Execute()
	if err != nil {
		sendMetric(err.Error())
		_, _ = fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}

	sendMetric()
}

func buildCommands() *cobra.Command {
	userHomeDir := api.UserHomeDir()
	ritchieHomeDir := api.RitchieHomeDir()

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
	treeManager := tree.NewTreeManager(ritchieHomeDir, repoLister, api.CoreCmds)
	credSettings := credential.NewSettings(fileManager, dirManager, userHomeDir)
	autocompleteGen := autocomplete.NewGenerator(treeManager)
	credResolver := envcredential.NewResolver(credFinder)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver
	tutorialFinder := rtutorial.NewFinder(ritchieHomeDir, fileManager)
	tutorialSetter := rtutorial.NewSetter(ritchieHomeDir, fileManager)
	tutorialFindSetter := rtutorial.NewFindSetter(ritchieHomeDir, tutorialFinder, tutorialSetter)
	formBuildMake := builder.NewBuildMake()
	formBuildBat := builder.NewBuildBat(fileManager)
	formBuildDocker := builder.NewBuildDocker()
	formulaLocalBuilder := builder.NewBuildLocal(ritchieHomeDir, dirManager, fileManager, treeGen)

	postRunner := runner.NewPostRunner(fileManager, dirManager)
	inputManager := runner.NewInput(envResolvers, fileManager, inputList, inputText, inputBool, inputPassword)
	formulaSetup := runner.NewPreRun(ritchieHomeDir, formBuildMake, formBuildDocker, formBuildBat, dirManager, fileManager)
	formulaRunner := runner.NewFormulaRunner(postRunner, inputManager, formulaSetup, fileManager, ctxFinder)

	formulaCreator := creator.NewCreator(treeManager, dirManager, fileManager, tplManager)
	formulaWorkspace := fworkspace.New(ritchieHomeDir, fileManager)

	watchManager := watcher.New(formulaLocalBuilder, dirManager)
	createBuilder := formula.NewCreateBuilder(formulaCreator, formulaLocalBuilder)

	upgradeManager := upgrade.DefaultManager{Updater: upgrade.DefaultUpdater{}}
	defaultUpgradeResolver := version.DefaultVersionResolver{
		StableVersionUrl: cmd.StableVersionUrl,
		FileUtilService:  fileutil.DefaultService{},
		HttpClient:       &http.Client{Timeout: 1 * time.Second},
	}
	defaultUrlFinder := upgrade.DefaultUrlFinder{}
	rootCmd := cmd.NewRootCmd(ritchieHomeDir, dirManager, tutorialFinder)

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	initCmd := cmd.NewInitCmd(repoAdder, githubRepo, tutorialFinder, inputList, fileManager)
	listCmd := cmd.NewListCmd()
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()
	buildCmd := cmd.NewBuildCmd()
	upgradeCmd := cmd.NewUpgradeCmd(defaultUpgradeResolver, upgradeManager, defaultUrlFinder)
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
	listRepoCmd := cmd.NewListRepoCmd(repoLister, tutorialFinder)
	deleteRepoCmd := cmd.NewDeleteRepoCmd(repoLister, inputList, repoDeleter)
	setPriorityCmd := cmd.NewSetPriorityCmd(inputList, inputInt, repoLister, repoPrioritySetter)
	autocompleteZsh := cmd.NewAutocompleteZsh(autocompleteGen)
	autocompleteBash := cmd.NewAutocompleteBash(autocompleteGen)
	autocompleteFish := cmd.NewAutocompleteFish(autocompleteGen)
	autocompletePowerShell := cmd.NewAutocompletePowerShell(autocompleteGen)

	createFormulaCmd := cmd.NewCreateFormulaCmd(userHomeDir, createBuilder, tplManager, formulaWorkspace, inputText, inputTextValidator, inputList, tutorialFinder)
	buildFormulaCmd := cmd.NewBuildFormulaCmd(userHomeDir, formulaLocalBuilder, formulaWorkspace, watchManager, dirManager, inputText, inputList, tutorialFinder)

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash, autocompleteFish, autocompletePowerShell)
	addCmd.AddCommand(addRepoCmd)
	updateCmd.AddCommand(updateRepoCmd)
	createCmd.AddCommand(createFormulaCmd)
	deleteCmd.AddCommand(deleteCtxCmd, deleteRepoCmd)
	listCmd.AddCommand(listRepoCmd)
	listCmd.AddCommand(listCredentialCmd)
	setCmd.AddCommand(setCredentialCmd, setCtxCmd, setPriorityCmd)
	showCmd.AddCommand(showCtxCmd)
	buildCmd.AddCommand(buildFormulaCmd)

	formulaCmd := cmd.NewFormulaCommand(api.CoreCmds, treeManager, formulaRunner)
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

func sendMetric(err ...string) {
	metricEnable := metric.NewChecker(stream.NewFileManager())
	if metricEnable.Check() {
		var collectData metric.APIData
		metricManager := metric.NewHttpSender(metric.ServerRestURL, http.DefaultClient)
		userIdManager := metric.NewUserIdGenerator()
		data := metric.NewDataCollector(userIdManager)

		collectData, _ = data.Collect(cmd.Version, err...)
		metricManager.Send(collectData)
	}
}
