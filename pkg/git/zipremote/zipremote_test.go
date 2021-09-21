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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepoInfo(t *testing.T) {
	type in struct {
		url   string
		token string
	}
	tests := []struct {
		name string
		in   in
		want DefaultRepoInfo
	}{
		{
			name: "Run with success",
			in: in{
				url:   "https://provider.com/username/ritchie-formulas/src/master/",
				token: "some_token",
			},
			want: DefaultRepoInfo{url: "https://provider.com/username/ritchie-formulas/src/master/"},
		},
		{
			name: "Return err when the URL is incorrect",
			in: in{
				url:   "",
				token: "some_token",
			},
			want: DefaultRepoInfo{},
		},
	}

	for _, tt := range tests {
		got := NewRepoInfo(tt.in.url, tt.in.token)

		assert.Equal(t, tt.want, got)
	}
}

func TestTagsUrl(t *testing.T) {
	const want = ""
	repoInfo := NewRepoInfo("https://provider.com/username/ritchie-formulas/src/master/", "some_token")
	tagsUrl := repoInfo.TagsUrl()

	assert.Equal(t, want, tagsUrl)
}

func TestZipUrl(t *testing.T) {
	const want = "https://provider.com/username/ritchie-formulas/get/1.0.0.zip"
	repoInfo := NewRepoInfo("https://provider.com/username/ritchie-formulas/get/1.0.0.zip", "some_token")
	zipUrl := repoInfo.ZipUrl("1.0.0")

	assert.Equal(t, want, zipUrl)
}

func TestLatestTagUrl(t *testing.T) {
	const want = ""
	repoInfo := NewRepoInfo("https://provider.com/username/ritchie-formulas/src/master/", "some_token")
	latestTagsUrl := repoInfo.LatestTagUrl()

	assert.Equal(t, want, latestTagsUrl)
}

func TestTokenHeader(t *testing.T) {
	const want = ""
	repoInfo := NewRepoInfo("https://provider.com/username/ritchie-formulas/src/master/", "some_token")
	tokenHeader := repoInfo.TokenHeader()

	assert.Equal(t, want, tokenHeader)
}

func TestToken(t *testing.T) {
	const want = ""
	repoInfo := NewRepoInfo("https://provider.com/username/ritchie-formulas/src/master/", "some_token")
	token := repoInfo.Token()

	assert.Equal(t, want, token)
}
