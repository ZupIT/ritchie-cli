package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"

	"k8s.io/kubectl/pkg/util/templates"

	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator"

	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credsingle"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/env/envcredential"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/watcher"
	fworkspace "github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/security/secsingle"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/session/sesssingle"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/workspace"
)

func main() {
	rootCmd := buildCommands()
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}
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
	inputMultiline := prompt.NewSurveyMultiline()

	// deps
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	treeGen := tree.NewGenerator(dirManager, fileManager)
	repoAdder := repo.NewAdder(ritchieHomeDir, http.DefaultClient, treeGen, dirManager, fileManager)
	repoLister := repo.NewLister(ritchieHomeDir, fileManager)

	sessionManager := session.NewManager(ritchieHomeDir)
	workspaceManager := workspace.NewChecker(ritchieHomeDir)
	ctxFinder := rcontext.NewFinder(ritchieHomeDir)
	ctxSetter := rcontext.NewSetter(ritchieHomeDir, ctxFinder)
	ctxRemover := rcontext.NewRemover(ritchieHomeDir, ctxFinder)
	ctxFindSetter := rcontext.NewFindSetter(ritchieHomeDir, ctxFinder, ctxSetter)
	ctxFindRemover := rcontext.NewFindRemover(ritchieHomeDir, ctxFinder, ctxRemover)
	sessionValidator := sesssingle.NewValidator(sessionManager)
	passphraseManager := secsingle.NewPassphraseManager(sessionManager)
	credSetter := credsingle.NewSetter(ritchieHomeDir, ctxFinder, sessionManager)
	credFinder := credsingle.NewFinder(ritchieHomeDir, ctxFinder, sessionManager)
	treeManager := tree.NewTreeManager(ritchieHomeDir, repoLister, api.SingleCoreCmds)
	autocompleteGen := autocomplete.NewGenerator(treeManager)
	credResolver := envcredential.NewResolver(credFinder)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver

	inputManager := runner.NewInputManager(envResolvers, inputList, inputText, inputBool, inputPassword)
	formulaSetup := runner.NewDefaultSingleSetup(ritchieHomeDir, http.DefaultClient)

	defaultPreRunner := runner.NewDefaultPreRunner(formulaSetup)
	dockerPreRunner := runner.NewDockerPreRunner(formulaSetup)

	postRunner := runner.NewPostRunner()

	defaultRunner := runner.NewDefaultRunner(defaultPreRunner, postRunner, inputManager)
	dockerRunner := runner.NewDockerRunner(dockerPreRunner, postRunner, inputManager)

	formulaCreator := creator.NewCreator(treeManager, dirManager, fileManager)
	formulaWorkspace := fworkspace.New(ritchieHomeDir, fileManager)
	formulaBuilder := builder.New(ritchieHomeDir, dirManager, fileManager)
	watchManager := watcher.New(formulaBuilder, dirManager)
	createBuilder := formula.NewCreateBuilder(formulaCreator, formulaBuilder)

	upgradeManager := upgrade.DefaultManager{Updater: upgrade.DefaultUpdater{}}
	defaultUpgradeResolver := version.DefaultVersionResolver{
		StableVersionUrl: cmd.StableVersionUrl,
		FileUtilService:  fileutil.DefaultService{},
		HttpClient:       &http.Client{Timeout: 1 * time.Second},
	}
	defaultUrlFinder := upgrade.DefaultUrlFinder{}
	rootCmd := cmd.NewSingleRootCmd(workspaceManager, sessionValidator)

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	cleanCmd := cmd.NewCleanCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	initCmd := cmd.NewSingleInitCmd(inputPassword, passphraseManager)
	listCmd := cmd.NewListCmd()
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()
	buildCmd := cmd.NewBuildCmd()
	upgradeCmd := cmd.NewUpgradeCmd(api.Single, defaultUpgradeResolver, upgradeManager, defaultUrlFinder)

	// level 2
	setCredentialCmd := cmd.NewSingleSetCredentialCmd(
		credSetter,
		inputText,
		inputBool,
		inputList,
		inputPassword,
		inputMultiline)
	deleteCtxCmd := cmd.NewDeleteContextCmd(ctxFindRemover, inputBool, inputList)
	setCtxCmd := cmd.NewSetContextCmd(ctxFindSetter, inputText, inputList)
	showCtxCmd := cmd.NewShowContextCmd(ctxFinder)
	addRepoCmd := cmd.NewAddRepoCmd(http.DefaultClient, repoAdder, inputText, inputPassword, inputURL, inputList, inputBool, inputInt)
	autocompleteZsh := cmd.NewAutocompleteZsh(autocompleteGen)
	autocompleteBash := cmd.NewAutocompleteBash(autocompleteGen)

	createFormulaCmd := cmd.NewCreateFormulaCmd(userHomeDir, createBuilder, formulaWorkspace, inputText, inputTextValidator, inputList)
	buildFormulaCmd := cmd.NewBuildFormulaCmd(userHomeDir, formulaBuilder, formulaWorkspace, watchManager, dirManager, inputText, inputList)

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash)
	addCmd.AddCommand(addRepoCmd)
	createCmd.AddCommand(createFormulaCmd)
	deleteCmd.AddCommand(deleteCtxCmd)
	setCmd.AddCommand(setCredentialCmd, setCtxCmd)
	showCmd.AddCommand(showCtxCmd)
	buildCmd.AddCommand(buildFormulaCmd)

	formulaCmd := cmd.NewFormulaCommand(api.SingleCoreCmds, treeManager, defaultRunner, dockerRunner)
	if err := formulaCmd.Add(rootCmd); err != nil {
		panic(err)
	}

	groups := templates.CommandGroups{
		{
			Message: api.CoreCmdsDesc,
			Commands: []*cobra.Command{
				addCmd,
				autocompleteCmd,
				cleanCmd,
				createCmd,
				deleteCmd,
				initCmd,
				listCmd,
				setCmd,
				showCmd,
				updateCmd,
				buildCmd,
				upgradeCmd,
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
