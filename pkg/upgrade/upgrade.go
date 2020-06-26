package upgrade

import (
	"fmt"
	"runtime"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type UrlFinder interface {
	Url(edition api.Edition, resolver version.Resolver) string
}

type DefaultUrlFinder struct {}

func (duf DefaultUrlFinder) Url(edition api.Edition, resolver version.Resolver) string {
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
