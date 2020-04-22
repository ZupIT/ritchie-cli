package rcontext

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestRemove(t *testing.T) {
	tmp := os.TempDir()
	var write stream.FileWriter
	var finder Finder
	fileReadExister := stream.NewReadExister(stream.NewFileReader(), stream.NewFileExister())
	setter := NewSetter(tmp, NewFinder(tmp, fileReadExister), stream.NewFileWriter())
	setter.Set(dev)
	setter.Set(qa)

	type in struct {
		ctx    string
		finder *finderMock
		write  *writeUtilMock
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
			name: "dev context",
			in: &in{
				ctx: dev,
			},
			out: &out{
				want: ContextHolder{Current: qa, All: []string{qa}},
				err:  nil,
			},
		},
		{
			name: "current context",
			in: &in{
				ctx: CurrentCtx + qa,
			},
			out: &out{
				want: ContextHolder{All: []string{}},
				err:  nil,
			},
		},
		{
			name: "error find context",
			in: &in{
				ctx: dev,
				finder: &finderMock{
					ctx: ContextHolder{},
					err: errors.New("error find context"),
				},
			},
			out: &out{
				want: ContextHolder{},
				err:  errors.New("error find context"),
			},
		},
		{
			name: "error write context",
			in: &in{
				ctx: dev,
				write: &writeUtilMock{
					err: errors.New("write context error"),
				},
			},
			out: &out{
				want: ContextHolder{},
				err:  errors.New("write context error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			if in.write != nil {
				write = in.write
			} else {
				write = stream.NewFileWriter()
			}
			if in.finder != nil {
				finder = in.finder
			} else {
				finder = NewFinder(tmp, fileReadExister)
			}

			remover := NewRemover(tmp, finder, write)
			findRemover := NewFindRemover(finder, remover)
			got, err := findRemover.Remove(in.ctx)
			out := tt.out
			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Remove(%s) got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Remove(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
