package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestWorkspaceManager_Add(t *testing.T) {
	tmpDir := os.TempDir()
	fileManager := stream.NewFileManager()
	workspaceFile := fmt.Sprintf(workspacesPattern, tmpDir)
	if err := fileManager.Remove(workspaceFile); err != nil {
		t.Error(err)
	}

	type in struct {
		workspace   Workspace
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
				workspace: Workspace{
					Name: "zup",
					Dir:  "home/user/go/src/github.com/ZupIT/ritchie-formulas",
				},
				fileManager: fileManager,
			},
			out: nil,
		},
		{
			name: "success edit",
			in: in{
				workspace: Workspace{
					Name: "commons",
					Dir:  "home/user/go/src/github.com/ZupIT/ritchie-formulas-commons",
				},
				fileManager: fileManager,
			},
			out: nil,
		},
		{
			name: "read not found",
			in: in{
				workspace: Workspace{
					Name: "commons",
					Dir:  "home/user/go/src/github.com/ZupIT/ritchie-formulas-commons",
				},
				fileManager: fileManagerMock{exist: true, readErr: errors.New("not found file")},
			},
			out: errors.New("not found file"),
		},
		{
			name: "unmarshal error",
			in: in{
				workspace: Workspace{
					Name: "commons",
					Dir:  "home/user/go/src/github.com/ZupIT/ritchie-formulas-commons",
				},
				fileManager: fileManagerMock{exist: true, read: []byte("error")},
			},
			out: errors.New("invalid character 'e' looking for beginning of value"),
		},
		{
			name: "write error",
			in: in{
				workspace: Workspace{
					Name: "commons",
					Dir:  "home/user/go/src/github.com/ZupIT/ritchie-formulas-commons",
				},
				fileManager: fileManagerMock{exist: false, writeErr: errors.New("write file error")},
			},
			out: errors.New("write file error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in

			workspace := New(tmpDir, in.fileManager)
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
	workspaceFile := fmt.Sprintf(workspacesPattern, tmpDir)

	type in struct {
		workspaces  *Workspaces
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
				workspaces:  &Workspaces{"commons": "/home/user/ritchie-formulas"},
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

			workspace := New(tmpDir, in.fileManager)
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
