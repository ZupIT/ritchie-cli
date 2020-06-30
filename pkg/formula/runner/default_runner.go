package runner

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type DefaultRunner struct {
	formula.PreRunner
	formula.PostRunner
	formula.InputRunner
}

func NewDefaultRunner(preRunner formula.PreRunner, postRunner formula.PostRunner, inRunner formula.InputRunner) DefaultRunner {
	return DefaultRunner{preRunner, postRunner, inRunner}
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
	outputEnv := fmt.Sprintf(formula.EnvPattern, formula.OutputEnv, setup.TmpOutputDir)
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

func printAndValidOutputDir(setup formula.Setup) string {

	files, err := ioutil.ReadDir(setup.TmpOutputDir)
	if err != nil {
		return prompt.Red("Fail to read output dir")
	}
	fOutputs := map[string]string{}

	resolveKey := func(name string) string { return strings.ToUpper(name) }

	if len(files) != len(setup.Config.Outputs) {
		return prompt.Red("Output dir size is different of outputs array in config.json")
	}

	for _, f := range files {
		fName := fmt.Sprintf("%s/%s", setup.TmpOutputDir, f.Name())
		key := resolveKey(f.Name())
		b, err := ioutil.ReadFile(fName)
		if err != nil {
			return prompt.Red("fail to read file: " + fName)
		}
		fOutputs[key] = string(b)
	}

	var result string
	for _, o := range setup.Config.Outputs {
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
