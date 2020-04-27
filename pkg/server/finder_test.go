package server

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	fileManager := stream.NewFileManager()
	finder := NewFinder(tmp, fileManager)
	setter := NewSetter(tmp, fileManager)

	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "empty server",
			in:   "",
			out:  "",
		},
		{
			name: "existing server",
			in:   "http://localhost/mocked",
			out:  "http://localhost/mocked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			out := tt.out

			if in != "" {
				_ = setter.Set(in)
			}

			got, err := finder.Find()
			if err != nil {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, nil)
			}
			if got != out {
				t.Errorf("Find(%s) got %v, want %v", tt.name, got, out)
			}
		})
	}
}
