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
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestInputManager(t *testing.T) {

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
        "name": "sample_list2",
        "type": "list",
        "default": "in1",
        "items": [
            "in_list1",
            "in_list2",
            "in_list3",
            "in_listN"
        ],
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
	ritHome := filepath.Join(os.TempDir(), "inputs")
	ritInvalidHome := filepath.Join(ritHome, "invalid")
	_ = os.Mkdir(ritHome, os.ModePerm)
	defer os.RemoveAll(ritHome)

	tests := []struct {
		name            string
		ritHome         string
		cacheContents   string
		credResolverErr error
		inputBoolErr    error
		inputPassErr    error
		inputTextErr    error
		inputListErr    error
		expectedError   string
	}{
		{
			name:    "success prompt",
			ritHome: ritHome,
		},
		{
			name:          "error unmarshal load items",
			ritHome:       ritHome,
			cacheContents: "error",
			expectedError: "invalid character 'e' looking for beginning of value",
		},
		{
			name:    "cache file doesn't exist success",
			ritHome: ritHome,
		},
		{
			name:          "persist cache file write error",
			ritHome:       ritInvalidHome,
			expectedError: mocks.FileNotFoundError(fmt.Sprintf(CachePattern, ritInvalidHome, strings.ToUpper("SAMPLE_TEXT"))),
		},
		{
			name:            "error env resolver prompt",
			ritHome:         ritHome,
			credResolverErr: errors.New("credential not found"),
			expectedError:   "credential not found",
		},
		{
			name:          "error input bool",
			ritHome:       ritHome,
			inputBoolErr:  errors.New("bool error"),
			expectedError: "bool error",
		},
		{
			name:          "error input pass",
			ritHome:       ritHome,
			inputPassErr:  errors.New("pass error"),
			expectedError: "pass error",
		},
		{
			name:          "error input text",
			ritHome:       ritHome,
			inputTextErr:  errors.New("text error"),
			expectedError: "text error",
		},
		{
			name:          "error input list",
			ritHome:       ritHome,
			inputListErr:  errors.New("list error"),
			expectedError: "list error",
		},
	}

	inPath := &mocks.InputPathMock{}
	inPath.On("Read", "Type : ").Return("", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheFile := fmt.Sprintf(CachePattern, tt.ritHome, strings.ToUpper("SAMPLE_TEXT"))
			if tt.cacheContents == "" {
				_ = os.Remove(cacheFile)
			} else {
				_ = ioutil.WriteFile(cacheFile, []byte(tt.cacheContents), os.ModePerm)
			}

			iText := &mocks.InputTextMock{}
			iText.On("Text", mock.Anything, mock.Anything, mock.Anything).Return("text value", nil)
			iTextValidator := &mocks.InputTextValidatorMock{}
			iTextValidator.On("Text", mock.Anything, mock.Anything, mock.Anything).Return("validator value", nil)
			iList := &mocks.InputListMock{}
			iList.On("List", mock.Anything, mock.Anything, mock.Anything).Return("list value", tt.inputListErr)
			iBool := &mocks.InputBoolMock{}
			iBool.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(true, tt.inputBoolErr)
			iPass := &mocks.InputPasswordMock{}
			iPass.On("Password", mock.Anything, mock.Anything, mock.Anything).Return("pass value", tt.inputPassErr)
			iMultiselect := &mocks.InputMultiselectMock{}
			iMultiselect.On("Multiselect", mock.Anything).Return([]string{"multiselect value"}, nil)
			iTextDefault := &mocks.InputDefaultTextMock{}
			iTextDefault.On("Text", mock.Anything, mock.Anything, mock.Anything).Return("default value", tt.inputTextErr)
			credResover := &mocks.CredResolverMock{}
			credResover.On("Resolve", mock.Anything).Return("resolver value", tt.credResolverErr)

			setup := formula.Setup{
				Config: formula.Config{
					Inputs: inputs,
				},
				FormulaPath: tt.ritHome,
			}

			inputManager := NewInputManager(
				credResover,
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

			if got != nil {
				assert.EqualError(t, got, tt.expectedError)
			} else {
				assert.Empty(t, tt.expectedError)
				expected := []string{
					"SAMPLE_TEXT=default value",
					"SAMPLE_TEXT=default value",
					"SAMPLE_TEXT_2=default value",
					"SAMPLE_LIST=list value",
					"SAMPLE_LIST2=list value",
					"SAMPLE_BOOL=true",
					"SAMPLE_PASSWORD=pass value",
					"TEST_RESOLVER=resolver value",
				}
				assert.Equal(t, expected, cmd.Env)
			}
		})
	}
}

