package server

import (
	"os"
	"reflect"
	"testing"
)

func TestNewSetter(t *testing.T) {
	NewSetter(os.TempDir())
}

func TestSet(t *testing.T) {

	s := NewSetter(os.TempDir())

	type in struct {
		serverURL string
	}

	type out struct {
		err  error
		want string
	}

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "empty serverUrl",
			in:   &in {
				serverURL: "",
			},
			out: &out{
				err:  nil,
			},
		},
		{
			name: "existing serverUrl",
			in: &in{
				serverURL: "http://localhost/mocked",
			},
			out: &out{
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out
			got := s.Set(in.serverURL)
			if !reflect.DeepEqual(out.err, got) {
				t.Errorf("Set(%s) got %v, want %v", in.serverURL , got, nil)
			}
		})
	}
}