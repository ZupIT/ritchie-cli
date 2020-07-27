package cmd

import "github.com/spf13/cobra"

func NewBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "build SUB_COMMAND",
		Short:   "Build formulas",
		Long:    "Build formula with latest changes.",
		Example: "rit build formula",
	}
}
