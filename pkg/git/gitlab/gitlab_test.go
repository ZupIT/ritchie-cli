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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepoInfo(t *testing.T) {
	want := DefaultRepoInfo{
		host:  "gitlab.com",
		owner: "username",
		repo:  "ritchie-formulas",
		token: "some_token",
	}
	got := NewRepoInfo("https://gitlab.com/username/ritchie-formulas", "some_token")

	assert.Equal(t, want, got)
}

func TestTagsUrl(t *testing.T) {
	const want = "https://gitlab.com/api/v4/projects/username%2Fritchie-formulas/releases"
	repoInfo := NewRepoInfo("https://gitlab.com/username/ritchie-formulas", "some_token")
	tagsUrl := repoInfo.TagsUrl()

	assert.Equal(t, want, tagsUrl)
}

func TestZipUrl(t *testing.T) {
	const want = "https://gitlab.com/api/v4/projects/username%2Fritchie-formulas/repository/archive.zip?sha=1.0.0"
	repoInfo := NewRepoInfo("https://gitlab.com/username/ritchie-formulas", "some_token")
	zipUrl := repoInfo.ZipUrl("1.0.0")

	assert.Equal(t, want, zipUrl)
}

func TestLatestTagUrl(t *testing.T) {
	const want = "https://gitlab.com/api/v4/projects/username%2Fritchie-formulas/releases?per_page=1&page=1"
	repoInfo := NewRepoInfo("https://gitlab.com/username/ritchie-formulas", "some_token")
	latestTagUrl := repoInfo.LatestTagUrl()

	assert.Equal(t, want, latestTagUrl)
}

func TestTokenHeader(t *testing.T) {
	const want = "some_token"
	repoInfo := NewRepoInfo("https://gitlab.com/username/ritchie-formulas", "some_token")
	tokenHeader := repoInfo.TokenHeader()

	assert.Equal(t, want, tokenHeader)
}

func TestToken(t *testing.T) {
	const want = "some_token"
	repoInfo := NewRepoInfo("https://gitlab.com/username/ritchie-formulas", "some_token")
	token := repoInfo.Token()

	assert.Equal(t, want, token)
}
