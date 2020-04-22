package rcontext

import (
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	dev = "dev"
	qa  = "qa"
)

func TestMain(m *testing.M) {
	remover := stream.NewFileRemover(stream.NewFileExister())
	cleanCtx(remover)
	e := m.Run()
	os.Exit(e)
}

func cleanCtx(file stream.FileRemover) {
	_ = file.Remove(fmt.Sprintf(ContextPath, os.TempDir()))
}
