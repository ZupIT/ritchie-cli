package rcontext

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestRemove(t *testing.T) {
	tmp := os.TempDir()
	finder := NewFinder(tmp)
	setter := NewSetter(tmp, finder)
	remover := NewRemover(tmp, finder)

	_, err := setter.Set(dev)
	if err != nil {
		fmt.Sprintf("Error in Set")
		return
	}
	_, err = setter.Set(qa)
	if err != nil {
		fmt.Sprintf("Error in Set")
		return
	}

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
			name: "dev context",
			in:   dev,
			out: &out{
				want: ContextHolder{Current: qa, All: []string{qa}},
				err:  nil,
			},
		},
		{
			name: "current context",
			in:   CurrentCtx + qa,
			out: &out{
				want: ContextHolder{All: []string{}},
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := tt.in
			out := tt.out

			got, err := remover.Remove(in)
			if err != nil {
				t.Errorf("Remove(%s) got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Remove(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
