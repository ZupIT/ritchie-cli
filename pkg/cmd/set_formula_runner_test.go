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
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestSetFormulaRunnerCmd(t *testing.T) {
	tmpDir := os.TempDir()
	ritHome := filepath.Join(tmpDir, "runner")
	defer os.RemoveAll(ritHome)

	fileManager := stream.NewFileManager()
	configManager := runner.NewConfigManager(ritHome, fileManager)
	runnerFile := filepath.Join(ritHome, runner.FileName)

	var tests = []struct {
		name       string
		args       []string
		inputStdin string
		runner     string
		listErr    error
		err        error
	}{
		{
			name:   "success prompt set formula run",
			args:   []string{},
			runner: formula.LocalRun.String(),
		},
		{
			name:    "error on list",
			args:    []string{},
			listErr: errors.New("list error"),
			err:     errors.New("list error"),
		},
		{
			name: "success on flags",
			args: []string{"--runner=local"},
		},
		{
			name: "fail when missing flag",
			args: []string{"--runner="},
			err:  errors.New(missingFlagText(runnerFlagName)),
		},
		{
			name: "fail for wrong flag name",
			args: []string{"--runner=invalid"},
			err:  ErrInvalidRunType,
		},
		{
			name:       "success on stdin",
			args:       []string{},
			inputStdin: "{\"runType\": \"local\"}\n",
		},
		{
			name:       "fail with stdin wrong runner",
			args:       []string{},
			inputStdin: "{\"runType\": \"invalid\"}\n",
			err:        ErrInvalidRunType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Remove(runnerFile)
			inputList := &mocks.InputListMock{}
			inputList.On(
				"List", "Select a default formula run type", mock.Anything, mock.Anything,
			).Return(tt.runner, tt.listErr)

			cmd := NewSetFormulaRunnerCmd(configManager, inputList)
			cmd.SetArgs(tt.args)

			cmd.PersistentFlags().Bool("stdin", tt.inputStdin != "", "input by stdin")

			newReader := strings.NewReader(tt.inputStdin)
			cmd.SetIn(newReader)

			err := cmd.Execute()

			if err != nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, tt.err)
				assert.FileExists(t, runnerFile)
			}
		})
	}
}
