package builder

import (
	"errors"
	"fmt"
	"os"
	"testing"

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

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(workspacePath)
	_ = streams.Unzip("../../../testdata/ritchie-formulas-test.zip", workspacePath)

	type in struct {
		fileManager stream.FileCopyExistLister
		dirManager  stream.DirCreateListCopier
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
			},
			want: nil,
		},
		{
			name: "list dir error",
			in: in{
				fileManager: fileManager,
				dirManager:  dirManagerMock{listErr: errors.New("error to list dir")},
			},
			want: errors.New("error to list dir"),
		},
		{
			name: "create dir error",
			in: in{
				fileManager: fileManager,
				dirManager:  dirManagerMock{createErr: errors.New("error to create dir")},
			},
			want: errors.New("error to create dir"),
		},
		{
			name: "copy so dir error",
			in: in{
				fileManager: fileManager,
				dirManager:  dirManagerMock{data: []string{"linux"}, copyErr: errors.New("error to copy dir")},
			},
			want: errors.New("error to copy dir"),
		},
		{
			name: "copy commons dir error",
			in: in{
				fileManager: fileManager,
				dirManager:  dirManagerMock{data: []string{"commons"}, copyErr: errors.New("error to copy dir")},
			},
			want: errors.New("error to copy dir"),
		},
		{
			name: "list files error",
			in: in{
				fileManager: fileManagerMock{listErr: errors.New("error to list files")},
				dirManager:  dirManager,
			},
			want: errors.New("error to list files"),
		},
		{
			name: "copy files error",
			in: in{
				fileManager: fileManagerMock{copyErr: errors.New("error to copy files")},
				dirManager:  dirManager,
			},
			want: errors.New("error to copy files"),
		},
	}

	for _, tt := range testes {
		t.Run(tt.name, func(t *testing.T) {
			builderManager := New(ritHome, tt.in.dirManager, tt.in.fileManager)
			got := builderManager.Build(workspacePath, formulaPath)

			if got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Build(%s) got %v, want %v", tt.name, got, tt.want)
			}

			if tt.want == nil {
				hasRitchieHome := dirManager.Exists(ritHome)
				if !hasRitchieHome {
					t.Errorf("Build(%s) did not create the Ritchie home directory", tt.name)
				}

				treeLocalFile := fmt.Sprintf("%s/repo/local/tree.json", ritHome)
				hasTreeLocalFile := fileManager.Exists(treeLocalFile)
				if !hasTreeLocalFile {
					t.Errorf("Build(%s) did not copy the tree local file", tt.name)
				}

				formulaFiles := fmt.Sprintf("%s/formulas/testing/formula/bin", ritHome)
				files, err := fileManager.List(formulaFiles)
				if err == nil && len(files) != 7 {
					t.Errorf("Build(%s) did not copy formulas files", tt.name)
				}

				configFile := fmt.Sprintf("%s/formulas/testing/formula/config.json", ritHome)
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
	data    []string
	listErr error
	copyErr error
	exist   bool
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
