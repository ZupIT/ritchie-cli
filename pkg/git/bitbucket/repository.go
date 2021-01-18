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
package bitbucket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
)

type RepoManager struct {
	client *http.Client
}

type value struct {
	Name string `json:"name"`
}

type bitbucketTags struct {
	Values []value `json:"values"`
}

func NewRepoManager(client *http.Client) RepoManager {
	return RepoManager{client: client}
}

func (re RepoManager) Zipball(info git.RepoInfo, version string) (io.ReadCloser, error) {
	zipUrl := info.ZipUrl(version)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, zipUrl, nil)
	if err != nil {
		return nil, err
	}

	if info.Token() != "" {
		req.Header.Add(headers.Authorization, fmt.Sprintf("Bearer %s", info.Token()))
	}

	resp, err := re.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		all, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(all)
		return nil, errors.New(resp.Status)
	}

	return resp.Body, nil
}

func (re RepoManager) Tags(info git.RepoInfo) (git.Tags, error) {
	apiUrl := info.TagsUrl()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiUrl, nil)
	if err != nil {
		return git.Tags{}, err
	}

	if info.Token() != "" {
		req.Header.Add(headers.Authorization, fmt.Sprintf("Bearer %s", info.Token()))
	}

	res, err := re.client.Do(req)
	if err != nil {
		return git.Tags{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("There was an error adding the repository, status: %d - %s.", res.StatusCode, http.StatusText(res.StatusCode))
		return git.Tags{}, errors.New(errorMessage)
	}

	var bTags bitbucketTags
	if err := json.NewDecoder(res.Body).Decode(&bTags); err != nil {
		return git.Tags{}, err
	}

	var tags git.Tags
	for _, v := range bTags.Values {
		tag := git.Tag{Name: v.Name}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (re RepoManager) LatestTag(info git.RepoInfo) (git.Tag, error) {
	apiUrl := info.LatestTagUrl()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiUrl, nil)
	if err != nil {
		return git.Tag{}, err
	}

	if info.Token() != "" {
		req.Header.Add(headers.Authorization, fmt.Sprintf("Bearer %s", info.Token()))
	}

	res, err := re.client.Do(req)
	if err != nil {
		return git.Tag{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return git.Tag{}, err
		}
		return git.Tag{}, errors.New(string(b))
	}

	var bTags bitbucketTags
	if err := json.NewDecoder(res.Body).Decode(&bTags); err != nil {
		return git.Tag{}, err
	}

	if len(bTags.Values) == 0 {
		return git.Tag{}, errors.New("release not found")
	}

	latestTag := git.Tag{Name: bTags.Values[len(bTags.Values)-1].Name}

	return latestTag, nil
}
