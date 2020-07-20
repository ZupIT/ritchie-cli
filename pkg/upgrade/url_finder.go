package upgrade

import (
	"fmt"
	"runtime"

	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type UrlFinder interface {
	Url(resolver version.Resolver) string
}

type DefaultUrlFinder struct{}

func (duf DefaultUrlFinder) Url(resolver version.Resolver) string {
	stableVersion, err := resolver.StableVersion()
	if err != nil {
		return ""
	}

	upgradeUrl := fmt.Sprintf(upgradeUrlFormat, stableVersion, runtime.GOOS)

	if runtime.GOOS == "windows" {
		upgradeUrl += ".exe"
	}

	return upgradeUrl
}
