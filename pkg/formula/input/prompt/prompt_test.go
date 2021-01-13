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

package prompt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
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
			"name": "sample_text_3",
			"type": "path",
			"label": "Type : ",
	"required": true
	},
	{
		"name": "sample_text_4",
		"type": "path",
		"label": "Type : ",
		"cache": {
				"active": true,
				"qty": 6,
				"newLabel": "Type new value. "
			},
	"tutorial": "Add a text for this field."
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
	_ = os.Setenv("SAMPLE_TEXT", "someValue")

	setup := formula.Setup{
		Config: formula.Config{
			Inputs: inputs,
		},
		FormulaPath: os.TempDir(),
	}
	fileManager := stream.NewFileManager()

	type in struct {
		iText          inputMock
		iTextValidator inputTextValidatorMock
		iTextDefault   inputTextDefaultMock
		iList          inputMock
		iBool          inputMock
		iPass          inputMock
		inType         api.TermInputType
		creResolver    credential.Resolver
		file           stream.FileWriteReadExister
	}

	tests := []struct {
		name          string
		in            in
		want          error
		expectedError string
	}{
		{
			name: "success prompt",
			in: in{
				iText:          inputMock{text: ""},
				iTextValidator: inputTextValidatorMock{},
				iTextDefault:   inputTextDefaultMock{"test", nil},
				iList:          inputMock{text: "Type new value?"},
				iBool:          inputMock{boolean: false},
				iPass:          inputMock{text: "******"},
				inType:         api.Prompt,
				creResolver:    envResolverMock{in: "test"},
				file:           fileManager,
			},
			want:          nil,
			expectedError: "",
		},
		{
			name: "error read file load items",
			in: in{
				iText:          inputMock{text: DefaultCacheNewLabel},
				iTextValidator: inputTextValidatorMock{},
				iTextDefault:   inputTextDefaultMock{"test", nil},
				iList:          inputMock{text: "test"},
				iBool:          inputMock{boolean: false},
				iPass:          inputMock{text: "******"},
				inType:         api.Prompt,
				creResolver:    envResolverMock{in: "test"},
				file:           fileManagerMock{rErr: errors.New("error to read file"), exist: true},
			},
			want:          errors.New(""),
			expectedError: "error to read file",
		},
		{
			name: "error unmarshal load items",
			in: in{
				iText:          inputMock{text: DefaultCacheNewLabel},
				iTextValidator: inputTextValidatorMock{},
				iTextDefault:   inputTextDefaultMock{"test", nil},
				iList:          inputMock{text: "test"},
				iBool:          inputMock{boolean: false},
				iPass:          inputMock{text: "******"},
				inType:         api.Prompt,
				creResolver:    envResolverMock{in: "test"},
				file:           fileManagerMock{rBytes: []byte("error"), exist: true},
			},
			want:          errors.New(""),
			expectedError: "invalid character 'e' looking for beginning of value",
		},
		{
			name: "cache file doesn't exist success",
			in: in{
				iText:          inputMock{text: DefaultCacheNewLabel},
				iTextValidator: inputTextValidatorMock{},
				iTextDefault:   inputTextDefaultMock{"test", nil},
				iList:          inputMock{text: "test"},
				iBool:          inputMock{boolean: false},
				iPass:          inputMock{text: "******"},
				inType:         api.Prompt,
				creResolver:    envResolverMock{in: "test"},
				file:           fileManagerMock{exist: false},
			},
			want:          nil,
			expectedError: "",
		},
		{
			name: "cache file doesn't exist error file write",
			in: in{
				iText:          inputMock{text: DefaultCacheNewLabel},
				iTextValidator: inputTextValidatorMock{},
				iTextDefault:   inputTextDefaultMock{"test", nil},
				iList:          inputMock{text: "test"},
				iBool:          inputMock{boolean: false},
				iPass:          inputMock{text: "******"},
				inType:         api.Prompt,
				creResolver:    envResolverMock{in: "test"},
				file:           fileManagerMock{wErr: errors.New("error to write file"), exist: false},
			},
			want:          errors.New(""),
			expectedError: "error to write file",
		},
		{
			name: "persist cache file write error",
			in: in{
				iText:          inputMock{text: DefaultCacheNewLabel},
				iTextValidator: inputTextValidatorMock{},
				iTextDefault:   inputTextDefaultMock{"test", nil},
				iList:          inputMock{text: "test"},
				iBool:          inputMock{boolean: false},
				iPass:          inputMock{text: "******"},
				inType:         api.Prompt,
				creResolver:    envResolverMock{in: "test"},
				file:           fileManagerMock{wErr: errors.New("error to write file"), rBytes: []byte(`["in_list1","in_list2"]`), exist: true},
			},
			want:          nil,
			expectedError: "",
		},
		{
			name: "error env resolver prompt",
			in: in{
				iText:          inputMock{text: DefaultCacheNewLabel},
				iTextValidator: inputTextValidatorMock{},
				iTextDefault:   inputTextDefaultMock{"test", nil},
				iList:          inputMock{text: "test"},
				iBool:          inputMock{boolean: false},
				iPass:          inputMock{text: "******"},
				inType:         api.Prompt,
				creResolver:    envResolverMock{in: "test", err: errors.New("credential not found")},
				file:           fileManager,
			},
			want:          errors.New(""),
			expectedError: "credential not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iText := tt.in.iText
			iTextValidator := tt.in.iTextValidator
			iTextDefault := tt.in.iTextDefault
			iList := tt.in.iList
			iBool := tt.in.iBool
			iPass := tt.in.iPass
			iMultiselect := inputMock{}
			inPath := InputPathMock{}

			inputManager := NewInputManager(
				tt.in.creResolver,
				tt.in.file,
				iList,
				iText,
				iTextValidator,
				iTextDefault,
				iBool,
				iPass,
				iMultiselect,
				inPath,
			)
			cmd := &exec.Cmd{}
			got := inputManager.Inputs(cmd, setup, nil)

			if tt.expectedError != "" {
				assert.EqualError(t, got, tt.expectedError)
			}

			if tt.want == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInputManager_ConditionalInputs(t *testing.T) {

	inputJson := `[
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
        "name": "sample_text",
        "type": "text",
        "label": "Type : ",
		"default": "test",
		"condition": {
			"variable": "%s",
			"operator": "%s",
			"value":    "in_list1"
		}
    }
]`

	fileManager := stream.NewFileManager()

	type in struct {
		variable string
		operator string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "equal conditional",
			in: in{
				variable: "sample_list",
				operator: "==",
			},
			want: nil,
		},
		{
			name: "not equal conditional",
			in: in{
				variable: "sample_list",
				operator: "!=",
			},
			want: nil,
		},
		{
			name: "greater than conditional",
			in: in{
				variable: "sample_list",
				operator: ">",
			},
			want: nil,
		},
		{
			name: "greater than or equal to conditional",
			in: in{
				variable: "sample_list",
				operator: ">=",
			},
			want: nil,
		},
		{
			name: "less than conditional",
			in: in{
				variable: "sample_list",
				operator: "<",
			},
			want: nil,
		},
		{
			name: "less than or equal to conditional",
			in: in{
				variable: "sample_list",
				operator: "<=",
			},
			want: nil,
		},
		{
			name: "wrong operator conditional",
			in: in{
				variable: "sample_list",
				operator: "eq",
			},
			want: errors.New("config.json: conditional operator eq not valid. Use any of (==, !=, >, >=, <, <=)"),
		},
		{
			name: "non-existing variable conditional",
			in: in{
				variable: "non_existing",
				operator: "==",
			},
			want: errors.New("config.json: conditional variable non_existing not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs []formula.Input
			_ = json.Unmarshal([]byte(fmt.Sprintf(inputJson, tt.in.variable, tt.in.operator)), &inputs)

			setup := formula.Setup{
				Config: formula.Config{
					Inputs: inputs,
				},
				FormulaPath: os.TempDir(),
			}

			iText := inputMock{text: DefaultCacheNewLabel}
			iTextValidator := inputTextValidatorMock{}
			iTextDefault := inputTextDefaultMock{}
			iList := inputMock{text: "in_list1"}
			iBool := inputMock{boolean: false}
			iPass := inputMock{text: "******"}
			iMultiselect := inputMock{}
			inPath := InputPathMock{}

			inputManager := NewInputManager(
				envResolverMock{},
				fileManager,
				iList,
				iText,
				iTextValidator,
				iTextDefault,
				iBool,
				iPass,
				iMultiselect,
				inPath,
			)

			cmd := &exec.Cmd{}
			got := inputManager.Inputs(cmd, setup, nil)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInputManager_RegexType(t *testing.T) {
	type in struct {
		inputJson      string
		inText         inputMock
		iTextValidator inputTextValidatorMock
		iTextDefault   inputTextDefaultMock
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "Success regex test",
			in: in{
				inputJson: `[
					    {
					        "name": "sample_text",
					        "type": "text",
									"label": "Type : ",
									"pattern": {
										"regex": "a|b",
										"mismatchText": "mismatch"
									}
					    }
					]`,
				inText:         inputMock{text: "a"},
				iTextValidator: inputTextValidatorMock{str: "a"},
			},
			want: nil,
		},
		{
			name: "Failed regex test",
			in: in{
				inputJson: `[
					    {
					        "name": "sample_text",
					        "type": "text",
									"label": "Type : ",
									"pattern": {
										"regex": "c|d",
										"mismatchText": "mismatch"
									}
					    }
					]`,
				inText:         inputMock{text: "a"},
				iTextValidator: inputTextValidatorMock{str: "a"},
			},
			want: errors.New("mismatch"),
		},
		{
			name: "Success regex test",
			in: in{
				inputJson: `[
					    {
					        "name": "sample_text",
					        "type": "text",
									"label": "Type : ",
									"pattern": {
										"regex": "abcc",
										"mismatchText": "mismatch"
									}
					    }
					]`,
				inText:         inputMock{text: "abcc"},
				iTextValidator: inputTextValidatorMock{str: "abcc"},
			},
			want: nil,
		},
	}

	fileManager := stream.NewFileManager()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs []formula.Input
			_ = json.Unmarshal([]byte(tt.in.inputJson), &inputs)

			setup := formula.Setup{
				Config: formula.Config{
					Inputs: inputs,
				},
				FormulaPath: os.TempDir(),
			}

			iText := tt.in.inText
			iTextValidator := tt.in.iTextValidator
			iTextDefault := tt.in.iTextDefault
			iList := inputMock{text: "in_list1"}
			iBool := inputMock{boolean: false}
			iPass := inputMock{text: "******"}
			iMultiselect := inputMock{}
			inPath := InputPathMock{}

			inputManager := NewInputManager(
				envResolverMock{},
				fileManager,
				iList,
				iText,
				iTextValidator,
				iTextDefault,
				iBool,
				iPass,
				iMultiselect,
				inPath,
			)

			cmd := &exec.Cmd{}
			got := inputManager.Inputs(cmd, setup, nil)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInputManager_DynamicInputs(t *testing.T) {
	type in struct {
		inputJson      string
		inText         inputMock
		iTextValidator inputTextValidatorMock
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "Success dynamic input test",
			in: in{
				inputJson: `[
					     {
      						"label": "Choose your repository ",
      						"name": "repo_list",
      						"type": "dynamic",
      						"requestInfo": {
      						  "url":"https://api.github.com/orgs/zupIt/repos",
      						  "jsonPath":"$..full_name"
      					 	}
    					}
					]`,
				inText:         inputMock{text: "a"},
				iTextValidator: inputTextValidatorMock{str: "a"},
			},
			want: nil,
		},
		{
			name: "fail dynamic input when http status is not ok",
			in: in{
				inputJson: `[
					     {
      						"label": "Choose your repository ",
      						"name": "repo_list",
      						"type": "dynamic",
      						"requestInfo": {
      						  "url":"https://github.com/ZupIT/ritchie-cli/issuesa",
      						  "jsonPath":"$..full_name"
      					 	}
    					}
					]`,
				inText:         inputMock{text: "a"},
				iTextValidator: inputTextValidatorMock{str: "a"},
			},
			want: errors.New("dynamic list request got http status 404 expecting some 2xx range"),
		},
		{
			name: "fail dynamic input when jsonpath is wrong",
			in: in{
				inputJson: `[
					     {
      						"label": "Choose your repository ",
      						"name": "repo_list",
      						"type": "dynamic",
      						"requestInfo": {
      						  "url":"https://api.github.com/orgs/ZupIT/repos",
      						  "jsonPath":"$.[*]full_name"
      					 	}
    					}
					]`,
				inText:         inputMock{text: "a"},
				iTextValidator: inputTextValidatorMock{str: "a"},
			},
			want: errors.New(`unexpected "[" while scanning JSON select expected Ident, "." or "*"`),
		},
		{
			name: "fail dynamic input when config.json url is empty",
			in: in{
				inputJson: `[
					     {
      						"label": "Choose your repository ",
      						"name": "repo_list",
      						"type": "dynamic",
      						"requestInfo": {
      						  "url":"",
      						  "jsonPath":"$.[*]full_name"
      					 	}
    					}
					]`,
				inText:         inputMock{text: "a"},
				iTextValidator: inputTextValidatorMock{str: "a"},
			},
			want: errors.New(`unsupported protocol scheme ""`),
		},
	}

	fileManager := stream.NewFileManager()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs []formula.Input
			_ = json.Unmarshal([]byte(tt.in.inputJson), &inputs)

			setup := formula.Setup{
				Config: formula.Config{
					Inputs: inputs,
				},
				FormulaPath: os.TempDir(),
			}

			iText := tt.in.inText
			iTextValidator := tt.in.iTextValidator
			iTextDefault := inputTextDefaultMock{}
			iList := inputMock{text: "in_list1"}
			iBool := inputMock{boolean: false}
			iPass := inputMock{text: "******"}
			iMultiselect := inputMock{}
			inPath := InputPathMock{}

			inputManager := NewInputManager(
				envResolverMock{},
				fileManager,
				iList,
				iText,
				iTextValidator,
				iTextDefault,
				iBool,
				iPass,
				iMultiselect,
				inPath,
			)

			cmd := &exec.Cmd{}
			got := inputManager.Inputs(cmd, setup, nil)

			if tt.want != nil {
				assert.NotNil(t, got)
			}

			if tt.want == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInputManager_Multiselect(t *testing.T) {
	type in struct {
		inputJSON     string
		inMultiselect inputMock
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success multiselect input test",
			in: in{
				inputJSON: `[
					{
						"name": "sample_multiselect",
						"type": "multiselect",
						"items": [
							"item_1",
							"item_2",
							"item_3",
							"item_4"
						],
						"label": "Choose one or more items: ",
						"required": true,
						"tutorial": "Select one or more items for this field."
					}
				]`,
				inMultiselect: inputMock{items: []string{"item_1", "item_2"}},
			},
			want: nil,
		},
		{
			name: "fail multiselect input test",
			in: in{
				inputJSON: `[
					{
						"name": "sample_multiselect",
						"type": "multiselect",
						"items": [],
						"label": "Choose one or more items: ",
						"required": true,
						"tutorial": "Select one or more items for this field."
					}
				]`,
				inMultiselect: inputMock{items: []string{}},
			},
			want: fmt.Errorf(EmptyItems, "sample_multiselect"),
		},
	}

	fileManager := stream.NewFileManager()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs []formula.Input
			_ = json.Unmarshal([]byte(tt.in.inputJSON), &inputs)

			setup := formula.Setup{
				Config: formula.Config{
					Inputs: inputs,
				},
				FormulaPath: os.TempDir(),
			}

			iText := tt.in.inMultiselect
			iTextValidator := inputTextValidatorMock{}
			iTextDefault := inputTextDefaultMock{}
			iList := inputMock{text: "in_list1"}
			iBool := inputMock{boolean: false}
			iPass := inputMock{text: "******"}
			iMultiselect := inputMock{items: []string{}}
			inPath := InputPathMock{}

			inputManager := NewInputManager(
				envResolverMock{},
				fileManager,
				iList,
				iText,
				iTextValidator,
				iTextDefault,
				iBool,
				iPass,
				iMultiselect,
				inPath,
			)

			cmd := &exec.Cmd{}
			got := inputManager.Inputs(cmd, setup, nil)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInputManager_DefaultFlag(t *testing.T) {
	inputJson := `[
		{
			"name": "sample_text",
			"type": "text",
			"label": "Type : ",
			"default": "test"
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
		inType      api.TermInputType
		creResolver credential.Resolver
		file        stream.FileWriteReadExister
	}

	tests := []struct {
		name string
		in   in
	}{
		{
			name: "success prompt",
			in: in{
				inType:      api.Prompt,
				creResolver: envResolverMock{in: "test"},
				file:        fileManager,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputManager := NewInputManager(
				tt.in.creResolver,
				tt.in.file,
				inputMock{},
				inputMock{},
				inputTextValidatorMock{},
				inputTextDefaultMock{},
				inputMock{},
				inputMock{},
				inputMock{},
				InputPathMock{},
			)

			cmd := &exec.Cmd{}
			flags := pflag.NewFlagSet("default", 0)
			flags.Bool("default", true, "default")

			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := inputManager.Inputs(cmd, setup, flags)

			_ = w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			assert.Nil(t, err)
			assert.Contains(t, string(out), "Added sample_text by default: test")
		})
	}
}

type inputTextValidatorMock struct {
	str string
}

func (i inputTextValidatorMock) Text(name string, validate func(interface{}) error, helper ...string) (string, error) {
	return i.str, validate(i.str)
}

type inputMock struct {
	text    string
	boolean bool
	items   []string
	err     error
}

func (i inputMock) List(string, []string, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Text(string, bool, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Bool(string, []string, ...string) (bool, error) {
	return i.boolean, i.err
}

func (i inputMock) Password(string, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Multiselect(formula.Input) ([]string, error) {
	return i.items, i.err
}

type inputTextDefaultMock struct {
	text string
	err  error
}

func (i inputTextDefaultMock) Text(formula.Input) (string, error) {
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
	rBytes   []byte
	rErr     error
	wErr     error
	aErr     error
	mErr     error
	rmErr    error
	lErr     error
	exist    bool
	listNews []string
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

func (fi fileManagerMock) Move(oldPath, newPath string, files []string) error {
	return fi.mErr
}

func (fi fileManagerMock) Remove(path string) error {
	return fi.rmErr
}

func (fi fileManagerMock) ListNews(oldPath, newPath string) ([]string, error) {
	return fi.listNews, fi.lErr
}

type InputPathMock struct{}

func (InputPathMock) Read(text string) (string, error) {
	return text, nil
}
