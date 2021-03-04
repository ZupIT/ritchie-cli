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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	repoUrl = "http://github.com/username/repo-name"
	token   = "some_token"
	version = "0.0.3"
)

func TestLatestTagUrl(t *testing.T) {
	const want = "https://api.github.com/repos/username/repo-name/releases/latest"
	repoInfo := NewRepoInfo(repoUrl, token)
	latestTagsUrl := repoInfo.LatestTagUrl()

	assert.Equal(t, want, latestTagsUrl)
}

func TestTagsUrl(t *testing.T) {
	const want = "https://api.github.com/repos/username/repo-name/releases"
	repoInfo := NewRepoInfo(repoUrl, token)
	tagsUrl := repoInfo.TagsUrl()

	assert.Equal(t, want, tagsUrl)
}

func TestTokenHeader(t *testing.T) {
	const want = "token some_token"
	repoInfo := NewRepoInfo(repoUrl, token)
	tokenHeader := repoInfo.TokenHeader()

	assert.Equal(t, want, tokenHeader)
}

func TestZipUrl(t *testing.T) {
	const want = "https://api.github.com/repos/username/repo-name/zipball/0.0.3"
	repoInfo := NewRepoInfo(repoUrl, token)
	zipUrl := repoInfo.ZipUrl(version)

	assert.Equal(t, want, zipUrl)
}

func TestToken(t *testing.T) {
	const want = "some_token"
	repoInfo := NewRepoInfo(repoUrl, token)
	token := repoInfo.Token()

	assert.Equal(t, want, token)
}
