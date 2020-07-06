package cmd

import "github.com/spf13/cobra"

var descUpdateLong = `
This command consists of multiple subcommands to interact with ritchie.

It can be used to update formulas, repositories and etc.
`

// NewUpdateCmd create a new update instance
func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update SUBCOMMAND",
		Short: "Update objects (repo and others)",
		Long:  descUpdateLong,
	}
}
