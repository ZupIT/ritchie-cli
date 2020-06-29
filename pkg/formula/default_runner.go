package formula

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

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
	outputEnv := fmt.Sprintf(EnvPattern, OutputEnv, setup.tmpOutputDir)
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

	fmt.Print(printAndValidOutputDir(setup))

	if err := d.PostRun(setup, false); err != nil {
		return err
	}

	return nil
}

func printAndValidOutputDir(setup Setup) string {

	files, err := ioutil.ReadDir(setup.tmpOutputDir)
	if err != nil {
		return prompt.Red("Fail to read output dir")
	}
	fOutputs := map[string]string{}

	resolveKey := func(name string) string { return strings.ToUpper(name) }

	if len(files) != len(setup.config.Outputs) {
		return prompt.Red("Output dir not have all the outputs files")
	}

	for _, f := range files {
		fName := fmt.Sprintf("%s/%s", setup.tmpOutputDir, f.Name())
		key := resolveKey(f.Name())
		b, err := ioutil.ReadFile(fName)
		if err != nil {
			return prompt.Red("fail to read file: " + fName)
		}
		fOutputs[key] = string(b)
	}

	var result string
	for _, o := range setup.config.Outputs {
		key := resolveKey(o.Name)
		v, exist := fOutputs[key]
		if !exist {
			return prompt.Red("file:" + key + " not found in output dir")
		}
		if o.Print {
			result += fmt.Sprintf("%s=%s\n", key, v)
		}
	}
	return result
}
