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

package tree

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/mock"
)

func TestTree(t *testing.T) {
	type in struct {
		RepositoryLister formula.RepositoryLister
		Repos            git.Repositories
		Core             bool
		// wantMergedTree   formula.Tree
	}
	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name:     "someRepo1",
								Provider: "Github",
								Url:      "https://github.com/owner/repo",
								Token:    "token",
							},
						}, nil
					},
				},
				Repos: mock.DefaultGitRepositoryMock,
				Core:  false,
			},
			wantErr: false,
			// wantMergedTree:,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// repoProviders := formula.NewRepoProviders()
			// repoProviders.Add("Github", formula.Git{Repos: tt.in.Repos, NewRepoInfo: github.NewRepoInfo})

			// commands := api.Commands{
			// 	{
			// 		Id:     "root_mock",
			// 		Parent: "root",
			// 		Usage:  "mock",
			// 		Help:   "mock for add",
			// 	},
			// 	{
			// 		Id:      "root_mock_test",
			// 		Parent:  "root_mock",
			// 		Usage:   "test",
			// 		Help:    "test for add",
			// 		Formula: true,
			// 	},
			// }

			// newTree := NewTreeManager("any", tt.in.RepositoryLister, commands, repoProviders)
			// mergedTree := newTree.MergedTree(tt.in.Core)
			// if mergedTree != tt.wantMergedTree {
			// 	t.Errorf("NewTreeManager_MergedTree() mergedTree = %v, want mergedTree %v", ermergedTreer, tt.wantErwantMergedTreer)
			// }
		})
	}
}

type RepositoryListerCustomMock struct {
	list func() (formula.Repos, error)
}

func (m RepositoryListerCustomMock) List() (formula.Repos, error) {
	return m.list()
}
