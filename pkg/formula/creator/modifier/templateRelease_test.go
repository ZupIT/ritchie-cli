package modifier

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/git/github"
)

func TestTemplateRelease(t *testing.T) {
	repoInfo := github.NewRepoInfo(TemplateFormulasRepoURL, "")
	githubRepo := github.NewRepoManager(http.DefaultClient)
	tag, _ := githubRepo.LatestTag(repoInfo)

	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "modify with success",
			args: args{
				b: []byte(`{tag}`),
			},
			want: []byte(tag.Name),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TemplateRelease{}
			if got := tr.modify(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("modify() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
