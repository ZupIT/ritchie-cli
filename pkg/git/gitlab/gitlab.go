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
package gitlab

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/git"
)

const (
	// https://gitlab.com/api/v4/projects/kaduartur%2Fritchie-formulas-kadu/repository/archive.zip?sha=1.0.1
	ZipUrlPattern  = "https://%s/api/v4/projects/%s/repository/archive.zip?sha=%s"
	TagsUrlPattern = "https://%s/api/v4/projects/%s/releases"
)

type DefaultRepoInfo struct {
	host  string
	owner string
	repo  string
	token string
}

// NewRepoInfo returns the RepoInfo built by repository url
// Repository url e.g. https://gitlab.com/{{owner}}/{{repo}}
func NewRepoInfo(url string, token string) git.RepoInfo {
	split := strings.Split(url, "/")
	repo := split[len(split)-1]
	owner := split[len(split)-2]
	host := split[len(split)-3]

	return DefaultRepoInfo{
		host:  host,
		owner: owner,
		repo:  repo,
		token: token,
	}
}

// ZipUrl returns the Gitlab API URL for download zipball repository
// e.g. https://yourhost/{{owner}}/{{repo}}/-/archive/{{tag-version}}/{{repo}}-{{tag-version}}.zip
func (in DefaultRepoInfo) ZipUrl(version string) string {
	id := url.QueryEscape(path.Join(in.owner, in.repo))
	return fmt.Sprintf(ZipUrlPattern, in.host, id, version)
}

// TagsUrl returns the Gitlab API URL for get all tags
// e.g. https://yourhost/api/v4/projects/{{owner}}%2F{{repo}}/releases
func (in DefaultRepoInfo) TagsUrl() string {
	id := url.QueryEscape(path.Join(in.owner, in.repo))
	return fmt.Sprintf(TagsUrlPattern, in.host, id)
}

// Deprecated: Gitlab API does not implement the latest tag URL
func (in DefaultRepoInfo) LatestTagUrl() string {
	return ""
}

// TokenHeader returns the Authorization value formatted for Gitlab API integration
// e.g. "f39c5aca-858f-4a04-9ca3-5104d02b9c56"
func (in DefaultRepoInfo) TokenHeader() string {
	return in.token
}

// Deprecated: Uses TokenHeader() function
func (in DefaultRepoInfo) Token() string {
	return in.token
}
