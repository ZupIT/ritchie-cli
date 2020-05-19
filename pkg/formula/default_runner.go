package formula

import (
	"os"
	"os/exec"
)

type DefaultRunner struct {
	PreRunner
	InputRunner
}

func NewDefaultRunner(preRunner PreRunner, inRunner InputRunner) DefaultRunner {
	return DefaultRunner{preRunner, inRunner}
}

func (d DefaultRunner) Run(def Definition) error {
	setup, err := d.PreRun(def)
	if err != nil {
		return err
	}

	cmd := exec.Command(setup.tmpBinFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := d.Inputs(cmd, setup.formulaPath, &setup.config, false); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if err := PostRun(setup, false); err != nil {
		return err
	}

	return nil
}
