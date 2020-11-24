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

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func NewShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "show SUB_COMMAND",
		Short:     "Show env and formula-runner default",
		Long:      "Show current env and formula-runner default",
		Example:   "rit show env",
		ValidArgs: []string{"env", "formula-runner"},
		Args:      cobra.OnlyValidArgs,
	}

	deprecatedMsg := fmt.Sprintf(
		`you can now use the "%v" command for the same purpose as the "%v" command.`,
		prompt.Bold("rit show env"),
		prompt.Bold("rit show context"),
	)

	DeprecateCmd(cmd, "context", deprecatedMsg)

	return cmd
}
