package repo

import (
	"errors"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestUpdateManager_Update(t *testing.T) {

	type fields struct {
		ritHome string
		repo    formula.RepositoryListCreator
		tree    formula.TreeGenerator
		file    stream.FileWriter
	}
	type args struct {
		name    formula.RepoName
		version formula.RepoVersion
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Return err when listRepos fail",
			fields: fields{
				ritHome: "",
				repo: repositoryListCreatorCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{}, errors.New("some error")
					},
				},
				tree: treeGeneratorCustomMock{},
				file: FileWriteCreatorReadExistRemover{},
			},
			args: args{
				name:    "any_name",
				version: "any_version",
			},
			wantErr: true,
		},
		{
			name: "Return err when listRepos is empty",
			fields: fields{
				ritHome: "",
				repo: repositoryListCreatorCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "any_repo",
							},
						}, nil
					},
				},
				tree: treeGeneratorCustomMock{},
				file: FileWriteCreatorReadExistRemover{},
			},
			args: args{
				name:    "not_a_repo_added_name",
				version: "any_version",
			},
			wantErr: true,
		},
		{
			name: "Return err when Create fail",
			fields: fields{
				ritHome: "",
				repo: repositoryListCreatorCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "any_repo",
							},
						}, nil
					},
					create: func(repo formula.Repo) error {
						return errors.New("some error")
					},
				},
				tree: treeGeneratorCustomMock{},
				file: FileWriteCreatorReadExistRemover{},
			},
			args: args{
				name:    "any_repo",
				version: "any_version",
			},
			wantErr: true,
		},
		{
			name: "Return err when write fail",
			fields: fields{
				ritHome: "",
				repo: repositoryListCreatorCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "any_repo",
							},
						}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{},
				file: FileWriteCreatorReadExistRemover{
					write: func(path string, content []byte) error {
						return errors.New("some error")
					},
				},
			},
			args: args{
				name:    "any_repo",
				version: "any_version",
			},
			wantErr: true,
		},
		{
			name: "Return err when generate fail",
			fields: fields{
				ritHome: "",
				repo: repositoryListCreatorCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "any_repo",
							},
						}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					generate: func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, errors.New("some error")
					},
				},
				file: FileWriteCreatorReadExistRemover{
					write: func(path string, content []byte) error {
						return nil
					},
				},
			},
			args: args{
				name:    "any_repo",
				version: "any_version",
			},
			wantErr: true,
		},
		{
			name: "Return err when fail to write tree",
			fields: fields{
				ritHome: "",
				repo: repositoryListCreatorCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "any_repo",
							},
						}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					generate: func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					write: func(path string, content []byte) error {
						if strings.Contains(path, "tree.json") {
							return errors.New("some error")
						}
						return nil
					},
				},
			},
			args: args{
				name:    "any_repo",
				version: "any_version",
			},
			wantErr: true,
		},
		{
			name: "Run with success",
			fields: fields{
				ritHome: "",
				repo: repositoryListCreatorCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "any_repo",
							},
						}, nil
					},
					create: func(repo formula.Repo) error {
						return nil
					},
				},
				tree: treeGeneratorCustomMock{
					generate: func(repoPath string) (formula.Tree, error) {
						return formula.Tree{}, nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					write: func(path string, content []byte) error {
						return nil
					},
				},
			},
			args: args{
				name:    "any_repo",
				version: "any_version",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			up := NewUpdater(
				tt.fields.ritHome,
				tt.fields.repo,
				tt.fields.tree,
				tt.fields.file,
			)
			if err := up.Update(tt.args.name, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type repositoryListCreatorCustomMock struct {
	create func(repo formula.Repo) error
	list   func() (formula.Repos, error)
}

func (m repositoryListCreatorCustomMock) Create(repo formula.Repo) error {
	return m.create(repo)
}

func (m repositoryListCreatorCustomMock) List() (formula.Repos, error) {
	return m.list()
}
