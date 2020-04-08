package sliceutil

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"world", true},
		{"notfound", false},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := Contains([]string{"world", "earth", "universe"}, tt.in)
			if got != tt.out {
				t.Errorf("Contains got %v, want %v", got, tt.out)
			}
		})
	}
}

func TestContainsCmd(t *testing.T) {
	tests := []struct {
		in  api.Command
		out bool
	}{
		{api.Command{Parent: "root_set", Usage: "credential"}, true},
		{api.Command{Parent: "root", Usage: "notfound"}, false},
	}

	for _, tt := range tests {
		path := fmt.Sprintf("%s_%s", tt.in.Parent, tt.in.Usage)
		t.Run(path, func(t *testing.T) {
			got := ContainsCmd(api.CoreCmds, tt.in)
			if got != tt.out {
				t.Errorf("ContainsCmd got %v, want %v", got, tt.out)
			}
		})
	}
}
