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

type showContextCmd struct {
	env.Finder
}

func NewShowContextCmd(f env.Finder) *cobra.Command {
	s := showContextCmd{f}

	return &cobra.Command{
		Use:       "context",
		Short:     "Show current context",
		Example:   "rit show context",
		RunE:      s.runFunc(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
}

func (s showContextCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, err := s.Find()
		if err != nil {
			return err
		}

		if ctx.Current == "" {
			ctx.Current = env.Default
		}

		prompt.Info(fmt.Sprintf("Current context: %s \n", ctx.Current))
		return nil
	}
}
