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

import "github.com/spf13/cobra"

const descListLong = `
This command consists of multiple subcommands to interact with ritchie.

It can be used to list repositories or credentials.
`

// NewListCmd create a new list instance
func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list SUBCOMMAND",
		Short:   "List repositories or credentials",
		Long:    descListLong,
		Example: "rit list repo, rit list credential",
	}
}
