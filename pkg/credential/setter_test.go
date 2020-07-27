package credential

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

var (
	githubCred = Detail{Service: "github"}
	ctxFinder =   rcontext.FindManager{CtxFile: ""}
)

func TestSet(t *testing.T) {
	tmp := os.TempDir()
	setter := NewSetter(tmp, ctxFinder)
	tests := []struct {
		name string
		in   Detail
		out  error
	}{
		{
			name: "github credential",
			in:   githubCred,
			out:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setter.Set(tt.in)
			if got != tt.out {
				t.Errorf("Set(%s) got %v, want %v", tt.name, got, tt.out)
			}
		})
	}
}

