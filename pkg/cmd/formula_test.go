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

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestFormulaCommand_Add(t *testing.T) {
	treeMock := treeMock{
		tree: formula.Tree{
			Commands: api.Commands{
				{
					Id:     "root_mock",
					Parent: "root",
					Usage:  "mock",
					Help:   "mock for add",
				},
				{
					Id:      "root_mock_test",
					Parent:  "root_mock",
					Usage:   "test",
					Help:    "test for add",
					Formula: true,
				},
			},
		},
	}

	type in struct {
		execMock FormulaExecutorMock
		args     []string
	}

	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "success default",
			in: in{
				execMock: FormulaExecutorMock{},
				args:     []string{"mock", "test"},
			},
		},
		{
			name: "success docker",
			in: in{
				execMock: FormulaExecutorMock{},
				args:     []string{"mock", "test", "--docker"},
			},
		},
		{
			name: "success local",
			in: in{
				execMock: FormulaExecutorMock{},
				args:     []string{"mock", "test", "--local"},
			},
		},
		{
			name: "success stdin",
			in: in{
				execMock: FormulaExecutorMock{},
				args:     []string{"mock", "test", "--stdin"},
			},
		},
		{
			name: "invalid flags",
			in: in{
				execMock: FormulaExecutorMock{},
				args:     []string{"mock", "test", "--local", "--docker"},
			},
			wantErr: true,
		},
		{
			name: "formula exec error",
			in: in{
				execMock: FormulaExecutorMock{err: errors.New("error to execute formula")},
				args:     []string{"mock", "test"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formulaCmd := NewFormulaCommand(api.CoreCmds, treeMock, tt.in.execMock, fileReaderMock{data: []byte(inputJson)})
			rootCmd := &cobra.Command{Use: "rit"}
			rootCmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			got := formulaCmd.Add(rootCmd)
			if got != nil {
				t.Errorf("Add got %v, want nil", got)
			}

			rootCmd.SetArgs(tt.in.args)

			if err := rootCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("furmula_exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

const inputJson = `{
  "inputs": [
    {
      "name": "sample_text",
      "type": "text",
      "label": "Type : ",
      "default": "test"
    },
    {
      "name": "sample_text_2",
      "type": "text",
      "label": "Type : ",
      "required": true
    },
    {
      "name": "sample_list",
      "type": "text",
      "default": "in1",
      "items": [
        "in_list1",
        "in_list2",
        "in_list3",
        "in_listN"
      ],
      "cache": {
        "active": true,
        "qty": 3,
        "newLabel": "Type new value?"
      },
      "label": "Pick your : ",
      "tutorial": "Select an item for this field."
    },
    {
      "name": "sample_bool",
      "type": "bool",
      "default": "false",
      "items": [
        "false",
        "true"
      ],
      "label": "Pick: ",
      "tutorial": "Select true or false for this field."
    },
    {
      "name": "sample_password",
      "type": "password",
      "label": "Pick: ",
      "tutorial": "Add a secret password for this field."
    },
    {
      "name": "test_resolver",
      "type": "CREDENTIAL_TEST"
    }
  ]
}`

type fileReaderMock struct {
	data []byte
	err  error
}

func (f fileReaderMock) Read(filepath string) ([]byte, error) {
	return f.data, f.err
}
