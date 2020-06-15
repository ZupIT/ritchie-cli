package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	versionUtil "github.com/ZupIT/ritchie-cli/pkg/version"
	"github.com/ZupIT/ritchie-cli/pkg/workspace"

	"github.com/spf13/cobra"
)

const (
	versionMsg          = "%s (%s)\n  Build date: %s\n  Built with: %s\n"
	cmdUse              = "rit"
	cmdShortDescription = "rit is a NoOps CLI"
	cmdDescription      = `A CLI that developers can build and operate
your applications without help from the infra staff.
Complete documentation is available at https://github.com/ZupIT/ritchie-cli`
)

var (
	// Version contains the current version.
	Version = "dev"
	// BuildDate contains a string with the build date.
	BuildDate = "unknown"

	// MsgInit error message for init cmd
	MsgInit = "To start using rit, you need to initialize rit first.\nCommand: rit init"
	// MsgSession error message for session not initialized
	MsgSession = "To use this command, you need to start a session first.\nCommand: rit login"

	// Url to get Rit Stable Version
	StableVersionUrl = "https://commons-repo.ritchiecli.io/stable.txt"

	singleWhitelist = []string{
		fmt.Sprint(cmdUse),
		fmt.Sprintf("%s help", cmdUse),
		fmt.Sprintf("%s completion zsh", cmdUse),
		fmt.Sprintf("%s completion bash", cmdUse),
		fmt.Sprintf("%s init", cmdUse),
		fmt.Sprintf("%s upgrade", cmdUse),
	}

	teamWhitelist = []string{
		fmt.Sprint(cmdUse),
		fmt.Sprintf("%s login", cmdUse),
		fmt.Sprintf("%s logout", cmdUse),
		fmt.Sprintf("%s help", cmdUse),
		fmt.Sprintf("%s completion zsh", cmdUse),
		fmt.Sprintf("%s completion bash", cmdUse),
		fmt.Sprintf("%s init", cmdUse),
		fmt.Sprintf("%s upgrade", cmdUse),
	}

	upgradeValidationWhiteList = []string{
		fmt.Sprintf("%s upgrade", cmdUse),
	}
)

type singleRootCmd struct {
	workspaceChecker workspace.Checker
	sessionValidator session.Validator
}

type teamRootCmd struct {
	workspaceChecker workspace.Checker
	serverFinder     server.Finder
	sessionValidator session.Validator
}

// NewSingleRootCmd creates the root command for single edition.
func NewSingleRootCmd(wc workspace.Checker, sv session.Validator) *cobra.Command {
	o := &singleRootCmd{
		workspaceChecker: wc,
		sessionValidator: sv,
	}

	cmd := &cobra.Command{
		Use:                cmdUse,
		Version:            version(api.Single),
		Short:              cmdShortDescription,
		Long:               cmdDescription,
		PersistentPreRunE:  o.PreRunFunc(),
		PersistentPostRunE: o.PostRunFunc(),
		RunE:               runHelp,
		SilenceErrors:      true,
		TraverseChildren:   true,
	}
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	return cmd
}

// NewTeamRootCmd creates the root command for team edition.
func NewTeamRootCmd(wc workspace.Checker,
	sf server.Finder,
	sv session.Validator) *cobra.Command {
	o := &teamRootCmd{
		workspaceChecker: wc,
		serverFinder:     sf,
		sessionValidator: sv,
	}

	cmd := &cobra.Command{
		Use:                cmdUse,
		Version:            version(api.Team),
		Short:              cmdShortDescription,
		Long:               cmdDescription,
		PersistentPreRunE:  o.PreRunFunc(),
		PersistentPostRunE: o.PostRunFunc(),
		RunE:               runHelp,
		SilenceErrors:      true,
	}
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	return cmd
}

func (o *singleRootCmd) PreRunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if err := o.workspaceChecker.Check(); err != nil {
			return err
		}

		if isWhitelist(singleWhitelist, cmd) {
			return nil
		}

		if err := o.sessionValidator.Validate(); err != nil {
			fmt.Println(MsgInit)
			os.Exit(0)
		}

		return nil
	}
}

func (o *teamRootCmd) PreRunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if err := o.workspaceChecker.Check(); err != nil {
			return err
		}

		if isWhitelist(teamWhitelist, cmd) {
			return nil
		}

		cfg, err := o.serverFinder.Find()
		if err != nil {
			return err
		} else if cfg.URL == "" {
			fmt.Println(MsgInit)
			os.Exit(0)
		}

		if err := o.sessionValidator.Validate(); err != nil {
			fmt.Println(MsgSession)
			os.Exit(0)
		}

		return nil
	}
}

func (o *singleRootCmd) PostRunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		verifyNewVersion(cmd)
		return nil
	}
}

func (o *teamRootCmd) PostRunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		verifyNewVersion(cmd)
		return nil
	}
}

func verifyNewVersion(cmd *cobra.Command) {
	if !isWhitelist(upgradeValidationWhiteList, cmd) {
		resolver := versionUtil.DefaultVersionResolver{
			CurrentVersion:   Version,
			StableVersionUrl: StableVersionUrl,
			FileUtilService:  fileutil.DefaultFileUtilService{},
		}
		versionUtil.VerifyNewVersion(resolver, os.Stdout)
	}
}

func isWhitelist(whitelist []string, cmd *cobra.Command) bool {
	return sliceutil.Contains(whitelist, cmd.CommandPath())
}

func version(edition api.Edition) string {
	return fmt.Sprintf(versionMsg, Version, edition, BuildDate, runtime.Version())
}

func runHelp(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
