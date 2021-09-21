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
package zipremote

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/git"
)

type RepoManager struct {
	client *http.Client
}

func NewRepoManager(client *http.Client) RepoManager {
	return RepoManager{client: client}
}

func (re RepoManager) Zipball(info git.RepoInfo, version string) (io.ReadCloser, error) {
	zipUrl := info.ZipUrl(version)
	res, err := re.performRequest(zipUrl)
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
	return git.Tags{}, nil
}

func (re RepoManager) LatestTag(info git.RepoInfo) (git.Tag, error) {
	return git.Tag{}, nil
}

func (re RepoManager) performRequest(url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := re.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
