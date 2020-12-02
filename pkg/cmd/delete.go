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

// NewDeleteCmd create a new delete instance.
func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "delete SUBCOMMAND",
		Short:     "Delete env, repositories, formulas and workspaces",
		Long:      "Delete env, repositories, formulas and workspaces",
		Example:   "rit delete env",
		ValidArgs: []string{"env", "formula", "repo", "workspace", "credential"},
		Args:      cobra.OnlyValidArgs,
	}

	deprecatedMsg := fmt.Sprintf(
		`you can now use the "%v" command for the same purpose as the "%v" command.`,
		prompt.Bold("rit delete env"),
		prompt.Bold("rit delete context"),
	)

	DeprecateCmd(cmd, "context", deprecatedMsg)

	return cmd
}
