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
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	tmpDir  = os.TempDir()
	ritHome = filepath.Join(tmpDir, ".rit-tree")

	coreCmds = []api.Command{
		{Parent: "root", Usage: "add"},
		{Parent: "root_add", Usage: "repo"},
		{Parent: "root", Usage: "metrics"},
	}
)

func TestMergedTree(t *testing.T) {
	defer cleanRitHome()

	treeLocalDir := filepath.Join(ritHome, "repos", "local")
	treeSomeRepoDir := filepath.Join(ritHome, "repos", "someRepo")

	createDir(ritHome)
	createDir(treeLocalDir)
	createDir(treeSomeRepoDir)

	localTree := formula.Tree{Commands: []api.Command{
		{Parent: "root", Usage: "jedi-list"},
		{Parent: "root_jedi-list", Usage: "add"},
	}}

	someRepoTree := formula.Tree{Commands: []api.Command{
		{Parent: "root", Usage: "pokemon-list"},
		{Parent: "root_pokemon-list", Usage: "add"},
	}}

	addTreeLocal(treeLocalDir, localTree)
	addTreeLocal(treeSomeRepoDir, someRepoTree)

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
			}, nil
		},
	}

	githubRepo := github.NewRepoManager(http.DefaultClient)
	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: githubRepo, NewRepoInfo: github.NewRepoInfo})

	newTree := NewTreeManager(ritHome, repoLister, coreCmds)
	mergedTree := newTree.MergedTree(true)

	if !isSameFormulaTree(mergedTree, expectedTreeComplete) {
		t.Errorf("NewTreeManager_MergedTree() mergedTree = %v, expectedTree = %v", mergedTree, expectedTreeComplete)
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
		"someRepo1": {
			Commands: api.Commands{},
		},
	}

	type in struct {
		repo formula.RepositoryLister
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
								Name:     "someRepo1",
								Provider: "Github",
								Url:      "https://github.com/owner/repo",
								Token:    "token",
							},
						}, nil
					},
				},
			},
			wantErr:      false,
			expectedTree: expectedTreeComplete,
		},
		{
			name: "return error when repository lister resturns error",
			in: in{
				repo: repositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, errFoo
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
			newTree := NewTreeManager(ritHome, in.repo, coreCmds)

			tree, err := newTree.Tree()

			if !isSameTree(tree, tt.expectedTree) {
				t.Errorf("NewTreeManager_Tree() tree = %v", tree)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTreeManager_Tree() error = %v", err)
			}
		})
	}
}

func cleanRitHome() {
	_ = os.RemoveAll(ritHome)
}

func createDir(path string) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	_ = dirManager.Remove(path)
	_ = dirManager.Create(path)
}

func addTreeLocal(dest string, tree formula.Tree) {
	fileManager := stream.NewFileManager()

	treeJSON, _ := json.MarshalIndent(tree, "", "\t")

	treeLocalFile := filepath.Join(dest, "tree.json")
	_ = fileManager.Write(treeLocalFile, treeJSON)
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
	for i, v := range formula.Commands {
		commandsExists := expected.Commands[i] != api.Command{}
		if !commandsExists {
			return false
		}
		if !isSameCommand(v, expected.Commands[i]) {
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

type repositoryListerCustomMock struct {
	list func() (formula.Repos, error)
}

func (m repositoryListerCustomMock) List() (formula.Repos, error) {
	return m.list()
}
