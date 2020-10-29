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

	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(workspacePath)

	zipFile := filepath.Join("..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, workspacePath)

	type in struct {
		formulaPath string
		fileManager stream.FileWriteReadExister
		dirManager  stream.DirCreateListCopyRemover
		tree        formula.TreeGenerator
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
				fileManager: fileManager,
				dirManager:  dirManager,
				tree:        defaultTreeManager,
			},
			want: nil,
		},
		{
			name: "success build without build.sh",
			in: in{
				formulaPath: filepath.Join(tmpDir, "ritchie-formulas-test", "testing", "without-build-sh"),
				fileManager: fileManager,
				dirManager:  dirManager,
				tree:        defaultTreeManager,
			},
			want: nil,
		},
		{
			name: "create dir error",
			in: in{
				formulaPath: formulaPath,
				fileManager: fileManager,
				dirManager:  dirManagerMock{createErr: errors.New("error to create dir")},
				tree:        defaultTreeManager,
			},
			want: errors.New("error to create dir"),
		},
		{
			name: "copy so dir error",
			in: in{
				formulaPath: formulaPath,
				fileManager: fileManager,
				dirManager:  dirManagerMock{data: []string{"linux"}, copyErr: errors.New("error to copy dir")},
				tree:        defaultTreeManager,
			},
			want: errors.New("error to copy dir"),
		},
		{
			name: "copy commons dir error",
			in: in{
				formulaPath: formulaPath,
				fileManager: fileManager,
				dirManager:  dirManagerMock{data: []string{"commons"}, copyErr: errors.New("error to copy dir")},
				tree:        defaultTreeManager,
			},
			want: errors.New("error to copy dir"),
		},
		{
			name: "dir remove error",
			in: in{
				formulaPath: formulaPath,
				fileManager: fileManager,
				dirManager:  dirManagerMock{data: []string{"commons"}, removeErr: errors.New("error to remove dir")},
				tree:        defaultTreeManager,
			},
			want: errors.New("error to remove dir"),
		},
		{
			name: "tree generate error",
			in: in{
				formulaPath: formulaPath,
				fileManager: fileManager,
				dirManager:  dirManager,
				tree:        treeGenerateMock{err: errors.New("error to generate tree")},
			},
			want: errors.New("error to generate tree"),
		},
		{
			name: "write tree error",
			in: in{
				formulaPath: formulaPath,
				fileManager: fileManagerMock{writeErr: errors.New("error to write tree")},
				dirManager:  dirManager,
				tree:        defaultTreeManager,
			},
			want: errors.New("error to write tree"),
		},
		{
			name: "chdir error",
			in: in{
				formulaPath: "invalid",
				fileManager: fileManager,
				dirManager:  dirManager,
				tree:        defaultTreeManager,
			},
			want: errors.New("chdir invalid: no such file or directory"),
		},
	}

	for _, tt := range testes {
		t.Run(tt.name, func(t *testing.T) {
			builderManager := NewBuildLocal(ritHome, tt.in.dirManager, tt.in.fileManager, tt.in.tree)
			info := formula.BuildInfo{FormulaPath: tt.in.formulaPath, Workspace: formula.Workspace{Dir: workspacePath}}
			got := builderManager.Build(info)

			if (tt.want == nil && got != nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Build(%s) got %v, want %v", tt.name, got, tt.want)
			}

			if tt.want == nil {
				hasRitchieHome := dirManager.Exists(ritHome)
				if !hasRitchieHome {
					t.Errorf("Build(%s) did not create the Ritchie home directory", tt.name)
				}

				treeLocalFile := filepath.Join(ritHome, "repos", "local", "tree.json")
				hasTreeLocalFile := fileManager.Exists(treeLocalFile)
				if !hasTreeLocalFile {
					t.Errorf("Build(%s) did not copy the tree local file", tt.name)
				}

				formulaFiles := filepath.Join(ritHome, "repos", "local", "testing", "formula", "bin")
				files, err := fileManager.List(formulaFiles)
				if err == nil && len(files) != 4 {
					t.Errorf("Build(%s) did not generate bin files", tt.name)
				}

				configFile := filepath.Join(ritHome, "repos", "local", "testing", "formula", "config.json")
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

type fileManagerMock struct {
	data     []byte
	readErr  error
	exist    bool
	writeErr error
}

func (f fileManagerMock) Read(string) ([]byte, error) {
	return f.data, f.readErr
}

func (f fileManagerMock) Exists(string) bool {
	return f.exist
}

func (f fileManagerMock) Write(string, []byte) error {
	return f.writeErr
}

type treeGenerateMock struct {
	tree formula.Tree
	err  error
}

func (t treeGenerateMock) Generate(repoPath string) (formula.Tree, error) {
	return t.tree, t.err
}
