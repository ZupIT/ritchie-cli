package cmd

import (
	"github.com/spf13/cobra"
)

// NewAddCmd create a new add instance
func NewAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "add SUBCOMMAND",
		Short:   "Add repositories ",
		Long:    "Add a new formula repository to ritchie-cli.",
		Example: "rit add repo",
	}
}
