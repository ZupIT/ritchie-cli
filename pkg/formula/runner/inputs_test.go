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

package runner

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
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
        }
    },
 	{
        "name": "sample_text",
        "type": "text",
        "label": "Type : ",
		"default": "test"
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
        "label": "Pick your : "
    },
    {
        "name": "sample_bool",
        "type": "bool",
        "default": "false",
        "items": [
            "false",
            "true"
        ],
		"condition": {
			"variable": "sample_list",
			"operator": "==",
			"value": "in_list1"
		},
        "label": "Pick: "
    },
    {
        "name": "sample_password",
        "type": "password",
        "label": "Pick: "
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
		iText       inputMock
		iList       inputMock
		iBool       inputMock
		iPass       inputMock
		inType      api.TermInputType
		creResolver env.Resolvers
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
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Stdin,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				stdin:       `{"sample_text":"test_text","sample_list":"test_list","sample_bool": false}`,
				file:        fileManager,
			},
			want: nil,
		},
		{
			name: "error stdin",
			in: in{
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Stdin,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				stdin:       `"sample_text"`,
				file:        fileManager,
			},
			want: stdin.ErrInvalidInput,
		},
		{
			name: "success prompt",
			in: in{
				iText:       inputMock{text: ""},
				iList:       inputMock{text: "Type new value?"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Prompt,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				file:        fileManager,
			},
			want: nil,
		},
		{
			name: "success conditional prompt",
			in: in{
				iText:       inputMock{text: ""},
				iList:       inputMock{text: "in_list1"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Prompt,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				file:        fileManager,
			},
			want: nil,
		},
		{
			name: "error read file load items",
			in: in{
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Prompt,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				file:        fileManagerMock{rErr: errors.New("error to read file"), exist: true},
			},
			want: errors.New("error to read file"),
		},
		{
			name: "error unmarshal load items",
			in: in{
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Prompt,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				file:        fileManagerMock{rBytes: []byte("error"), exist: true},
			},
			want: errors.New("invalid character 'e' looking for beginning of value"),
		},
		{
			name: "cache file doesn't exist success",
			in: in{
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Prompt,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				file:        fileManagerMock{exist: false},
			},
			want: nil,
		},
		{
			name: "cache file doesn't exist error file write",
			in: in{
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Prompt,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				file:        fileManagerMock{wErr: errors.New("error to write file"), exist: false},
			},
			want: errors.New("error to write file"),
		},
		{
			name: "persist cache file write error",
			in: in{
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Prompt,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				file:        fileManagerMock{wErr: errors.New("error to write file"), rBytes: []byte(`["in_list1","in_list2"]`), exist: true},
			},
			want: nil,
		},
		{
			name: "error unknown prompt",
			in: in{
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.TermInputType(3),
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}},
				file:        fileManager,
			},
			want: ErrInputNotRecognized,
		},
		{
			name: "error env resolver prompt",
			in: in{
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Prompt,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test", err: errors.New("credential not found")}},
				file:        fileManager,
			},
			want: errors.New("credential not found"),
		},
		{
			name: "error env resolver stdin",
			in: in{
				iText:       inputMock{text: DefaultCacheNewLabel},
				iList:       inputMock{text: "test"},
				iBool:       inputMock{boolean: false},
				iPass:       inputMock{text: "******"},
				inType:      api.Stdin,
				stdin:       `{"sample_text":"test_text","sample_list":"test_list","sample_bool": false}`,
				creResolver: env.Resolvers{"CREDENTIAL": envResolverMock{in: "test", err: errors.New("credential not found")}},
				file:        fileManager,
			},
			want: errors.New("credential not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iText := tt.in.iText
			iList := tt.in.iList
			iBool := tt.in.iBool
			iPass := tt.in.iPass

			inputManager := NewInput(tt.in.creResolver, tt.in.file, iList, iText, iBool, iPass)

			cmd := &exec.Cmd{}
			if tt.in.inType == api.Stdin {
				cmd.Stdin = strings.NewReader(tt.in.stdin)
			}

			got := inputManager.Inputs(cmd, setup, tt.in.inType)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Inputs(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

type inputMock struct {
	text    string
	boolean bool
	err     error
}

func (i inputMock) List(string, []string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Text(string, bool, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Bool(string, []string) (bool, error) {
	return i.boolean, i.err
}

func (i inputMock) Password(string) (string, error) {
	return i.text, i.err
}

type envResolverMock struct {
	in  string
	err error
}

func (e envResolverMock) Resolve(string) (string, error) {
	return e.in, e.err
}

type fileManagerMock struct {
	rBytes []byte
	rErr   error
	wErr   error
	aErr   error
	exist  bool
}

func (fi fileManagerMock) Write(string, []byte) error {
	return fi.wErr
}

func (fi fileManagerMock) Read(string) ([]byte, error) {
	return fi.rBytes, fi.rErr
}

func (fi fileManagerMock) Exists(string) bool {
	return fi.exist
}

func (fi fileManagerMock) Append(path string, content []byte) error {
	return fi.aErr
}
