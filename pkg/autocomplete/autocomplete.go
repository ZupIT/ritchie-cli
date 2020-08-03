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

package autocomplete

import "github.com/spf13/cobra"

type (
	BashCommand struct {
		LastCommand string
		RootCommand string
		Commands    string
		Level       int
	}

	CompletionCommand struct {
		Content []string
		Before  string
	}
)

type Generator interface {
	Generate(s ShellName, cmd *cobra.Command) (string, error)
}

type ShellName string

func (s ShellName) String() string {
	return string(s)
}
