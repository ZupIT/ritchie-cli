package cmd

import "github.com/spf13/cobra"

// CommandRunnerFunc represents that runner func for commands
type CommandRunnerFunc func(cmd *cobra.Command, args []string) error

// RunFuncE delegates to stdinFunc if --stdin flag is passed otherwise delegates to promptFunc
func RunFuncE(stdinFunc, promptFunc CommandRunnerFunc) CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		stdin, err := cmd.Flags().GetBool("stdin")
		if err != nil {
			return err
		}

		if stdin {
			return stdinFunc(cmd, args)
		}
		return promptFunc(cmd, args)
	}
}
