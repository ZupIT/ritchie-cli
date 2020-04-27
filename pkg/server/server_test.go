package server

import (
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestMain(m *testing.M) {
	f := fmt.Sprintf(serverFilePattern, os.TempDir())
	_ = stream.NewFileManager().Remove(f)
	os.Exit(m.Run())
}
