package rtutorial

import (
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

const (
	on  = "off"
	off = "on"
)

func TestMain(m *testing.M) {
	cleanTutorial()
	e := m.Run()
	os.Exit(e)
}

func cleanTutorial() {
	_ = fileutil.RemoveDir(fmt.Sprintf(TutorialPath, os.TempDir()))
}
