package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type UpgradeCmd struct {
	upgrade.Manager
	resolver version.Resolver
	upgrade.UrlFinder
}

func NewUpgradeCmd(r version.Resolver, m upgrade.Manager, uf upgrade.UrlFinder) *cobra.Command {

	u := UpgradeCmd{
		Manager:  m,
		resolver: r,
		UrlFinder: uf,
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
		err := u.resolver.UpdateCache()
		if err != nil {
			return prompt.NewError(err.Error()+"\n")
		}
		upgradeUrl := u.Url(u.resolver)
		err = u.Run(upgradeUrl)
		if err != nil {
			return prompt.NewError(err.Error()+"\n")
		}
		prompt.Success("Rit upgraded with success")
		return nil
	}
}
