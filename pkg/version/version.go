/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package version

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

const (
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

func (r DefaultVersionResolver) UpdateCache() error {
	cachePath := filepath.Join(api.RitchieHomeDir(), stableVersionFileCache)
	stableVersion, err := requestStableVersion(r.StableVersionUrl, r.HttpClient)
	if err != nil {
		return err
	}

	err = saveCache(stableVersion, cachePath, r.FileUtilService)
	return err
}

func (r DefaultVersionResolver) StableVersion() (string, error) {
	cachePath := filepath.Join(api.RitchieHomeDir(), stableVersionFileCache)
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
	if currentVersion != stableVersion && currentVersion != "" {
		return MsgRitUpgrade
	}
	return ""
}
