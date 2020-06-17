package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
)

type UpgradeCmd struct {
	upgradeUrl string
	upgrade.Manager
}

func NewUpgradeCmd(upgradeUrl string, manager upgrade.Manager) *cobra.Command {

	u := UpgradeCmd{
		upgradeUrl: upgradeUrl,
		Manager:    manager,
	}

	return &cobra.Command{
		Use:   "upgrade",
		Short: "Update rit version",
		Long:  `Update rit version to last stable version.`,
		RunE:  u.runFunc(),
	}
}

func (u UpgradeCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		err := u.Run(u.upgradeUrl)
		if err != nil {
			prompt.Error(err.Error())
			return err
		}
		prompt.Success("Rit upgraded with success")
		return nil
	}
}
