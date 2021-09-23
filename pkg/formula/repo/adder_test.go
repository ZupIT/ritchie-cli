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
	"io"
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

var (
	testDataPath = filepath.Join("..", "..", "..", "testdata", "repos")
)

func TestAdd(t *testing.T) {
	ritHome := filepath.Join(os.TempDir(), ".rit_add_repo")

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
					Name:     "test",
					Version:  "1.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
			},
		},
		{
			name: "success add same repo",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "default",
					Version:  "1.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
			},
		},
		{
			name: "error create repo",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "test",
					Version:  "1.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				mock:          true,
				latestTag:     "2.0.0",
				createRepoErr: errors.New("error to create repo"),
			},
			want: errors.New("error to create repo"),
		},
		{
			name: "error list repos",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "test",
					Version:  "1.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				mock:        true,
				latestTag:   "2.0.0",
				listRepoErr: errors.New("error to list repos"),
			},
			want: errors.New("error to list repos"),
		},
		{
			name: "error write repos",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "test",
					Version:  "1.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				mock:         true,
				latestTag:    "2.0.0",
				writeRepoErr: errors.New("error to write repos"),
			},
			want: errors.New("error to write repos"),
		},
		{
			name: "error tree generation",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "test",
					Version:  "1.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				mock:       true,
				latestTag:  "2.0.0",
				treeGenErr: errors.New("error to generate tree.json"),
			},
			want: errors.New("error to generate tree.json"),
		},
		{
			name: "error invalid repo",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "test",
					Version:  "1.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				mock:      true,
				latestTag: "2.0.0",
				treeGen:   formula.Tree{},
			},
			want: ErrInvalidRepo,
		},
		{
			name: "error to delete invalid repo",
			in: in{
				ritHome: ritHome,
				repo: formula.Repo{
					Provider: "Github",
					Name:     "test",
					Version:  "1.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				mock:          true,
				latestTag:     "2.0.0",
				deleteRepoErr: errors.New("error to delete invalid repo"),
				treeGen:       formula.Tree{},
			},
			want: errors.New("error to delete invalid repo"),
		},
		{
			name: "error to write tree.json",
			in: in{
				ritHome: "invalid",
				repo: formula.Repo{
					Provider: "Github",
					Name:     "test",
					Version:  "1.0.0",
					Url:      "https://github.com/ZupIT/ritchie-cli",
					Priority: 0,
				},
				mock:      true,
				latestTag: "2.0.0",
				treeGen: formula.Tree{
					Commands: api.Commands{
						"root_test": api.Command{},
					},
				},
			},
			want: errors.New("open invalid/repos/test/tree.json: no such file or directory"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := filepath.Join(ritHome, "repos", tt.in.repo.Name.String())
			var repoAdder formula.RepositoryAdder
			if !tt.in.mock {
				repoAdder = addRepoSetup(ritHome, repoPath)
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

				repoAdder = NewAdder(tt.in.ritHome, repoManagerMock, treeManager)
			}

			got := repoAdder.Add(tt.in.repo)

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
				assert.Equal(t, formula.RepoVersion("2.0.0"), repo.LatestVersion)
				assert.NotEmpty(t, repo.Cache)
				assert.FileExists(t, reposPath)
				assert.FileExists(t, filepath.Join(repoPath, "tree.json"))
			}
		})
	}
}

func addRepoSetup(ritHome, repoPath string) formula.RepositoryAdder {
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
			LatestVersion: "2.0.0",
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
	githubRepo.On("LatestTag", mock.Anything).Return(git.Tag{Name: "2.0.0"}, nil)

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

	return NewAdder(ritHome, repoManager, treeGen)
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

func TestIsTemplateRepo(t *testing.T) {
	formulasPath := filepath.Join(testDataPath, "formulas")
	tplFormPath := filepath.Join(testDataPath, "tplForm")
	tplPath := filepath.Join(testDataPath, "commons")
	tplNotValidPath := filepath.Join(testDataPath, "tplNotValid")

	type in struct {
		path string
	}

	type want struct {
		check bool
		err   error
	}

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			name: "is not a template Repo 1",
			in: in{
				path: formulasPath,
			},
			want: want{
				check: false,
				err:   nil,
			},
		},
		{
			name: "is not a template Repo 2",
			in: in{
				path: tplFormPath,
			},
			want: want{
				check: false,
				err:   nil,
			},
		},
		{
			name: "is a template repo",
			in: in{
				path: tplPath,
			},
			want: want{
				check: true,
				err:   nil,
			},
		},
		{
			name: "error: has templates but is not valid",
			in: in{
				path: tplNotValidPath,
			},
			want: want{
				check: true,
				err:   ErrInvalidRepo,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isTemplateRepo(tt.in.path)
			assert.Equal(t, tt.want.check, got)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestIsValidTemplateRepo(t *testing.T) {
	tplPath := filepath.Join(testDataPath, "commons", "templates")
	tplLangPath := filepath.Join(testDataPath, "tplLang", "templates")
	tplNotValidPath := filepath.Join(testDataPath, "tplNotValid", "templates")

	type in struct {
		path string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "is template repo",
			in: in{
				path: tplPath,
			},
			want: nil,
		},
		{
			name: "error: is not a valid template repo",
			in: in{
				path: tplLangPath,
			},
			want: ErrInvalidRepo,
		},
		{
			name: "error: does not have templates",
			in: in{
				path: tplNotValidPath,
			},
			want: ErrInvalidRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidTemplateRepo(tt.in.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHasTemplates(t *testing.T) {
	tplPath := filepath.Join(testDataPath, "tplValid", "templates", "create_formula", "languages")

	type in struct {
		lang string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "has templates 1",
			in: in{
				lang: "rust",
			},
			want: nil,
		},
		{
			name: "has templates 2",
			in: in{
				lang: "csharp",
			},
			want: nil,
		},
		{
			name: "error: does not have templates",
			in: in{
				lang: "crystal",
			},
			want: ErrInvalidRepo,
		},
		{
			name: "error: has a no valid template 1",
			in: in{
				lang: "go",
			},
			want: ErrInvalidRepo,
		},
		{
			name: "error: has a no valid template 2",
			in: in{
				lang: "perl",
			},
			want: ErrInvalidRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasTemplates(tplPath, tt.in.lang)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsValidTemplate(t *testing.T) {
	tplPath := filepath.Join(testDataPath, "tplValid", "templates", "create_formula", "languages")

	type in struct {
		lang string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "is valid template",
			in: in{
				lang: "rust",
			},
			want: nil,
		},
		{
			name: "error on path",
			in: in{
				lang: "java8",
			},
			want: ErrInvalidRepo,
		},
		{
			name: "error: does not have the 'build.bat' file",
			in: in{
				lang: "ruby",
			},
			want: ErrInvalidRepo,
		},
		{
			name: "error: does not have the 'build.sh' file",
			in: in{
				lang: "go",
			},
			want: ErrInvalidRepo,
		},
		{
			name: "error: does not have the 'build.bat' and the 'build.sh' files",
			in: in{
				lang: "java11",
			},
			want: ErrInvalidRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := filepath.Join(tplPath, tt.in.lang, "helloWorld")
			got := isValidTemplate(fullPath)
			assert.Equal(t, tt.want, got)
		})
	}
}
