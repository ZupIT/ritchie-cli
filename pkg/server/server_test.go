package server

import (
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

func TestMain(m *testing.M) {
	f := fmt.Sprintf(serverFilePattern, os.TempDir())
	fileutil.RemoveFile(f)
	os.Exit(m.Run())
}
