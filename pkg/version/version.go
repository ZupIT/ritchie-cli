package version

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

var (
	// MsgUpgrade error message to inform user to upgrade rit version
	MsgRitUpgrade = "\nWarning: Rit have a new stable version.\nPlease run: rit upgrade"
	// stableVersionFileCache is the file name to cache stableVersion
	stableVersionFileCache = "stable-version-cache.json"
)

type DefaultVersionResolver struct {
	StableVersionUrl string
	FileUtilService  fileutil.Service
	HttpClient       *http.Client
}

type stableVersionCache struct {
	StableVersion string `json:"stableVersion"`
	ExpiresAt     int64  `json:"expiresAt"`
}

func (r DefaultVersionResolver) StableVersion() (string, error) {
	cachePath := api.RitchieHomeDir() + "/" + stableVersionFileCache
	cacheData, err := r.FileUtilService.ReadFile(cachePath)
	cache := &stableVersionCache{}
	if err == nil {
		err = json.Unmarshal(cacheData, cache)
	}

	if err != nil || cache.ExpiresAt <= time.Now().Unix() {

		request, err := http.NewRequest(http.MethodGet, r.StableVersionUrl, nil)
		if err != nil {
			return "", err
		}
		response, err := r.HttpClient.Do(request)
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
			ExpiresAt:     time.Now().Add(time.Hour * 10).Unix(),
		}

		newCacheJson, err := json.Marshal(newCache)
		if err == nil {
			_ = r.FileUtilService.WriteFilePerm(cachePath, newCacheJson, 0600)
		}

		return stableVersion, nil

	} else {
		return cache.StableVersion, nil
	}
}

func (r DefaultVersionResolver) StableVersionForCmd() (string, error){
	request, err := http.NewRequest(http.MethodGet, r.StableVersionUrl, nil)
	if err != nil {
		return "", err
	}
	response, err := r.HttpClient.Do(request)
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
		ExpiresAt:     time.Now().Add(time.Hour * 10).Unix(),
	}

	newCacheJson, err := json.Marshal(newCache)

	if err == nil {
		cachePath := api.RitchieHomeDir() + "/" + stableVersionFileCache
		_ = r.FileUtilService.WriteFilePerm(cachePath, newCacheJson, 0600)
	}

	return stableVersion, nil
}

func VerifyNewVersion(resolve Resolver, currentVersion string) string {
	stableVersion, err := resolve.StableVersion()
	if err != nil {
		return ""
	}
	if currentVersion != stableVersion {
		return MsgRitUpgrade
	}
	return ""
}
