package cmd

import "github.com/spf13/cobra"

const descListLong = `
This command consists of multiple subcommands to interact with ritchie.

It can be used to list formulas, repositories and other objects.
`

// NewListCmd create a new list instance
func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list SUBCOMMAND",
		Short: "list objects ",
		Long:  descListLong,
	}
}
