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

func TestFind(t *testing.T) {
	type in struct {
		tutorial string
		holder   TutorialHolder
	}

	type out struct {
		err       error
		want      TutorialHolder
		waitError bool
	}

	err := errors.New("some error")

	tests := []struct {
		name            string
		in              *in
		out             *out
		FileReadExister stream.FileReadExister
	}{
		{
			name: "With no tutorial file",
			out: &out{
				want:      TutorialHolder{Current: "enabled"},
				err:       nil,
				waitError: false,
			},
			FileReadExister: sMocks.FileReadExisterCustomMock{
				ExistsMock: func(path string) bool {
					return false
				},
			},
		},
		{
			name: "With existing tutorial file",
			out: &out{
				want:      TutorialHolder{Current: "disabled"},
				err:       nil,
				waitError: false,
			},
			FileReadExister: sMocks.FileReadExisterCustomMock{
				ReadMock: func(path string) ([]byte, error) {
					return []byte("{\"tutorial\":\"disabled\"}"), nil
				},
				ExistsMock: func(path string) bool {
					return true
				},
			},
		},
		{
			name: "Error reading the tutorial file",
			out: &out{
				want:      TutorialHolder{Current: "enabled"},
				err:       err,
				waitError: true,
			},
			FileReadExister: sMocks.FileReadExisterCustomMock{
				ReadMock: func(path string) ([]byte, error) {
					return []byte(""), err
				},
				ExistsMock: func(path string) bool {
					return true
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := os.TempDir()
			tmpTutorial := fmt.Sprintf(TutorialPath, tmp)
			defer os.RemoveAll(tmpTutorial)

			finder := NewFinder(tmp, tt.FileReadExister)

			out := tt.out
			got, err := finder.Find()
			if err != nil && !tt.out.waitError {
				t.Errorf("Set(%s) - Execution error - got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Set(%s) - Error in the expected response -  got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
