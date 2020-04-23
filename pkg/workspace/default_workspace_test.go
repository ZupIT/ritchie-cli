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
	fileManager := stream.NewFileManager()
	dirCreater := stream.NewDirCreater()
	workman := NewChecker(home, dirCreater, fileManager)
	if err := workman.Check(); err != nil {
		t.Errorf("Check got %v, want %v", err, nil)
	}
}
