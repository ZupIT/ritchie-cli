package builder

import (
	"errors"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const buildFile = "build.bat"

var ErrBuildFormulaBuildBat = errors.New("failed building formula with build.bat, verify your repository")

type BatManager struct{}

func NewBuildBat() formula.BatBuilder {
	return BatManager{}
}

func (ba BatManager) Build(formulaPath string) error {
	if err := os.Chdir(formulaPath); err != nil {
		return err
	}

	cmd := exec.Command(buildFile)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return ErrBuildFormulaBuildBat
	}

	return nil
}
