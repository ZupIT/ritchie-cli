package repo

import (
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewListUpdater(t *testing.T) {

	ritHome := os.TempDir()
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	repoList := NewLister(ritHome, fileManager)
	repoCreator := NewCreator(ritHome, github.NewRepoManager(http.DefaultClient), dirManager, fileManager)
	repoListCreator := NewListCreator(repoList, repoCreator)
	treeGenerator := tree.NewGenerator(dirManager, fileManager)
	repoUpdate := NewUpdater(ritHome, repoListCreator, treeGenerator, fileManager)

	type args struct {
		repoList   formula.RepositoryLister
		repoUpdate formula.RepositoryUpdater
	}
	tests := []struct {
		name string
		args args
		want formula.RepositoryListUpdater
	}{
		{
			name: "Build with success",
			args: args{
				repoList:   repoList,
				repoUpdate: repoUpdate,
			},
			want: ListUpdateManager{
				RepositoryLister:  repoList,
				RepositoryUpdater: repoUpdate,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewListUpdater(tt.args.repoList, tt.args.repoUpdate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewListUpdater() = %v, want %v", got, tt.want)
			}
		})
	}
}
