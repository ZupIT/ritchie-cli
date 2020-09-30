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
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
)

var (
	tmpDir  = os.TempDir()
	ritHome = filepath.Join(tmpDir, ".rit-tree")

	coreCmds = []api.Command{
		{Parent: "root", Usage: "add"},
		{Parent: "root_add", Usage: "repo"},
		{Parent: "root", Usage: "metrics"},
	}

	someRepoTree = formula.Tree{Commands: []api.Command{
		{Parent: "root", Usage: "pokemon-list"},
		{Parent: "root_pokemon-list", Usage: "add"},
	}}
)

func TestMergedTree(t *testing.T) {
	defer cleanRitHome()

	errFoo := errors.New("some error")
	localTree := formula.Tree{Commands: []api.Command{
		{Parent: "root", Usage: "jedi-list"},
		{Parent: "root_jedi-list", Usage: "add"},
	}}

	expectedTreeComplete := formula.Tree{
		Commands: api.Commands{
			coreCmds[0],
			coreCmds[1],
			coreCmds[2],
			{
				Parent: "root",
				Usage:  "jedi-list",
				Repo:   "local",
			},
			{
				Parent: "root_jedi-list",
				Usage:  "add",
				Repo:   "local",
			},
			{
				Parent: "root",
				Usage:  "pokemon-list",
				Repo:   "someRepo",
			},
			{
				Parent: "root_pokemon-list",
				Usage:  "add",
				Repo:   "someRepo",
			},
		},
	}

	repoLister := repositoryListerCustomMock{
		list: func() (formula.Repos, error) {
			return formula.Repos{
				{
					Name:     "someRepo",
					Provider: "Github",
					Url:      "https://github.com/owner/repo",
					Token:    "token",
				},
				{
					Name:     "otherRepo",
					Provider: "Github",
					Url:      "https://github.com/owner/noame",
					Token:    "",
				},
			}, nil
		},
	}

	githubRepo := github.NewRepoManager(http.DefaultClient)
	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: githubRepo, NewRepoInfo: github.NewRepoInfo})

	fileManager := FileReadExisterMock{
		exists: func(path string) bool {
			isLocalRepo := strings.Contains(path, "/local/tree.json")
			isSomeRepo := strings.Contains(path, "/someRepo/tree.json")
			isOtherRepo := strings.Contains(path, "/otherRepo/tree.json")

			if isLocalRepo || isSomeRepo || isOtherRepo {
				return true
			}
			return false
		},
		read: func(path string) ([]byte, error) {
			isLocalRepo := strings.Contains(path, "/local/tree.json")
			isSomeRepo := strings.Contains(path, "/someRepo/tree.json")
			isOtherRepo := strings.Contains(path, "/otherRepo/tree.json")

			if isLocalRepo {
				return []byte(getStringOfTree(localTree)), nil
			}
			if isSomeRepo {
				return []byte(getStringOfTree(someRepoTree)), nil
			}
			if isOtherRepo {
				return []byte("any"), errFoo
			}
			return []byte("some data"), nil
		},
	}

	newTree := NewTreeManager(ritHome, repoLister, coreCmds, fileManager)
	mergedTree := newTree.MergedTree(true)

	if !isSameFormulaTree(mergedTree, expectedTreeComplete) {
		t.Errorf("NewTreeManager_MergedTree() \n\tmergedTree = %v\n\texpectedTree = %v", mergedTree, expectedTreeComplete)
	}
}

