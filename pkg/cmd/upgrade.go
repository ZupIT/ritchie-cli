package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type UpgradeCmd struct {
	edition api.Edition
	upgrade.Manager
	resolver version.Resolver
}

func NewUpgradeCmd(e api.Edition, r version.Resolver, manager upgrade.Manager) *cobra.Command {

	u := UpgradeCmd{
		edition: e,
		Manager:  manager,
		resolver: r,
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
		upgradeUrl := upgrade.Url(u.edition, u.resolver)
		err := u.Run(upgradeUrl)
		if err != nil {
			return fmt.Errorf(prompt.Red, err.Error()+"\n")
		}
		prompt.Success("Rit upgraded with success")
		return nil
	}
}
