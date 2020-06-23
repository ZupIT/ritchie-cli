package sliceutil

import (
	"fmt"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
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
		{api.Command{Parent: "root", Usage: "add"}, true},
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

func TestRemove(t *testing.T) {
	type in struct {
		slice  []string
		remove string
	}

	tests := []struct {
		name string
		in   in
		out  int
	}{
		{
			name: "success",
			in: in{
				slice:  []string{"test_1", "test_2", "test_3"},
				remove: "test_2",
			},
			out: 2,
		},
		{
			name: "not remove any",
			in: in{
				slice:  []string{"test_1", "test_2", "test_3"},
				remove: "test_0",
			},
			out: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Remove(tt.in.slice, tt.in.remove)

			if tt.out != len(got) {
				t.Errorf("Remove(%s) got %v, want %v", tt.name, len(got), tt.out)
			}
		})
	}

}
