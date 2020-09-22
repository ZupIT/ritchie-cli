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

package mock

import (
	"io"

	"github.com/ZupIT/ritchie-cli/pkg/git"
)

var DefaultGitRepositoryMock = GitRepositoryMock{
	latestTag: func(info git.RepoInfo) (git.Tag, error) {
		return git.Tag{}, nil
	},
	tags: func(info git.RepoInfo) (git.Tags, error) {
		return git.Tags{git.Tag{Name: "1.0.0"}}, nil
	},
	zipball: func(info git.RepoInfo, version string) (io.ReadCloser, error) {
		return nil, nil
	},
}

type GitRepositoryMock struct {
	zipball   func(info git.RepoInfo, version string) (io.ReadCloser, error)
	tags      func(info git.RepoInfo) (git.Tags, error)
	latestTag func(info git.RepoInfo) (git.Tag, error)
}

func (m GitRepositoryMock) Zipball(info git.RepoInfo, version string) (io.ReadCloser, error) {
	return m.zipball(info, version)
}

func (m GitRepositoryMock) Tags(info git.RepoInfo) (git.Tags, error) {
	return m.tags(info)
}

func (m GitRepositoryMock) LatestTag(info git.RepoInfo) (git.Tag, error) {
	return m.latestTag(info)
}
