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
)

func Test_showFormulaRunnerCmd_runPrompt(t *testing.T) {
	type in struct {
		config formula.ConfigRunner
	}
	var tests = []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "success show formula run type",
			in: in{
				config: ConfigRunnerMock{
					runType: formula.Local,
				},
			},
			wantErr: false,
		},
		{
			name: "error to find config",
			in: in{
				config: ConfigRunnerMock{
					findErr: errors.New("error to create config"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewShowFormulaRunnerCmd(tt.in.config)
			o.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := o.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("set credential command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
