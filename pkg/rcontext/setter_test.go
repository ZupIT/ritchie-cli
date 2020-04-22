package rcontext

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestSet(t *testing.T) {
	var finder Finder
	var writer stream.FileWriter
	tmp := os.TempDir()

	type in struct {
		ctx           string
		finderMock    *finderMock
		writeUtilMock *writeUtilMock
	}

	type out struct {
		want ContextHolder
		err  error
	}

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "new dev context",
			in: &in{
				ctx: dev,
			},
			out: &out{
				want: ContextHolder{Current: dev, All: []string{dev}},
				err:  nil,
			},
		},
		{
			name: "no duplicate context",
			in: &in{
				ctx: dev,
			},
			out: &out{
				want: ContextHolder{Current: dev, All: []string{dev}},
				err:  nil,
			},
		},
		{
			name: "new qa context",
			in: &in{
				ctx: qa,
			},
			out: &out{
				want: ContextHolder{Current: qa, All: []string{dev, qa}},
				err:  nil,
			},
		},
		{
			name: "default context",
			in: &in{
				ctx: DefaultCtx,
			},
			out: &out{
				want: ContextHolder{Current: "", All: []string{dev, qa}},
				err:  nil,
			},
		},
		{
			name: "error to read context",
			in: &in{
				ctx: DefaultCtx,
				finderMock: &finderMock{
					ctx: ContextHolder{},
					err: errors.New("error to read file"),
				},
			},
			out: &out{
				want: ContextHolder{},
				err:  errors.New("error to read file"),
			},
		},
		{
			name: "error to write context",
			in: &in{
				ctx:           DefaultCtx,
				writeUtilMock: &writeUtilMock{err: errors.New("error to write context")},
			},
			out: &out{
				want: ContextHolder{},
				err:  errors.New("error to write context"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			if in.writeUtilMock != nil {
				writer = in.writeUtilMock
			} else {
				writer = stream.NewFileWriter()
			}

			if in.finderMock != nil {
				finder = in.finderMock
			} else {
				fileReader := stream.NewFileReader()
				fileExister := stream.NewFileExister()
				finder = NewFinder(tmp, stream.NewReadExister(fileReader, fileExister))
			}

			setter := NewSetter(tmp, finder, writer)
			findSetter := NewFindSetter(finder, setter)

			out := tt.out

			got, err := findSetter.Set(in.ctx)
			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Set(%s) got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Set(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}

type writeUtilMock struct {
	err error
}

func (w writeUtilMock) Write(string, []byte) error {
	return w.err
}
