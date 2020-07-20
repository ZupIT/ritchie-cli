package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/security/otp"

	"k8s.io/kubectl/pkg/util/templates"

	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator"

	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula/watcher"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credteam"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/env/envcredential"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	fworkspace "github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
	"github.com/ZupIT/ritchie-cli/pkg/metrics"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
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
	inputText := prompt.NewSurveyText()
	inputTextValidator := prompt.NewSurveyTextValidator()
	inputInt := prompt.NewSurveyInt()
	inputBool := prompt.NewSurveyBool()
	inputPassword := prompt.NewSurveyPassword()
	inputList := prompt.NewSurveyList()
	inputURL := prompt.NewSurveyURL()
	inputMultiline := prompt.NewSurveyMultiline()

	// deps
	sessionManager := session.NewManager(ritchieHomeDir)
	workspaceManager := workspace.NewChecker(ritchieHomeDir)
	ctxFinder := rcontext.NewFinder(ritchieHomeDir)
	ctxSetter := rcontext.NewSetter(ritchieHomeDir, ctxFinder)
	ctxRemover := rcontext.NewRemover(ritchieHomeDir, ctxFinder)
	ctxFindSetter := rcontext.NewFindSetter(ritchieHomeDir, ctxFinder, ctxSetter)
	ctxFindRemover := rcontext.NewFindRemover(ritchieHomeDir, ctxFinder, ctxRemover)
	serverFinder := server.NewFinder(ritchieHomeDir)
	serverSetter := server.NewSetter(ritchieHomeDir, makeHttpClientIgnoreSsl())
	serverFindSetter := server.NewFindSetter(serverFinder, serverSetter)

	httpClient := makeHttpClient(serverFinder)
	repoManager := repo.NewTeamRepoManager(ritchieHomeDir, serverFinder, httpClient, sessionManager)
	repoLoader := repo.NewTeamLoader(serverFinder, httpClient, sessionManager, repoManager)
	sessionValidator := sessteam.NewValidator(sessionManager)
	loginManager := secteam.NewLoginManager(
		serverFinder,
		httpClient,
		sessionManager)
	logoutManager := secteam.NewLogoutManager(sessionManager)
	credSetter := credteam.NewSetter(serverFinder, httpClient, sessionManager, ctxFinder)
	credFinder := credteam.NewFinder(serverFinder, httpClient, sessionManager, ctxFinder)
	credSettings := credteam.NewSettings(serverFinder, httpClient, sessionManager, ctxFinder)
	treeManager := tree.NewTreeManager(ritchieHomeDir, repoManager, api.TeamCoreCmds)
	autocompleteGen := autocomplete.NewGenerator(treeManager)
	credResolver := envcredential.NewResolver(credFinder)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver

	inputManager := runner.NewInputManager(envResolvers, inputList, inputText, inputBool, inputPassword)
	formulaSetup := runner.NewDefaultTeamSetup(ritchieHomeDir, httpClient, sessionManager)

	defaultPreRunner := runner.NewDefaultPreRunner(formulaSetup)
	dockerPreRunner := runner.NewDockerPreRunner(formulaSetup)
	postRunner := runner.NewPostRunner()

	defaultRunner := runner.NewDefaultRunner(defaultPreRunner, postRunner, inputManager)
	dockerRunner := runner.NewDockerRunner(dockerPreRunner, postRunner, inputManager, ctxFinder)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	formulaCreator := creator.NewCreator(treeManager, dirManager, fileManager)
	formulaWorkspace := fworkspace.New(ritchieHomeDir, fileManager)
	formulaBuilder := builder.New(ritchieHomeDir, dirManager, fileManager)
	watchManager := watcher.New(formulaBuilder, dirManager)
	createBuilder := formula.NewCreateBuilder(formulaCreator, formulaBuilder)

	upgradeManager := upgrade.DefaultManager{Updater: upgrade.DefaultUpdater{}}
	uhc := makeHttpClient(serverFinder)
	uhc.Timeout = 1 * time.Second
	defaultUpgradeResolver := version.DefaultVersionResolver{
		StableVersionUrl: cmd.StableVersionUrl,
		FileUtilService:  fileutil.DefaultService{},
		HttpClient:       uhc,
	}
	defaultUrlFinder := upgrade.DefaultUrlFinder{}

	otpResolver := otp.NewOtpResolver(httpClient)

	// commands
	rootCmd := cmd.NewTeamRootCmd(workspaceManager, serverFinder, sessionValidator)

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	cleanCmd := cmd.NewCleanCmd()
	initCmd := cmd.NewTeamInitCmd(
		inputText,
		inputPassword,
		inputURL,
		inputBool,
		serverFindSetter,
		loginManager,
		repoLoader,
		otpResolver,
	)
	listCmd := cmd.NewListCmd()
	loginCmd := cmd.NewLoginCmd(inputText, inputPassword, loginManager, repoLoader, serverFinder, otpResolver)
	logoutCmd := cmd.NewLogoutCmd(logoutManager)
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()
	buildCmd := cmd.NewBuildCmd()
	upgradeCmd := cmd.NewUpgradeCmd(api.Team, defaultUpgradeResolver, upgradeManager, defaultUrlFinder)

	// level 2
	setCredentialCmd := cmd.NewTeamSetCredentialCmd(
		credSetter,
		credSettings,
		inputText,
		inputBool,
		inputList,
		inputPassword,
		inputMultiline)
	deleteCtxCmd := cmd.NewDeleteContextCmd(ctxFindRemover, inputBool, inputList)
	setCtxCmd := cmd.NewSetContextCmd(ctxFindSetter, inputText, inputList)
	showCtxCmd := cmd.NewShowContextCmd(ctxFinder)
	addRepoCmd := cmd.NewAddRepoCmd(repoManager, inputText, inputURL, inputInt, inputBool)
	deleteRepoCmd := cmd.NewDeleteRepoCmd(repoManager, inputList, inputBool)
	listRepoCmd := cmd.NewListRepoCmd(repoManager)
	updateRepoCmd := cmd.NewUpdateRepoCmd(repoManager)
	autocompleteZsh := cmd.NewAutocompleteZsh(autocompleteGen)
	autocompleteBash := cmd.NewAutocompleteBash(autocompleteGen)
	autocompleteFish := cmd.NewAutocompleteFish(autocompleteGen)
	autocompletePowerShell := cmd.NewAutocompletePowerShell(autocompleteGen)

	createFormulaCmd := cmd.NewCreateFormulaCmd(userHomeDir, createBuilder, formulaWorkspace, inputText, inputTextValidator, inputList)
	buildFormulaCmd := cmd.NewBuildFormulaCmd(userHomeDir, formulaBuilder, formulaWorkspace, watchManager, dirManager, inputText, inputList)
	cleanFormulasCmd := cmd.NewCleanFormulasCmd()

	autocompleteCmd.AddCommand(autocompleteZsh, autocompleteBash, autocompleteFish, autocompletePowerShell)
	addCmd.AddCommand(addRepoCmd)
	createCmd.AddCommand(createFormulaCmd)
	deleteCmd.AddCommand(deleteRepoCmd, deleteCtxCmd)
	cleanCmd.AddCommand(cleanFormulasCmd)
	listCmd.AddCommand(listRepoCmd)
	setCmd.AddCommand(setCredentialCmd, setCtxCmd)
	showCmd.AddCommand(showCtxCmd)
	updateCmd.AddCommand(updateRepoCmd)
	buildCmd.AddCommand(buildFormulaCmd)

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
				createCmd,
				deleteCmd,
				cleanCmd,
				initCmd,
				listCmd,
				loginCmd,
				logoutCmd,
				setCmd,
				showCmd,
				buildCmd,
				updateCmd,
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

	sendMetrics(sessionManager, serverFinder)

	return rootCmd
}

