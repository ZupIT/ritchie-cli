package runner

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/formula"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

type DefaultRunner struct {
	formula.PreRunner
	formula.PostRunner
	formula.InputRunner
	output formula.OutputRunner
}

func NewDefaultRunner(
	preRunner formula.PreRunner,
	postRunner formula.PostRunner,
	inRunner formula.InputRunner,
	outRunner formula.OutputRunner,
) DefaultRunner {
	return DefaultRunner{preRunner, postRunner, inRunner, outRunner}
}

func (d DefaultRunner) Run(def formula.Definition, inputType api.TermInputType) error {
	setup, err := d.PreRun(def)
	if err != nil {
		return err
	}

	cmd := exec.Command(setup.TmpBinFilePath)

	cmd.Env = os.Environ()
	pwdEnv := fmt.Sprintf(formula.EnvPattern, formula.PwdEnv, setup.Pwd)
	cPwdEnv := fmt.Sprintf(formula.EnvPattern, formula.CPwdEnv, setup.Pwd)
	cmd.Env = append(cmd.Env, pwdEnv)
	cmd.Env = append(cmd.Env, cPwdEnv)

	d.output.PrepareEnv(cmd, setup)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := d.Inputs(cmd, setup, inputType); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if err := d.output.ValidAndPrint(setup); err != nil {
		return err
	}

	if err := d.PostRun(setup, false); err != nil {
		return err
	}

	return nil
}
