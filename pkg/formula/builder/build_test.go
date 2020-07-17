package builder

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestBuild(t *testing.T) {
	tmpDir := os.TempDir()
	workspacePath := fmt.Sprintf("%s/ritchie-formulas-test", tmpDir)
	formulaPath := fmt.Sprintf("%s/ritchie-formulas-test/testing/formula", tmpDir)
	ritHome := fmt.Sprintf("%s/.my-rit", os.TempDir())
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	defaultTreeManagerMock := tree.NewGenerator(dirManager, fileManager)

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(ritHome)
	_ = dirManager.Create(workspacePath)
	_ = streams.Unzip("../../../testdata/ritchie-formulas-test.zip", workspacePath)

	type in struct {
		fileManager stream.FileCopyExistListerWriter
		dirManager  stream.DirCreateListCopier
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
				fileManager: fileManager,
				dirManager:  dirManager,
				tree:        defaultTreeManagerMock,
			},
			want: nil,
		},
		{
			name: "list dir error",
			in: in{
				fileManager: fileManager,
				dirManager:  dirManagerMock{listErr: errors.New("error to list dir")},
				tree:        defaultTreeManagerMock,
			},
			want: errors.New("error to list dir"),
		},
		{
			name: "create dir error",
			in: in{
				fileManager: fileManager,
				dirManager:  dirManagerMock{createErr: errors.New("error to create dir")},
				tree:        defaultTreeManagerMock,
			},
			want: errors.New("error to create dir"),
		},
		{
			name: "copy so dir error",
			in: in{
				fileManager: fileManager,
				dirManager:  dirManagerMock{data: []string{"linux"}, copyErr: errors.New("error to copy dir")},
				tree:        defaultTreeManagerMock,
			},
			want: errors.New("error to copy dir"),
		},
		{
			name: "copy commons dir error",
			in: in{
				fileManager: fileManager,
				dirManager:  dirManagerMock{data: []string{"commons"}, copyErr: errors.New("error to copy dir")},
				tree:        defaultTreeManagerMock,
			},
			want: errors.New("error to copy dir"),
		},
		{
			name: "list files error",
			in: in{
				fileManager: fileManagerMock{listErr: errors.New("error to list files")},
				dirManager:  dirManager,
				tree:        defaultTreeManagerMock,
			},
			want: errors.New("error to list files"),
		},
		{
			name: "copy files error",
			in: in{
				fileManager: fileManagerMock{copyErr: errors.New("error to copy files")},
				dirManager:  dirManager,
				tree:        defaultTreeManagerMock,
			},
			want: errors.New("error to copy files"),
		},
	}

	for _, tt := range testes {
		t.Run(tt.name, func(t *testing.T) {
			builderManager := New(ritHome, tt.in.dirManager, tt.in.fileManager, tt.in.tree)
			got := builderManager.Build(workspacePath, formulaPath)

			if (tt.want == nil && got != nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Build(%s) got %v, want %v", tt.name, got, tt.want)
			}

			if tt.want == nil {
				hasRitchieHome := dirManager.Exists(ritHome)
				if !hasRitchieHome {
					t.Errorf("Build(%s) did not create the Ritchie home directory", tt.name)
				}

				treeLocalFile := fmt.Sprintf("%s/repos/local/tree.json", ritHome)
				hasTreeLocalFile := fileManager.Exists(treeLocalFile)
				if !hasTreeLocalFile {
					t.Errorf("Build(%s) did not copy the tree local file", tt.name)
				}

				formulaFiles := fmt.Sprintf("%s/repos/local/testing/formula/bin", ritHome)
				files, err := fileManager.List(formulaFiles)
				if err == nil && len(files) != 4 {
					t.Errorf("Build(%s) did not generate bin files", tt.name)
				}

				configFile := fmt.Sprintf("%s/repos/local/testing/formula/config.json", ritHome)
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

type fileManagerMock struct {
	data     []string
	listErr  error
	copyErr  error
	exist    bool
	writeErr error
}

func (f fileManagerMock) List(string) ([]string, error) {
	return f.data, f.listErr
}

func (f fileManagerMock) Copy(string, string) error {
	return f.copyErr
}

func (f fileManagerMock) Exists(string) bool {
	return f.exist
}

func (f fileManagerMock) Write(path string, content []byte) error {
	return f.writeErr
}