func TestConditionalInputs(t *testing.T) {

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

	inPath := &mocks.InputPathMock{}
	inPath.On("Read", "Type : ").Return("", nil)

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

			iTextDefault := &mocks.InputDefaultTextMock{}
			iTextDefault.On("Text", mock.Anything, mock.Anything, mock.Anything).Return("default value", nil)
			iList := &mocks.InputListMock{}
			iList.On("List", mock.Anything, mock.Anything, mock.Anything).Return("list value", nil)

			inputManager := NewInputManager(
				&mocks.CredResolverMock{},
				iList,
				&mocks.InputTextMock{},
				&mocks.InputTextValidatorMock{},
				iTextDefault,
				&mocks.InputBoolMock{},
				&mocks.InputPasswordMock{},
				&mocks.InputMultiselectMock{},
				inPath,
			)

			cmd := &exec.Cmd{}
			got := inputManager.Inputs(cmd, setup, nil)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRegexType(t *testing.T) {
	tests := []struct {
		name      string
		inputJson string
		inText    string
		want      error
	}{
		{
			name: "Success regex test",
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
			inText: "a",
			want:   nil,
		},
		{
			name: "Failed regex test",
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
			inText: "a",
			want:   errors.New("mismatch"),
		},
		{
			name: "Success regex test",
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
			inText: "abcc",
			want:   nil,
		},
	}

	inPath := &mocks.InputPathMock{}
	inPath.On("Read", "Type : ").Return("", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs []formula.Input
			_ = json.Unmarshal([]byte(tt.inputJson), &inputs)

			setup := formula.Setup{
				Config: formula.Config{
					Inputs: inputs,
				},
				FormulaPath: os.TempDir(),
			}

			iText := &mocks.InputTextMock{}
			iText.On("Text", mock.Anything, mock.Anything, mock.Anything).Return(tt.inText, nil)
			iTextValidator := &mocks.InputTextValidatorMock{}
			iTextValidator.On("Text", mock.Anything, mock.Anything, mock.Anything).Return(tt.inText, nil)

			inputManager := NewInputManager(
				&mocks.CredResolverMock{},
				&mocks.InputListMock{},
				iText,
				iTextValidator,
				&mocks.InputDefaultTextMock{},
				&mocks.InputBoolMock{},
				&mocks.InputPasswordMock{},
				&mocks.InputMultiselectMock{},
				inPath,
			)

			cmd := &exec.Cmd{}
			got := inputManager.Inputs(cmd, setup, nil)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDynamicInputs(t *testing.T) {
	tests := []struct {
		name      string
		inputJson string
		want      error
	}{
		{
			name: "Success dynamic input test",
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
			want: nil,
		},
		{
			name: "fail dynamic input when http status is not ok",
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
			want: errors.New("dynamic list request got http status 404 expecting some 2xx range"),
		},
		{
			name: "fail dynamic input when jsonpath is wrong",
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
			want: errors.New(`unexpected "[" while scanning JSON select expected Ident, "." or "*"`),
		},
		{
			name: "fail dynamic input when config.json url is empty",
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
			want: errors.New(`unsupported protocol scheme ""`),
		},
	}

	inPath := &mocks.InputPathMock{}
	inPath.On("Read", "Type : ").Return("", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs []formula.Input
			_ = json.Unmarshal([]byte(tt.inputJson), &inputs)

			setup := formula.Setup{
				Config: formula.Config{
					Inputs: inputs,
				},
				FormulaPath: os.TempDir(),
			}

			iText := &mocks.InputTextMock{}
			iText.On("Text", mock.Anything, mock.Anything, mock.Anything).Return("a", nil)
			iTextValidator := &mocks.InputTextValidatorMock{}
			iTextValidator.On("Text", mock.Anything, mock.Anything, mock.Anything).Return("a", nil)
			iList := &mocks.InputListMock{}
			iList.On("List", mock.Anything, mock.Anything, mock.Anything).Return("list value", nil)

			inputManager := NewInputManager(
				&mocks.CredResolverMock{},
				iList,
				iText,
				iTextValidator,
				&mocks.InputDefaultTextMock{},
				&mocks.InputBoolMock{},
				&mocks.InputPasswordMock{},
				&mocks.InputMultiselectMock{},
				inPath,
			)

			cmd := &exec.Cmd{}
			got := inputManager.Inputs(cmd, setup, nil)

			if tt.want != nil {
				assert.NotNil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMultiselect(t *testing.T) {
	tests := []struct {
		name             string
		inputJSON        string
		multiselectValue []string
		want             error
	}{
		{
			name: "success multiselect input test",
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
			multiselectValue: []string{"item_1", "item_2"},
			want:             nil,
		},
		{
			name: "success multiselect input test when the required field is not sent",
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
						"tutorial": "Select one or more items for this field."
					}
				]`,
			multiselectValue: []string{"item_1", "item_2"},
			want:             nil,
		},
		{
			name: "fail multiselect input test",
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
			multiselectValue: []string{},
			want:             fmt.Errorf(EmptyItems, "sample_multiselect"),
		},
	}

	inPath := &mocks.InputPathMock{}
	inPath.On("Read", "Type : ").Return("", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs []formula.Input
			_ = json.Unmarshal([]byte(tt.inputJSON), &inputs)

			setup := formula.Setup{
				Config: formula.Config{
					Inputs: inputs,
				},
				FormulaPath: os.TempDir(),
			}

			iMultiselect := &mocks.InputMultiselectMock{}
			iMultiselect.On("Multiselect", mock.Anything).Return(tt.multiselectValue, nil)

			inputManager := NewInputManager(
				&mocks.CredResolverMock{},
				&mocks.InputListMock{},
				&mocks.InputTextMock{},
				&mocks.InputTextValidatorMock{},
				&mocks.InputDefaultTextMock{},
				&mocks.InputBoolMock{},
				&mocks.InputPasswordMock{},
				iMultiselect,
				inPath,
			)

			cmd := &exec.Cmd{}
			got := inputManager.Inputs(cmd, setup, nil)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultFlag(t *testing.T) {
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

	t.Run("success prompt", func(t *testing.T) {
		inputManager := NewInputManager(
			&mocks.CredResolverMock{},
			&mocks.InputListMock{},
			&mocks.InputTextMock{},
			&mocks.InputTextValidatorMock{},
			&mocks.InputDefaultTextMock{},
			&mocks.InputBoolMock{},
			&mocks.InputPasswordMock{},
			&mocks.InputMultiselectMock{},
			&mocks.InputPathMock{},
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

func TestEmptyList(t *testing.T) {
	inputJson := `[
		{
			"name": "sample_list",
			"type": "list",
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

	t.Run("success prompt", func(t *testing.T) {
		inputManager := NewInputManager(
			&mocks.CredResolverMock{},
			&mocks.InputListMock{},
			&mocks.InputTextMock{},
			&mocks.InputTextValidatorMock{},
			&mocks.InputDefaultTextMock{},
			&mocks.InputBoolMock{},
			&mocks.InputPasswordMock{},
			&mocks.InputMultiselectMock{},
			&mocks.InputPathMock{},
		)

		cmd := &exec.Cmd{}
		got := inputManager.Inputs(cmd, setup, nil)

		assert.Equal(t, fmt.Errorf(EmptyItems, "sample_list"), got)
	})
}
