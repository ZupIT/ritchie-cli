package server

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	finder := NewFinder(tmp)

	tests := []struct {
		name string
		in   Config
		out  Config
	}{
		{
			name: "empty server",
			in:   Config{},
			out:  Config{},
		},
		{
			name: "existing server",
			in:   Config{Organization: "org", URL: "http://localhost/mocked"},
			out:  Config{Organization: "org", URL: "http://localhost/mocked"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out

			if in.URL != "" {
				b, _ := json.Marshal(in)
				_ = fileutil.WriteFile(finder.serverFile, b)
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
