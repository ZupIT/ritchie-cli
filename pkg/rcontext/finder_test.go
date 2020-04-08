package rcontext

import (
	"os"
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	finder := NewFinder(tmp)
	setter := NewSetter(tmp, finder)

	type in struct {
		ctx    string
		holder ContextHolder
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			if in != nil {
				setter.Set(in.ctx)
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
