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
	"errors"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestWorkspaceManager_Add(t *testing.T) {
	cleanForm()
	fullDir := createFullDir()

	tmpDir := os.TempDir()
	fileManager := stream.NewFileManager()
	workspaceFile := path.Join(tmpDir, formula.WorkspacesFile)
	if err := fileManager.Remove(workspaceFile); err != nil {
		t.Error(err)
	}

	type in struct {
		workspace   formula.Workspace
		fileManager stream.FileWriteReadExister
	}

	tests := []struct {
		name string
		in   in
		out  error
	}{
		{
			name: "success create",
			in: in{
				workspace: formula.Workspace{
					Name: "zup",
					Dir:  fullDir,
				},
				fileManager: fileManager,
			},
			out: nil,
		},
		{
			name: "success edit",
			in: in{
				workspace: formula.Workspace{
					Name: "commons",
					Dir:  fullDir,
				},
				fileManager: fileManager,
			},
			out: nil,
		},
		{
			name: "invalid workspace",
			in: in{
				workspace: formula.Workspace{
					Name: "zup",
					Dir:  "home/user/go/src/github.com/ZupIT/ritchie-formulas-commons",
				},
				fileManager: fileManager,
			},
			out: ErrInvalidWorkspace,
		},
		{
			name: "read not found",
			in: in{
				workspace: formula.Workspace{
					Name: "commons",
					Dir:  fullDir,
				},
				fileManager: fileManagerMock{exist: true, readErr: errors.New("not found file")},
			},
			out: errors.New("not found file"),
		},
		{
			name: "unmarshal error",
			in: in{
				workspace: formula.Workspace{
					Name: "commons",
					Dir:  fullDir,
				},
				fileManager: fileManagerMock{exist: true, read: []byte("error")},
			},
			out: errors.New("invalid character 'e' looking for beginning of value"),
		},
		{
			name: "write error",
			in: in{
				workspace: formula.Workspace{
					Name: "commons",
					Dir:  fullDir,
				},
				fileManager: fileManagerMock{exist: true, read: []byte("{\"name\":\"name\"}"), writeErr: errors.New("write file error")},
			},
			out: errors.New("write file error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in

			workspace := New(tmpDir, tmpDir, in.fileManager)
			got := workspace.Add(in.workspace)

			if got != nil && got.Error() != tt.out.Error() {
				t.Errorf("Add(%s) got %v, out %v", tt.name, got, tt.out)
			}
		})
	}
}

func TestManager_List(t *testing.T) {
	tmpDir := os.TempDir()
	fileManager := stream.NewFileManager()
	workspaceFile := path.Join(tmpDir, formula.WorkspacesFile)

	type in struct {
		workspaces  *formula.Workspaces
		fileManager stream.FileWriteReadExister
	}

	type out struct {
		listSize int
		error    error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "success list",
			in: in{
				workspaces:  &formula.Workspaces{"commons": "/home/user/ritchie-formulas"},
				fileManager: fileManager,
			},
			out: out{
				listSize: 1,
				error:    nil,
			},
		},
		{
			name: "not exist file",
			in: in{
				workspaces:  nil,
				fileManager: fileManager,
			},
			out: out{
				listSize: 0,
				error:    nil,
			},
		},
		{
			name: "read not found",
			in: in{
				fileManager: fileManagerMock{exist: true, readErr: errors.New("not found file")},
			},
			out: out{
				listSize: 0,
				error:    errors.New("not found file"),
			},
		},
		{
			name: "unmarshal error",
			in: in{
				fileManager: fileManagerMock{exist: true, read: []byte("error")},
			},
			out: out{
				listSize: 0,
				error:    errors.New("invalid character 'e' looking for beginning of value"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out

			_ = fileManager.Remove(workspaceFile)
			if in.workspaces != nil {
				content, _ := json.Marshal(in.workspaces)
				_ = fileManager.Write(workspaceFile, content)
			}

			workspace := New(tmpDir, tmpDir, in.fileManager)
			got, err := workspace.List()

			if err != nil && err.Error() != out.error.Error() {
				t.Errorf("List(%s) got err %v, out err %v", tt.name, err, out.error)
			}

			if len(got) != out.listSize {
				t.Errorf("List(%s) got size %v, out size %v", tt.name, len(got), out.listSize)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	cleanForm()
	fullDir := createFullDir()

	tmpDir := os.TempDir()
	fileManager := stream.NewFileManager()
	workspaceFile := path.Join(tmpDir, formula.WorkspacesFile)
	if err := fileManager.Remove(workspaceFile); err != nil {
		t.Error(err)
	}

	type in struct {
		workspace   formula.Workspace
		fileManager stream.FileWriteReadExister
	}

	tests := []struct {
		name string
		in   in
		out  error
	}{
		{
			name: "valid",
			in: in{
				workspace: formula.Workspace{
					Name: "zup",
					Dir:  fullDir,
				},
				fileManager: fileManager,
			},
			out: nil,
		},
		{
			name: "invalid workspace",
			in: in{
				workspace: formula.Workspace{
					Name: "zup",
					Dir:  "/home/user/invalid-workspace",
				},
				fileManager: fileManager,
			},
			out: ErrInvalidWorkspace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in

			workspace := New(tmpDir, in.fileManager)
			got := workspace.Validate(in.workspace)

			if got != nil && got.Error() != tt.out.Error() {
				t.Errorf("Validate(%s) got %v, out %v", tt.name, got, tt.out)
			}
		})
	}

}

func cleanForm() {
	_ = fileutil.RemoveDir(filepath.Join(os.TempDir(), "my-custom-repo"))
}

func createFullDir() string {
	dir := filepath.Join(os.TempDir(), "my-custom-repo")
	_ = fileutil.CreateDirIfNotExists(dir, os.ModePerm)

	return dir
}

type fileManagerMock struct {
	exist    bool
	read     []byte
	readErr  error
	writeErr error
}

func (f fileManagerMock) Exists(string) bool {
	return f.exist
}

func (f fileManagerMock) Read(string) ([]byte, error) {
	return f.read, f.readErr
}

func (f fileManagerMock) Write(string, []byte) error {
	return f.writeErr
}
