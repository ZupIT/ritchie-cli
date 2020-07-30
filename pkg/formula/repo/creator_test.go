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
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestCreate(t *testing.T) {

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	type in struct {
		ritHome string
		github  github.Repositories
		dir     stream.DirCreateListCopyRemover
		file    stream.FileWriteCreatorReadExistRemover
		repo    formula.Repo
	}
	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "TestCreateManager_Create_with_success")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Remove(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					return ritHomePath
				}(),
				github: GitRepositoryMock{
					zipball: func(info github.RepoInfo, version string) (io.ReadCloser, error) {
						data, _ := fileManager.Read("../../../testdata/ritchie-formulas.zip")
						return ioutil.NopCloser(bytes.NewReader(data)), nil
					},
				},
				dir:  dirManager,
				file: fileManager,
				repo: formula.Repo{
					Name:     "testing_repo",
					Version:  "0.0.3",
					Url:      "https://github.com/viniciussousazup/ritchie-formulas/releases",
					Token:    "",
					Priority: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "Return err when zipball fail",
			in: in{
				ritHome: "",
				github: GitRepositoryMock{
					zipball: func(info github.RepoInfo, version string) (io.ReadCloser, error) {
						return nil, errors.New("some error")
					},
				},
				dir:  dirManager,
				file: fileManager,
				repo: formula.Repo{
					Name:     "testing_repo",
					Version:  "0.0.3",
					Url:      "https://github.com/viniciussousazup/ritchie-formulas/releases",
					Token:    "",
					Priority: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when dir remove fail",
			in: in{
				ritHome: "",
				github: GitRepositoryMock{
					zipball: func(info github.RepoInfo, version string) (io.ReadCloser, error) {
						data, _ := fileManager.Read("../../../testdata/ritchie-formulas.zip")
						return ioutil.NopCloser(bytes.NewReader(data)), nil
					},
				},
				dir: DirCreateListCopyRemoverCustomMock{
					remove: func(dir string) error {
						return errors.New("some error")
					},
				},
				file: fileManager,
				repo: formula.Repo{
					Name:     "testing_repo",
					Version:  "0.0.3",
					Url:      "https://github.com/viniciussousazup/ritchie-formulas/releases",
					Token:    "",
					Priority: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when dir create fail",
			in: in{
				ritHome: "",
				github: GitRepositoryMock{
					zipball: func(info github.RepoInfo, version string) (io.ReadCloser, error) {
						data, _ := fileManager.Read("../../../testdata/ritchie-formulas.zip")
						return ioutil.NopCloser(bytes.NewReader(data)), nil
					},
				},
				dir: DirCreateListCopyRemoverCustomMock{
					remove: func(dir string) error {
						return nil
					},
					create: func(dir string) error {
						return errors.New("some error")
					},
				},
				file: fileManager,
				repo: formula.Repo{
					Name:     "testing_repo",
					Version:  "0.0.3",
					Url:      "https://github.com/viniciussousazup/ritchie-formulas/releases",
					Token:    "",
					Priority: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when file create fail",
			in: in{
				ritHome: "",
				github: GitRepositoryMock{
					zipball: func(info github.RepoInfo, version string) (io.ReadCloser, error) {
						data, _ := fileManager.Read("../../../testdata/ritchie-formulas.zip")
						return ioutil.NopCloser(bytes.NewReader(data)), nil
					},
				},
				dir: DirCreateListCopyRemoverCustomMock{
					remove: func(dir string) error {
						return nil
					},
					create: func(dir string) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					create: func(path string, data io.ReadCloser) error {
						return errors.New("some error")
					},
				},
				repo: formula.Repo{
					Name:     "testing_repo",
					Version:  "0.0.3",
					Url:      "https://github.com/viniciussousazup/ritchie-formulas/releases",
					Token:    "",
					Priority: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when file remove fail",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "TestCreateManager_Create_fail_remove_File")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Remove(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					return ritHomePath
				}(),
				github: GitRepositoryMock{
					zipball: func(info github.RepoInfo, version string) (io.ReadCloser, error) {
						data, _ := fileManager.Read("../../../testdata/ritchie-formulas.zip")
						return ioutil.NopCloser(bytes.NewReader(data)), nil
					},
				},
				dir: dirManager,
				file: FileWriteCreatorReadExistRemover{
					create: fileManager.Create,
					remove: func(path string) error {
						return errors.New("some error")
					},
				},
				repo: formula.Repo{
					Name:     "testing_repo",
					Version:  "0.0.3",
					Url:      "https://github.com/viniciussousazup/ritchie-formulas/releases",
					Token:    "",
					Priority: 0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := NewCreator(
				tt.in.ritHome,
				tt.in.github,
				tt.in.dir,
				tt.in.file,
			)
			if err := cr.Create(tt.in.repo); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				newRepoPath := filepath.Join(tt.in.ritHome, "repos", tt.in.repo.Name.String())
				readMePath := filepath.Join(newRepoPath, "README.md")
				if !fileManager.Exists(readMePath) {
					t.Errorf("ReadMe not exist on path %s ", readMePath)
				}

			}
		})
	}
}

type GitRepositoryMock struct {
	zipball   func(info github.RepoInfo, version string) (io.ReadCloser, error)
	tags      func(info github.RepoInfo) (github.Tags, error)
	latestTag func(info github.RepoInfo) (github.Tag, error)
}

func (m GitRepositoryMock) Zipball(info github.RepoInfo, version string) (io.ReadCloser, error) {
	return m.zipball(info, version)
}

func (m GitRepositoryMock) Tags(info github.RepoInfo) (github.Tags, error) {
	return m.tags(info)
}

func (m GitRepositoryMock) LatestTag(info github.RepoInfo) (github.Tag, error) {
	return m.latestTag(info)
}
