package credsingle

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	setter := NewSetter(tmp, ctxFinder, sessManager)
	err := setter.Set(githubCred)
	if err != nil {
		fmt.Sprintln("Error in Set")
		return
	}
	finder := NewFinder(tmp, ctxFinder, sessManager)

	type out struct {
		cred credential.Detail
		err  error
	}

	tests := []struct {
		name string
		in   string
		out  out
	}{
		{
			name: "github",
			in:   "github",
			out: out{
				cred: githubCred,
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.out
			got, err := finder.Find(tt.in)

			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, out.err)
			}

			if !reflect.DeepEqual(out.cred, got) {
				t.Errorf("Find(%s) got %v, want %v", tt.name, got, out.cred)
			}
		})
	}
}
