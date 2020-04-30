package server

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	finder := NewFinder(tmp)
	setter := NewSetter(tmp)

	type out struct {
		status int
		err    error
	}

	tests := []struct {
		name string
		in   string
		out  out
	}{
		{
			name: "empty server",
			in:   "",
			out: out{
				status: 404,
			},
		},
		{
			name: "existing server",
			in:   "http://localhost/mocked",
			out: out{
				status: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			out := tt.out

			var body []byte
			if in != "" {
				body, _ = json.Marshal(&in)
				server := mockServer(out.status, body)
				defer server.Close()
				err := setter.Set(server.URL)
				if err != nil {
					fmt.Sprintln("Error in set")
					return
				}
			} else {
				err := setter.Set("")
				if err != nil {
					fmt.Sprintln("Error in set")
					return
				}
			}

			got, err := finder.Find()
			if err != nil {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, nil)
			}
			if got != "" && out.status != 200 {
				t.Errorf("Find(%s) got %v, want HttpStatus %v", tt.name, out.status, 200)
			}
		})
	}
}
