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
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

var (
	tmpDir  = os.TempDir()
	ritHome = filepath.Join(tmpDir, ".rit-tree")

	defaultGitRepositoryMock = GitRepositoryMock{
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
)

func TestMergedTree(t *testing.T) {
	defer os.Remove(ritHome)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	workspacePath := filepath.Join(ritHome, "repos", "someRepo1")
	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(workspacePath)
	defer os.Remove(workspacePath)

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

	newTree := NewTreeManager(ritHome, repoLister, api.CoreCmds)
	mergedTree := newTree.MergedTree(true)

	nullTree := formula.Tree{Commands: []api.Command{}}
	if len(mergedTree.Commands) == len(nullTree.Commands) {
		t.Errorf("NewTreeManager_MergedTree() mergedTree = %v", mergedTree)
	}
}

type repositoryListerCustomMock struct {
	list func() (formula.Repos, error)
}

func (m repositoryListerCustomMock) List() (formula.Repos, error) {
	return m.list()
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
