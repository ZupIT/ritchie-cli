package credential

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestFind(t *testing.T) {

	fileManager := stream.NewFileManager()
	tmp := os.TempDir()
	setter := NewSetter(tmp, ctxFinder)
	_ = setter.Set(githubCred)

	type out struct {
		cred Detail
		err  error
	}

	type in struct {
		homePath  string
		ctxFinder rcontext.Finder
		file      stream.FileReader
		provider  string
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "Run with success",
			in: in{
				homePath:  tmp,
				ctxFinder: ctxFinder,
				file:      fileManager,
				provider:  githubCred.Service,
			},
			out: out{
				cred: githubCred,
				err:  nil,
			},
		},
		{
			name: "Return err when file not exist",
			in: in{
				homePath:  tmp,
				ctxFinder: ctxFinder,
				file:      fileManager,
				provider:  "aws",
			},
			out: out{
				cred: Detail{},
				err:  errors.New(prompt.Red(fmt.Sprintf(errNotFoundTemplate, "aws"))),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.out
			finder := NewFinder(tt.in.homePath, tt.in.ctxFinder, tt.in.file)
			got, err := finder.Find(tt.in.provider)
			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, out.err)
			}

			if !reflect.DeepEqual(out.cred, got) {
				t.Errorf("Find(%s) got %v, want %v", tt.name, got, out.cred)
			}
		})
	}
}
