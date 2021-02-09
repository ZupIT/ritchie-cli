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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

func TestTutorialCmd(t *testing.T) {
	ritPath := filepath.Join(os.TempDir(), "tutorial")
	_ = os.Mkdir(ritPath, os.ModePerm)
	defer os.RemoveAll(ritPath)

	tutorialFinder := rtutorial.NewFinder(ritPath)
	tutorialSetter := rtutorial.NewSetter(ritPath)
	tutorialFindSetter := rtutorial.NewFindSetter(tutorialFinder, tutorialSetter)

	tests := []struct {
		name    string
		ritHome string
		args    []string
		listErr error
		err     error
	}{
		{
			name:    "success prompt",
			ritHome: ritPath,
			args:    []string{},
		},
		{
			name:    "fail list",
			ritHome: ritPath,
			args:    []string{},
			listErr: errors.New("list error"),
			err:     errors.New("list error"),
		},
		{
			name:    "success flags",
			ritHome: ritPath,
			args:    []string{"--enabled=true"},
		},
		{
			name:    "fail empty flag",
			ritHome: ritPath,
			args:    []string{"--enabled="},
			err:     errors.New("invalid argument \"\" for \"--enabled\" flag: strconv.ParseBool: parsing \"\": invalid syntax"),
		},
		{
			name:    "fail set",
			ritHome: ritPath,
			args:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputList := &mocks.InputListMock{}
			inputList.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tutorialStatusEnabled, tt.listErr)
			cmd := NewTutorialCmd(inputList, tutorialFindSetter)

			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if err != nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.Nil(t, tt.err)
				assert.FileExists(t, filepath.Join(tt.ritHome, rtutorial.TutorialFile))
			}
		})
	}
}

func TestNewTutorialStdin(t *testing.T) {
	cmd := NewTutorialCmd(inputListMock{}, TutorialFindSetterMock{})
	cmd.PersistentFlags().Bool("stdin", true, "input by stdin")
	cmd.SetArgs([]string{})

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

func TestTutorialRunOnlyStdin(t *testing.T) {
	t.Run("Error when readJson returns err", func(t *testing.T) {
		wantErr := true

		cmdStdin := NewTutorialCmd(inputListMock{}, TutorialFindSetterMock{})
		cmdStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

		newReader := strings.NewReader("{\"tutorial\": 1}\n")
		cmdStdin.SetIn(newReader)

		if err := cmdStdin.Execute(); (err != nil) != wantErr {
			t.Errorf("cmd_runStdin() error = %v, wantErr %v", err, wantErr)
		}
	})
}
