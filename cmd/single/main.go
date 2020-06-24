package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"

	"k8s.io/kubectl/pkg/util/templates"

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
	inputInt := prompt.NewSurveyInt()
	inputBool := prompt.NewSurveyBool()
	inputPassword := prompt.NewSurveyPassword()
	inputList := prompt.NewSurveyList()
	inputURL := prompt.NewSurveyURL()

	// deps
	sessionManager := session.NewManager(ritchieHomeDir)
	workspaceManager := workspace.NewChecker(ritchieHomeDir)
	ctxFinder := rcontext.NewFinder(ritchieHomeDir)
	ctxSetter := rcontext.NewSetter(ritchieHomeDir, ctxFinder)
	ctxRemover := rcontext.NewRemover(ritchieHomeDir, ctxFinder)
	ctxFindSetter := rcontext.NewFindSetter(ritchieHomeDir, ctxFinder, ctxSetter)
	ctxFindRemover := rcontext.NewFindRemover(ritchieHomeDir, ctxFinder, ctxRemover)
	repoManager := formula.NewSingleRepoManager(ritchieHomeDir, http.DefaultClient, sessionManager)
	repoLoader := formula.NewSingleLoader(cmd.CommonsRepoURL, repoManager)
	sessionValidator := sesssingle.NewValidator(sessionManager)
	passphraseManager := secsingle.NewPassphraseManager(sessionManager)
	credSetter := credsingle.NewSetter(ritchieHomeDir, ctxFinder, sessionManager)
	credFinder := credsingle.NewFinder(ritchieHomeDir, ctxFinder, sessionManager)
	treeManager := formula.NewTreeManager(ritchieHomeDir, repoManager, api.SingleCoreCmds)
	autocompleteGen := autocomplete.NewGenerator(treeManager)
	credResolver := envcredential.NewResolver(credFinder)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver

	inputManager := formula.NewInputManager(envResolvers, inputList, inputText, inputBool)
	formulaSetup := formula.NewDefaultSingleSetup(ritchieHomeDir, http.DefaultClient)

	defaultPreRunner := formula.NewDefaultPreRunner(formulaSetup)
	dockerPreRunner := formula.NewDockerPreRunner(formulaSetup)

	postRunner := formula.NewPostRunner()

	defaultRunner := formula.NewDefaultRunner(defaultPreRunner, postRunner, inputManager)
	dockerRunner := formula.NewDockerRunner(dockerPreRunner, postRunner, inputManager)

	formulaCreator := formula.NewCreator(userHomeDir, treeManager)

	upgradeManager := upgrade.DefaultManager{Updater: upgrade.DefaultUpdater{}}
	defaultUpgradeResolver := version.DefaultVersionResolver{
		StableVersionUrl: cmd.StableVersionUrl,
		FileUtilService:  fileutil.DefaultService{},
		HttpClient:       &http.Client{Timeout: 1 * time.Second},
	}
	upgradeUrl := upgrade.UpgradeUrl(api.Single, defaultUpgradeResolver)

	rootCmd := cmd.NewSingleRootCmd(workspaceManager, sessionValidator)

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	cleanCmd := cmd.NewCleanCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	initCmd := cmd.NewSingleInitCmd(inputPassword, passphraseManager, repoLoader)
	listCmd := cmd.NewListCmd()
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()
	buildCmd := cmd.NewBuildCmd()
	upgradeCmd := cmd.NewUpgradeCmd(upgradeUrl, upgradeManager)

	// level 2
	setCredentialCmd := cmd.NewSingleSetCredentialCmd(
		credSetter,
		inputText,
		inputBool,
		inputList,
		inputPassword)
	deleteCtxCmd := cmd.NewDeleteContextCmd(ctxFindRemover, inputBool, inputList)
	setCtxCmd := cmd.NewSetContextCmd(ctxFindSetter, inputText, inputList)
	showCtxCmd := cmd.NewShowContextCmd(ctxFinder)
	addRepoCmd := cmd.NewAddRepoCmd(repoManager, inputText, inputURL, inputInt, inputBool)
	cleanRepoCmd := cmd.NewCleanRepoCmd(repoManager, inputText)
	deleteRepoCmd := cmd.NewDeleteRepoCmd(repoManager, inputList, inputBool)
	listRepoCmd := cmd.NewListRepoCmd(repoManager)
	updateRepoCmd := cmd.NewUpdateRepoCmd(repoManager)
	autocompleteZsh := cmd.NewAutocompleteZsh(autocompleteGen)
	autocompleteBash := cmd.NewAutocompleteBash(autocompleteGen)
	createFormulaCmd := cmd.NewCreateFormulaCmd(formulaCreator, inputText, inputList, inputBool)
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	formulaWorkspace := fworkspace.New(ritchieHomeDir, fileManager)
	formulaBuilder := builder.New(ritchieHomeDir, dirManager, fileManager)
	watchManager := watcher.New(formulaBuilder, dirManager)
	buildFormulaCmd := cmd.NewBuildFormulaCmd(userHomeDir, formulaWorkspace, formulaBuilder, watchManager, dirManager, inputText, inputList)

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash)
	addCmd.AddCommand(addRepoCmd)
	cleanCmd.AddCommand(cleanRepoCmd)
	createCmd.AddCommand(createFormulaCmd)
	deleteCmd.AddCommand(deleteRepoCmd, deleteCtxCmd)
	listCmd.AddCommand(listRepoCmd)
	setCmd.AddCommand(setCredentialCmd, setCtxCmd)
	showCmd.AddCommand(showCtxCmd)
	updateCmd.AddCommand(updateRepoCmd)
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
