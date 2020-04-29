package formula

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

const (
	fCmdExists    = "rit add repo"
	fCmdCorrect   = "rit scaffold generate test"
	fCmdIncorrect = "git scaffold generate testing"
)

type repoListerMock struct{}

func (repoListerMock) List() ([]Repository, error) {
	return []Repository{}, nil
}

func cleanForm() {
	_ = fileutil.RemoveDir(fmt.Sprintf(FormCreatePathPattern, os.TempDir()))
}

func TestCreator(t *testing.T) {
	cleanForm()
	treeMan := NewTreeManager("../../testdata", repoListerMock{}, api.SingleCoreCmds)

	type in struct {
		fCmd string
	}

	type out struct {
		err error
	}

	creator := NewCreator(fmt.Sprintf(FormCreatePathPattern, os.TempDir()), treeMan)

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "command exists",
			in: &in{
				fCmd: fCmdExists,
			},
			out: &out{
				err: errors.New("this command already exists"),
			},
		},
		{
			name: "command correct",
			in: &in{
				fCmd: fCmdCorrect,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command incorrect",
			in: &in{
				fCmd: fCmdIncorrect,
			},
			out: &out{
				err: errors.New("the formula's command needs to start with \"rit\" [ex.: rit group verb <noun>]"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out

			got := creator.Create(in.fCmd)
			if got != nil && got.Error() != out.err.Error() {
				t.Errorf("Create(%s) got %v, want %v", tt.name, got, out.err)
			}
		})
	}
}
