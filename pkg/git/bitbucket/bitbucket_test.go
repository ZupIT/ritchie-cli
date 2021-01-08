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
	"reflect"
	"testing"
)

func TestNewRepoInfo(t *testing.T) {
	want := DefaultRepoInfo{
		host:  "bitbucket.org",
		owner: "username",
		repo:  "ritchie-formulas",
		token: "gHexna7h7CQWafwNYSXp",
	}
	got := NewRepoInfo("https://bitbucket.org/username/ritchie-formulas/src/master/", "gHexna7h7CQWafwNYSXp")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewRepoInfo() = %v, want %v", got, want)
	}
}

func TestTagsUrl(t *testing.T) {
	const want = "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/refs/tags"
	repoInfo := NewRepoInfo("https://bitbucket.org/username/ritchie-formulas/src/master/", "gHexna7h7CQWafwNYSXp")
	tagsUrl := repoInfo.TagsUrl()

	if !reflect.DeepEqual(tagsUrl, want) {
		t.Errorf("NewRepoInfo() = %v, want %v", "got", want)
	}
}

func TestZipUrl(t *testing.T) {
	const want = "https://bitbucket.org/username/ritchie-formulas/get/1.0.0.zip"
	repoInfo := NewRepoInfo("https://bitbucket.org/username/ritchie-formulas/src/master/", "gHexna7h7CQWafwNYSXp")
	tagsUrl := repoInfo.ZipUrl("1.0.0")

	if !reflect.DeepEqual(tagsUrl, want) {
		t.Errorf("NewRepoInfo() = %v, want %v", tagsUrl, want)
	}
}
