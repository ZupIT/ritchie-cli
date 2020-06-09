package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type UpgradeCmd struct {
}

func NewUpgradeCmd() *cobra.Command {

	u := UpgradeCmd{}

	return &cobra.Command{
		Use:   "upgrade",
		Short: "Update rit version",
		Long:  `Update rit version to last stable version.`,
		RunE:  u.runFunc(),
	}
}

func (u UpgradeCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hello World")
		return nil
	}
}
