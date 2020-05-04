package server

import (
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

func TestMain(m *testing.M) {
	f := fmt.Sprintf(serverFilePattern, os.TempDir())
	err := fileutil.RemoveFile(f)
	if err != nil {
		fmt.Sprintln("Error in remove file")
		return
	}
	os.Exit(m.Run())
}
