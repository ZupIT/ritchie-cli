package main

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"net/http"
	"os"
	"os/user"
	"time"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credsingle"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credteam"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/env/envcredential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/metrics"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/security/secsingle"
	"github.com/ZupIT/ritchie-cli/pkg/security/secteam"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/session/sesssingle"
	"github.com/ZupIT/ritchie-cli/pkg/session/sessteam"
	"github.com/ZupIT/ritchie-cli/pkg/workspace"
)

const (
	ritchieHomePattern = "%s/.rit"
)

func main() {
	homePath := ritchieHomePath()
	rootCmd := buildCommands(homePath)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}
}

func ritchieHomePath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(ritchieHomePattern, usr.HomeDir)
}

func buildCommands(ritchieHomePath string) *cobra.Command {
	sessionManager := session.NewManager(ritchieHomePath)
	workspaceManager := workspace.NewChecker(ritchieHomePath)
	ctxFinder := rcontext.NewFinder(ritchieHomePath)
	ctxSetter := rcontext.NewSetter(ritchieHomePath, ctxFinder)
	ctxRemover := rcontext.NewRemover(ritchieHomePath, ctxFinder)
	ctxFindSetter := rcontext.NewFindSetter(ritchieHomePath, ctxFinder, ctxSetter)
	ctxFindRemover := rcontext.NewFindRemover(ritchieHomePath, ctxFinder, ctxRemover)
	repoManager := formula.NewRepoManager(ritchieHomePath, http.DefaultClient, sessionManager)

	var sessionValidator session.Validator
	var loginManager security.LoginManager
	var logoutManager security.LogoutManager
	var userManager security.UserManager
	var credSetter credential.Setter
	var credFinder credential.Finder
	var credSettings credential.Settings
	var coreCmds []api.Command

	switch env.Edition {
	case env.Single:
		coreCmds = api.SingleCoreCmds
		sessionValidator = sesssingle.NewValidator(sessionManager)
		loginManager = secsingle.NewLoginManager(sessionManager)
		credSetter = credsingle.NewSetter(ritchieHomePath, ctxFinder, sessionManager)
		credFinder = credsingle.NewFinder(ritchieHomePath, ctxFinder, sessionManager)
	case env.Team:
		coreCmds = api.TeamCoreCmds
		sessionValidator = sessteam.NewValidator(sessionManager)
		loginManager = secteam.NewLoginManager(ritchieHomePath, env.ServerURL, security.OAuthProvider, http.DefaultClient, sessionManager)
		logoutManager = secteam.NewLogoutManager(security.OAuthProvider, sessionManager)
		userManager = secteam.NewUserManager(env.ServerURL, http.DefaultClient, sessionManager)
		credSetter = credteam.NewSetter(env.ServerURL, http.DefaultClient, sessionManager, ctxFinder)
		credFinder = credteam.NewFinder(env.ServerURL, http.DefaultClient, sessionManager, ctxFinder)
		credSettings = credteam.NewSettings(env.ServerURL, http.DefaultClient, sessionManager, ctxFinder)
	default:
		panic("The env.Edition is required on build")
	}

	treeManager := formula.NewTreeManager(ritchieHomePath, repoManager, coreCmds)
	autocompleteGen := autocomplete.NewGenerator(treeManager)

	credResolver := envcredential.NewResolver(credFinder)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver

	formulaRunner := formula.NewRunner(ritchieHomePath, envResolvers, http.DefaultClient, treeManager)
	formulaCreator := formula.NewCreator(ritchieHomePath, treeManager)

	rootCmd := cmd.NewRootCmd(workspaceManager, loginManager, sessionValidator)

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	cleanCmd := cmd.NewCleanCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	listCmd := cmd.NewListCmd()
	loginCmd := cmd.NewLoginCmd(loginManager, repoManager)
	logoutCmd := cmd.NewLogoutCmd(logoutManager)
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()

	// level 2
	setCredentialCmd := cmd.NewSetCredentialCmd(credSetter, credSettings)
	createUserCmd := cmd.NewCreateUserCmd(userManager)
	deleteUserCmd := cmd.NewDeleteUserCmd(userManager)
	deleteCtxCmd := cmd.NewDeleteContextCmd(ctxFindRemover)
	setCtxCmd := cmd.NewSetContextCmd(ctxFindSetter)
	showCtxCmd := cmd.NewShowContextCmd(ctxFinder)
	addRepoCmd := cmd.NewAddRepoCmd(repoManager)
	cleanRepoCmd := cmd.NewCleanRepoCmd(repoManager)
	deleteRepoCmd := cmd.NewDeleteRepoCmd(repoManager)
	listRepoCmd := cmd.NewListRepoCmd(repoManager)
	updateRepoCmd := cmd.NewUpdateRepoCmd(repoManager)
	autocompleteZsh := cmd.NewAutocompleteZsh(autocompleteGen)
	autocompleteBash := cmd.NewAutocompleteBash(autocompleteGen)
	createFormulaCmd := cmd.NewCreateFormulaCmd(formulaCreator)

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash)
	addCmd.AddCommand(addRepoCmd)
	cleanCmd.AddCommand(cleanRepoCmd)
	createCmd.AddCommand(createUserCmd, createFormulaCmd)
	deleteCmd.AddCommand(deleteUserCmd, deleteRepoCmd, deleteCtxCmd)
	listCmd.AddCommand(listRepoCmd)
	setCmd.AddCommand(setCredentialCmd, setCtxCmd)
	showCmd.AddCommand(showCtxCmd)
	updateCmd.AddCommand(updateRepoCmd)

	rootCmd.AddCommand(addCmd, autocompleteCmd, cleanCmd, createCmd, deleteCmd, listCmd, loginCmd, logoutCmd, setCmd, showCmd, updateCmd)
	if env.Edition == env.Single {
		createCmd.RemoveCommand(createUserCmd)
		deleteCmd.RemoveCommand(deleteUserCmd)
		rootCmd.RemoveCommand(loginCmd, logoutCmd)
	}

	formulaCmd := cmd.NewFormulaCommand(coreCmds, treeManager, formulaRunner)
	err := formulaCmd.Add(rootCmd)
	if err != nil {
		panic(err)
	}

	sendMetrics(sessionManager)

	return rootCmd
}

func sendMetrics(sessionManager session.DefaultManager) {
	if env.Edition == env.Team {
		metricsManager := metrics.NewSender(env.ServerURL, &http.Client{Timeout: 2 * time.Second}, sessionManager)
		go metricsManager.SendCommand()
	}
}
