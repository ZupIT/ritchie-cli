package formula

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

type DefaultRunner struct {
	PreRunner
	PostRunner
	InputRunner
}

func NewDefaultRunner(preRunner PreRunner, postRunner PostRunner, inRunner InputRunner) DefaultRunner {
	return DefaultRunner{preRunner, postRunner, inRunner}
}

func (d DefaultRunner) Run(def Definition, inputType api.TermInputType) error {
	setup, err := d.PreRun(def)
	if err != nil {
		return err
	}

	cmd := exec.Command(setup.tmpBinFilePath)

	cmd.Env = os.Environ()
	pwdEnv := fmt.Sprintf(EnvPattern, PwdEnv, setup.pwd)
	cPwdEnv := fmt.Sprintf(EnvPattern, CPwdEnv, setup.pwd)
	cmd.Env = append(cmd.Env, pwdEnv)
	cmd.Env = append(cmd.Env, cPwdEnv)

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

	if err := d.PostRun(setup, false); err != nil {
		return err
	}

	return nil
}
