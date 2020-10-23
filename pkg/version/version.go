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
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	// MsgUpgrade error message to inform user to upgrade rit version
	MsgRitUpgrade = "\nWarning: Rit has a new stable version.\nPlease run: rit upgrade"
	// stableVersionFileCache is the file name to cache stableVersion
	stableVersionFileCache = "stable-version-cache.json"

	StableVersionUrl = "https://commons-repo.ritchiecli.io/stable.txt"
)

type Manager struct {
	stableUrl string
	file      stream.FileWriteReadExister
}

var _ Resolver = Manager{}

func NewManager(
	stableVersionUrl string,
	file stream.FileWriteReadExister,
) Manager {
	return Manager{
		stableUrl: stableVersionUrl,
		file:      file,
	}
}

type stableVersionCache struct {
	Stable    string `json:"stableVersion"`
	ExpiresAt int64  `json:"expiresAt"`
}

func (m Manager) StableVersion() (string, error) {
	cachePath := filepath.Join(
		api.RitchieHomeDir(),
		stableVersionFileCache)

	shouldRequestAndSave := false
	cache := &stableVersionCache{}

	if !m.file.Exists(cachePath) {
		shouldRequestAndSave = true
	} else {
		cacheData, err := m.file.Read(cachePath)
		if err != nil {
			return "", err
		}

		if err = json.Unmarshal(cacheData, cache); err != nil {
			return "", err
		}

		if cache.ExpiresAt <= time.Now().Unix() {
			shouldRequestAndSave = true
		}
	}

	if shouldRequestAndSave {
		stableVersion, err := requestStableVersion(m.stableUrl)
		if err != nil {
			return "", err
		}
		if err := saveCache(stableVersion, cachePath, m.file); err != nil {
			return "", err
		}

		return stableVersion, nil
	}

	return cache.Stable, nil
}

func (m Manager) VerifyNewVersion(current, installed string) string {
	if current != installed && current != "" {
		return MsgRitUpgrade
	}
	return ""
}

func (m Manager) UpdateCache() error {
	cachePath := filepath.Join(api.RitchieHomeDir(), stableVersionFileCache)

	stableVersion, err := requestStableVersion(m.stableUrl)
	if err != nil {
		return err
	}

	if err = saveCache(stableVersion, cachePath, m.file); err != nil {
		return err
	}
	return nil
}

func requestStableVersion(stableVersionURL string) (string, error) {
	response, err := http.Get(stableVersionURL)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status is not %d", http.StatusOK)
	}

	stableVersionBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	stableVersion := string(stableVersionBytes)
	stableVersion = strings.ReplaceAll(stableVersion, "\n", "")
	return stableVersion, nil
}

func saveCache(
	stableVersion string,
	cachePath string,
	file stream.FileWriteReadExister,
) error {
	newCache := stableVersionCache{
		Stable:    stableVersion,
		ExpiresAt: time.Now().Add(time.Hour * 10).Unix(),
	}

	newCacheJSON, err := json.Marshal(newCache)
	if err != nil {
		return err
	}

	if err := file.Write(cachePath, newCacheJSON); err != nil {
		return err
	}
	return nil
}
