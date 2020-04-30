package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credsingle"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/env/envcredential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/security/secsingle"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/session/sesssingle"
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

	//prompt
	inputText := prompt.NewInputText()
	inputInt := prompt.NewInputInt()
	inputBool := prompt.NewInputBool()
	inputPassword := prompt.NewInputPassword()
	inputList := prompt.NewInputList()
	inputURL := prompt.NewInputURL()

	// deps
	sessionManager := session.NewManager(ritchieHomeDir)
	workspaceManager := workspace.NewChecker(ritchieHomeDir)
	ctxFinder := rcontext.NewFinder(ritchieHomeDir)
	ctxSetter := rcontext.NewSetter(ritchieHomeDir, ctxFinder)
	ctxRemover := rcontext.NewRemover(ritchieHomeDir, ctxFinder)
	ctxFindSetter := rcontext.NewFindSetter(ritchieHomeDir, ctxFinder, ctxSetter)
	ctxFindRemover := rcontext.NewFindRemover(ritchieHomeDir, ctxFinder, ctxRemover)
	repoManager := formula.NewSingleRepoManager(ritchieHomeDir, http.DefaultClient, sessionManager)
	sessionValidator := sesssingle.NewValidator(sessionManager)
	loginManager := secsingle.NewLoginManager(sessionManager)
	credSetter := credsingle.NewSetter(ritchieHomeDir, ctxFinder, sessionManager)
	credFinder := credsingle.NewFinder(ritchieHomeDir, ctxFinder, sessionManager)
	treeManager := formula.NewTreeManager(ritchieHomeDir, repoManager, api.SingleCoreCmds)
	autocompleteGen := autocomplete.NewGenerator(treeManager)
	credResolver := envcredential.NewResolver(credFinder)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver
	formulaRunner := formula.NewRunner(
		ritchieHomeDir,
		envResolvers,
		http.DefaultClient,
		treeManager,
		inputList,
		inputText,
		inputBool)
	formulaCreator := formula.NewCreator(userHomeDir, treeManager)

	//commands
	rootCmd := cmd.NewSingleRootCmd(
		workspaceManager,
		loginManager,
		repoManager,
		sessionValidator,
		api.Single,
		inputText,
		inputPassword)

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	cleanCmd := cmd.NewCleanCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	listCmd := cmd.NewListCmd()
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()

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
	addRepoCmd := cmd.NewAddRepoCmd(repoManager, inputText, inputURL, inputInt)
	cleanRepoCmd := cmd.NewCleanRepoCmd(repoManager, inputText)
	deleteRepoCmd := cmd.NewDeleteRepoCmd(repoManager, inputText)
	listRepoCmd := cmd.NewListRepoCmd(repoManager)
	updateRepoCmd := cmd.NewUpdateRepoCmd(repoManager)
	autocompleteZsh := cmd.NewAutocompleteZsh(autocompleteGen)
	autocompleteBash := cmd.NewAutocompleteBash(autocompleteGen)
	createFormulaCmd := cmd.NewCreateFormulaCmd(formulaCreator, inputText)

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash)
	addCmd.AddCommand(addRepoCmd)
	cleanCmd.AddCommand(cleanRepoCmd)
	createCmd.AddCommand(createFormulaCmd)
	deleteCmd.AddCommand(deleteRepoCmd, deleteCtxCmd)
	listCmd.AddCommand(listRepoCmd)
	setCmd.AddCommand(setCredentialCmd, setCtxCmd)
	showCmd.AddCommand(showCtxCmd)
	updateCmd.AddCommand(updateRepoCmd)

	rootCmd.AddCommand(
		addCmd,
		autocompleteCmd,
		cleanCmd,
		createCmd,
		deleteCmd,
		listCmd,
		setCmd,
		showCmd,
		updateCmd)

	formulaCmd := cmd.NewFormulaCommand(api.SingleCoreCmds, treeManager, formulaRunner)
	if err := formulaCmd.Add(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}
