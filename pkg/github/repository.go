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

package github

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
)

type RepoManager struct {
	client *http.Client
}

func NewRepoManager(client *http.Client) RepoManager {
	return RepoManager{client: client}
}

func (re RepoManager) Zipball(info RepoInfo, version string) (io.ReadCloser, error) {
	zipUrl := info.ZipUrl(version)
	req, err := http.NewRequest(http.MethodGet, zipUrl, nil)
	if err != nil {
		return nil, err
	}

	if info.Token() != "" {
		authToken := info.TokenHeader()
		req.Header.Add(headers.Authorization, authToken)
	}

	req.Header.Add(headers.Accept, "application/vnd.github.v3+json")
	resp, err := re.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (re RepoManager) Tags(info RepoInfo) (Tags, error) {
	apiUrl := info.TagsUrl()
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return Tags{}, err
	}

	if info.Token() != "" {
		authToken := info.TokenHeader()
		req.Header.Add(headers.Authorization, authToken)
	}

	req.Header.Add(headers.Accept, "application/vnd.github.v3+json")
	res, err := re.client.Do(req)
	if err != nil {
		return Tags{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return Tags{}, err
		}
		return Tags{}, errors.New(string(b))
	}

	var tags Tags
	if err := json.NewDecoder(res.Body).Decode(&tags); err != nil {
		return Tags{}, err
	}

	return tags, nil
}

func (re RepoManager) LatestTag(info RepoInfo) (Tag, error) {
	apiUrl := info.LatestTagUrl()
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return Tag{}, err
	}

	if info.Token() != "" {
		authToken := info.TokenHeader()
		req.Header.Add(headers.Authorization, authToken)
	}

	req.Header.Add(headers.Accept, "application/vnd.github.v3+json")
	res, err := re.client.Do(req)
	if err != nil {
		return Tag{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return Tag{}, err
		}
		return Tag{}, errors.New(string(b))
	}

	var tag Tag
	if err := json.NewDecoder(res.Body).Decode(&tag); err != nil {
		return Tag{}, err
	}

	return tag, nil
}
