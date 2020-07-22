package rtutorial

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	type in struct {
		tutorial string
		holder   TutorialHolder
	}

	type out struct {
		err  error
		want TutorialHolder
	}

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "empty tutorial",
			in:   nil,
			out: &out{
				want: TutorialHolder{Current: "on"},
				err:  nil,
			},
		},
		{
			name: "off tutorial",
			in: &in{
				tutorial: "off",
				holder:   TutorialHolder{Current: "off"},
			},
			out: &out{
				want: TutorialHolder{Current: "off"},
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := os.TempDir()
			defer os.RemoveAll(tmp)

			finder := NewFinder(tmp)
			setter := NewSetter(tmp)

			in := tt.in
			if in != nil {
				_, err := setter.Set(in.tutorial)
				if err != nil {
					fmt.Sprintln("Error in Set")
					return
				}
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
