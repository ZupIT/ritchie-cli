package rtutorial

import (
	"os"
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	type out struct {
		want TutorialHolder
		err  error
	}
	type fieldsMock struct {
		read   func(path string) ([]byte, error)
		exists func(path string) bool
	}

	tests := []struct {
		name string
		in   string
		out  *out
	}{
		{
			name: "new on tutorial",
			in:   "on",
			out: &out{
				want: TutorialHolder{Current: "on"},
				err:  nil,
			},
		},
		{
			name: "new off tutorial",
			in:   "off",
			out: &out{
				want: TutorialHolder{Current: "off"},
				err:  nil,
			},
		},
		{
			name: "default tutorial",
			in:   DefaultTutorial,
			out: &out{
				want: TutorialHolder{Current: "on"},
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := os.TempDir()
			defer os.RemoveAll(tmp)

			setter := NewSetter(tmp)

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
