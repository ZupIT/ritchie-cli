package repo

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestAdd(t *testing.T) {

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	type in struct {
		ritHome string
		creator formula.RepositoryCreator
		tree    formula.TreeGenerator
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
			name: "Run with success, when repository json not exist",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-adder-test-success")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Remove(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					return ritHomePath
				}(),
				creator: repositoryCreatorCustomMock{
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, nil
					},
				},
				file: fileManager,
				dir:  dirManager,
				repo: formula.Repo{
					Name:     "some_repo_name",
					Priority: 10,
					Token:    "",
					Url:      "https://github.com/someUser/ritchie-formulas",
					Version:  "2.0",
				},
			},
			wantErr: false,
		},
		{
			name: "Run with success, when repository json exist",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-adder-test-success")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Remove(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos", "some_repo_name"))
					repoFileData := `
					[
							{
								"name": "commons",
								"version": "v2.0.0",
								"url": "https://github.com/kaduartur/ritchie-formulas",
								"priority": 0
							}
					]
					`
					_ = fileManager.Write(filepath.Join(ritHomePath, "repos", reposFileName), []byte(repoFileData))
					return ritHomePath
				}(),
				creator: repositoryCreatorCustomMock{
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, nil
					},
				},
				file: fileManager,
				dir:  dirManager,
				repo: formula.Repo{
					Name:     "some_repo_name",
					Priority: 10,
					Token:    "",
					Url:      "https://github.com/someUser/ritchie-formulas",
					Version:  "2.0",
				},
			},

			wantErr: false,
		},
		{
			name: "Return err when RepositoryCreator fail",
			in: in{
				creator: repositoryCreatorCustomMock{
					create: func(repo formula.Repo) error {
						return errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when file read fail",
			in: in{
				creator: repositoryCreatorCustomMock{
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					read: func(path string) ([]byte, error) {
						return nil, errors.New("some error")
					},
					exists: func(path string) bool {
						return true
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when fail to parse json",
			in: in{
				creator: repositoryCreatorCustomMock{
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					read: func(path string) ([]byte, error) {
						return []byte("not json data"), nil
					},
					exists: func(path string) bool {
						return true
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when saveRepo fail",
			in: in{
				creator: repositoryCreatorCustomMock{
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					exists: func(path string) bool {
						return false
					},
					write: func(path string, content []byte) error {
						return errors.New("some error")
					},
				},
				dir: DirCreateListCopyRemoverCustomMock{
					create: func(dir string) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when tree Generate fail",
			in: in{
				creator: repositoryCreatorCustomMock{
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					exists: func(path string) bool {
						return false
					},
					write: func(path string, content []byte) error {
						return nil
					},
				},
				dir: DirCreateListCopyRemoverCustomMock{
					create: func(dir string) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad := NewAdder(
				tt.in.ritHome,
				tt.in.creator,
				tt.in.tree,
				tt.in.dir,
				tt.in.file,
			)
			if err := ad.Add(tt.in.repo); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				treePath := filepath.Join(tt.in.ritHome, "repos", tt.in.repo.Name.String(), "tree.json")
				if !fileManager.Exists(treePath) {
					t.Errorf("Tree with path %s not exist.", treePath)
				}
				repoJsonPath := filepath.Join(tt.in.ritHome, "repos", reposFileName)
				if !fileManager.Exists(repoJsonPath) {
					t.Errorf("RepoJsonPath with path %s not exist.", repoJsonPath)
				}
			}
		})
	}
}

type repositoryCreatorCustomMock struct {
	create func(repo formula.Repo) error
}

func (m repositoryCreatorCustomMock) Create(repo formula.Repo) error {
	return m.create(repo)
}

type treeGeneratorCustomMock struct {
	generate func(repoPath string) (formula.Tree, error)
}

func (m treeGeneratorCustomMock) Generate(repoPath string) (formula.Tree, error) {
	return m.generate(repoPath)
}

type FileWriteCreatorReadExistRemover struct {
	write  func(path string, content []byte) error
	create func(path string, data io.ReadCloser) error
	read   func(path string) ([]byte, error)
	exists func(path string) bool
	remove func(path string) error
}

func (m FileWriteCreatorReadExistRemover) Write(path string, content []byte) error {
	return m.write(path, content)
}

func (m FileWriteCreatorReadExistRemover) Create(path string, data io.ReadCloser) error {
	return m.create(path, data)
}

func (m FileWriteCreatorReadExistRemover) Read(path string) ([]byte, error) {
	return m.read(path)
}

func (m FileWriteCreatorReadExistRemover) Exists(path string) bool {
	return m.exists(path)
}

func (m FileWriteCreatorReadExistRemover) Remove(path string) error {
	return m.remove(path)
}

type DirCreateListCopyRemoverCustomMock struct {
	create func(dir string) error
	list   func(dir string, hiddenDir bool) ([]string, error)
	copy   func(src, dst string) error
	remove func(dir string) error
}

func (m DirCreateListCopyRemoverCustomMock) Create(dir string) error {
	return m.create(dir)
}

func (m DirCreateListCopyRemoverCustomMock) List(dir string, hiddenDir bool) ([]string, error) {
	return m.list(dir, hiddenDir)
}

func (m DirCreateListCopyRemoverCustomMock) Copy(src, dst string) error {
	return m.copy(src, dst)
}

func (m DirCreateListCopyRemoverCustomMock) Remove(dir string) error {
	return m.remove(dir)
}
