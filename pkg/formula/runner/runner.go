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

package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/mattn/go-isatty"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	dockerCmd = "docker"
	envFile   = ".env"
)

type RunManager struct {
	formula.PostRunner
	formula.InputRunner
	formula.PreRunner
	file stream.FileWriteExistAppender
}

func NewFormulaRunner(
	postRun formula.PostRunner,
	input formula.InputRunner,
	preRun formula.PreRunner,
	file stream.FileWriteExistAppender,
) formula.Runner {
	return RunManager{
		PostRunner:  postRun,
		InputRunner: input,
		PreRunner:   preRun,
		file:        file,
	}
}

func (ru RunManager) Run(def formula.Definition, inputType api.TermInputType, local bool) error {
	setup, err := ru.PreRun(def, local)
	if err != nil {
		return err
	}

	var isDocker bool
	var cmd *exec.Cmd
	if local || setup.ContainerId == "" {
		cmd, err = ru.runLocal(setup, inputType)
		if err != nil {
			return err
		}
	} else {
		cmd, err = ru.runDocker(setup, inputType)
		if err != nil {
			return err
		}

		isDocker = true
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := ru.PostRun(setup, isDocker); err != nil {
		return err
	}

	return nil
}

func (ru RunManager) runDocker(setup formula.Setup, inputType api.TermInputType) (*exec.Cmd, error) {
	volume := fmt.Sprintf("%s:/app", setup.Pwd)
	var args []string
	if isatty.IsTerminal(os.Stdout.Fd()) {
		args = []string{"run", "-it", "--env-file", envFile, "-v", volume, "--name", setup.ContainerId, setup.ContainerId}
	} else {
		args = []string{"run", "--env-file", envFile, "-v", volume, "--name", setup.ContainerId, setup.ContainerId}
	}

	cmd := exec.Command(dockerCmd, args...) // Run command "docker run -env-file .env -v "$(pwd):/app" --name (randomId) (randomId)"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := ru.Inputs(cmd, setup, inputType); err != nil {
		return nil, err
	}

	if err := ru.setDockerEnvs(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

func (ru RunManager) runLocal(setup formula.Setup, inputType api.TermInputType) (*exec.Cmd, error) {
	formulaRun := path.Join(setup.TmpDir, setup.BinName)
	cmd := exec.Command(formulaRun)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()
	pwdEnv := fmt.Sprintf(formula.EnvPattern, formula.PwdEnv, setup.Pwd)
	cPwdEnv := fmt.Sprintf(formula.EnvPattern, formula.CPwdEnv, setup.Pwd)
	cmd.Env = append(cmd.Env, pwdEnv)
	cmd.Env = append(cmd.Env, cPwdEnv)

	if err := ru.Inputs(cmd, setup, inputType); err != nil {
		return nil, err
	}

	return cmd, nil
}

func (ru RunManager) setDockerEnvs(cmd *exec.Cmd) error {
	pwdEnv := fmt.Sprintf(formula.EnvPattern, formula.PwdEnv, "/app")
	cPwdEnv := fmt.Sprintf(formula.EnvPattern, formula.CPwdEnv, "/app")
	cmd.Env = append(cmd.Env, pwdEnv)
	cmd.Env = append(cmd.Env, cPwdEnv)

	for _, e := range cmd.Env { // Create a file named .env and add the environment variable inName=inValue
		if !ru.file.Exists(envFile) {
			if err := ru.file.Write(envFile, []byte(e+"\n")); err != nil {
				return err
			}
			continue
		}
		if err := ru.file.Append(envFile, []byte(e+"\n")); err != nil {
			return err
		}
	}

	return nil
}
