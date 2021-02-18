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
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
)

type RepoManager struct {
	client *http.Client
}

func NewRepoManager(client *http.Client) RepoManager {
	return RepoManager{client: client}
}

func (re RepoManager) Zipball(info git.RepoInfo, version string) (io.ReadCloser, error) {
	zipURL := info.ZipUrl(version)
	res, err := re.performRequest(info, zipURL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		all, _ := ioutil.ReadAll(res.Body)
		return nil, errors.New(res.Status + "-" + string(all))
	}

	return res.Body, nil
}

func (re RepoManager) Tags(info git.RepoInfo) (git.Tags, error) {
	apiURL := info.TagsUrl()
	res, err := re.performRequest(info, apiURL)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(res.Body)
		return nil, errors.New(res.Status + "-" + string(all))
	}

	var tags git.Tags
	if err := json.NewDecoder(res.Body).Decode(&tags); err != nil {
		return git.Tags{}, err
	}

	return tags, nil
}

func (re RepoManager) LatestTag(info git.RepoInfo) (git.Tag, error) {
	apiURL := info.LatestTagUrl()
	res, err := re.performRequest(info, apiURL)
	if err != nil {
		return git.Tag{}, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(res.Body)
		return git.Tag{}, errors.New(res.Status + "-" + string(all))
	}

	githubTag := struct {
		Name               string `json:"tag_name"`
		ReleaseDescription string `json:"body"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&githubTag); err != nil {
		return git.Tag{}, err
	}
	tag := git.Tag{
		Name:        githubTag.Name,
		Description: githubTag.ReleaseDescription,
	}
	return tag, nil
}

func (re RepoManager) performRequest(info git.RepoInfo, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if info.Token() != "" {
		authToken := info.TokenHeader()
		req.Header.Add(headers.Authorization, authToken)
	}

	req.Header.Add(headers.Accept, "application/vnd.github.v3+json")
	res, err := re.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
