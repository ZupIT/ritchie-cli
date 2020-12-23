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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestAdd(t *testing.T) {
	treeWithOneCommand := formula.Tree{Commands: api.Commands{"": {Parent: "", Usage: "", Help: "", LongHelp: "", Formula: false, Repo: ""}}}

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	type in struct {
		ritHome  string
		repoMock formula.RepositoryListWriteCreator
		deleter  formula.RepositoryDeleter
		tree     formula.TreeGenerator
		file     stream.FileWriteCreatorReadExistRemover
		repo     formula.Repo
	}
	tests := []struct {
		name    string
		in      in
		wantErr error
	}{
		{
			name: "Run with success, when repository json not exist",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-adder-test-success")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Remove(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					return ritHomePath
				}(),
				repoMock: repoListWriteCreatorMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
					write: func(repos formula.Repos) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return treeWithOneCommand, nil
					},
				},
				file: fileManager,
				repo: formula.Repo{
					Name:     "some_repo_name",
					Priority: 10,
					Token:    "",
					Url:      "https://github.com/someUser/ritchie-formulas",
					Version:  "2.0",
				},
			},
		},
		{
			name: "Run with success, when repository json exist",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-adder-test-success")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Remove(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					repoFileData := `
					[
							{
								"name": "commons",
								"version": "v2.0.0",
								"url": "https://github.com/kaduartur/ritchie-formulas",
								"priority": 0
							}
					]
					`
					_ = fileManager.Write(filepath.Join(ritHomePath, "repos", reposFileName), []byte(repoFileData))
					return ritHomePath
				}(),
				repoMock: repoListWriteCreatorMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Provider: "Local",
								Name:     "repo-local",
								Version:  "0.0.0",
								Priority: 0,
								IsLocal:  true,
							},
						}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
					write: func(repos formula.Repos) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return treeWithOneCommand, nil
					},
				},
				file: fileManager,
				repo: formula.Repo{
					Name:     "some_repo_name",
					Priority: 10,
					Token:    "",
					Url:      "https://github.com/someUser/ritchie-formulas",
					Version:  "2.0",
				},
			},
		},
		{
			name: "Return err when RepositoryCreator fail",
			in: in{
				repoMock: repoListWriteCreatorMock{
					create: func(repo formula.Repo) error {
						return errors.New("error to create repo")
					},
				},
			},
			wantErr: errors.New("error to create repo"),
		},
		{
			name: "Return err when repos list fail",
			in: in{
				repoMock: repoListWriteCreatorMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, errors.New("error to read repos file")
					},
					create: func(repo formula.Repo) error {
						return nil
					},
					write: func(repos formula.Repos) error {
						return nil
					},
				},
			},
			wantErr: errors.New("error to read repos file"),
		},
		{
			name: "Return err when saveRepo fail",
			in: in{
				repoMock: repoListWriteCreatorMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
					write: func(repos formula.Repos) error {
						return errors.New("error to write repo file")
					},
				},
			},
			wantErr: errors.New("error to write repo file"),
		},
		{
			name: "Return err when tree Generate fail",
			in: in{
				repoMock: repoListWriteCreatorMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
					write: func(repos formula.Repos) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					exists: func(path string) bool {
						return false
					},
					write: func(path string, content []byte) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, errors.New("error to generate tree.json")
					},
				},
			},
			wantErr: errors.New("error to generate tree.json"),
		},
		{
			name: "Return err when write tree.json",
			in: in{
				repoMock: repoListWriteCreatorMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
					write: func(repos formula.Repos) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					write: func(path string, content []byte) error {
						return errors.New("error to write tree.json")
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return treeWithOneCommand, nil
					},
				},
			},
			wantErr: errors.New("error to write tree.json"),
		},
		{
			name: "Return err when tree.json has no commands",
			in: in{
				repoMock: repoListWriteCreatorMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
					write: func(repos formula.Repos) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, nil
					},
				},
				deleter: repositoryDeleterMock{
					deleteMock: func(repoName formula.RepoName) error {
						return nil
					},
				},
			},
			wantErr: ErrInvalidRepo,
		},
		{
			name: "Return err when Delete fail",
			in: in{
				repoMock: repoListWriteCreatorMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
					write: func(repos formula.Repos) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, nil
					},
				},
				deleter: repositoryDeleterMock{
					deleteMock: func(repoName formula.RepoName) error {
						return errors.New("error to delete repository")
					},
				},
			},
			wantErr: errors.New("error to delete repository"),
		},
		{
			name: "success when add a new empty local repository",
			in: in{
				repo: formula.Repo{
					Provider: "Local",
					Name:     formula.RepoName("my-repo"),
					Version:  "0.0.0",
					Url:      "local repository",
					Priority: 0,
					IsLocal:  true,
				},
				repoMock: repoListWriteCreatorMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
					write: func(repos formula.Repos) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					write: func(path string, content []byte) error {
						return nil
					},
				},
				deleter: repositoryDeleterMock{
					deleteMock: func(repoName formula.RepoName) error {
						return nil
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad := NewAdder(
				tt.in.ritHome,
				tt.in.repoMock,
				tt.in.deleter,
				tt.in.tree,
				tt.in.file,
			)

			got := ad.Add(tt.in.repo)
			assert.Equal(t, tt.wantErr, got)
		})
	}
}

type treeGeneratorCustomMock struct {
	generate func(repoPath string) (formula.Tree, error)
}

func (m treeGeneratorCustomMock) Generate(repoPath string) (formula.Tree, error) {
	return m.generate(repoPath)
}

type FileWriteCreatorReadExistRemover struct {
	write  func(path string, content []byte) error
	create func(path string, data io.ReadCloser) error
	read   func(path string) ([]byte, error)
	exists func(path string) bool
	remove func(path string) error
}

func (m FileWriteCreatorReadExistRemover) Write(path string, content []byte) error {
	return m.write(path, content)
}

func (m FileWriteCreatorReadExistRemover) Create(path string, data io.ReadCloser) error {
	return m.create(path, data)
}

func (m FileWriteCreatorReadExistRemover) Read(path string) ([]byte, error) {
	return m.read(path)
}

func (m FileWriteCreatorReadExistRemover) Exists(path string) bool {
	return m.exists(path)
}

func (m FileWriteCreatorReadExistRemover) Remove(path string) error {
	return m.remove(path)
}

type DirCreateListCopyRemoverCustomMock struct {
	create func(dir string) error
	list   func(dir string, hiddenDir bool) ([]string, error)
	copy   func(src, dst string) error
	remove func(dir string) error
}

func (m DirCreateListCopyRemoverCustomMock) Create(dir string) error {
	return m.create(dir)
}

func (m DirCreateListCopyRemoverCustomMock) List(dir string, hiddenDir bool) ([]string, error) {
	return m.list(dir, hiddenDir)
}

func (m DirCreateListCopyRemoverCustomMock) Copy(src, dst string) error {
	return m.copy(src, dst)
}

func (m DirCreateListCopyRemoverCustomMock) Remove(dir string) error {
	return m.remove(dir)
}

type repositoryDeleterMock struct {
	deleteMock func(repoName formula.RepoName) error
}

func (c repositoryDeleterMock) Delete(repoName formula.RepoName) error {
	return c.deleteMock(repoName)
}
