package cmd

import "github.com/spf13/cobra"

type CommandRunnerFunc func(cmd *cobra.Command, args []string) error
