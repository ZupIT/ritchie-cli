package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"k8s.io/kubectl/pkg/util/templates"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/server"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credteam"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/env/envcredential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/metrics"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/security/secteam"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/session/sessteam"
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
	inputText := prompt.NewInputText()
	inputInt := prompt.NewInputInt()
	inputBool := prompt.NewInputBool()
	inputEmail := prompt.NewInputEmail()
	inputPassword := prompt.NewInputPassword()
	inputList := prompt.NewInputList()
	inputURL := prompt.NewInputURL()

	// deps
	sessionManager := session.NewManager(ritchieHomeDir)
	workspaceManager := workspace.NewChecker(ritchieHomeDir)
	serverFinder := server.NewFinder(ritchieHomeDir)
	serverValidator := server.NewValidator(serverFinder)
	ctxFinder := rcontext.NewFinder(ritchieHomeDir)
	ctxSetter := rcontext.NewSetter(ritchieHomeDir, ctxFinder)
	ctxRemover := rcontext.NewRemover(ritchieHomeDir, ctxFinder)
	ctxFindSetter := rcontext.NewFindSetter(ritchieHomeDir, ctxFinder, ctxSetter)
	ctxFindRemover := rcontext.NewFindRemover(ritchieHomeDir, ctxFinder, ctxRemover)
	serverSetter := server.NewSetter(ritchieHomeDir, http.DefaultClient)
	repoManager := formula.NewTeamRepoManager(ritchieHomeDir, serverFinder, http.DefaultClient, sessionManager)
	repoLoader := formula.NewTeamLoader(serverFinder, http.DefaultClient, sessionManager, repoManager)
	sessionValidator := sessteam.NewValidator(sessionManager)
	loginManager := secteam.NewLoginManager(
		ritchieHomeDir,
		serverFinder,
		security.OAuthProvider,
		http.DefaultClient,
		sessionManager)
	logoutManager := secteam.NewLogoutManager(security.OAuthProvider, sessionManager, serverFinder)
	userManager := secteam.NewUserManager(serverFinder, http.DefaultClient, sessionManager)
	credSetter := credteam.NewSetter(serverFinder, http.DefaultClient, sessionManager, ctxFinder)
	credFinder := credteam.NewFinder(serverFinder, http.DefaultClient, sessionManager, ctxFinder)
	credSettings := credteam.NewSettings(serverFinder, http.DefaultClient, sessionManager, ctxFinder)
	treeManager := formula.NewTreeManager(ritchieHomeDir, repoManager, api.TeamCoreCmds)
	autocompleteGen := autocomplete.NewGenerator(treeManager)
	credResolver := envcredential.NewResolver(credFinder)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver

	inputManager := formula.NewInputManager(envResolvers, inputList, inputText, inputBool)
	formulaSetup := formula.NewDefaultTeamSetup(ritchieHomeDir, http.DefaultClient, sessionManager)

	defaultPreRunner := formula.NewDefaultPreRunner(formulaSetup)
	dockerPreRunner := formula.NewDockerPreRunner(formulaSetup)
	postRunner := formula.NewPostRunner()

	defaultRunner := formula.NewDefaultRunner(defaultPreRunner, postRunner, inputManager)
	dockerRunner := formula.NewDockerRunner(dockerPreRunner, postRunner, inputManager)

	formulaCreator := formula.NewCreator(userHomeDir, treeManager)

	// commands
	rootCmd := cmd.NewTeamRootCmd(
		workspaceManager,
		loginManager,
		repoLoader,
		serverValidator,
		sessionValidator,
		api.Team,
		inputText,
		inputPassword)

	rootCmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	cleanCmd := cmd.NewCleanCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	listCmd := cmd.NewListCmd()
	loginCmd := cmd.NewLoginCmd(loginManager, repoLoader, inputText)
	logoutCmd := cmd.NewLogoutCmd(logoutManager)
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()

	// level 2
	setCredentialCmd := cmd.NewTeamSetCredentialCmd(
		credSetter,
		credSettings,
		inputText,
		inputBool,
		inputList,
		inputPassword)
	createUserCmd := cmd.NewCreateUserCmd(userManager, inputText, inputEmail, inputPassword)
	deleteUserCmd := cmd.NewDeleteUserCmd(userManager, inputBool, inputText)
	deleteCtxCmd := cmd.NewDeleteContextCmd(ctxFindRemover, inputBool, inputList)
	setCtxCmd := cmd.NewSetContextCmd(ctxFindSetter, inputText, inputList)
	setServerCmd := cmd.NewSetServerCmd(serverSetter, inputURL)
	showCtxCmd := cmd.NewShowContextCmd(ctxFinder)
	addRepoCmd := cmd.NewAddRepoCmd(repoManager, inputText, inputURL, inputInt, inputBool)
	cleanRepoCmd := cmd.NewCleanRepoCmd(repoManager, inputText)
	deleteRepoCmd := cmd.NewDeleteRepoCmd(repoManager, inputList, inputBool)
	listRepoCmd := cmd.NewListRepoCmd(repoManager)
	updateRepoCmd := cmd.NewUpdateRepoCmd(repoManager)
	autocompleteZsh := cmd.NewAutocompleteZsh(autocompleteGen)
	autocompleteBash := cmd.NewAutocompleteBash(autocompleteGen)
	createFormulaCmd := cmd.NewCreateFormulaCmd(formulaCreator, inputText, inputList)

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash)
	addCmd.AddCommand(addRepoCmd)
	cleanCmd.AddCommand(cleanRepoCmd)
	createCmd.AddCommand(createUserCmd, createFormulaCmd)
	deleteCmd.AddCommand(deleteUserCmd, deleteRepoCmd, deleteCtxCmd)
	listCmd.AddCommand(listRepoCmd)
	setCmd.AddCommand(setCredentialCmd, setCtxCmd, setServerCmd)
	showCmd.AddCommand(showCtxCmd)
	updateCmd.AddCommand(updateRepoCmd)

	formulaCmd := cmd.NewFormulaCommand(api.TeamCoreCmds, treeManager, defaultRunner, dockerRunner)
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
				listCmd,
				loginCmd,
				logoutCmd,
				setCmd,
				showCmd,
				updateCmd,
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

	sendMetrics(sessionManager, serverFinder)

	return rootCmd
}

func sendMetrics(sm session.DefaultManager, sf server.Finder) {
	hc := &http.Client{Timeout: 2 * time.Second}
	metricsManager := metrics.NewSender(hc, sf, sm)
	go metricsManager.SendCommand()
}
