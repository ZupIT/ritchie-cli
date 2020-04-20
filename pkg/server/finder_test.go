package server

import (
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	finder := NewFinder(tmp)
	setter := NewSetter(tmp)

	type in struct {
		serverUrl string
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
			name: "empty server",
			in:   nil,
			out: &out{
				want: "",
				err:  nil,
			},
		},
		{
			name: "existing server",
			in: &in{
				serverUrl: "http://localhost/mocked",
			},
			out: &out{
				want: "http://localhost/mocked",
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			if in != nil {
				setter.Set(in.serverUrl)
			} else {
				fileutil.WriteFile(tmp+"/server", []byte(""))
			}

			out := tt.out
			got, err := finder.Find()
			if err != nil {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Find(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}