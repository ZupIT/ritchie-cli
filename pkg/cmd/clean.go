package cmd

import (
"github.com/spf13/cobra"
)

var descCleanLong = `
This command consists of multiple subcommands to interact with ritchie.

It can be used to clean formulas from your current ritchie build
`

// NewCleanCmd create a new clean instance
func NewCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean SUBCOMMAND",
		Short: "Clean properties",
		Long:  descCleanLong,
	}
}

