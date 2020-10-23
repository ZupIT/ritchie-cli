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
package git

import "io"

type Tag struct {
	Name string `json:"tag_name"`
}

type Tags []Tag

func (t Tags) Names() []string {
	tags := make([]string, 0, len(t))
	for i := range t {
		tags = append(tags, t[i].Name)
	}

	return tags
}

type RepoInfo interface {
	ZipUrl(version string) string
	TagsUrl() string
	LatestTagUrl() string
	TokenHeader() string
	Token() string
}

type Repositories interface {
	Zipball(info RepoInfo, version string) (io.ReadCloser, error)
	Tags(info RepoInfo) (Tags, error)
	LatestTag(info RepoInfo) (Tag, error)
}
