package cmd

import (
	"errors"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/workspace"
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

	whitelist = []string{
		fmt.Sprintf("%s login", cmdUse),
		fmt.Sprintf("%s logout", cmdUse),
		fmt.Sprintf("%s completion zsh", cmdUse),
		fmt.Sprintf("%s completion bash", cmdUse),
	}
)

type rootCmd struct {
	workspaceManager workspace.Checker
	loginManager     security.LoginManager
	sessionValidator session.Validator
}

// NewRootCmd creates the root for all ritchie commands.
func NewRootCmd(wm workspace.Checker, l security.LoginManager, sv session.Validator) *cobra.Command {
	o := &rootCmd{wm, l, sv}

	return &cobra.Command{
		Use:               cmdUse,
		Version:           version(),
		Short:             cmdShortDescription,
		Long:              cmdDescription,
		PersistentPreRunE: o.PreRunFunc(),
		SilenceErrors:     true,
	}
}

func version() string {
	return fmt.Sprintf(versionMsg, Version, env.Edition, BuildDate, runtime.Version())
}

func (o *rootCmd) PreRunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if err := o.workspaceManager.Check(); err != nil {
			return err
		}

		if err := o.checkSession(cmd.CommandPath()); err != nil {
			return err
		}

		return nil
	}
}

func (o *rootCmd) checkSession(commandPath string) error {
	if sliceutil.Contains(whitelist, commandPath) {
		return nil
	}

	err := o.sessionValidator.Validate()
	if err != nil {
		fmt.Print("To use this command, you need to start a session on Ritchie\n\n")
		secret, err := sessionPrompt()
		if err != nil {
			return err
		}

		if err := o.loginManager.Login(secret); err != nil {
			return err
		}

		fmt.Println("Session created successfully!")
		os.Exit(0)
	}

	return nil
}

func sessionPrompt() (security.Passcode, error) {
	var passcode string
	var err error

	switch env.Edition {
	case env.Single:
		passcode, err = prompt.Password("Define a passphrase for the session: ")
	case env.Team:
		passcode, err = prompt.String("Enter your organization: ", true)
	default:
		err = errors.New("invalid Ritchie build, no edition defined")
	}

	if err != nil {
		return "", err
	}

	return security.Passcode(passcode), nil
}
