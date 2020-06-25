package upgrade

import (
	"fmt"
	"runtime"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

func Url(edition api.Edition, resolver version.Resolver) string {
	stableVersion, err := resolver.StableVersion(true)
	if err != nil {
		return ""
	}

	upgradeUrl := fmt.Sprintf(upgradeUrlFormat, stableVersion, runtime.GOOS, edition)

	if runtime.GOOS == "windows" {
		upgradeUrl += ".exe"
	}

	return upgradeUrl
}
