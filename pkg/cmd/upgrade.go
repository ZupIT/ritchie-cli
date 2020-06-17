package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version/version_util"
)

const (
	upgradeUrlFormat = "https://commons-repo.ritchiecli.io/%s/%s/%s/rit"
)

type UpgradeCmd struct {
	upgradeUrl string
	upgrade    upgrade.Upgrade
}

func UpgradeUrl(edition api.Edition, resolver version_util.Resolver) string {
	stableVersion, err := resolver.GetStableVersion()
	if err != nil {
		return ""
	}

	upgradeUrl := fmt.Sprintf(upgradeUrlFormat, stableVersion, runtime.GOOS, edition)

	if runtime.GOOS == "windows" {
		upgradeUrl += ".exe"
	}

	return upgradeUrl
}

func NewUpgradeCmd(upgradeUrl string, upgrade upgrade.Upgrade) *cobra.Command {

	u := UpgradeCmd{
		upgradeUrl: upgradeUrl,
		upgrade:    upgrade,
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

		if u.upgradeUrl == "" {
			prompt.Error("Fail to resolve upgrade url.")
			return errors.New("fail to resolve upgrade url")
		}

		resp, err := http.Get(u.upgradeUrl)
		if err != nil {
			prompt.Error(fmt.Sprintf("Fail to download new version.\nErr:%s\n", err))
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			prompt.Error(fmt.Sprintf("Fail to download new version.\nStatus:%d\n", resp.StatusCode))
			return fmt.Errorf("upgradeUrl return status:%d", resp.StatusCode)
		}

		err = u.upgrade.Apply(resp.Body, update.Options{})
		if err != nil {
			prompt.Error(fmt.Sprintf("Fail to upgrade new version.\nErr:%s\n", err))
			return err
		}
		prompt.Success("Rit upgraded with success\n")
		return nil
	}
}
