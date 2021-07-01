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

package builder

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestBuild(t *testing.T) {
	workspacePath := filepath.Join(tmpDir, "ritchie-formulas-test")
	formulaPath := filepath.Join(tmpDir, "ritchie-formulas-test", "testing", "formula")
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	defaultTreeManager := tree.NewGenerator(dirManager, fileManager)

	repoProviders := formula.NewRepoProviders()
	repoCreator := repo.NewCreator(ritHome, repoProviders, dirManager, fileManager)
	repoLister := repo.NewLister(ritHome, fileManager)
	repoWriter := repo.NewWriter(ritHome, fileManager)
	repoListWrite := repo.NewListWriter(repoLister, repoWriter)
	repoDeleter := repo.NewDeleter(ritHome, repoListWrite, dirManager)
	repoDetail := repo.NewDetail(repoProviders)
	repoListWriteCreator := repo.NewCreateWriteListDetailDeleter(repoLister, repoCreator, repoWriter, repoDetail, repoDeleter)
	repoAdder := repo.NewAdder(ritHome, repoListWriteCreator, defaultTreeManager)

	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(workspacePath)

	zipFile := filepath.Join("..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, workspacePath)

	type in struct {
		formulaPath string
		dirManager  stream.DirCreateListCopyRemover
		repo        formula.RepositoryAdder
	}

	testes := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				formulaPath: formulaPath,
				dirManager:  dirManager,
				repo:        repoAdder,
			},
			want: nil,
		},
		{
			name: "success build without build.sh",
			in: in{
				formulaPath: filepath.Join(tmpDir, "ritchie-formulas-test", "testing", "without-build-sh"),
				dirManager:  dirManager,
				repo:        repoAdder,
			},
			want: nil,
		},
		{
			name: "create dir error",
			in: in{
				formulaPath: formulaPath,
				dirManager:  dirManagerMock{createErr: errors.New("error to create dir")},
			},
			want: errors.New("error to create dir"),
		},
		{
			name: "copy workspace dir error",
			in: in{
				formulaPath: formulaPath,
				dirManager:  dirManagerMock{data: []string{"linux"}, copyErr: errors.New("error to copy dir")},
			},
			want: errors.New("error to copy dir"),
		},
		{
			name: "dir remove error",
			in: in{
				formulaPath: formulaPath,
				dirManager:  dirManagerMock{data: []string{"commons"}, removeErr: errors.New("error to remove dir")},
				repo:        repoAdder,
			},
			want: errors.New("error to remove dir"),
		},
		{
			name: "repo add error",
			in: in{
				formulaPath: formulaPath,
				dirManager:  dirManager,
				repo:        repoAdderMock{err: errors.New("error to add repo")},
			},
			want: errors.New("error to add repo"),
		},
	}

	for _, tt := range testes {
		t.Run(tt.name, func(t *testing.T) {
			builderManager := NewBuildLocal(ritHome, tt.in.dirManager, tt.in.repo)
			info := formula.BuildInfo{FormulaPath: tt.in.formulaPath, Workspace: formula.Workspace{Name: "repo", Dir: workspacePath}}
			got := builderManager.Build(info)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Build(%s) got %v, want %v", tt.name, got, tt.want)
			}

			if tt.want == nil {
				hasRitchieHome := dirManager.Exists(ritHome)
				if !hasRitchieHome {
					t.Errorf("Build(%s) did not create the Ritchie home directory", tt.name)
				}

				treeLocalFile := filepath.Join(ritHome, "repos", "local-repo", "tree.json")
				hasTreeLocalFile := fileManager.Exists(treeLocalFile)
				if !hasTreeLocalFile {
					t.Errorf("Build(%s) did not copy the tree local file", tt.name)
				}

				formulaFiles := filepath.Join(ritHome, "repos", "local-repo", "testing", "formula", "bin")
				files, err := fileManager.List(formulaFiles)
				if err == nil && len(files) != 4 {
					t.Errorf("Build(%s) did not generate bin files", tt.name)
				}

				configFile := filepath.Join(ritHome, "repos", "local-repo", "testing", "formula", "config.json")
				hasConfigFile := fileManager.Exists(configFile)
				if !hasConfigFile {
					t.Errorf("Build(%s) did not copy formula config", tt.name)
				}
			}
		})
	}
}

type dirManagerMock struct {
	data      []string
	createErr error
	listErr   error
	copyErr   error
	removeErr error
}

func (d dirManagerMock) Create(string) error {
	return d.createErr
}

func (d dirManagerMock) List(string, bool) ([]string, error) {
	return d.data, d.listErr
}

func (d dirManagerMock) Copy(string, string) error {
	return d.copyErr
}

func (d dirManagerMock) Remove(string) error {
	return d.removeErr
}

type repoAdderMock struct {
	err error
}

func (r repoAdderMock) Add(repo formula.Repo) error {
	return r.err
}
