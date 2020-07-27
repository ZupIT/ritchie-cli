package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func Test_listRepoCmd_runFunc(t *testing.T) {
	finderTutorial := rtutorial.NewFinder(os.TempDir(), stream.NewFileManager())
	type in struct {
		RepositoryLister formula.RepositoryLister
	}
	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "someRepo1",
							},
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Run with success with more than 1 repo",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "someRepo1",
							},
							{
								Name: "someRepo2",
							},
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Return err when list fail",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return nil, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lr := NewListRepoCmd(tt.in.RepositoryLister, finderTutorial)
			lr.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := lr.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("setCredentialCmd_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type RepositoryListerCustomMock struct {
	list func() (formula.Repos, error)
}

func (m RepositoryListerCustomMock) List() (formula.Repos, error) {
	return m.list()
}
