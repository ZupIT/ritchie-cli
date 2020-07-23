package cmd

import (
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/github"
)

func TestNewSingleInitCmd(t *testing.T) {
	cmd := NewInitCmd(defaultRepoAdderMock, defaultGitRepositoryMock)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	if cmd == nil {
		t.Errorf("NewInitCmd got %v", cmd)
		return
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func Test_initCmd_runPrompt(t *testing.T) {
	type fields struct {
		repo formula.RepositoryAdder
		git  github.Repositories
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Run With Success",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
			},
			wantErr: false,
		},
		{
			name: "Fail when call git.LatestTag",
			fields: fields{
				repo: defaultRepoAdderMock,
				git: GitRepositoryMock{
					latestTag: func(info github.RepoInfo) (github.Tag, error) {
						return github.Tag{}, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Fail when call repo.Add",
			fields: fields{
				repo: repoListerAdderCustomMock{
					add: func(d formula.Repo) error {
						return errors.New("some error")
					},
				},
				git: defaultGitRepositoryMock,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewInitCmd(tt.fields.repo, tt.fields.git)
			o.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := o.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("init_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
