package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type stableVersionCacheMock struct {
	Stable    string `json:"stableVersion"`
	ExpiresAt int64  `json:"expiresAt"`
}

func buildStableBodyMock(expiresAt int64) []byte {
	cache := stableVersionCacheMock{
		Stable:    "2.0.4",
		ExpiresAt: expiresAt,
	}
	b, _ := json.Marshal(cache)
	return b
}

func Test_rootCmd(t *testing.T) {
	type in struct {
		dir  stream.DirCreateChecker
		file stream.FileWriteReadExistRemover
		vm   version.Manager
	}

	notExpiredCache := time.Now().Add(time.Hour).Unix()
	versionManager := version.NewManager(
		"any value",
		sMocks.FileWriteReadExisterCustomMock{
			ExistsMock: func(path string) bool {
				return true
			},
			ReadMock: func(path string) ([]byte, error) {
				return buildStableBodyMock(notExpiredCache), nil
			},
		},
	)

	var tests = []struct {
		name    string
		wantErr bool
		in      in
	}{
		{
			name:    "Run with success",
			wantErr: false,
			in: in{
				dir: DirManagerCustomMock{
					exists: func(dir string) bool {
						return true
					},
					create: func(dir string) error {
						return nil
					},
				},
				file: stream.NewFileManager(),
				vm:   versionManager,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := os.TempDir()
			rootCmd := NewRootCmd(tmpDir, tt.in.dir, tt.in.file, TutorialFinderMock{}, tt.in.vm, nil, nil)

			if err := rootCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("root error = %v | error wanted: %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertContextToEnv(t *testing.T) {
	ctxFile := `{
  "current_context": "prod",
  "contexts": [
    "prod",
    "qa",
    "dev"
  ]
}`

	type in struct {
		ritHome string
		file    stream.FileWriteReadExistRemover
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				ritHome: os.TempDir(),
				file: sMocks.FileManagerMock{
					WriteFunc: func(path string, content []byte) error {
						return nil
					},
					ReadFunc: func(path string) ([]byte, error) {
						return []byte(ctxFile), nil
					},
					ExistsFunc: func(path string) bool {
						return true
					},
					RemoveFunc: func(path string) error {
						return nil
					},
				},
			},
			want: nil,
		},
		{
			name: "success when contexts file does not exist",
			in: in{
				ritHome: os.TempDir(),
				file: sMocks.FileManagerMock{
					ExistsFunc: func(path string) bool {
						return false
					},
				},
			},
			want: nil,
		},
		{
			name: "read contexts file error",
			in: in{
				ritHome: os.TempDir(),
				file: sMocks.FileManagerMock{
					ExistsFunc: func(path string) bool {
						return true
					},
					ReadFunc: func(path string) ([]byte, error) {
						return nil, errors.New("error to read contexts file")
					},
				},
			},
			want: errors.New("error to read contexts file"),
		},
		{
			name: "unmarshal error",
			in: in{
				ritHome: os.TempDir(),
				file: sMocks.FileManagerMock{
					ExistsFunc: func(path string) bool {
						return true
					},
					ReadFunc: func(path string) ([]byte, error) {
						return []byte("invalid"), nil
					},
				},
			},
			want: errors.New("invalid character 'i' looking for beginning of value"),
		},
		{
			name: "write envs file error",
			in: in{
				ritHome: os.TempDir(),
				file: sMocks.FileManagerMock{
					WriteFunc: func(path string, content []byte) error {
						return errors.New("error to write envs file")
					},
					ReadFunc: func(path string) ([]byte, error) {
						return []byte(ctxFile), nil
					},
					ExistsFunc: func(path string) bool {
						return true
					},
				},
			},
			want: errors.New("error to write envs file"),
		},
		{
			name: "remove contexts file error",
			in: in{
				ritHome: os.TempDir(),
				file: sMocks.FileManagerMock{
					WriteFunc: func(path string, content []byte) error {
						return nil
					},
					ReadFunc: func(path string) ([]byte, error) {
						return []byte(ctxFile), nil
					},
					ExistsFunc: func(path string) bool {
						return true
					},
					RemoveFunc: func(path string) error {
						return errors.New("error to remove contexts file")
					},
				},
			},
			want: errors.New("error to remove contexts file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := rootCmd{
				ritchieHome: tt.in.ritHome,
				file:        tt.in.file,
			}

			got := cmd.convertContextsFileToEnvsFile()
			if got != nil && got.Error() != tt.want.Error() {
				t.Fatalf("convertContextsFileToEnvsFile(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestConvertTree(t *testing.T) {
	type (
		repo struct {
			repos    formula.Repos
			listErr  error
			writeErr error
		}
		file struct {
			exist     bool
			readBytes []byte
			readErr   error
			writeErr  error
			removeErr error
		}
		mockTree struct {
			tree formula.Tree
			err  error
		}
	)

	type in struct {
		file file
		repo repo
		tree mockTree
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				tree: mockTree{
					tree: formula.Tree{
						Commands: api.Commands{
							"root_group": {
								Parent: "root",
								Usage:  "group",
								Help:   "group for add",
							},
							"root_group_verb": {
								Parent:  "root_group",
								Usage:   "verb",
								Help:    "verb for add",
								Formula: true,
							},
						},
					},
				},
				repo: repo{
					repos: formula.Repos{
						{
							Provider: "Github",
							Name:     "test1",
							Version:  "1.0.0",
							Priority: 0,
						},
						{
							Provider:    "Github",
							Name:        "test2",
							Version:     "1.0.0",
							Priority:    0,
							TreeVersion: tree.Version,
						},
					},
				},
			},
		},
		{
			name: "success without update",
			in: in{
				repo: repo{
					repos: formula.Repos{
						{
							Provider:    "Github",
							Name:        "test1",
							Version:     "1.0.0",
							Priority:    0,
							TreeVersion: tree.Version,
						},
						{
							Provider:    "Github",
							Name:        "test2",
							Version:     "1.0.0",
							Priority:    0,
							TreeVersion: tree.Version,
						},
					},
				},
			},
		},
		{
			name: "repo list error",
			in: in{
				repo: repo{
					listErr: errors.New("error to list repos"),
				},
			},
			want: errors.New("error to list repos"),
		},
		{
			name: "generate mockTree error",
			in: in{
				tree: mockTree{
					err: errors.New("error to generate mockTree"),
				},
				repo: repo{
					repos: formula.Repos{
						{
							Provider: "Github",
							Name:     "test1",
							Version:  "1.0.0",
							Priority: 0,
						},
						{
							Provider:    "Github",
							Name:        "test2",
							Version:     "1.0.0",
							Priority:    0,
							TreeVersion: tree.Version,
						},
					},
				},
			},
		},
		{
			name: "write mockTree.json error",
			in: in{
				file: file{
					writeErr: errors.New("error to write mockTree.json"),
				},
				repo: repo{
					repos: []formula.Repo{
						{
							Provider: "Github",
							Name:     "test1",
							Version:  "1.0.0",
							Priority: 0,
						},
						{
							Provider:    "Github",
							Name:        "test2",
							Version:     "1.0.0",
							Priority:    0,
							TreeVersion: tree.Version,
						},
					},
				},
			},
		},
		{
			name: "error to write repos",
			in: in{
				repo: repo{
					repos: []formula.Repo{
						{
							Provider: "Github",
							Name:     "test1",
							Version:  "1.0.0",
							Priority: 0,
						},
						{
							Provider:    "Github",
							Name:        "test2",
							Version:     "1.0.0",
							Priority:    0,
							TreeVersion: tree.Version,
						},
					},
					writeErr: errors.New("error to write repos"),
				},
			},
			want: errors.New("error to write repos"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			repoMock := new(mocks.RepoManager)
			repoMock.On("List").Return(in.repo.repos, in.repo.listErr)
			repoMock.On("Write", mock.Anything).Return(in.repo.writeErr)
			fileMock := new(mocks.FileManager)
			fileMock.On("Write", mock.Anything, mock.Anything).Return(in.file.writeErr)
			treeMock := new(mocks.TreeManager)
			treeMock.On("Generate", mock.Anything).Return(in.tree.tree, in.tree.err)

			cmd := rootCmd{
				file: fileMock,
				repo: repoMock,
				tree: treeMock,
			}

			got := cmd.convertTree()

			assert.Equal(t, tt.want, got)
		})
	}
}
