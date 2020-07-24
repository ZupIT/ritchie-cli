package rtutorial

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func TestSet(t *testing.T) {
	type out struct {
		want      TutorialHolder
		err       error
		waitError bool
	}
	err := errors.New("some error")

	tests := []struct {
		name       string
		in         string
		out        *out
		FileWriter stream.FileWriter
	}{
		{
			name: "Set enabled tutorial",
			in:   "enabled",
			out: &out{
				want:      TutorialHolder{Current: "enabled"},
				err:       nil,
				waitError: false,
			},
			FileWriter: sMocks.FileWriterCustomMock{
				WriteMock: func(path string, content []byte) error {
					return nil
				},
			},
		},
		{
			name: "Set disabled tutorial",
			in:   "disabled",
			out: &out{
				want:      TutorialHolder{Current: "disabled"},
				err:       nil,
				waitError: false,
			},
			FileWriter: sMocks.FileWriterCustomMock{
				WriteMock: func(path string, content []byte) error {
					return nil
				},
			},
		},
		{
			name: "Error writing the tutorial file",
			in:   DefaultTutorial,
			out: &out{
				want:      TutorialHolder{Current: DefaultTutorial},
				err:       err,
				waitError: true,
			},
			FileWriter: sMocks.FileWriterCustomMock{
				WriteMock: func(path string, content []byte) error {
					return err
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := os.TempDir()
			tmpTutorial := fmt.Sprintf(TutorialPath, tmp)
			defer os.RemoveAll(tmpTutorial)

			setter := NewSetter(tmp, tt.FileWriter)

			in := tt.in
			out := tt.out

			got, err := setter.Set(in)
			if err != nil && !tt.out.waitError {
				t.Errorf("Set(%s) - Execution error - got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Set(%s) - Error in the expected response -  got %v, want %v", tt.name, got, out.want)
			}
		})
	}

}
