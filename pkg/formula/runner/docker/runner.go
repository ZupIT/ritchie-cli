/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/mattn/go-isatty"
	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const (
	envFile   = ".env"
	dockerPwd = "/app"
)

var _ formula.Runner = RunManager{}

type RunManager struct {
	homeDir string
	env     env.Finder
	formula.InputResolver
	formula.PreRunner
}

func NewRunner(
	homeDir string,
	input formula.InputResolver,
	preRun formula.PreRunner,
	env env.Finder,
) formula.Runner {
	return RunManager{
		InputResolver: input,
		PreRunner:     preRun,
		env:           env,
		homeDir:       homeDir,
	}
}

func (ru RunManager) Run(def formula.Definition, inputType api.TermInputType, verbose bool, flags *pflag.FlagSet) error {
	setup, err := ru.PreRun(def)
	if err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(envFile); err != nil {
			return
		}
	}()

	cmd, err := ru.runDocker(setup, inputType, verbose, flags)
	if err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (ru RunManager) runDocker(setup formula.Setup, inputType api.TermInputType, verbose bool, flags *pflag.FlagSet) (*exec.Cmd, error) {
	volume := fmt.Sprintf("%s:/app", setup.Pwd)
	homeDirVolume := fmt.Sprintf("%s/.rit:/root/.rit", ru.homeDir)
	var args []string
	if isatty.IsTerminal(os.Stdout.Fd()) && inputType != api.Stdin {
		args = []string{
			"run",
			"--rm",
			"-it",
			"--env-file",
			envFile,
			"-v",
			volume,
			"-v",
			homeDirVolume,
			"--name",
			setup.ContainerId,
			setup.ContainerId,
		}
	} else {
		args = []string{
			"run",
			"--rm",
			"--env-file",
			envFile,
			"-v",
			volume,
			"-v",
			homeDirVolume,
			"--name",
			setup.ContainerId,
			setup.ContainerId,
		}
	}

	//nolint:gosec,lll
	cmd := exec.Command(dockerCmd, args...) // Run command "docker run -env-file .env -v "$(pwd):/app" --name (randomId) (randomId)"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	inputRunner, err := ru.InputResolver.Resolve(inputType)
	if err != nil {
		return nil, err
	}

	if err := inputRunner.Inputs(cmd, setup, flags); err != nil {
		return nil, err
	}

	if err := ru.setEnvs(cmd, dockerPwd, verbose); err != nil {
		return nil, err
	}

	return cmd, nil
}

func (ru RunManager) setEnvs(cmd *exec.Cmd, pwd string, verbose bool) error {
	envHolder, err := ru.env.Find()
	if err != nil {
		return err
	}

	dockerEnv := fmt.Sprintf(formula.EnvPattern, formula.DockerExecutionEnv, "true")
	pwdEnv := fmt.Sprintf(formula.EnvPattern, formula.PwdEnv, pwd)
	ctxEnv := fmt.Sprintf(formula.EnvPattern, formula.CtxEnv, envHolder.Current)
	env := fmt.Sprintf(formula.EnvPattern, formula.Env, envHolder.Current)
	verboseEnv := fmt.Sprintf(formula.EnvPattern, formula.VerboseEnv, strconv.FormatBool(verbose))
	cmd.Env = append(cmd.Env, pwdEnv, ctxEnv, verboseEnv, dockerEnv, env)

	envs := strings.Builder{}
	for _, e := range cmd.Env {
		envs.WriteString(e + "\n")
	}

	// Create a file named .env and add the environment variable inName=inValue
	if err := ioutil.WriteFile(envFile, []byte(envs.String()), os.ModePerm); err != nil {
		return err
	}

	return nil
}
