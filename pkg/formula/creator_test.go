package formula

import (
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
	fCmdCorrectPHP 	  = "rit scaffold generate test-php"
	fCmdCorrectPython = "rit scaffold generate test-python"
	fCmdCorrectShell  = "rit scaffold generate test-shell"
	fCmdIncorrect     = "git scaffold generate testing"
	langGo            = "Go"
	langJava          = "Java"
	langNode          = "Node"
	langPHP           = "PHP"
	langPython        = "Python"
	langShell         = "Shell"
)

type repoListerMock struct{}

func (repoListerMock) List() ([]Repository, error) {
	return []Repository{}, nil
}

func cleanForm() {
	_ = fileutil.RemoveDir(fmt.Sprintf(FormCreatePathPattern, os.TempDir()))
	_ = fileutil.RemoveDir(os.TempDir() + "/customRepo")
	_ = fileutil.RemoveDir(os.TempDir() + "/customRepoMakefile")
	_ = fileutil.RemoveDir(os.TempDir() + "/customRepoTreejson")
}

func createDirWithMakefile() string {
	dir := os.TempDir() + "/customRepoMakefile"
	_ = fileutil.CreateDirIfNotExists(dir, os.ModePerm)
	makefilePath := fmt.Sprintf("%s/%s", dir, Makefile)
	_ = fileutil.CreateFileIfNotExist(makefilePath, []byte(""))
	return dir
}

func createDirWithTree() string {
	dir := os.TempDir() + "/customRepoTreejson"
	treeJsonDir := fmt.Sprintf("%s/%s", dir, "tree")
	treeJsonFile := fmt.Sprintf(TreeCreatePathPattern, dir)
	_ = fileutil.CreateDirIfNotExists(dir, os.ModePerm)
	_ = fileutil.CreateDirIfNotExists(treeJsonDir, os.ModePerm)
	_ = fileutil.CreateFileIfNotExist(treeJsonFile, []byte(""))
	return dir
}

func createFullDir() string {
	dir := os.TempDir() + "/customRepo"
	treeJsonDir := fmt.Sprintf("%s/%s", dir, "tree")
	treeJsonFile := fmt.Sprintf(TreeCreatePathPattern, dir)
	makefilePath := fmt.Sprintf("%s/%s", dir, Makefile)
	_ = fileutil.CreateDirIfNotExists(dir, os.ModePerm)
	_ = fileutil.CreateDirIfNotExists(treeJsonDir, os.ModePerm)
	makefile, _ := fileutil.ReadFile("../../testdata/Makefile")
	_ = fileutil.CreateFileIfNotExist(makefilePath, makefile)
	_ = fileutil.CreateFileIfNotExist(treeJsonFile, []byte("{}"))

	return dir
}

func TestCreator(t *testing.T) {
	cleanForm()

	makefileDir := createDirWithMakefile()
	jsonDir := createDirWithTree()
	fullDir := createFullDir()

	treeMan := NewTreeManager("../../testdata", repoListerMock{}, api.SingleCoreCmds)

	type out struct {
		err error
	}

	creator := NewCreator(fmt.Sprintf(FormCreatePathPattern, os.TempDir()), treeMan)
	tests := []struct {
		name string
		in   *Create
		out  *out
	}{
		{
			name: "command exists",
			in: &Create{
				FormulaCmd: fCmdExists,
				Lang: langGo,
			},
			out: &out{
				err: ErrRepeatedCommand,
			},
		},
		{
			name: "command correct-go",
			in: &Create{
				FormulaCmd: fCmdCorrectGo,
				Lang: langGo,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command correct-java",
			in: &Create{
				FormulaCmd: fCmdCorrectJava,
				Lang: langJava,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command correct-node",
			in: &Create{
				FormulaCmd: fCmdCorrectNode,
				Lang: langNode,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command correct-php",
			in: &Create{
				FormulaCmd: fCmdCorrectPHP,
				Lang: langPHP,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command correct-python",
			in: &Create{
				FormulaCmd: fCmdCorrectPython,
				Lang: langPython,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command correct-shell",
			in: &Create{
				FormulaCmd: fCmdCorrectShell,
				Lang: langShell,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "command incorrect",
			in: &Create{
				FormulaCmd: fCmdIncorrect,
				Lang: langGo,
			},
			out: &out{
				err: ErrDontStartWithRit,
			},
		},
		{
			name: "command to custom repo with missing packge.json",
			in: &Create{
				FormulaCmd:          fCmdCorrectGo,
				Lang:          langGo,
				LocalRepoDir: makefileDir,
			},
			out: &out{
				err: ErrTreeJsonNotFound,
			},
		},
		{
			name: "command to custom repo with missing Makefile",
			in: &Create{
				FormulaCmd:          fCmdCorrectGo,
				Lang:          langGo,
				LocalRepoDir: jsonDir,
			},
			out: &out{
				err: ErrMakefileNotFound,
			},
		},
		{
			name: "command to custom repo correct",
			in: &Create{
				FormulaCmd:          fCmdCorrectGo,
				Lang:          langGo,
				LocalRepoDir: fullDir,
			},
			out: &out{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out
			_, got := creator.Create(*in)
			if got != nil && got.Error() != out.err.Error() {
				t.Errorf("Create(%s) got %v, want %v", tt.name, got, out.err)
			}
		})
	}
}
