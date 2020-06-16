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
	"github.com/ZupIT/ritchie-cli/pkg/version/versionutil"
)

const (
	upgradeUrlFormat = "https://commons-repo.ritchiecli.io/%s/%s/%s/rit"
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

func UpgradeUrl(edition api.Edition, resolver versionutil.Resolver) string {
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
			fmt.Printf(prompt.Error, fmt.Sprintf("Fail to download new version.\nErr:%s\n", err))
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Printf(prompt.Error, fmt.Sprintf("Fail to download new version.\nStatus:%d\n", resp.StatusCode))
			return fmt.Errorf("upgradeUrl return status:%d", resp.StatusCode)
		}

		err = u.upgradeUtil.Apply(resp.Body, update.Options{})
		if err != nil {
			fmt.Printf(prompt.Error, fmt.Sprintf("Fail to upgrade new version.\nErr:%s\n", err))
			return err
		}
		fmt.Printf(prompt.Success, "Rit upgrated with success\n")
		return nil
	}
}
