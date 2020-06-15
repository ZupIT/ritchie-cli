package version

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	// MsgUpgrade error message to inform user to upgrade rit version
	MsgRitUpgrade = "\nWarning: Rit have a new stable version.\nPlease run: rit upgrade\n"
)

type Resolver interface {
	GetCurrentVersion() (string, error)
	GetStableVersion() (string, error)
}

type DefaultVersionResolver struct {
	CurrentVersion   string
	StableVersionUrl string
}

func (r DefaultVersionResolver) GetCurrentVersion() (string, error) {
	return r.CurrentVersion, nil
}

func (r DefaultVersionResolver) GetStableVersion() (string, error) {

	response, err := http.Get(r.StableVersionUrl)
	if err != nil {
		return "", err
	}
	stableVersionBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	stableVersion := string(stableVersionBytes)
	stableVersion = strings.ReplaceAll(stableVersion, "\n", "")

	return stableVersion, nil
}

func VerifyNewVersion(resolve Resolver, writer io.Writer) {
	stableVersion, err := resolve.GetStableVersion()
	if err != nil {
		return
	}
	currentVersion, err := resolve.GetCurrentVersion()
	if err != nil {
		return
	}
	if currentVersion != stableVersion {
		fmt.Fprintf(writer, prompt.Warning, MsgRitUpgrade)
	}
}
