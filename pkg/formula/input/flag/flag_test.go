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

package flag

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input"
)

func TestInputs(t *testing.T) {
	var inputs []formula.Input
	_ = json.Unmarshal([]byte(inputJson), &inputs)

	setup := formula.Setup{
		Config: formula.Config{
			Inputs: inputs,
		},
		FormulaPath: os.TempDir(),
	}

	type in struct {
		creResolver      env.Resolvers
		defaultFlagValue string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success flags",
			in: in{
				creResolver:      env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				defaultFlagValue: "text",
			},
			want: nil,
		},
		{
			name: "error flags empty",
			in: in{
				creResolver:     env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
			},
			want: errors.New("this flags cannot be empty [--sample_text_cache, --sample_text_2, --sample_password]"),
		},
		{
			name: "error env resolver",
			in: in{
				creResolver:     env.Resolvers{"CREDENTIAL": envResolverMock{in: "test", err: errors.New("credential not found")}},
				defaultFlagValue: "text",
			},
			want: errors.New("credential not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputManager := NewInputManager(tt.in.creResolver)

			cmd := &exec.Cmd{}
			flags := pflag.NewFlagSet("test", 0)

			for _, in := range inputs {
				switch in.Type {
				case input.TextType, input.PassType:
					flags.String(in.Name, tt.in.defaultFlagValue, in.Tutorial)
				case input.BoolType:
					flags.Bool(in.Name, false, in.Tutorial)
				}
			}

			got := inputManager.Inputs(cmd, setup, flags)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Inputs(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

type promptMock struct {
	err error
}

func (p promptMock) Inputs(cmd *exec.Cmd, setup formula.Setup, flags *pflag.FlagSet) error {
	return p.err
}

type envResolverMock struct {
	in  string
	err error
}

func (e envResolverMock) Resolve(string) (string, error) {
	return e.in, e.err
}

const inputJson = `[
    {
        "name": "sample_text_cache",
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
