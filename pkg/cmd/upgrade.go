package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	versionUtil "github.com/ZupIT/ritchie-cli/pkg/version"
)

const (
	UpgradeUrlFormat = "https://commons-repo.ritchiecli.io/%s/%s/%s/rit"
)

type UpgradeUtil interface {
	Apply(reader io.Reader, opts update.Options) error
}

type DefaultUpgradeUtil struct{}

func (u DefaultUpgradeUtil) Apply(reader io.Reader, opts update.Options) error {
	return update.Apply(reader, opts)
}

type UpgradeCmd struct {
	upgradeUrl  string
	upgradeUtil UpgradeUtil
}

func GetUpgradeUrl(edition api.Edition, resolver versionUtil.Resolver) string {
	stableVersion, err := resolver.GetStableVersion()
	if err != nil {
		return ""
	}
	return fmt.Sprintf(UpgradeUrlFormat, stableVersion, runtime.GOOS, edition)
}

func NewUpgradeCmd(upgradeUrl string, upgradeUtil UpgradeUtil) *cobra.Command {

	u := UpgradeCmd{
		upgradeUrl:  upgradeUrl,
		upgradeUtil: upgradeUtil,
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
			fmt.Printf(prompt.Error, "Fail to resolve upgrade url.\n")
			return errors.New("fail to resolve upgrade url")
		}

		resp, err := http.Get(u.upgradeUrl)
		if err != nil {
			fmt.Printf(prompt.Error, "Fail to download new version.\n")
			return err
		}
		defer resp.Body.Close()

		err = u.upgradeUtil.Apply(resp.Body, update.Options{})
		if err != nil {
			fmt.Printf(prompt.Error, "Fail to upgrade new version.\n")
			return err
		}
		fmt.Printf(prompt.Success, "Rit upgrated with success\n")
		return nil
	}
}
