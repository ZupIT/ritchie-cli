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

package workspace

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestWorkspaceManagerAdd(t *testing.T) {
	cleanForm()
	fullDir := createFullDir()

	tmpDir := path.Join(os.TempDir(), "workspace-add")
	_ = os.Mkdir(tmpDir, os.ModePerm)
	defer os.RemoveAll(tmpDir)

	dirManager := stream.NewDirManager(stream.NewFileManager())
	workspaceFile := path.Join(tmpDir, formula.WorkspacesFile)
	workspaceBrokenPath := path.Join(tmpDir, "broken")
	workspaceNonExistingPath := filepath.Join(tmpDir, "non-existing-dir")
	err := os.Mkdir(workspaceBrokenPath, os.ModePerm)
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(workspaceBrokenPath, formula.WorkspacesFile), []byte("error"), os.ModePerm)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		workspacePath string
		workspace     formula.Workspace
		outErr        string
	}{
		{
			name:          "success create",
			workspacePath: tmpDir,
			workspace: formula.Workspace{
				Name: "zup",
				Dir:  fullDir,
			},
		},
		{
			name:          "success create with trailing separator",
			workspacePath: tmpDir,
			workspace: formula.Workspace{
				Name: "zup2",
				Dir:  fullDir + string(filepath.Separator),
			},
		},
		{
			name:          "success edit",
			workspacePath: tmpDir,
			workspace: formula.Workspace{
				Name: "commons",
				Dir:  fullDir,
			},
		},
		{
			name:          "invalid workspace",
			workspacePath: workspaceFile,
			workspace: formula.Workspace{
				Name: "zup",
				Dir:  "home/user/go/src/github.com/ZupIT/ritchie-formulas-commons",
			},
			outErr: ErrInvalidWorkspace.Error(),
		},
		{
			name:          "unmarshal error",
			workspacePath: workspaceBrokenPath,
			workspace: formula.Workspace{
				Name: "commons",
				Dir:  fullDir,
			},
			outErr: "invalid character 'e' looking for beginning of value",
		},
		{
			name:          "write error",
			workspacePath: workspaceNonExistingPath,
			workspace: formula.Workspace{
				Name: "commons",
				Dir:  fullDir,
			},
			outErr: mocks.FileNotFoundError(filepath.Join(workspaceNonExistingPath, formula.WorkspacesFile)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localBuilder := &mocks.LocalBuilderMock{}
			localBuilder.On("Init", mock.AnythingOfType("string"), tt.workspace.Name).Return("", nil)

			workspace := New(tt.workspacePath, tt.workspacePath, dirManager, localBuilder)
			got := workspace.Add(tt.workspace)

			if got != nil {
				assert.EqualError(t, got, tt.outErr)
			} else {
				assert.Empty(t, tt.outErr)

				file, err := ioutil.ReadFile(path.Join(tt.workspacePath, formula.WorkspacesFile))
				assert.NoError(t, err)
				workspaces := formula.Workspaces{}
				err = json.Unmarshal(file, &workspaces)
				assert.NoError(t, err)
				pathName := workspaces[tt.workspace.Name]
				assert.Contains(t, tt.workspace.Dir, pathName)
			}
		})
	}
}

func TestManagerDelete(t *testing.T) {
	cleanForm()
	fullDir := createFullDir()

	tmpDir := path.Join(os.TempDir(), "workspace-delete")
	_ = os.Mkdir(tmpDir, os.ModePerm)
	defer os.RemoveAll(tmpDir)

	workspaceFile := path.Join(tmpDir, formula.WorkspacesFile)
	fileNonExistentPath := path.Join(tmpDir, "non-existent")
	dirManager := stream.NewDirManager(stream.NewFileManager())

	tests := []struct {
		name          string
		workspacePath string
		workspace     formula.Workspace
		outErr        string
	}{
		{
			name:          "success delete",
			workspacePath: tmpDir,
			workspace: formula.Workspace{
				Name: "zup",
				Dir:  fullDir,
			},
		},
		{
			name:          "invalid workspace",
			workspacePath: tmpDir,
			workspace: formula.Workspace{
				Name: "zup-not-exists",
				Dir:  "home/user/go/src/github.com/ZupIT/ritchie-formulas-commons",
			},
			outErr: ErrInvalidWorkspace.Error(),
		},
		{
			name:          "write file error",
			workspacePath: fileNonExistentPath,
			workspace: formula.Workspace{
				Name: "Default",
				Dir:  fullDir,
			},
			outErr: mocks.FileNotFoundError(filepath.Join(fileNonExistentPath, formula.WorkspacesFile)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ioutil.WriteFile(workspaceFile, []byte(`{"zup": "some/dir/path"}`), os.ModePerm)
			assert.NoError(t, err)

			localBuilder := &mocks.LocalBuilderMock{}

			workspace := New(tt.workspacePath, tt.workspacePath, dirManager, localBuilder)
			got := workspace.Delete(tt.workspace)

			if got != nil {
				assert.EqualError(t, got, tt.outErr)
			} else {
				assert.Empty(t, tt.outErr)

				file, err := ioutil.ReadFile(path.Join(tt.workspacePath, formula.WorkspacesFile))
				assert.NoError(t, err)
				workspaces := formula.Workspaces{}
				err = json.Unmarshal(file, &workspaces)
				assert.NoError(t, err)
				_, exists := workspaces[tt.workspace.Name]
				assert.False(t, exists)
			}
		})
	}
}

