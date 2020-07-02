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

	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"

	"k8s.io/kubectl/pkg/util/templates"

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
	serverFinder := server.NewFinder(ritchieHomeDir)
	serverSetter := server.NewSetter(ritchieHomeDir, makeHttpClientIgnoreSsl())
	serverFindSetter := server.NewFindSetter(serverFinder, serverSetter)

	httpClient := makeHttpClient(serverFinder)
	repoManager := formula.NewTeamRepoManager(ritchieHomeDir, serverFinder, httpClient, sessionManager)
	repoLoader := formula.NewTeamLoader(serverFinder, httpClient, sessionManager, repoManager)
	sessionValidator := sessteam.NewValidator(sessionManager)
	loginManager := secteam.NewLoginManager(
		serverFinder,
		httpClient,
		sessionManager)
	logoutManager := secteam.NewLogoutManager(sessionManager)
	credSetter := credteam.NewSetter(serverFinder, httpClient, sessionManager, ctxFinder)
	credFinder := credteam.NewFinder(serverFinder, httpClient, sessionManager, ctxFinder)
	credSettings := credteam.NewSettings(serverFinder, httpClient, sessionManager, ctxFinder)

	treeManager := formula.NewTreeManager(ritchieHomeDir, repoManager, api.TeamCoreCmds)
	autocompleteGen := autocomplete.NewGenerator(treeManager)
	credResolver := envcredential.NewResolver(credFinder)
	envResolvers := make(env.Resolvers)
	envResolvers[env.Credential] = credResolver

	inputManager := formula.NewInputManager(envResolvers, inputList, inputText, inputBool, inputPassword)
	formulaSetup := formula.NewDefaultTeamSetup(ritchieHomeDir, httpClient, sessionManager)

	defaultPreRunner := formula.NewDefaultPreRunner(formulaSetup)
	dockerPreRunner := formula.NewDockerPreRunner(formulaSetup)
	postRunner := formula.NewPostRunner()

	defaultRunner := formula.NewDefaultRunner(defaultPreRunner, postRunner, inputManager)
	dockerRunner := formula.NewDockerRunner(dockerPreRunner, postRunner, inputManager)

	formulaCreator := formula.NewCreator(userHomeDir, treeManager)

	upgradeManager := upgrade.DefaultManager{Updater: upgrade.DefaultUpdater{}}
	uhc := makeHttpClient(serverFinder)
	uhc.Timeout =  1 * time.Second
	defaultUpgradeResolver := version.DefaultVersionResolver{
		StableVersionUrl: cmd.StableVersionUrl,
		FileUtilService:  fileutil.DefaultService{},
		HttpClient:       uhc,
	}
	upgradeUrl := upgrade.UpgradeUrl(api.Team, defaultUpgradeResolver)

	// commands
	rootCmd := cmd.NewTeamRootCmd(workspaceManager, serverFinder, sessionValidator)

	// level 1
	autocompleteCmd := cmd.NewAutocompleteCmd()
	addCmd := cmd.NewAddCmd()
	cleanCmd := cmd.NewCleanCmd()
	createCmd := cmd.NewCreateCmd()
	deleteCmd := cmd.NewDeleteCmd()
	initCmd := cmd.NewTeamInitCmd(inputText, inputPassword, inputURL, inputBool, serverFindSetter, loginManager, repoLoader)
	listCmd := cmd.NewListCmd()
	loginCmd := cmd.NewLoginCmd(inputText, inputPassword, loginManager, repoLoader, serverFinder)
	logoutCmd := cmd.NewLogoutCmd(logoutManager)
	setCmd := cmd.NewSetCmd()
	showCmd := cmd.NewShowCmd()
	updateCmd := cmd.NewUpdateCmd()
	buildCmd := cmd.NewBuildCmd()
	upgradeCmd := cmd.NewUpgradeCmd(upgradeUrl, upgradeManager)

	// level 2
	setCredentialCmd := cmd.NewTeamSetCredentialCmd(
		credSetter,
		credSettings,
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
		fmt.Println(fmt.Errorf(prompt.Red, "error load cli config, try run \"rit init\""))
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
		if addr == pAddr {
			if err != nil {
				return c, err
			}
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
