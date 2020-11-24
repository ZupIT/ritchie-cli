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

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type showEnvCmd struct {
	env env.Finder
}

func NewShowEnvCmd(f env.Finder) *cobra.Command {
	s := showEnvCmd{f}

	return &cobra.Command{
		Use:       "env",
		Short:     "Show current env",
		Example:   "rit show env",
		RunE:      s.runFunc(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
}

func (s showEnvCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		envHolder, err := s.env.Find()
		if err != nil {
			return err
		}

		if envHolder.Current == "" {
			envHolder.Current = env.Default
		}

		fmt.Printf("Current env: %v\n", prompt.Bold(envHolder.Current))
		return nil
	}
}
