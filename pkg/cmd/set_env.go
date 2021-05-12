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

package cmd

import (
	"errors"
	"os"
	"reflect"

	renv "github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const (
	newEnv                = "Type the new env?"
	successMsg            = "Set env successful!"
	setEnvFlagName        = "env"
	setEnvFlagDescription = "Env name to set"
)

var setEnvFlags = flags{
	{
		name:        setEnvFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: setEnvFlagDescription,
	},
}

// setEnvCmd type for clean repo command.
type setEnvCmd struct {
	env renv.FindSetter
	prompt.InputText
	prompt.InputList
}

// setEnv type for stdin json decoder.
type setEnv struct {
	Env string `json:"env"`
}

func NewSetEnvCmd(
	fs renv.FindSetter,
	it prompt.InputText,
	il prompt.InputList,
) *cobra.Command {
	s := setEnvCmd{fs, it, il}

	cmd := &cobra.Command{
		Use:       "env",
		Short:     "Set env",
		Example:   "rit set env",
		RunE:      RunFuncE(s.runStdin(), s.runCmd()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	cmd.LocalFlags()

	addReservedFlags(cmd.Flags(), setEnvFlags)

	return cmd
}

func (s setEnvCmd) runCmd() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		env, err := s.resolveInput(cmd)
		if err != nil {
			return err
		} else if env == "" {
			return nil
		}

		if _, err := s.env.Set(env); err != nil {
			return err
		}

		prompt.Success(successMsg)
		return nil
	}

}

func (s setEnvCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		sc := setEnv{}

		err := stdin.ReadJson(os.Stdin, &sc)
		if err != nil {
			return err
		}

		if _, err := s.env.Set(sc.Env); err != nil {
			return err
		}

		prompt.Success(successMsg)
		return nil
	}
}

func (s *setEnvCmd) resolveInput(cmd *cobra.Command) (string, error) {
	if IsFlagInput(cmd) {
		return s.resolveFlags(cmd)
	}
	return s.resolvePrompt()
}

func (s *setEnvCmd) resolvePrompt() (string, error) {
	envHolder, err := s.env.Find()
	if err != nil {
		return "", err
	}

	envHolder.All = append(envHolder.All, renv.Default)
	envHolder.All = append(envHolder.All, newEnv)
	env, err := s.List("All:", envHolder.All)
	if err != nil {
		return "", err
	}

	if env == newEnv {
		env, err = s.Text("New env: ", true)
		if err != nil {
			return "", err
		}
	}

	return env, nil
}

func (s *setEnvCmd) resolveFlags(cmd *cobra.Command) (string, error) {
	env, err := cmd.Flags().GetString(setEnvFlagName)
	if err != nil {
		return "", err
	} else if env == "" {
		return "", errors.New(missingFlagText(setEnvFlagName))
	}

	return env, nil
}
