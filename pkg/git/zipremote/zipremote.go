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
	"github.com/ZupIT/ritchie-cli/pkg/git"
)

type DefaultRepoInfo struct {
	url string
}

// NewRepoInfo returns the RepoInfo built by repository url
// Repository url e.g. https://github.com/{{owner}}/{{repo}}/archive/refs/tags/{{tag-version}}.zip
func NewRepoInfo(url string, token string) git.RepoInfo {
	return DefaultRepoInfo{url}
}

// ZipUrl returns the URL for download zipball repository
// e.g. https://github.com/{{owner}}/{{repo}}/archive/refs/tags/{{tag-version}}.zip
func (in DefaultRepoInfo) ZipUrl(version string) string {
	return in.url
}

// TagsUrl returns the URL for get all tags
func (in DefaultRepoInfo) TagsUrl() string {
	return ""
}

// LatestTagUrl returns the URL for get latest tag release
func (in DefaultRepoInfo) LatestTagUrl() string {
	return ""
}

// TokenHeader returns the Authorization value formatted
func (in DefaultRepoInfo) TokenHeader() string {
	return ""
}

func (in DefaultRepoInfo) Token() string {
	return ""
}
