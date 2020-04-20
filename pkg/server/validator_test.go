package server

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestNewValidator(t *testing.T) {
	tmp := os.TempDir()
	finder := NewFinder(tmp)
	NewValidator(finder)
}

func TestValidator(t *testing.T) {
	tmp := os.TempDir()
	finder := NewFinder(tmp)
	setter := NewSetter(tmp)

	type in struct {
		serverFinder Finder
		serverUrl string
	}

	type out struct {
		err  error
	}

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "empty serverUrl",
			in:   &in {
				serverFinder: finder,
				serverUrl: "",
			},
			out: &out{
				err:  fmt.Errorf("No server URL found ! Please set a server URL."),
			},
		},
		{
			name: "existing serverUrl",
			in: &in{
				serverFinder: finder,
				serverUrl: "http://localhost/mocked",
			},
			out: &out{
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			setter.Set(in.serverUrl)
			validator := NewValidator(in.serverFinder)

			out := tt.out
			err := validator.Validate(); if err != nil {
				if !reflect.DeepEqual(out.err.Error(), err.Error()) {
					t.Errorf("Find(%s) got %v, want %v", tt.name, err.Error(), out.err)
				}
			}
		})
	}
}