func sendMetrics(sm session.DefaultManager, sf server.Finder) {
	hc := makeHttpClient(sf)
	hc.Timeout = 2 * time.Second
	metricsManager := metrics.NewSender(hc, sf, sm)
	go metricsManager.SendCommand()
}

func makeHttpClient(finder server.Finder) *http.Client {
	c, err := finder.Find()
	if err != nil {
		fmt.Println(prompt.NewError("error load cli config, try run \"rit init\""))
		os.Exit(1)
	}
	client := &http.Client{}
	client.Transport = &http.Transport{
		DialTLSContext: makeDialer(c.PinningKey, c.PinningAddr, true),
	}
	return client
}

type Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

func makeDialer(pKey, pAddr string, skipCAVerification bool) Dialer {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		c, err := tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: skipCAVerification})
		if err != nil {
			return c, err
		}
		if addr == pAddr {
			connState := c.ConnectionState()
			keyPinValid := false
			for _, peerCert := range connState.PeerCertificates {
				der, err := x509.MarshalPKIXPublicKey(peerCert.PublicKey)
				if err != nil {
					return nil, err
				}
				uEnc := base64.StdEncoding.EncodeToString(der)
				if uEnc == pKey {
					keyPinValid = true
				}
			}
			if !keyPinValid {
				return nil, errors.New("certificate of server not valid")
			}
		}
		return c, nil
	}
}

func makeHttpClientIgnoreSsl() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}
