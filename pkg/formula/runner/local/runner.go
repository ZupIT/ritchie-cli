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

package local

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
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
		if err := ru.PostRun(setup, false); err != nil {
			prompt.Error(err.Error())
			return
		}
	}()

	formulaRun := filepath.Join(setup.TmpDir, setup.BinName)
	cmd := exec.Command(formulaRun)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := ru.setEnvs(cmd, setup.Pwd, verbose); err != nil {
		return err
	}

	inputRunner, err := ru.InputResolver.Resolve(inputType)
	if err != nil {
		return err
	}

	if err := inputRunner.Inputs(cmd, setup, flags); err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	metric.RepoName = def.RepoName

	return nil
}

func (ru RunManager) setEnvs(cmd *exec.Cmd, pwd string, verbose bool) error {
	ctx, err := ru.ctx.Find()
	if err != nil {
		return err
	}

	if ctx.Current != "" {
		prompt.Info(
			fmt.Sprintf("Formula running on context: %s\n", prompt.Cyan(ctx.Current)),
		)
	}

	cmd.Env = os.Environ()
	dockerEnv := fmt.Sprintf(formula.EnvPattern, formula.DockerExecutionEnv, "false")
	pwdEnv := fmt.Sprintf(formula.EnvPattern, formula.PwdEnv, pwd)
	ctxEnv := fmt.Sprintf(formula.EnvPattern, formula.CtxEnv, ctx.Current)
	verboseEnv := fmt.Sprintf(formula.EnvPattern, formula.VerboseEnv, strconv.FormatBool(verbose))
	cmd.Env = append(cmd.Env, pwdEnv, ctxEnv, verboseEnv, dockerEnv)

	return nil
}
