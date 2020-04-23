package credsingle

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestSet(t *testing.T) {
	tmp := os.TempDir()
	file := stream.NewFileManager()
	dir := stream.NewDirCreater()
	setter := NewSetter(tmp, ctxFinder, sessManager, dir, file)

	tests := []struct {
		name string
		in   credential.Detail
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
