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

package stdin

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestInputManager_Inputs(t *testing.T) {

	inputJson := `[
    {
        "name": "sample_text",
        "type": "text",
        "label": "Type : ",
        "cache": {
            "active": true,
            "qty": 6,
            "newLabel": "Type new value. "
        },
		"tutorial": "Add a text for this field."
    },
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
]`
	var inputs []formula.Input
	_ = json.Unmarshal([]byte(inputJson), &inputs)

	setup := formula.Setup{
		Config: formula.Config{
			Inputs: inputs,
		},
		FormulaPath: os.TempDir(),
	}

	fileManager := stream.NewFileManager()

	type in struct {
		creResolver credential.Resolver
		file        stream.FileWriteReadExister
		stdin       string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success stdin",
			in: in{
				creResolver: envResolverMock{in: "test"},
				stdin:       `{"sample_text":"test_text","sample_list":"test_list","sample_bool": false}`,
				file:        fileManager,
			},
			want: nil,
		},
		{
			name: "error stdin",
			in: in{
				creResolver: envResolverMock{in: "test"},
				stdin:       `"sample_text"`,
				file:        fileManager,
			},
			want: stdin.ErrInvalidInput,
		},
		{
			name: "error env resolver stdin",
			in: in{
				stdin:       `{"sample_text":"test_text","sample_list":"test_list","sample_bool": false}`,
				creResolver: envResolverMock{in: "test", err: errors.New("credential not found")},
				file:        fileManager,
			},
			want: errors.New("credential not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			inputManager := NewInputManager(tt.in.creResolver)

			cmd := &exec.Cmd{}
			cmd.Stdin = strings.NewReader(tt.in.stdin)

			got := inputManager.Inputs(cmd, setup, nil)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Inputs(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

type envResolverMock struct {
	in  string
	err error
}

func (e envResolverMock) Resolve(string) (string, error) {
	return e.in, e.err
}
