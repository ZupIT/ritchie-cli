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
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const newCtx = "Type new context?"

// setContextCmd type for clean repo command
type setContextCmd struct {
	rcontext.FindSetter
	prompt.InputText
	prompt.InputList
}

// setContext type for stdin json decoder
type setContext struct {
	Context string `json:"context"`
}

func NewSetContextCmd(
	fs rcontext.FindSetter,
	it prompt.InputText,
	il prompt.InputList) *cobra.Command {
	s := setContextCmd{fs, it, il}

	cmd := &cobra.Command{
		Use:       "context",
		Short:     "Set context",
		Example:   "rit set context",
		RunE:      RunFuncE(s.runStdin(), s.runPrompt()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	cmd.LocalFlags()

	return cmd
}

func (s setContextCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctxHolder, err := s.Find()
		if err != nil {
			return err
		}

		ctxHolder.All = append(ctxHolder.All, rcontext.DefaultCtx)
		ctxHolder.All = append(ctxHolder.All, newCtx)
		ctx, err := s.List("All:", ctxHolder.All, "")
		if err != nil {
			return err
		}

		if ctx == newCtx {
			ctx, err = s.Text("New context: ", true, "")
			if err != nil {
				return err
			}
		}

		if _, err := s.Set(ctx); err != nil {
			return err
		}

		prompt.Success("Set context successful!")
		return nil
	}

}

func (s setContextCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		sc := setContext{}

		err := stdin.ReadJson(os.Stdin, &sc)
		if err != nil {
			return err
		}

		if _, err := s.Set(sc.Context); err != nil {
			return err
		}

		prompt.Success("Set context successful!")
		return nil
	}
}
