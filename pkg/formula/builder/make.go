package builder

import (
	"errors"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

var ErrBuildFormulaMakefile = errors.New("failed building formula with make, verify your repository")

type MakeManager struct {}

func NewBuildMake() formula.MakeBuilder {
	return MakeManager{}
}

func (ma MakeManager) Build(formulaPath string) error {
	if err := os.Chdir(formulaPath); err != nil {
		return err
	}
	cmd := exec.Command("make", "build")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return ErrBuildFormulaMakefile
	}

	return nil
}