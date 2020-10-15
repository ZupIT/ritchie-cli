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
	"os"
	"os/exec"
	"strconv"

	"github.com/mattn/go-isatty"
	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	envFile   = ".env"
	dockerPwd = "/app"
)

var _ formula.Runner = RunManager{}

type RunManager struct {
	formula.PostRunner
	formula.InputResolver
	formula.PreRunner
	file    stream.FileWriteExistAppender
	ctx     rcontext.Finder
	homeDir string
}

func NewRunner(
	postRun formula.PostRunner,
	input formula.InputResolver,
	preRun formula.PreRunner,
	file stream.FileWriteExistAppender,
	ctx rcontext.Finder,
	homeDir string,
) formula.Runner {
	return RunManager{
		PostRunner:    postRun,
		InputResolver: input,
		PreRunner:     preRun,
		file:          file,
		ctx:           ctx,
		homeDir:       homeDir,
	}
}

func (ru RunManager) Run(def formula.Definition, inputType api.TermInputType, verbose bool, flags *pflag.FlagSet) error {
	setup, err := ru.PreRun(def)
	if err != nil {
		return err
	}

	defer func() {
		if err := ru.PostRun(setup, true); err != nil {
			prompt.Error(err.Error())
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
		args = []string{"run", "--rm", "-it", "--env-file", envFile, "-v", volume, "-v", homeDirVolume, "--name", setup.ContainerId, setup.ContainerId}
	} else {
		args = []string{"run", "--rm", "--env-file", envFile, "-v", volume, "-v", homeDirVolume, "--name", setup.ContainerId, setup.ContainerId}
	}

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
	ctx, err := ru.ctx.Find()
	if err != nil {
		return err
	}

	pwdEnv := fmt.Sprintf(formula.EnvPattern, formula.PwdEnv, pwd)
	ctxEnv := fmt.Sprintf(formula.EnvPattern, formula.CtxEnv, ctx.Current)
	verboseEnv := fmt.Sprintf(formula.EnvPattern, formula.VerboseEnv, strconv.FormatBool(verbose))
	cmd.Env = append(cmd.Env, pwdEnv, ctxEnv, verboseEnv)

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
