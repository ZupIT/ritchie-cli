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
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func Test_newUpdateRepoCmd(t *testing.T) {
	type in struct {
		config formula.ConfigRunner
		input  prompt.InputList
	}
	var tests = []struct {
		name       string
		in         in
		wantErr    bool
		inputStdin string
	}{
		{
			name: "success set formula run",
			in: in{
				input: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return formula.LocalRun.String(), nil
					},
				},
				config: ConfigRunnerMock{},
			},
			wantErr:    false,
			inputStdin: "{\"runType\": \"local\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newUpdateRepoPrompt := NewUpdateRepoCmd(tt.in.config, tt.in.input)
			newUpdateRepoStdin := NewUpdateRepoCmd(tt.in.config, tt.in.input)

			newUpdateRepoPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")
			newUpdateRepoStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

			newReader := strings.NewReader(tt.inputStdin)
			newUpdateRepoStdin.SetIn(newReader)

			if err := newUpdateRepoPrompt.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("new update repo type prompt command error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := newUpdateRepoStdin.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("new update repo type stdin command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
