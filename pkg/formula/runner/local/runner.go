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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var _ formula.Runner = RunManager{}

type RunManager struct {
	formula.PostRunner
	formula.InputResolver
	formula.PreRunner
	file    stream.FileWriteExistAppender
	env     env.Finder
	homeDir string
}

func NewRunner(
	postRun formula.PostRunner,
	input formula.InputResolver,
	preRun formula.PreRunner,
	file stream.FileWriteExistAppender,
	env env.Finder,
	homeDir string,
) formula.Runner {
	return RunManager{
		PostRunner:    postRun,
		InputResolver: input,
		PreRunner:     preRun,
		file:          file,
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
		if err := ru.PostRun(setup, false); err != nil {
			prompt.Error(err.Error())
			return
		}
	}()

	formulaRun := filepath.Join(setup.TmpDir, setup.BinName)
	cmd := exec.Command(formulaRun)
	cmd.Stdin = os.Stdin
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

	out, _ := cmd.StdoutPipe()
	done := make(chan struct{})
	scanner := bufio.NewScanner(out)

	output := []string{}
	if out != nil {
		go func() {
			if scanner != nil {
				for scanner.Scan() {
					line := scanner.Text()

					if strings.Contains(line,
						"::output") {
						output = append(output, line)
					} else {
						fmt.Println(line)
					}
				}
			}

			done <- struct{}{}
		}()
	} else {
		close(done)
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	<-done

	sanitizeOutput := sanitizeData(output)
	flattenOutput := flattenData(sanitizeOutput)
	transformOutput := transformData(flattenOutput)
	if len(transformOutput) > 0 {
		output := filepath.Join(setup.BinPath, "output.json")
		testJson, _ := json.MarshalIndent(transformOutput, "", "\t")
		ioutil.WriteFile(output, testJson, os.ModePerm)
	}

	metric.RepoName = def.RepoName

	return nil
}

func (ru RunManager) setEnvs(cmd *exec.Cmd, pwd string, verbose bool) error {
	envHolder, err := ru.env.Find()
	if err != nil {
		return err
	}

	if envHolder.Current != "" {
		prompt.Info(
			fmt.Sprintf("Formula running on env: %s\n", prompt.Cyan(envHolder.Current)),
		)
	}

	cmd.Env = os.Environ()
	dockerEnv := fmt.Sprintf(formula.EnvPattern, formula.DockerExecutionEnv, "false")
	pwdEnv := fmt.Sprintf(formula.EnvPattern, formula.PwdEnv, pwd)
	ctxEnv := fmt.Sprintf(formula.EnvPattern, formula.CtxEnv, envHolder.Current)
	env := fmt.Sprintf(formula.EnvPattern, formula.Env, envHolder.Current)
	verboseEnv := fmt.Sprintf(formula.EnvPattern, formula.VerboseEnv, strconv.FormatBool(verbose))
	cmd.Env = append(cmd.Env, pwdEnv, ctxEnv, verboseEnv, dockerEnv, env)

	return nil
}

func sanitizeData(data []string) []string {
	var sanitizeData []string
	for i := range data {
		output := strings.Split(data[i], " ")[1:]
		newOutput := strings.Join(output, " ")
		sanitizeData = append(sanitizeData, newOutput)
	}

	return sanitizeData
}

func flattenData(data []string) []string {
	var flattenData []string
	for i := range data {
		element := strings.Split(data[i], " ")
		for j := range element {
			flattenData = append(flattenData, element[j])
		}
	}

	return flattenData
}

func transformData(data []string) map[string]string {
	transformData := make(map[string]string)
	for i := range data {
		test := strings.Split(data[i], "=")
		key := test[0]
		value := test[1]
		transformData[key] = value
	}

	return transformData
}
