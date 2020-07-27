package rcontext

import (
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestSet(t *testing.T) {
	tmp := os.TempDir()
	file := stream.NewFileManager()
	finder := NewFinder(tmp, file)
	setter := NewSetter(tmp, finder)

	type out struct {
		want ContextHolder
		err  error
	}

	tests := []struct {
		name string
		in   string
		out  *out
	}{
		{
			name: "new dev context",
			in:   dev,
			out: &out{
				want: ContextHolder{Current: dev, All: []string{dev}},
				err:  nil,
			},
		},
		{
			name: "no duplicate context",
			in:   dev,
			out: &out{
				want: ContextHolder{Current: dev, All: []string{dev}},
				err:  nil,
			},
		},
		{
			name: "new qa context",
			in:   qa,
			out: &out{
				want: ContextHolder{Current: qa, All: []string{dev, qa}},
				err:  nil,
			},
		},
		{
			name: "default context",
			in:   DefaultCtx,
			out: &out{
				want: ContextHolder{Current: "", All: []string{dev, qa}},
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			out := tt.out

			got, err := setter.Set(in)
			if err != nil {
				t.Errorf("Set(%s) got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Set(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}

}
