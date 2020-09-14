package version

import (
	"encoding/json"
	"errors"
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
	stableVersionFileCache      = "stable-version-cache.json"
	errUnexpectedResponseMethod = errors.New("Unexpected response method")
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

func (r DefaultVersionResolver) UpdateCache() error {
	cachePath := api.RitchieHomeDir() + "/" + stableVersionFileCache

	stableVersion, err := requestStableVersion(r.StableVersionUrl, r.HttpClient)
	if err != nil {
		return err
	}

	err = saveCache(stableVersion, cachePath, r.FileUtilService)
	return err
}

func (r DefaultVersionResolver) StableVersion() (string, error) {
	cachePath := api.RitchieHomeDir() + "/" + stableVersionFileCache
	cacheData, err := r.FileUtilService.ReadFile(cachePath)
	cache := &stableVersionCache{}

	if err == nil {
		err = json.Unmarshal(cacheData, cache)
	}

	if err != nil || cache.ExpiresAt <= time.Now().Unix() {
		stableVersion, err := requestStableVersion(r.StableVersionUrl, r.HttpClient)
		if err != nil {
			return "", err
		}
		err = saveCache(stableVersion, cachePath, r.FileUtilService)
		if err != nil {
			return "", err
		}
		return stableVersion, nil
	} else {
		return cache.StableVersion, nil
	}
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

	if response.StatusCode != http.StatusOK {
		return "", errUnexpectedResponseMethod
	}

	stableVersionBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	stableVersion := string(stableVersionBytes)
	stableVersion = strings.ReplaceAll(stableVersion, "\n", "")
	return stableVersion, nil
}

func saveCache(stableVersion string, cachePath string, fileUtilService fileutil.Service) error {
	newCache := stableVersionCache{
		StableVersion: stableVersion,
		ExpiresAt:     time.Now().Add(time.Hour * 10).Unix(),
	}

	newCacheJson, err := json.Marshal(newCache)
	if err != nil {
		return err
	}
	err = fileUtilService.WriteFilePerm(cachePath, newCacheJson, 0600)
	return err
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
