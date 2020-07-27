package rcontext

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMock "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()

	type in struct {
		holder          ContextHolder
		FileReadExister stream.FileReadExister
	}

	type out struct {
		err  error
		want ContextHolder
	}

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "default context and existing ctx file",
			in: &in{
				holder: ContextHolder{Current: ""},
				FileReadExister: sMock.FileReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"current_context\":\"default\"}"), nil
					},
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
			out: &out{
				want: ContextHolder{Current: "default"},
				err:  nil,
			},
		},
		{
			name: "default context and missing ctx file",
			in: &in{
				holder: ContextHolder{Current: ""},
				FileReadExister: sMock.FileReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"current_context\":\"default\"}"), nil
					},
					ExistsMock: func(path string) bool {
						return false
					},
				},
			},
			out: &out{
				want: ContextHolder{Current: ""},
				err:  nil,
			},
		},
		{
			name: "default context and error on read file",
			in: &in{
				holder: ContextHolder{Current: ""},
				FileReadExister: sMock.FileReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte(""), errors.New("error reading file")
					},
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
			out: &out{
				want: ContextHolder{Current: ""},
				err:  errors.New("error reading file"),
			},
		},
		{
			name: "default context and incorrect json",
			in: &in{
				holder: ContextHolder{Current: "default"},
				FileReadExister: sMock.FileReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte(""), nil
					},
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
			out: &out{
				want: ContextHolder{Current: ""},
				err:  errors.New("unexpected end of JSON input"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := NewFinder(tmp, tt.in.FileReadExister)
			out := tt.out
			got, err := finder.Find()
			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Find(%s) - Execution error - got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Find(%s) - Error in the expected response - got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
