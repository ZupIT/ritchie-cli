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
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/git"
)

const (
	ZipUrlPattern       = "https://bitbucket.org/%s/%s/get/%s.zip"
	TagsUrlPattern      = "https://api.bitbucket.org/2.0/repositories/%s/%s/refs/tags"
	LatestTagUrlPattern = "https://api.bitbucket.org/2.0/repositories/%s/%s/refs/tags?sort=target.date"
)

type DefaultRepoInfo struct {
	host  string
	owner string
	repo  string
	token string
}

// NewRepoInfo returns the RepoInfo built by repository url
// Repository url e.g. https://bitbucket.org/{{owner}}/{{repo}}/src/master/
func NewRepoInfo(url string, token string) git.RepoInfo {
	split := strings.Split(url, "/")
	if len(split) < 8 {
		return DefaultRepoInfo{}
	}

	repo := split[len(split)-4]
	owner := split[len(split)-5]
	host := split[len(split)-6]

	return DefaultRepoInfo{
		host:  host,
		owner: owner,
		repo:  repo,
		token: token,
	}
}

// ZipUrl returns the Bitbucket API URL for download zipball repository
// e.g. https://bitbucket.org/{{owner}}/{{repo}}/get/{{tag-version}}.zip
func (in DefaultRepoInfo) ZipUrl(version string) string {
	return fmt.Sprintf(ZipUrlPattern, in.owner, in.repo, version)
}

// TagsUrl returns the Bitbucket API URL for get all tags
// e.g. https://api.bitbucket.org/2.0/repositories/{{owner}}/{{repo}}/refs/tags
func (in DefaultRepoInfo) TagsUrl() string {
	return fmt.Sprintf(TagsUrlPattern, in.owner, in.repo)
}

// LatestTagUrl returns the Bitbucket API URL for get latest tag release
// e.g. https://api.bitbucket.org/2.0/repositories/{{owner}}/{{repo}}/refs/tags?sort=target.date
func (in DefaultRepoInfo) LatestTagUrl() string {
	return fmt.Sprintf(LatestTagUrlPattern, in.owner, in.repo)
}

// TokenHeader returns the Authorization value formatted for Bitbucket API integration
// e.g. "m6ioSP4o4q6tmiinXHf9KOHxQsjbvShS-zMullcRsxiJtMHIqr2tuHpJZbl-UXpm7E-1meNlAqKzmORTxyoNAnXWZlCFPsvJpQj4evGtafuH4NBBgRrQ_Mc3"
func (in DefaultRepoInfo) TokenHeader() string {
	return in.token
}

func (in DefaultRepoInfo) Token() string {
	return in.token
}
