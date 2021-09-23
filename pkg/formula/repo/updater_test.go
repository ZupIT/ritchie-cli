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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestUpdateManager_Update(t *testing.T) {
	ritHome := filepath.Join(os.TempDir(), ".rit_update_repo")

	type in struct {
		ritHome       string
		repo          formula.Repo
		mock          bool
		latestTag     string
		createRepoErr error
		listRepos     formula.Repos
		listRepoErr   error
		writeRepoErr  error
		deleteRepoErr error
		treeGen       formula.Tree
		treeGenErr    error
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "default",
					Version:  "3.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
			},
		},
		{
			name: "error to list repos",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "default",
					Version:  "3.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				listRepoErr: errors.New("error to list repos"),
				mock:        true,
			},
			want: errors.New("error to list repos"),
		},
		{
			name: "repository not found",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "default",
					Version:  "3.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				listRepos: formula.Repos{},
				mock:      true,
			},
			want: errors.New("repository name \"default\" was not found"),
		},
		{
			name: "error local repo",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "default",
					Version:  "3.0.0",
					Priority: 0,
					IsLocal:  true,
				},
				listRepos: formula.Repos{
					formula.Repo{
						Name:    "default",
						IsLocal: true,
					},
				},
				mock: true,
			},
			want: ErrLocalRepo,
		},
		{
			name: "error to create repo",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "default",
					Version:  "3.0.0",
					Priority: 0,
				},
				listRepos: formula.Repos{
					formula.Repo{
						Name: "default",
					},
				},
				latestTag:     "3.0.0",
				createRepoErr: errors.New("error to create repo"),
				mock:          true,
			},
			want: errors.New("error to create repo"),
		},
		{
			name: "error to write repo",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "default",
					Version:  "3.0.0",
					Priority: 0,
				},
				listRepos: formula.Repos{
					formula.Repo{
						Name: "default",
					},
				},
				latestTag:    "3.0.0",
				writeRepoErr: errors.New("error to write repo"),
				mock:         true,
			},
			want: errors.New("error to write repo"),
		},
		{
			name: "error to generate tree.json",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "default",
					Version:  "3.0.0",
					Priority: 0,
				},
				listRepos: formula.Repos{
					formula.Repo{
						Name: "default",
					},
				},
				latestTag:  "3.0.0",
				treeGenErr: errors.New("error to generate tree"),
				mock:       true,
			},
			want: errors.New("error to generate tree"),
		},
		{
			name: "error to write tree.json",
			in: in{
				ritHome: "invalid",
				repo: formula.Repo{
					Provider: "Github",
					Name:     "default",
					Version:  "3.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				mock:      true,
				latestTag: "2.0.0",
				listRepos: formula.Repos{
					formula.Repo{
						Name: "default",
					},
				},
				treeGen: formula.Tree{
					Commands: api.Commands{
						"root_test": api.Command{},
					},
				},
			},
			want: errors.New("open invalid/repos/default/tree.json: no such file or directory"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := filepath.Join(ritHome, "repos", tt.in.repo.Name.String())
			var repoAdder formula.RepositoryUpdater
			if !tt.in.mock {
				repoAdder = updateRepoSetup(ritHome, repoPath)
				defer os.RemoveAll(ritHome)
			} else {
				repoManagerMock := &mocks.RepoManager{}
				repoManagerMock.On("LatestTag", mock.Anything).Return(tt.in.latestTag)
				repoManagerMock.On("Create", mock.Anything).Return(tt.in.createRepoErr)
				repoManagerMock.On("List").Return(tt.in.listRepos, tt.in.listRepoErr)
				repoManagerMock.On("Write", mock.Anything).Return(tt.in.writeRepoErr)
				repoManagerMock.On("Delete", mock.Anything).Return(tt.in.deleteRepoErr)

				treeManager := &mocks.TreeManager{}
				treeManager.On("Generate", mock.Anything).Return(tt.in.treeGen, tt.in.treeGenErr)

				repoAdder = NewUpdater(tt.in.ritHome, repoManagerMock, treeManager)
			}

			got := repoAdder.Update(tt.in.repo.Name, tt.in.repo.Version, tt.in.repo.Url)

			if got != nil {
				assert.EqualError(t, tt.want, got.Error())
			} else {
				assert.Nil(t, tt.want)
			}

			if !tt.in.mock {
				reposPath := filepath.Join(ritHome, "repos", "repositories.json")
				file, _ := ioutil.ReadFile(reposPath)

				var repos formula.Repos
				_ = json.Unmarshal(file, &repos)
				repo := repos[0]
				expectRepo := tt.in.repo

				assert.Equal(t, expectRepo.Provider, repo.Provider)
				assert.Equal(t, expectRepo.Name, repo.Name)
				assert.Equal(t, expectRepo.Version, repo.Version)
				assert.Equal(t, expectRepo.Url, repo.Url)
				assert.Equal(t, expectRepo.Priority, repo.Priority)
				assert.Equal(t, "v2", repo.TreeVersion)
				assert.Equal(t, formula.RepoVersion("3.0.0"), repo.LatestVersion)
				assert.NotEmpty(t, repo.Cache)
				assert.FileExists(t, reposPath)
				assert.FileExists(t, filepath.Join(repoPath, "tree.json"))
			}
		})
	}
}

func updateRepoSetup(ritHome, repoPath string) formula.RepositoryUpdater {
	_ = os.MkdirAll(filepath.Join(repoPath, "test", "test"), os.ModePerm)
	_ = ioutil.WriteFile(filepath.Join(repoPath, "test", "help.json"), []byte("{}"), os.ModePerm)

	defaultRepos := formula.Repos{
		{
			Provider:      "Github",
			Name:          "default",
			Version:       "1.0.0",
			Url:           "https://github.com/ZupIT/ritchie-cli",
			Priority:      0,
			TreeVersion:   "v2",
			LatestVersion: "3.0.0",
			Cache:         time.Now().Add(time.Hour),
		},
	}

	bytes, _ := json.Marshal(defaultRepos)
	reposPath := filepath.Join(ritHome, "repos", "repositories.json")
	_ = ioutil.WriteFile(reposPath, bytes, os.ModePerm)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	repoProviders := formula.NewRepoProviders()

	githubRepo := &mocks.GitRepositoryMock{}
	githubRepo.On("LatestTag", mock.Anything).Return(git.Tag{Name: "3.0.0"}, nil)

	gitProvider := formula.Git{Repos: githubRepo, NewRepoInfo: github.NewRepoInfo}
	repoProviders.Add("Github", gitProvider)

	repoCreator := &mocks.RepoManager{}
	repoCreator.On("Create", mock.Anything).Return(nil)

	repoLister := NewLister(ritHome, fileManager)
	repoWriter := NewWriter(ritHome, fileManager)
	repoDetail := NewDetail(repoProviders)
	repoListWriter := NewListWriter(repoLister, repoWriter)
	repoDeleter := NewDeleter(ritHome, repoListWriter, dirManager)
	repoManager := NewCreateWriteListDetailDeleter(repoLister, repoCreator, repoWriter, repoDetail, repoDeleter)
	treeGen := tree.NewGenerator(dirManager, fileManager)

	return NewUpdater(ritHome, repoManager, treeGen)
}
