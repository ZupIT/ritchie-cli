package server

import (
	"fmt"
	"os"
	"testing"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	finder := NewFinder(tmp)
	setter := NewSetter(tmp)

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
				err := setter.Set(in)
				if err != nil {
					fmt.Sprintln("Error in set")
					return
				}
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
