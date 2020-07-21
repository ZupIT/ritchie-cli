package repo

import (
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewListCreator(t *testing.T) {

	ritHome := os.TempDir()
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	repoList := NewLister(ritHome, fileManager)
	repoCreator := NewCreator(ritHome, github.NewRepoManager(http.DefaultClient), dirManager, fileManager)

	type args struct {
		repoList   formula.RepositoryLister
		repoCreate formula.RepositoryCreator
	}
	tests := []struct {
		name string
		args args
		want formula.RepositoryListCreator
	}{
		{
			name: "Build with success",
			args: args{
				repoList:   repoList,
				repoCreate: repoCreator,
			},
			want: ListCreateManager{
				RepositoryLister:  repoList,
				RepositoryCreator: repoCreator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewListCreator(tt.args.repoList, tt.args.repoCreate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewListCreator() = %v, want %v", got, tt.want)
			}
		})
	}
}
