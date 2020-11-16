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
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestMergedTree(t *testing.T) {
	// Setup
	fileManager := stream.NewFileManager()
	providers := formula.NewRepoProviders()

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

	// End setup

	type in struct {
		repo      formula.RepositoryLister
		file      stream.FileReadExister
		providers formula.RepoProviders
		core      bool
	}

	tests := []struct {
		name string
		in   in
		want formula.Tree
	}{
		{
			name: "success",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{repo1, repo2}, nil
					},
				},
				file:      fileManager,
				providers: providers,
			},
			want: expectedTree,
		},
		{
			name: "success with core commands",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{repo1, repo2}, nil
					},
				},
				file:      fileManager,
				providers: providers,
				core:      true,
			},
			want: expectedTreeWithCoreCmds,
		},
		{
			name: "return empty tree when invalid tree",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{repoInvalid}, nil
					},
				},
				file:      fileManager,
				providers: providers,
				core:      false,
			},
			want: formula.Tree{Version: treeVersion, Commands: api.Commands{}},
		},
		{
			name: "empty tree when tree.json does not exist",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{repo1}, nil
					},
				},
				file: FileReadExisterMock{
					exists: func(path string) bool {
						return false
					},
				},
				providers: providers,
				core:      false,
			},
			want: formula.Tree{Version: treeVersion, Commands: api.Commands{}},
		},
		{
			name: "read tree.json error",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{repo1}, nil
					},
				},
				file: FileReadExisterMock{
					exists: func(path string) bool {
						return true
					},
					read: func(path string) ([]byte, error) {
						return nil, errors.New("error to read file")
					},
				},
				providers: providers,
				core:      false,
			},
			want: formula.Tree{Version: treeVersion, Commands: api.Commands{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := NewTreeManager(ritHome, tt.in.repo, coreCmds, tt.in.file, tt.in.providers)

			got := tree.MergedTree(tt.in.core)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergedTree(%s) = %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTree(t *testing.T) {
	// Setup
	fileManager := stream.NewFileManager()
	providers := formula.NewRepoProviders()

	treeb1, _ := json.Marshal(tree1)
	treeb2, _ := json.Marshal(tree2)
	repo1Path := filepath.Join(ritHome, "repos", strings.ToLower(repo1.Name.String()), "tree.json")
	repo2Path := filepath.Join(ritHome, "repos", strings.ToLower(repo2.Name.String()), "tree.json")
	repo3Path := filepath.Join(ritHome, "repos", "invalid", "tree.json")

	_ = os.MkdirAll(filepath.Dir(repo1Path), os.ModePerm)
	_ = os.MkdirAll(filepath.Dir(repo2Path), os.ModePerm)
	_ = os.MkdirAll(filepath.Dir(repo3Path), os.ModePerm)

	_ = fileManager.Write(repo1Path, treeb1)
	_ = fileManager.Write(repo2Path, treeb2)
	_ = fileManager.Write(repo3Path, []byte("invalid"))

	// End setup

	type in struct {
		repo      formula.RepositoryLister
		file      stream.FileReadExister
		providers formula.RepoProviders
	}

	type want struct {
		treeByRepo map[formula.RepoName]formula.Tree
		err        error
	}

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "success",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{repo1, repo2}, nil
					},
				},
				file:      fileManager,
				providers: providers,
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
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, errors.New("repo list error")
					},
				},
				file:      fileManager,
				providers: providers,
			},
			want: want{
				err: errors.New("repo list error"),
			},
		},
		{
			name: "return repos with empty tree commands when tree.json does not exist",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{repo1, repo2}, nil
					},
				},
				file: FileReadExisterMock{
					exists: func(path string) bool {
						return false
					},
				},
				providers: providers,
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
			name: "read tree.json error",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{repo1, repo2}, nil
					},
				},
				file: FileReadExisterMock{
					exists: func(path string) bool {
						return true
					},
					read: func(path string) ([]byte, error) {
						return []byte("error"), errors.New("error to read tree.json")
					},
				},
				providers: providers,
			},
			want: want{
				err: errors.New("error to read tree.json"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := NewTreeManager(ritHome, tt.in.repo, coreCmds, tt.in.file, tt.in.providers)

			got, err := tree.Tree()

			if !reflect.DeepEqual(got, tt.want.treeByRepo) {
				t.Errorf("Tree(%s) got %v; want %v", tt.name, got, tt.want.treeByRepo)
			}

			if tt.want.err != nil && tt.want.err.Error() != err.Error() {
				t.Errorf("Tree(%s) got %v; want error %v", tt.name, err, tt.want.err)
			}
		})
	}
}

func BenchmarkMergedTree(b *testing.B) {
	fileManager := stream.NewFileManager()
	providers := formula.NewRepoProviders()

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

	repoMock := repositoryListerCustomMock{
		list: func() (formula.Repos, error) {
			return formula.Repos{repo1, repo2}, nil
		},
	}

	tree := NewTreeManager(ritHome, repoMock, coreCmds, fileManager, providers)

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
	}
	repo2 = formula.Repo{
		Name:     formula.RepoName("repo2"),
		Priority: 1,
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
		Version: treeVersion,
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
		Version: treeVersion,
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
		Version: treeVersion,
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

	expectedTreeWithCoreCmds = formula.Tree{
		Version: treeVersion,
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

type repositoryListerCustomMock struct {
	list func() (formula.Repos, error)
}

func (m repositoryListerCustomMock) List() (formula.Repos, error) {
	return m.list()
}

type FileReadExisterMock struct {
	read   func(path string) ([]byte, error)
	exists func(path string) bool
}

func (m FileReadExisterMock) Read(path string) ([]byte, error) {
	return m.read(path)
}

func (m FileReadExisterMock) Exists(path string) bool {
	return m.exists(path)
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
