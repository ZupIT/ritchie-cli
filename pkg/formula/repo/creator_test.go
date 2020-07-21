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

	type fields struct {
		ritHome string
		github  github.Repositories
		dir     stream.DirCreateListCopyRemover
		file    stream.FileWriteCreatorReadExistRemover
	}
	type args struct {
		repo formula.Repo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Run with success",
			fields: fields{
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
			},
			args: args{
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
			fields: fields{
				ritHome: "",
				github: GitRepositoryMock{
					zipball: func(info github.RepoInfo, version string) (io.ReadCloser, error) {
						return nil, errors.New("some error")
					},
				},
				dir:  dirManager,
				file: fileManager,
			},
			args: args{
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
			fields: fields{
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
			},
			args: args{
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
			fields: fields{
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
			},
			args: args{
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
			fields: fields{
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
			},
			args: args{
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
			fields: fields{
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
			},
			args: args{
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
				tt.fields.ritHome,
				tt.fields.github,
				tt.fields.dir,
				tt.fields.file,
			)
			if err := cr.Create(tt.args.repo); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				newRepoPath := filepath.Join(tt.fields.ritHome, "repos", tt.args.repo.Name.String())
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
