package workspace

import (
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	home string
)

func TestMain(m *testing.M) {
	home = fmt.Sprintf("%s/.rit", os.TempDir())
	os.Exit(m.Run())
}

func TestCheckWorkingDir(t *testing.T) {
	fileReader := stream.NewFileReader()
	fileWriter := stream.NewFileWriter()
	fileExister := stream.NewFileExister()
	fileRemover := stream.NewFileRemover(fileExister)
	fileManager := stream.NewFileManager(fileWriter, fileReader, fileExister, fileRemover)
	dirCreater := stream.NewDirCreater()
	workman := NewChecker(home, dirCreater, fileManager)
	if err := workman.Check(); err != nil {
		t.Errorf("Check got %v, want %v", err, nil)
	}
}