func TestManagerList(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), "workspace-list")
	_ = os.Mkdir(tmpDir, os.ModePerm)
	defer os.RemoveAll(tmpDir)

	dirManager := stream.NewDirManager(stream.NewFileManager())

	workspaceFile := path.Join(tmpDir, formula.WorkspacesFile)
	err := ioutil.WriteFile(workspaceFile, []byte(`{"zup": "/some/path"}`), os.ModePerm)
	assert.NoError(t, err)
	workspaceBrokenFile := path.Join(tmpDir, "broken", formula.WorkspacesFile)
	_ = os.Mkdir(path.Join(tmpDir, "broken"), os.ModePerm)
	err = ioutil.WriteFile(workspaceBrokenFile, []byte(`error`), os.ModePerm)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		workspacePath string
		listSize      int
		outErr        string
	}{
		{
			name:          "success list",
			workspacePath: tmpDir,
			listSize:      2,
		},
		{
			name:          "success on non existent file",
			workspacePath: path.Join(tmpDir, "non-existent"),
			listSize:      1,
		},
		{
			name:          "unmarshal error",
			workspacePath: path.Join(tmpDir, "broken"),
			outErr:        "invalid character 'e' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localBuilder := &mocks.LocalBuilderMock{}

			workspace := New(tt.workspacePath, tt.workspacePath, dirManager, localBuilder)
			got, err := workspace.List()

			if err != nil {
				assert.EqualError(t, err, tt.outErr)
			} else {
				assert.Empty(t, tt.outErr)
				assert.Equal(t, tt.listSize, len(got))
			}
		})
	}
}

func TestPreviousHash(t *testing.T) {
	dirManager := stream.NewDirManager(stream.NewFileManager())
	tmpDir := os.TempDir()
	ritHome := path.Join(tmpDir, ".rit")
	defer os.RemoveAll(ritHome)

	_ = os.MkdirAll(path.Join(ritHome, "hashes"), os.ModePerm)

	formulaPath := "my/formula"
	hashFile := path.Join(ritHome, "hashes", "my-formula.txt")
	hashValue := "somehash"
	err := ioutil.WriteFile(hashFile, []byte(hashValue), os.ModePerm)
	assert.NoError(t, err)
	formulaNonExistentPath := path.Join(tmpDir, "non-existent")

	tests := []struct {
		name     string
		homePath string
		outErr   string
	}{
		{
			name:     "shoud return hash file content on success",
			homePath: ritHome,
		},
		{
			name:     "shoud fail when file doesn't exist",
			homePath: formulaNonExistentPath,
			outErr:   mocks.FileNotFoundError(path.Join(formulaNonExistentPath, "hashes", "my-formula.txt")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localBuilder := &mocks.LocalBuilderMock{}

			workspace := New(tt.homePath, tt.homePath, dirManager, localBuilder)
			hash, err := workspace.PreviousHash(formulaPath)

			if err != nil {
				assert.EqualError(t, err, tt.outErr)
			} else {
				assert.Empty(t, tt.outErr)
				assert.Equal(t, hashValue, hash)
			}
		})
	}
}

func TestUpdateHash(t *testing.T) {
	dirManager := stream.NewDirManager(stream.NewFileManager())
	ritHome := path.Join(os.TempDir(), "update-hash")
	_ = os.Mkdir(ritHome, os.ModePerm)
	defer os.RemoveAll(ritHome)

	tests := []struct {
		name     string
		homePath string
	}{
		{
			name:     "should update the correct file",
			homePath: ritHome,
		},
		{
			name:     "should ignore dir creation errors",
			homePath: ritHome,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localBuilder := &mocks.LocalBuilderMock{}

			workspace := New(tt.homePath, tt.homePath, dirManager, localBuilder)
			err := workspace.UpdateHash("my/formula", "hash")
			assert.NoError(t, err)

			file, err := ioutil.ReadFile(path.Join(tt.homePath, "hashes", "my-formula.txt"))
			assert.NoError(t, err)
			assert.Equal(t, "hash", string(file))
		})
	}
}

func cleanForm() {
	_ = os.Remove(filepath.Join(os.TempDir(), "my-custom-repo"))
}

func createFullDir() string {
	dir := filepath.Join(os.TempDir(), "my-custom-repo")
	_ = os.MkdirAll(dir, os.ModePerm)

	return dir
}
