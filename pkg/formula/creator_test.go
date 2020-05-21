package formula

import (
	"errors"
	"fmt"
	"log"
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
	_ = fileutil.RemoveDir(os.TempDir() + "/customRepo")
}

func createDirWithMakefile() (string, error) {
	dir := os.TempDir() + "/customRepo"
	err := fileutil.CreateDirIfNotExists(dir, os.ModePerm)
	makefilePath := fmt.Sprintf("%s/%s", dir, Makefile)
	_, err = fileutil.CreateFileIfNotExist(makefilePath, []byte(""))
	return dir, err
}

func createDirWithTree() (string, error) {
	dir := os.TempDir() + "/customRepo"
	err := fileutil.CreateDirIfNotExists(dir, os.ModePerm)
	treeJsonPath := fmt.Sprintf("%s/%s", dir, TreeJson)
	_, err = fileutil.CreateFileIfNotExist(treeJsonPath, []byte(""))
	return dir, err
}

func createDir(dirname string, filename string, both bool) (string, error) {
	dir := os.TempDir() + dirname
	treeJsonDir := fmt.Sprintf("%s/%s", dir, "tree")
	treeJsonFile := fmt.Sprintf(TreeCreatePathPattern, dir)

	err := fileutil.CreateDirIfNotExists(dir, os.ModePerm)
	if both {
		err := fileutil.CreateDirIfNotExists(treeJsonDir, os.ModePerm)

		_, err = fileutil.CreateFileIfNotExist(treeJsonFile, []byte(""))
		makefilePath := fmt.Sprintf("%s/%s", dir, Makefile)
		_, err = fileutil.CreateFileIfNotExist(makefilePath, []byte(""))
		return dir, err
	}
	path := fmt.Sprintf("%s/%s", dir, filename)
	_, err = fileutil.CreateFileIfNotExist(path, []byte(""))
	return dir, err
}

func TestCreator(t *testing.T) {
	cleanForm()
	dir, err := createDirWithMakefile()

	fullDir,err := createDir("/customFullRepo","", true)

	fmt.Println(fullDir)
	if err != nil {
		log.Fatalf("Erro")
	}
	treeMan := NewTreeManager("../../testdata", repoListerMock{}, api.SingleCoreCmds)

	type in struct {
		fCmd          string
		lang          string
		customRepoDir string
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
		{
			name: "command to custom repo with missing packge.json",
			in: &in{
				fCmd:          fCmdCorrectGo,
				lang:          langGo,
				customRepoDir: dir,
			},
			out: &out{
				err: ErrTreeJsonNotFound,
			},
		},
		{
			name: "command to custom repo correct",
			in: &in{
				fCmd:          fCmdCorrectGo,
				lang:          langGo,
				customRepoDir: fullDir,
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
			_, got := creator.Create(in.fCmd, in.lang, in.customRepoDir)
			if got != nil && got.Error() != out.err.Error() {
				t.Errorf("Create(%s) got %v, want %v", tt.name, got, out.err)
			}
		})
	}
}
