package version

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	// MsgUpgrade error message to inform user to upgrade rit version
	MsgRitUpgrade = "\nWarning: Rit have a new stable version.\nPlease run: rit upgrade\n"
	// stableVersionFileCache is the file name to cache stableVersion
	stableVersionFileCache = "stable-version-cache.txt"
)

type Resolver interface {
	GetCurrentVersion() (string, error)
	GetStableVersion() (string, error)
}

type DefaultVersionResolver struct {
	CurrentVersion   string
	StableVersionUrl string
	FileUtilService fileutil.FileUtilService
}

type stableVersionCache struct {
	StableVersion string `json:"stableVersion"`
	ExpiresAt     time.Time `json:"expiresAt"`
}

func (r DefaultVersionResolver) GetCurrentVersion() (string, error) {
	return r.CurrentVersion, nil
}

func (r DefaultVersionResolver) GetStableVersion() (string, error) {

	cachePath := api.RitchieHomeDir() + "/" + stableVersionFileCache
	cacheData, err := r.FileUtilService.ReadFile(cachePath)
	cache := &stableVersionCache{}
	if err == nil {
		err = json.Unmarshal(cacheData, cache)
	}

	if err != nil || cache.ExpiresAt.Before(time.Now()) {

		api.RitchieHomeDir()

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

		newCache := stableVersionCache{
			StableVersion: stableVersion,
			ExpiresAt:     time.Now().Add(time.Hour * 10),
		}

		newCacheJson, err := json.Marshal(newCache)
		if err == nil {
			r.FileUtilService.WriteFilePerm(cachePath, newCacheJson, 0600)
		}

		return stableVersion, nil

	} else {
		return cache.StableVersion, nil
	}
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
