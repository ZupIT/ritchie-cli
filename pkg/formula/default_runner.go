package formula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
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
	outputEnv := fmt.Sprintf(EnvPattern, OutputEnv, setup.outputFilePath)
	cmd.Env = append(cmd.Env, pwdEnv)
	cmd.Env = append(cmd.Env, cPwdEnv)
	cmd.Env = append(cmd.Env, outputEnv)

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

	printOutEnvs(setup)

	if err := d.PostRun(setup, false); err != nil {
		return err
	}

	return nil
}

func printOutEnvs(setup Setup) {

	f, _ := os.Open(setup.outputFilePath)
	b, _ := ioutil.ReadAll(f)
	fOutputs := map[string]string{}
	fOutputsPrint := map[string]string{}
	if err := json.Unmarshal(b, &fOutputs); err != nil {
		prompt.Error("Fail to read json from output file")
		return
	}

	if len(fOutputs) != len(setup.config.Outputs) {
		prompt.Error("Output file return wrong size of outputs")
		return
	}
	for _, o := range setup.config.Outputs {
		v, exist := fOutputs[o.Name]
		if !exist {
			prompt.Error("Should return " + o.Name + " output on output file")
		}
		if o.Print == true {
			fOutputsPrint[o.Name] = v
		}
	}

	result, _ := json.Marshal(fOutputsPrint)
	fmt.Printf(string(result))

}
