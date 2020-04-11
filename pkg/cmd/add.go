package cmd

import (
	"github.com/spf13/cobra"
)

var descAddLong = `
This command consists of multiple subcommands to interact with ritchie.

It can be used to add formulas, repositories and other objects..
`

// NewAddCmd create a new add instance
func NewAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add SUBCOMMAND",
		Short: "add objects",
		Long:  descAddLong,
	}
}
