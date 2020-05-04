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
	fCmdExists        = "rit add repo"
	fCmdCorrectGo     = "rit scaffold generate test-go"
	fCmdCorrectJava   = "rit scaffold generate test-java"
	fCmdCorrectNode   = "rit scaffold generate test-node"
	fCmdCorrectPython = "rit scaffold generate test-python"
	fCmdCorrectShell  = "rit scaffold generate test-shell"
	fCmdIncorrect     = "git scaffold generate testing"
	langGo            = "Go"
	langJava          = "Java"
	langNode          = "Node"
	langPython        = "Python"
	langShell         = "Shell"
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
		lang string
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
				lang: langGo,
			},
			out: &out{
				err: errors.New("this command already exists"),
			},
		},
		{
			name: "command correct-go",
			in: &in{
				fCmd: fCmdCorrectGo,
				lang: langGo,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command correct-java",
			in: &in{
				fCmd: fCmdCorrectJava,
				lang: langJava,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command correct-node",
			in: &in{
				fCmd: fCmdCorrectNode,
				lang: langNode,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command correct-python",
			in: &in{
				fCmd: fCmdCorrectPython,
				lang: langPython,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command correct-shell",
			in: &in{
				fCmd: fCmdCorrectShell,
				lang: langShell,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command incorrect",
			in: &in{
				fCmd: fCmdIncorrect,
				lang: langGo,
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

			_, got := creator.Create(in.fCmd, in.lang)
			if got != nil && got.Error() != out.err.Error() {
				t.Errorf("Create(%s) got %v, want %v", tt.name, got, out.err)
			}
		})
	}
}
