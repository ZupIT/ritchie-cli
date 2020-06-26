package formula

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

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

	printOutEnvs(setup, cmd)

	if err := d.PostRun(setup, false); err != nil {
		return err
	}

	return nil
}

func printOutEnvs(setup Setup, cmd *exec.Cmd) {

	//Get From Env
	tEnv := "TESTE_ENV"
	println("tEnv:", os.Getenv(tEnv))
	println(OutputEnv,":", os.Getenv(OutputEnv))
	for _, e := range cmd.Env{
		k := strings.Split(e,"=")
		if k[0] == OutputEnv{
			fmt.Printf("%s=%s\n", k[0], k[1])
		}
		if k[0] == tEnv{
			fmt.Printf("%s=%s\n", k[0], k[1])
		}
	}

	//Get From File
	fOutputs := map[string]string{}
	f, _ := os.Open(setup.outputFilePath)
	b, _ := ioutil.ReadAll(f)
	for _, c := range strings.Split(string(b), ";") {
		l := strings.Split(c, "=")
		if len(l) == 2 {
			fOutputs[l[0]] = l[1]
		}
	}
	for _, o := range setup.config.Outputs {
		if _, exist := fOutputs[o.Name]; exist && o.Print == true {
			fmt.Printf("%s=%s\n", o.Name, fOutputs[o.Name])
		}
	}
}