func TestTree(t *testing.T) {
	defer cleanRitHome()

	errFoo := errors.New("some error")
	expectedTreeEmpty := map[string]formula.Tree{}
	expectedTreeComplete := map[string]formula.Tree{
		"CORE": {
			Commands: api.Commands{
				coreCmds[0],
				coreCmds[1],
				coreCmds[2],
			},
		},
		"LOCAL": {
			Commands: api.Commands{},
		},
		"someRepo": {
			Commands: api.Commands{
				someRepoTree.Commands[0],
				someRepoTree.Commands[1],
			},
		},
	}

	type in struct {
		repo formula.RepositoryLister
		file FileReadExisterMock
	}

	tests := []struct {
		name         string
		in           in
		wantErr      bool
		expectedTree map[string]formula.Tree
	}{
		{
			name: "run in sucess",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name:     "someRepo",
								Provider: "Github",
								Url:      "https://github.com/owner/repo",
								Token:    "token",
							},
						}, nil
					},
				},
				file: FileReadExisterMock{
					exists: func(path string) bool {
						if strings.Contains(path, "/someRepo/tree.json") {
							return true
						}
						return false
					},
					read: func(path string) ([]byte, error) {
						if strings.Contains(path, "/someRepo/tree.json") {
							return []byte(getStringOfTree(someRepoTree)), nil
						}
						return []byte("some data"), nil
					},
				},
			},
			wantErr:      false,
			expectedTree: expectedTreeComplete,
		},
		{
			name: "return error when repository lister returns error",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, errFoo
					},
				},
				file: FileReadExisterMock{
					exists: func(path string) bool {
						return false
					},
					read: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
				},
			},
			wantErr:      true,
			expectedTree: expectedTreeEmpty,
		},
		{
			name: "return error when local tree in read returns error",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, errFoo
					},
				},
				file: FileReadExisterMock{
					exists: func(path string) bool {
						return true
					},
					read: func(path string) ([]byte, error) {
						return []byte("some data"), errFoo
					},
				},
			},
			wantErr:      true,
			expectedTree: expectedTreeEmpty,
		},
		{
			name: "return error when local tree in read returns error",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, errFoo
					},
				},
				file: FileReadExisterMock{
					exists: func(path string) bool {
						return true
					},
					read: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
				},
			},
			wantErr:      true,
			expectedTree: expectedTreeEmpty,
		},
		{
			name: "return error when tree by repo returns error",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name:     "someRepo",
								Provider: "Github",
								Url:      "https://github.com/owner/repo",
							}}, nil
					},
				},
				file: FileReadExisterMock{
					exists: func(path string) bool {
						if strings.Contains(path, "/someRepo/tree.json") {
							return true
						}
						return false
					},
					read: func(path string) ([]byte, error) {
						if strings.Contains(path, "/someRepo/tree.json") {
							return []byte("some"), errFoo
						}
						return []byte("some data"), nil
					},
				},
			},
			wantErr:      true,
			expectedTree: expectedTreeEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			newTree := NewTreeManager(ritHome, in.repo, coreCmds, in.file)

			tree, err := newTree.Tree()

			if !isSameTree(tree, tt.expectedTree) {
				t.Errorf("NewTreeManager_Tree() \n\ttree = %v\n\texpected = %v", tree, tt.expectedTree)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTreeManager_Tree() \n\terror = %v\n\twantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func cleanRitHome() {
	_ = os.RemoveAll(ritHome)
}

func isSameTree(tree, expected map[string]formula.Tree) bool {
	if len(tree) != len(expected) {
		return false
	}
	for i, v := range tree {
		if !isSameFormulaTree(v, expected[i]) {
			return false
		}
	}
	return true
}

func isSameFormulaTree(formula, expected formula.Tree) bool {
	if len(formula.Commands) != len(expected.Commands) {
		return false
	}
	for i, v := range expected.Commands {
		commandsExists := formula.Commands[i] != api.Command{}
		if !commandsExists {
			return false
		}
		if !isSameCommand(v, formula.Commands[i]) {
			return false
		}
	}
	return true
}

func isSameCommand(command, expected api.Command) bool {
	var (
		idIsDiff       = command.Id != expected.Id
		parentIsDiff   = command.Parent != expected.Parent
		usageIsDiff    = command.Usage != expected.Usage
		helpIsDiff     = command.Help != expected.Help
		longHelpIsDiff = command.LongHelp != expected.LongHelp
		formulaIsDiff  = command.Formula != expected.Formula
		repoIsDiff     = command.Repo != expected.Repo
	)

	if idIsDiff || parentIsDiff || usageIsDiff || helpIsDiff || longHelpIsDiff || formulaIsDiff || repoIsDiff {
		return false
	}

	return true
}

func getStringOfTree(formula formula.Tree) string {
	bytes, _ := json.MarshalIndent(formula, "", "\t")
	return string(bytes)
}

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
