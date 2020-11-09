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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

// CommandRunnerFunc represents that runner func for commands.
type CommandRunnerFunc func(cmd *cobra.Command, args []string) error

// RunFuncE delegates to stdinFunc if --stdin flag is passed otherwise delegates to promptFunc.
func RunFuncE(stdinFunc, promptFunc CommandRunnerFunc) CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		s, err := cmd.Flags().GetBool(api.Stdin.ToLower())
		if err != nil {
			return err
		}

		exitsSdinEntry := stdin.ExistsEntry()

		inputStdinHasData := false
		if exitsSdinEntry {
			other := cmd.InOrStdin()
			inputStdinHasData = json.NewDecoder(other).More()
		}

		if s || inputStdinHasData {
			if s {
				fmt.Println("The flag --stdin is deprecated.\nIt's no longer needed for input via stdin.")
			}
			fmt.Println("s: ", s, "inputStdinHasData: ", inputStdinHasData)

			return stdinFunc(cmd, args)
		}
		fmt.Println(" -- prompt")
		return promptFunc(cmd, args)
	}
}
