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
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/git"
)

const (
	ZipUrlPattern       = "https://api.github.com/repos/%s/%s/zipball/%s"      //nolint:stylecheck
	TagsUrlPattern      = "https://api.github.com/repos/%s/%s/releases"        //nolint:stylecheck
	LatestTagUrlPattern = "https://api.github.com/repos/%s/%s/releases/latest" //nolint:stylecheck
	RitchieRepoURL      = "https://github.com/ZupIT/ritchie-cli"
)

type DefaultRepoInfo struct {
	owner string
	repo  string
	token string
}

// NewRepoInfo returns the RepoInfo built by repository url
// Repository url e.g. https://github.com/{{owner}}/{{repo}}
func NewRepoInfo(url string, token string) git.RepoInfo {
	split := strings.Split(url, "/")
	repo := split[len(split)-1]
	owner := split[len(split)-2]

	return DefaultRepoInfo{
		owner: owner,
		repo:  repo,
		token: token,
	}
}

// ZipUrl returns the GitHub API URL for download zipball repository
// e.g. https://api.github.com/repos/{{owner}}/{{repo}}/zipball/{{tag-version}}
func (in DefaultRepoInfo) ZipUrl(version string) string { //nolint:stylecheck
	return fmt.Sprintf(ZipUrlPattern, in.owner, in.repo, version)
}

// TagsUrl returns the GitHub API URL for get all tags
// e.g. https://api.github.com/repos/{{owner}}/{{repo}}/tags
func (in DefaultRepoInfo) TagsUrl() string { //nolint:stylecheck
	return fmt.Sprintf(TagsUrlPattern, in.owner, in.repo)
}

// LatestTagUrl returns the GitHub API URL for get latest tag release
// https://api.github.com/repos/:owner/:repo/releases/latest
func (in DefaultRepoInfo) LatestTagUrl() string { //nolint:stylecheck
	return fmt.Sprintf(LatestTagUrlPattern, in.owner, in.repo)
}

// TokenHeader returns the Authorization value formatted for Github API integration
// e.g. "token f39c5aca-858f-4a04-9ca3-5104d02b9c56"
func (in DefaultRepoInfo) TokenHeader() string {
	return fmt.Sprintf("token %s", in.token)
}

func (in DefaultRepoInfo) Token() string {
	return in.token
}
