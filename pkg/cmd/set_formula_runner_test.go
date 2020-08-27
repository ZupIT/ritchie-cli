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
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func Test_setFormulaRunnerCmd_runPrompt(t *testing.T) {
	type in struct {
		config formula.ConfigRunner
		input  prompt.InputList
	}
	var tests = []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "success set formula run",
			in: in{
				input: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return formula.Local.String(), nil
					},
				},
				config: ConfigRunnerMock{},
			},
			wantErr: false,
		},
		{
			name: "error to create config",
			in: in{
				input: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return formula.Local.String(), nil
					},
				},
				config: ConfigRunnerMock{
					createErr: errors.New("error to create config"),
				},
			},
			wantErr: true,
		},
		{
			name: "error to select run type",
			in: in{
				input: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return formula.Local.String(), errors.New("error to select run type")
					},
				},
				config: ConfigRunnerMock{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewSetFormulaRunnerCmd(tt.in.config, tt.in.input)
			o.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := o.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("set credential command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
