package rcontext

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"os"
	"testing"
)

const (
	dev = "dev"
	qa  = "qa"
)

func TestMain(m *testing.M) {
	cleanCtx()
	e := m.Run()
	os.Exit(e)
}

func cleanCtx() {
	_ = fileutil.RemoveDir(fmt.Sprintf(ContextPath, os.TempDir()))
}
