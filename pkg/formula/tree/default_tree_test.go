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
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
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

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	workspacePath := filepath.Join(ritHome, "repos", "someRepo1")
	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(workspacePath)

	formulaPath := filepath.Join(ritHome, "repos", "someRepo1", "testing", "formula")

	zipFile := filepath.Join("..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, workspacePath)

	defaultTreeManager := NewGenerator(dirManager, fileManager)
	builderManager := builder.NewBuildLocal(ritHome, dirManager, fileManager, defaultTreeManager)
	_ = builderManager.Build(workspacePath, formulaPath)

	repoLister := repositoryListerCustomMock{
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
	}

	githubRepo := github.NewRepoManager(http.DefaultClient)
	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: githubRepo, NewRepoInfo: github.NewRepoInfo})

	newTree := NewTreeManager(ritHome, repoLister, coreCmds)
	mergedTree := newTree.MergedTree(true)

	nullTree := formula.Tree{Commands: []api.Command{}}
	if len(mergedTree.Commands) == len(nullTree.Commands) {
		t.Errorf("NewTreeManager_MergedTree() mergedTree = %v", mergedTree)
	}
}

func TestTree(t *testing.T) {
	defer cleanRitHome()

	errFoo := errors.New("some error")
	expectedTreeEmpty := map[string]formula.Tree{}
	expectedTreeComplete := map[string]formula.Tree{
		"CORE": {
			Commands: api.Commands{
				{
					Id:       "",
					Parent:   "root",
					Usage:    "add",
					Help:     "",
					LongHelp: "",
					Formula:  false,
					Repo:     "",
				},
				{
					Id:       "",
					Parent:   "root_add",
					Usage:    "repo",
					Help:     "",
					LongHelp: "",
					Formula:  false,
					Repo:     "",
				},
				{
					Id:       "",
					Parent:   "root",
					Usage:    "metrics",
					Help:     "",
					LongHelp: "",
					Formula:  false,
					Repo:     "",
				},
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
	fileManager := stream.NewFileManager()
	_ = fileManager.Remove(ritHome)
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
