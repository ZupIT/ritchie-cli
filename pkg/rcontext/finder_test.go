package rcontext

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestFind(t *testing.T) {
	var finder Finder
	var setter Setter
	tmp := os.TempDir()
	fileManager := stream.NewFileManager()

	type in struct {
		ctx      string
		holder   ContextHolder
		fileMock *readUtilMock
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
			name: "empty context",
			in:   nil,
			out: &out{
				want: ContextHolder{},
				err:  nil,
			},
		},
		{
			name: "dev context",
			in: &in{
				ctx:    dev,
				holder: ContextHolder{Current: dev},
			},
			out: &out{
				want: ContextHolder{Current: dev, All: []string{dev}},
				err:  nil,
			},
		},
		{
			name: "read file error",
			in: &in{
				ctx:    dev,
				holder: ContextHolder{Current: dev},
				fileMock: &readUtilMock{
					exist: true,
					err:   errors.New("read error"),
					file:  nil,
				},
			},
			out: &out{
				want: ContextHolder{},
				err:  errors.New("read error"),
			},
		},
		{
			name: "file to context error",
			in: &in{
				ctx:    dev,
				holder: ContextHolder{Current: dev},
				fileMock: &readUtilMock{
					exist: true,
					err:   nil,
					file:  []byte("error"),
				},
			},
			out: &out{
				want: ContextHolder{},
				err:  errors.New("invalid character 'e' looking for beginning of value"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			if in != nil && in.fileMock != nil {
				finder = NewFinder(tmp, in.fileMock)
			} else {
				finder = NewFinder(tmp, fileManager)
			}

			if in != nil {
				setter = NewSetter(tmp, finder, fileManager)
				_, _ = setter.Set(in.ctx)
			}

			out := tt.out
			got, err := finder.Find()
			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Find(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}

type readUtilMock struct {
	exist bool
	err   error
	file  []byte
}

func (f readUtilMock) Exists(string) bool {
	return f.exist
}

func (f readUtilMock) Read(string) ([]byte, error) {
	return f.file, f.err
}

type finderMock struct {
	ctx ContextHolder
	err error
}

func (f finderMock) Find() (ContextHolder, error) {
	return f.ctx, f.err
}
