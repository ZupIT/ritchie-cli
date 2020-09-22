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
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

// deleteContextCmd type for clean repo command
type deleteContextCmd struct {
	rcontext.FindRemover
	prompt.InputBool
	prompt.InputList
}

// deleteContext type for stdin json decoder
type deleteContext struct {
	Context string `json:"context"`
}

func NewDeleteContextCmd(
	fr rcontext.FindRemover,
	ib prompt.InputBool,
	il prompt.InputList) *cobra.Command {
	d := deleteContextCmd{fr, ib, il}

	cmd := &cobra.Command{
		Use:       "context",
		Short:     "Delete context for credentials",
		Example:   "rit delete context",
		RunE:      RunFuncE(d.runStdin(), d.runPrompt()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	cmd.LocalFlags()

	return cmd
}

func (d deleteContextCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctxHolder, err := d.Find()
		if err != nil {
			return err
		}

		if len(ctxHolder.All) <= 0 {
			prompt.Error("You have no defined contexts")
			return nil
		}

		for i := range ctxHolder.All {
			if ctxHolder.All[i] == ctxHolder.Current {
				ctxHolder.All[i] = fmt.Sprintf("%s%s", rcontext.CurrentCtx, ctxHolder.Current)
			}
		}

		ctx, err := d.List("Contexts:", ctxHolder.All)
		if err != nil {
			return err
		}

		if b, err := d.Bool("Are you sure want to delete this context?", []string{"yes", "no"}); err != nil {
			return err
		} else if !b {
			return nil
		}

		if _, err := d.Remove(ctx); err != nil {
			return err
		}

		prompt.Success("Delete context successful!")
		return nil
	}
}

func (d deleteContextCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctxHolder, err := d.Find()
		if err != nil {
			return err
		}

		if len(ctxHolder.All) <= 0 {
			prompt.Error("You have no defined contexts")
			return nil
		}

		dc := deleteContext{}

		err = stdin.ReadJson(os.Stdin, &dc)
		if err != nil {
			return err
		}

		if _, err := d.Remove(dc.Context); err != nil {
			return err
		}

		prompt.Success("Delete context successful!")
		return nil
	}
}
