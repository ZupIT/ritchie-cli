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
	MsgRitUpgrade = "\nWarning: Rit has a new stable version.\nPlease run: rit upgrade"
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

func requestStableVersion(stableVersionUrl string, httpClient *http.Client) (string, error) {
	request, err := http.NewRequest(http.MethodGet, stableVersionUrl, nil)
	if err != nil {
		return "", err
	}
	response, err := httpClient.Do(request)
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

func updateCache(stableVersion string, cachePath string, fileUtilService fileutil.Service) {
	newCache := stableVersionCache{
		StableVersion: stableVersion,
		ExpiresAt:     time.Now().Add(time.Hour * 10).Unix(),
	}

	newCacheJson, err := json.Marshal(newCache)
	if err == nil {
		_ = fileUtilService.WriteFilePerm(cachePath, newCacheJson, 0600)
	}
}

func (r DefaultVersionResolver) StableVersion(fromCache bool) (string, error) {
	cachePath := api.RitchieHomeDir() + "/" + stableVersionFileCache

	if !fromCache {
		stableVersion, err := requestStableVersion(r.StableVersionUrl, r.HttpClient)
		if err != nil {
			return stableVersion, err
		}
		updateCache(stableVersion, cachePath, r.FileUtilService)
		return stableVersion, nil
	}

	cacheData, err := r.FileUtilService.ReadFile(cachePath)
	cache := &stableVersionCache{}
	if err == nil {
		err = json.Unmarshal(cacheData, cache)
	}

	if err != nil || cache.ExpiresAt <= time.Now().Unix() {

		stableVersion, err := requestStableVersion(r.StableVersionUrl, r.HttpClient)
		if err != nil {
			return stableVersion, err
		}
		updateCache(stableVersion, cachePath, r.FileUtilService)
		return stableVersion, nil

	} else {
		return cache.StableVersion, nil
	}
}

func VerifyNewVersion(resolve Resolver, currentVersion string) string {
	stableVersion, err := resolve.StableVersion(true)
	if err != nil {
		return ""
	}
	if currentVersion != stableVersion {
		return MsgRitUpgrade
	}
	return ""
}
