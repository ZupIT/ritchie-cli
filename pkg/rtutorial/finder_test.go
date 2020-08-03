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

var errReadingFile = errors.New("error reading file")

func TestFind(t *testing.T) {
	type out struct {
		err       error
		want      TutorialHolder
		wantError bool
	}

	tests := []struct {
		name string
		in   stream.FileReadExister
		out  *out
	}{
		{
			name: "With no tutorial file",
			in: sMocks.FileReadExisterCustomMock{
				ExistsMock: func(path string) bool {
					return false
				},
			},
			out: &out{
				want:      TutorialHolder{Current: "enabled"},
				err:       nil,
				wantError: false,
			},
		},
		{
			name: "With existing tutorial file",
			in: sMocks.FileReadExisterCustomMock{
				ReadMock: func(path string) ([]byte, error) {
					return []byte("{\"tutorial\":\"disabled\"}"), nil
				},
				ExistsMock: func(path string) bool {
					return true
				},
			},
			out: &out{
				want:      TutorialHolder{Current: "disabled"},
				err:       nil,
				wantError: false,
			},
		},
		{
			name: "Error reading the tutorial file",
			in: sMocks.FileReadExisterCustomMock{
				ReadMock: func(path string) ([]byte, error) {
					return []byte(""), errReadingFile
				},
				ExistsMock: func(path string) bool {
					return true
				},
			},
			out: &out{
				want:      TutorialHolder{Current: "enabled"},
				err:       errReadingFile,
				wantError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := os.TempDir()
			tmpTutorial := fmt.Sprintf(TutorialPath, tmp)
			defer os.RemoveAll(tmpTutorial)

			finder := NewFinder(tmp, tt.in)

			out := tt.out
			got, err := finder.Find()
			if err != nil && !tt.out.wantError {
				t.Errorf("%s - Execution error - got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("%s - Error in the expected response -  got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
