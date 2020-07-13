package cmd

import "github.com/spf13/cobra"

// NewShowCmd creates new cmd instance of show command
func NewShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show SUB_COMMAND",
		Short: "Show contexts",
		Long:  `Show contexts`,
	}
}
