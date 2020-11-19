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

package repo

import (
	"errors"
	"io"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
)

func TestNewDetailLatestTag(t *testing.T) {
	repoProviders := formula.NewRepoProviders()
	defaultGitRepositoryMock := GitRepositoryMock{
		tags: func(info git.RepoInfo) (git.Tags, error) {
			return git.Tags{git.Tag{}}, nil
		},
		zipball: func(info git.RepoInfo, version string) (io.ReadCloser, error) {
			return nil, nil
		},
	}

	type fields struct {
		repo          formula.Repo
		repoProviders formula.RepoProviders
		funcLatestTag func(info git.RepoInfo) (git.Tag, error)
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Return version",
			fields: fields{
				repoProviders: repoProviders,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "testing_repo",
					Version:  "0.0.0",
					Url:      "https://github.com/viniciussousazup/ritchie-formulas/releases",
					Token:    "",
					Priority: 0,
				},
				funcLatestTag: func(info git.RepoInfo) (git.Tag, error) {
					return git.Tag{"1.0.0"}, nil
				},
			},
			want: "1.0.0",
		},
		{
			name: "Return version nill when get latest returns error",
			fields: fields{
				repoProviders: repoProviders,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "testing_repo",
					Version:  "0.0.0",
					Url:      "https://github.com/viniciussousazup/ritchie-formulas/releases",
					Token:    "",
					Priority: 0,
				},
				funcLatestTag: func(info git.RepoInfo) (git.Tag, error) {
					return git.Tag{}, errors.New("some error")
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultGitRepositoryMock.latestTag = tt.fields.funcLatestTag
			repoProviders.Add("Github", formula.Git{Repos: defaultGitRepositoryMock, NewRepoInfo: github.NewRepoInfo})
			dm := NewDetail(tt.fields.repoProviders)

			tag := dm.LatestTag(tt.fields.repo)

			if tag != tt.want {
				t.Errorf("TestNewDetailLatestTag() receive = %v, expected %v", tag, tt.want)
			}
		})
	}
}
