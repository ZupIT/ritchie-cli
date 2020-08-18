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
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

func TestNewTutorialCmd(t *testing.T) {
	cmd := NewTutorialCmd("path/any", inputListMock{}, TutorialFindSetterMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	if cmd == nil {
		t.Errorf("NewTutorialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewTutorialStdin(t *testing.T) {
	cmd := NewTutorialCmd("path/any", inputListMock{}, TutorialFindSetterMock{})
	cmd.PersistentFlags().Bool("stdin", true, "input by stdin")

	input := "{\"tutorial\": \"enabled\"}\n"
	newReader := strings.NewReader(input)
	cmd.SetIn(newReader)

	if cmd == nil {
		t.Errorf("NewTutorialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func Test_tutorialCmd_runAnyEntry(t *testing.T) {
	var tutorialHolderEnabled, tutorialHolderDisabled rtutorial.TutorialHolder
	type fields struct {
		prompt.InputList
		tutorial rtutorial.FindSetter
	}

	tutorialHolderEnabled.Current = "enabled"
	tutorialHolderDisabled.Current = "disabled"

	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		inputStdin string
	}{
		{
			name: "Run With Success when set tutorial enabled",
			fields: fields{
				InputList: inputListCustomMock{name: "enabled"},
				tutorial:  TutorialFindSetterMock{},
			},
			wantErr:    false,
			inputStdin: "{\"tutorial\": \"enabled\"}\n",
		},
		{
			name: "Run With Success when set tutorial disabled",
			fields: fields{
				InputList: inputListCustomMock{name: "disabled"},
				tutorial:  TutorialFindSetterMock{},
			},
			wantErr:    false,
			inputStdin: "{\"tutorial\": \"disabled\"}\n",
		},
		{
			name: "Return error when set return error",
			fields: fields{
				InputList: inputListCustomMock{name: "enabled"},
				tutorial: TutorialFindSetterCustomMock{
					set: func(tutorial string) (rtutorial.TutorialHolder, error) {
						return tutorialHolderEnabled, errors.New("some error")
					},
					find: func() (rtutorial.TutorialHolder, error) {
						return tutorialHolderEnabled, nil
					},
				},
			},
			wantErr:    true,
			inputStdin: "{\"tutorial\": \"enabled\"}\n",
		},
	}
	for _, tt := range tests {
		initPrompt := NewTutorialCmd("path/any", tt.fields.InputList, tt.fields.tutorial)
		initStdin := NewTutorialCmd("path/any", tt.fields.InputList, tt.fields.tutorial)

		initPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")
		initStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

		newReader := strings.NewReader(tt.inputStdin)
		initStdin.SetIn(newReader)

		if err := initPrompt.Execute(); (err != nil) != tt.wantErr {
			t.Errorf("init_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
		}
		if err := initStdin.Execute(); (err != nil) != tt.wantErr {
			t.Errorf("init_runStdin() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func Test_initCmd_runOnlyPrompt(t *testing.T) {
	var tutorialHolderEnabled, tutorialHolderDisabled rtutorial.TutorialHolder
	type fields struct {
		prompt.InputList
		tutorial rtutorial.FindSetter
	}

	tutorialHolderEnabled.Current = "enabled"
	tutorialHolderDisabled.Current = "disabled"

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Return error when find return error",
			fields: fields{
				InputList: inputListCustomMock{name: "enabled"},
				tutorial: TutorialFindSetterCustomMock{
					find: func() (rtutorial.TutorialHolder, error) {
						return tutorialHolderEnabled, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Return error when list return error",
			fields: fields{
				InputList: inputListErrorMock{},
				tutorial:  TutorialFindSetterMock{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		initPrompt := NewTutorialCmd("path/any", tt.fields.InputList, tt.fields.tutorial)

		initPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")

		if err := initPrompt.Execute(); (err != nil) != tt.wantErr {
			t.Errorf("init_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func Test_initCmd_runOnlyStdin(t *testing.T) {
	t.Run("Error when readJson returns err", func(t *testing.T) {
		wantErr := true

		initStdin := NewTutorialCmd("path/any", inputListMock{}, TutorialFindSetterMock{})

		initStdin.PersistentFlags().Bool("stdin", true, "input by stdin")
		newReader := strings.NewReader("{\"tutorial\": 1}\n")
		initStdin.SetIn(newReader)

		if err := initStdin.Execute(); (err != nil) != wantErr {
			t.Errorf("init_runStdin() error = %v, wantErr %v", err, wantErr)
		}
	})
}
