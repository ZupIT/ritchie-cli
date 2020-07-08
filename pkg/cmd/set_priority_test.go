package cmd

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func TestSetRepoCmd_runFunc(t *testing.T) {
	type fields struct {
		InputList          prompt.InputList
		InputInt           prompt.InputInt
		RepoLister         formula.RepositoryLister
		RepoPrioritySetter formula.RepositoryPrioritySetter
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "run with success",
			fields: fields{
				InputList:  inputListMock{},
				RepoLister: repoListerNonEmptyMock{},
			},
			wantErr: false,
		},
		{
			name: "error on repoLister",
			fields: fields{
				InputList:  inputListMock{},
				RepoLister: repoListerErrorMock{},
			},
			wantErr: true,
		},
		{
			name: "error on inputList",
			fields: fields{
				InputList:  inputListErrorMock{},
				RepoLister: repoListerMock{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSetPriorityCmd(tt.fields.InputList, tt.fields.InputInt, tt.fields.RepoLister, tt.fields.RepoPrioritySetter)
			if err := s.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("runFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
