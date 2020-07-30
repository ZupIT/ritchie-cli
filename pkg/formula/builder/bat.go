package builder

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	buildFile   = "build.bat"
	msgBuildErr = "failed building formula with build.bat, verify your repository"
	errMsgFmt   = `%s
More about error: %s`
)

var ErrBuildFormulaBuildBat = errors.New(msgBuildErr)

type BatManager struct {
	file stream.FileExister
}

func NewBuildBat(file stream.FileExister) formula.BatBuilder {
	return BatManager{file: file}
}

func (ba BatManager) Build(formulaPath string) error {
	if err := os.Chdir(formulaPath); err != nil {
		return err
	}

	if !ba.file.Exists(buildFile) {
		return ErrBuildOnWindows
	}

	var stderr bytes.Buffer
	cmd := exec.Command(buildFile)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return ErrBuildFormulaBuildBat
	}

	if stderr.String() != "" {
		return fmt.Errorf(errMsgFmt, msgBuildErr, stderr.String())
	}

	return nil
}
