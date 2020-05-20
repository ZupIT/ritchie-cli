package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
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

	serverCheckerDesc  = "To use this command on the Team version, you need to inform the server URL first\n Command : rit set server\n"
	sessionCheckerDesc = "To use this command, you need to start a session on Ritchie\n\n"
)

var (
	// Version contains the current version.
	Version = "dev"
	// BuildDate contains a string with the build date.
	BuildDate = "unknown"

	whitelist = []string{
		fmt.Sprint(cmdUse),
		fmt.Sprintf("%s login", cmdUse),
		fmt.Sprintf("%s logout", cmdUse),
		fmt.Sprintf("%s help", cmdUse),
		fmt.Sprintf("%s completion zsh", cmdUse),
		fmt.Sprintf("%s completion bash", cmdUse),
		fmt.Sprintf("%s set server", cmdUse),
	}
)

type rootCmd struct {
	workspaceManager workspace.Checker
	loginManager     security.LoginManager
	repoLoader       formula.Loader
	serverValidator  server.Validator
	sessionValidator session.Validator
	edition          api.Edition
	prompt.InputText
	prompt.InputPassword
}

// NewSingleRootCmd creates the root command for single edition.
func NewSingleRootCmd(wm workspace.Checker,
	l security.LoginManager,
	r formula.Loader,
	sv session.Validator,
	e api.Edition,
	it prompt.InputText,
	ip prompt.InputPassword) *cobra.Command {
	o := &rootCmd{
		wm,
		l,
		r,
		nil,
		sv,
		e,
		it,
		ip,
	}

	return &cobra.Command{
		Use:               cmdUse,
		Version:           o.version(),
		Short:             cmdShortDescription,
		Long:              cmdDescription,
		PersistentPreRunE: o.PreRunFunc(),
		RunE:              runHelp,
		SilenceErrors:     true,
		TraverseChildren:  true,
	}
}

// NewTeamRootCmd creates the root command for team edition.
func NewTeamRootCmd(wm workspace.Checker,
	l security.LoginManager,
	r formula.Loader,
	srv server.Validator,
	sv session.Validator,
	e api.Edition,
	it prompt.InputText,
	ip prompt.InputPassword) *cobra.Command {
	o := &rootCmd{
		wm,
		l,
		r,
		srv,
		sv,
		e,
		it,
		ip,
	}

	return &cobra.Command{
		Use:               cmdUse,
		Version:           o.version(),
		Short:             cmdShortDescription,
		Long:              cmdDescription,
		PersistentPreRunE: o.PreRunFunc(),
		RunE:              runHelp,
		SilenceErrors:     true,
	}
}

func (o *rootCmd) PreRunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if err := o.workspaceManager.Check(); err != nil {
			return err
		}

		if sliceutil.Contains(whitelist, cmd.CommandPath()) {
			return nil
		}

		if err := o.checkServer(cmd.CommandPath()); err != nil {
			return err
		}

		if err := o.checkSession(cmd.CommandPath()); err != nil {
			return err
		}

		return nil
	}
}

func (o *rootCmd) checkServer(commandPath string) error {
	if o.edition == api.Team {
		if err := o.serverValidator.Validate(); err != nil {
			fmt.Print(serverCheckerDesc)
			os.Exit(0)
		}
	}
	return nil
}

func (o *rootCmd) checkSession(commandPath string) error {
	if err := o.sessionValidator.Validate(); err != nil {
		fmt.Print(sessionCheckerDesc)
		secret, err := o.sessionPrompt()
		if err != nil {
			return err
		}

		if err := o.loginManager.Login(secret); err != nil {
			return err
		}

		if err := o.repoLoader.Load(); err != nil {
			return err
		}

		fmt.Println("Session created successfully!")
		os.Exit(0)
	}

	return nil
}

func (o *rootCmd) sessionPrompt() (security.Passcode, error) {
	var passcode string
	var err error

	switch o.edition {
	case api.Single:
		passcode, err = o.Password("Define a passphrase for the session: ")
	case api.Team:
		passcode, err = o.Text("Enter your organization: ", true)
	default:
		err = errors.New("invalid Ritchie build, no edition defined")
	}

	if err != nil {
		return "", err
	}

	return security.Passcode(passcode), nil
}

func (o *rootCmd) version() string {
	return fmt.Sprintf(versionMsg, Version, o.edition, BuildDate, runtime.Version())
}

func runHelp(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
