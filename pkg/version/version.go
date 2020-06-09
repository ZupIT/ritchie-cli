package version

import (
	"fmt"
	"io"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	// MsgUpgrade error message to inform user to upgrade rit version
	MsgRitUpgrade = "\nWarning: Rit have a new version.\nPlease run: rit upgrade"
)

type Resolver interface {
	getCurrentVersion() string
	getStableVersion() string
}

type DefaultVersionResolver struct{
	CurrentVersion string
}

func (r DefaultVersionResolver) getCurrentVersion() string {
	return r.CurrentVersion
}

func (r DefaultVersionResolver) getStableVersion() string {
	return "bla"
}

func VerifyNewVersion(resolve Resolver, writer io.Writer) {
	stableVersion := resolve.getStableVersion()
	currentVersion := resolve.getCurrentVersion()
	if currentVersion != stableVersion {
		_, err := fmt.Fprintf(writer, prompt.Warning, MsgRitUpgrade)
		if err != nil {
			panic("Fail to Write MsgRitUpgrade")
		}
	}
}
