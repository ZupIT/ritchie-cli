package credential

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	stream "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

var (
	githubCred = Detail{Service: "github"}
	streamMock = stream.FileReadExisterCustomMock{
		ReadMock: func(path string) ([]byte, error) {
			return []byte("{\"current_context\":\"default\"}"), nil
		},
		ExistsMock: func(path string) bool {
			return true
		},
	}
	ctxFinder = rcontext.FindManager{CtxFile: "", File: streamMock}
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
