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
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const stdinWarning = "stdin commands are deprecated and will no longer be supported in future versions. Please use" +
	" flags for programatic formula execution"

const (
	nameFlagName     = "name"
	providerFlagName = "provider"
)

type flag struct {
	name        string
	shortName   string
	kind        reflect.Kind
	defValue    interface{}
	description string
}

type flags []flag

// CommandRunnerFunc represents that runner func for commands.
type CommandRunnerFunc func(cmd *cobra.Command, args []string) error

func missingFlagText(flagName string) string {
	return fmt.Sprintf("please provide a value for '%s'", flagName)
}

func addReservedFlags(flags *pflag.FlagSet, flagsToAdd flags) {
	for _, flag := range flagsToAdd {
		switch flag.kind { //nolint:exhaustive
		case reflect.String:
			flags.StringP(flag.name, flag.shortName, flag.defValue.(string), flag.description)
		case reflect.Bool:
			flags.BoolP(flag.name, flag.shortName, flag.defValue.(bool), flag.description)
		case reflect.Int:
			flags.IntP(flag.name, flag.shortName, flag.defValue.(int), flag.description)
		case reflect.Slice:
			flags.StringSliceP(flag.name, flag.shortName, []string{}, flag.description)
		default:
			warning := fmt.Sprintf("The %q type is not supported for the %q flag", flag.kind.String(), flag.name)
			prompt.Warning(warning)
		}
	}
}

// RunFuncE delegates to stdinFunc if --stdin flag is passed otherwise delegates to promptFunc.
func RunFuncE(stdinFunc, promptFunc CommandRunnerFunc) CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		stdin, err := cmd.Flags().GetBool(api.Stdin.ToLower())
		if err != nil {
			return err
		}

		if stdin {
			prompt.Warning(stdinWarning)
			return stdinFunc(cmd, args)
		}
		return promptFunc(cmd, args)
	}
}

func IsFlagInput(cmd *cobra.Command) bool {
	return cmd.Flags().NFlag() > 0
}

func DeprecateCmd(parentCmd *cobra.Command, deprecatedCmd, deprecatedMsg string) {
	command := &cobra.Command{
		Use:        deprecatedCmd,
		Deprecated: deprecatedMsg,
	}
	parentCmd.AddCommand(command)
}
