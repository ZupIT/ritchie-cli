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
	"fmt"
	"reflect"

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	envFlagName        = "env"
	envFlagDescription = "Env name to delete"
)

var (
	NoDefinedEnvsMsg    = "You have no defined envs"
	DeleteEnvSuccessMsg = "Delete env successful!"
	deleteEnvFlags      = flags{
		{
			name:        envFlagName,
			kind:        reflect.String,
			defValue:    "",
			description: envFlagDescription,
		},
	}
)

// deleteEnvCmd type for clean repo command.
type deleteEnvCmd struct {
	env env.FindRemover
	prompt.InputBool
	prompt.InputList
}

func NewDeleteEnvCmd(
	fr env.FindRemover,
	ib prompt.InputBool,
	il prompt.InputList,
) *cobra.Command {
	d := deleteEnvCmd{fr, ib, il}

	cmd := &cobra.Command{
		Use:       "env",
		Short:     "Delete env for credentials",
		Example:   "rit delete env",
		RunE:      d.runFormula(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), deleteEnvFlags)

	cmd.LocalFlags()

	return cmd
}

func (d *deleteEnvCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		env, err := d.resolveInput(cmd)
		if err != nil {
			return err
		}

		if env == "" {
			return nil
		}

		if _, err := d.env.Remove(env); err != nil {
			return err
		}

		prompt.Success(DeleteEnvSuccessMsg)
		return nil
	}
}

func (d *deleteEnvCmd) resolveInput(cmd *cobra.Command) (string, error) {
	if IsFlagInput(cmd) {
		return d.resolveFlags(cmd)
	}
	return d.resolvePrompt()
}

func (d *deleteEnvCmd) resolvePrompt() (string, error) {
	envHolder, err := d.env.Find()
	if err != nil {
		return "", err
	}

	if len(envHolder.All) == 0 {
		prompt.Error(NoDefinedEnvsMsg)
		return "", nil
	}

	for i := range envHolder.All {
		if envHolder.All[i] == envHolder.Current {
			envHolder.All[i] = fmt.Sprintf("%s%s", env.Current, envHolder.Current)
		}
	}

	envName, err := d.List("Envs:", envHolder.All)
	if err != nil {
		return "", err
	}

	proceed, err := d.Bool("Are you sure want to delete this env?", []string{"yes", "no"})
	if err != nil {
		return "", err
	}

	if !proceed {
		return "", nil
	}

	return envName, nil
}

func (d *deleteEnvCmd) resolveFlags(cmd *cobra.Command) (string, error) {
	env, err := cmd.Flags().GetString(envFlagName)
	if err != nil {
		return "", err
	}

	if env == "" {
		return "", errors.New("please provide a value for 'env'")
	}

	return env, nil
}
