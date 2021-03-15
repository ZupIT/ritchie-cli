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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestMergedTree(t *testing.T) {
	defaultTreeSetup()

	type repo struct {
		repos   formula.Repos
		listErr error
	}

	type in struct {
		repo    repo
		core    bool
	}

	tests := []struct {
		name string
		in   in
		want formula.Tree
	}{
		{
			name: "success",
			in: in{
				repo: repo{
					repos: formula.Repos{repo1, repo2},
				},
			},
			want: expectedTree,
		},
		{
			name: "success with core commands",
			in: in{
				repo: repo{
					repos: formula.Repos{repo1, repo2},
				},
				core: true,
			},
			want: expectedTreeWithCoreCmds,
		},
		{
			name: "return empty tree when invalid tree",
			in: in{
				repo: repo{
					repos: formula.Repos{repoInvalid},
				},
				core: false,
			},
			want: formula.Tree{
				Version:    Version,
				Commands:   api.Commands{},
				CommandsID: []api.CommandID{},
			},
		},
		{
			name: "empty tree when tree.json does not exist",
			in: in{
				repo: repo{
					repos: formula.Repos{repoInvalid},
				},
				core: false,
			},
			want: formula.Tree{
				Version:    Version,
				Commands:   api.Commands{},
				CommandsID: []api.CommandID{},
			},
		},
		{
			name: "unmarshal tree.json error",
			in: in{
				repo: repo{
					repos: formula.Repos{repoInvalid},
				},
				core: false,
			},
			want: formula.Tree{
				Version:    Version,
				Commands:   api.Commands{},
				CommandsID: []api.CommandID{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mocks.RepoManager)
			repoMock.On("List").Return(tt.in.repo.repos, tt.in.repo.listErr)
			repoMock.On("LatestTag", mock.Anything).Return("3.0.0")
			repoMock.On("Write", mock.Anything).Return(nil)

			tree := NewTreeManager(ritHome, repoMock, coreCmds)

			got := tree.MergedTree(tt.in.core)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTree(t *testing.T) {
	defaultTreeSetup()

	type (
		repo struct {
			repos   formula.Repos
			listErr error
		}
		in struct {
			ritHome string
			repo    repo
		}

		want struct {
			treeByRepo map[formula.RepoName]formula.Tree
			err        error
		}
	)

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "success",
			in: in{
				ritHome: ritHome,
				repo: repo{
					repos: formula.Repos{repo1, repo2},
				},
			},
			want: want{
				treeByRepo: map[formula.RepoName]formula.Tree{
					core:    {Commands: coreCmds},
					"repo1": tree1,
					"repo2": tree2,
				},
				err: nil,
			},
		},
		{
			name: "repo list error",
			in: in{
				ritHome: ritHome,
				repo: repo{
					repos:   formula.Repos{},
					listErr: errors.New("repo list error"),
				},
			},
			want: want{
				err: errors.New("repo list error"),
			},
		},
		{
			name: "return repos with empty tree commands when tree.json does not exist",
			in: in{
				ritHome: "/invalid",
				repo: repo{
					repos: formula.Repos{repo1, repo2},
				},
			},
			want: want{
				treeByRepo: map[formula.RepoName]formula.Tree{
					core:    {Commands: coreCmds},
					"repo1": {},
					"repo2": {},
				},
				err: nil,
			},
		},
		{
			name: "unmarshal tree.json error",
			in: in{
				ritHome: ritHome,
				repo: repo{
					repos: formula.Repos{repoInvalid},
				},
			},
			want: want{
				err: errors.New("invalid character 'i' looking for beginning of value"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mocks.RepoManager)
			repoMock.On("List").Return(tt.in.repo.repos, tt.in.repo.listErr)
			repoMock.On("LatestTag", mock.Anything).Return("3.0.0")
			repoMock.On("Write", mock.Anything).Return(nil)

			tree := NewTreeManager(tt.in.ritHome, repoMock, coreCmds)

			got, err := tree.Tree()

			assert.Equal(t, tt.want.treeByRepo, got)

			if tt.want.err != nil || err != nil {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func BenchmarkMergedTree(b *testing.B) {
	defaultTreeSetup()

	repoMock := new(mocks.RepoManager)
	repoMock.On("List").Return(formula.Repos{repo1, repo2}, nil)
	repoMock.On("LatestTag", mock.Anything).Return("3.0.0")
	repoMock.On("Write", mock.Anything).Return(nil)

	tree := NewTreeManager(ritHome, repoMock, coreCmds)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.MergedTree(false)
	}
}

var (
	tmpDir  = os.TempDir()
	ritHome = filepath.Join(tmpDir, ".rit-tree")

	repo1 = formula.Repo{
		Name:     formula.RepoName("repo1"),
		Priority: 0,
		Cache:    time.Now().Add(time.Hour),
	}
	repo2 = formula.Repo{
		Name:     formula.RepoName("repo2"),
		Priority: 1,
		Cache:    time.Now().Add(-time.Hour),
	}

	repoInvalid = formula.Repo{
		Name:     formula.RepoName("invalid"),
		Priority: 2,
	}

	coreCmds = api.Commands{
		"root_add":      {Parent: "root", Usage: "add"},
		"root_add_repo": {Parent: "root_add", Usage: "repo"},
	}

	tree1 = formula.Tree{
		Version: Version,
		Commands: api.Commands{
			"root_pokemon": api.Command{
				Parent:   "root",
				Usage:    "pokemon",
				Help:     "pokemon help",
				LongHelp: "pokemon help long",
				Formula:  false,
			},
			"root_pokemon_add": api.Command{
				Parent:   "root_pokemon",
				Usage:    "add-new-pokemon",
				Help:     "pokemon add-new-pokemon help",
				LongHelp: "pokemon add-new-pokemon help long",
				Formula:  true,
			},
		},
	}

	tree2 = formula.Tree{
		Version: Version,
		Commands: api.Commands{
			"root_star_wars": api.Command{
				Parent:   "root",
				Usage:    "star-wars",
				Help:     "star wars help",
				LongHelp: "star wars help long",
				Formula:  false,
			},
			"root_star_wars_list-jedis": api.Command{
				Parent:   "root_star_wars",
				Usage:    "list-jedis",
				Help:     "star wars list-jedis help",
				LongHelp: "star wars list-jedis help long",
				Formula:  true,
			},
		},
	}

	expectedTree = formula.Tree{
		Version: Version,
		CommandsID: []api.CommandID{
			"root_pokemon",
			"root_star_wars",
			"root_pokemon_add",
			"root_star_wars_list-jedis",
		},
		Commands: api.Commands{
			"root_pokemon": api.Command{
				Parent:   "root",
				Usage:    "pokemon",
				Help:     "pokemon help",
				LongHelp: "pokemon help long",
				Formula:  false,
				Repo:     "repo1",
			},
			"root_pokemon_add": api.Command{
				Parent:   "root_pokemon",
				Usage:    "add-new-pokemon",
				Help:     "pokemon add-new-pokemon help",
				LongHelp: "pokemon add-new-pokemon help long",
				Formula:  true,
				Repo:     "repo1",
			},
			"root_star_wars": api.Command{
				Parent:         "root",
				Usage:          "star-wars",
				Help:           "star wars help",
				LongHelp:       "star wars help long",
				Formula:        false,
				Repo:           "repo2",
				RepoNewVersion: "3.0.0",
			},
			"root_star_wars_list-jedis": api.Command{
				Parent:   "root_star_wars",
				Usage:    "list-jedis",
				Help:     "star wars list-jedis help",
				LongHelp: "star wars list-jedis help long",
				Formula:  true,
				Repo:     "repo2",
			},
		},
	}

	expectedTreeWithCoreCmds = formula.Tree{
		Version: Version,
		CommandsID: []api.CommandID{
			"root_pokemon",
			"root_star_wars",
			"root_pokemon_add",
			"root_star_wars_list-jedis",
		},
		Commands: api.Commands{
			"root_pokemon": api.Command{
				Parent:   "root",
				Usage:    "pokemon",
				Help:     "pokemon help",
				LongHelp: "pokemon help long",
				Formula:  false,
				Repo:     "repo1",
			},
			"root_pokemon_add": api.Command{
				Parent:   "root_pokemon",
				Usage:    "add-new-pokemon",
				Help:     "pokemon add-new-pokemon help",
				LongHelp: "pokemon add-new-pokemon help long",
				Formula:  true,
				Repo:     "repo1",
			},
			"root_star_wars": api.Command{
				Parent:         "root",
				Usage:          "star-wars",
				Help:           "star wars help",
				LongHelp:       "star wars help long",
				Formula:        false,
				Repo:           "repo2",
				RepoNewVersion: "3.0.0",
			},
			"root_star_wars_list-jedis": api.Command{
				Parent:   "root_star_wars",
				Usage:    "list-jedis",
				Help:     "star wars list-jedis help",
				LongHelp: "star wars list-jedis help long",
				Formula:  true,
				Repo:     "repo2",
			},
			"root_add": {
				Parent: "root",
				Usage:  "add",
			},
			"root_add_repo": {
				Parent: "root_add",
				Usage:  "repo",
			},
		},
	}
)

func defaultTreeSetup() {
	fileManager := stream.NewFileManager()

	tree1, _ := json.Marshal(tree1)
	tree2, _ := json.Marshal(tree2)
	repo1Path := filepath.Join(ritHome, "repos", strings.ToLower(repo1.Name.String()), "tree.json")
	repo2Path := filepath.Join(ritHome, "repos", strings.ToLower(repo2.Name.String()), "tree.json")
	repo3Path := filepath.Join(ritHome, "repos", "invalid", "tree.json")

	_ = os.MkdirAll(filepath.Dir(repo1Path), os.ModePerm)
	_ = os.MkdirAll(filepath.Dir(repo2Path), os.ModePerm)
	_ = os.MkdirAll(filepath.Dir(repo3Path), os.ModePerm)

	_ = fileManager.Write(repo1Path, tree1)
	_ = fileManager.Write(repo2Path, tree2)
	_ = fileManager.Write(repo3Path, []byte("invalid"))
}
